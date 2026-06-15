package cassandra

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
	"streaming-golang/internal/domain/timeexpr"
)

const defaultYearsAdjustment = 20

type CassandraQueryBuilder struct {
	tableMappings map[string]string
	keyspace      string
}

func NewCassandraQueryBuilder(tableMappings map[string]string, keyspaces ...string) *CassandraQueryBuilder {
	keyspace := "ts"
	if len(keyspaces) > 0 && strings.TrimSpace(keyspaces[0]) != "" {
		keyspace = strings.TrimSpace(keyspaces[0])
	}
	return &CassandraQueryBuilder{tableMappings: tableMappings, keyspace: keyspace}
}

func (b *CassandraQueryBuilder) BuildQueries(_ context.Context, command domain.Command) ([]domain.ExecutableQuery, error) {
	if len(command.Mappings) == 0 {
		return nil, apperr.New(apperr.Invalid, "cannot build Cassandra query without mappings")
	}

	queries := make([]domain.ExecutableQuery, 0)
	for _, mapping := range command.Mappings {
		if mapping.Source != domain.SourceCassandra {
			continue
		}

		table, ok := b.resolveTable(mapping)
		if !ok {
			return nil, apperr.New(apperr.Invalid, fmt.Sprintf("no Cassandra table mapping for %q", mapping.DataCategory))
		}

		statement, arguments, skip, err := b.buildStatement(table, mapping, command)
		if err != nil {
			return nil, err
		}
		if skip {
			continue
		}

		queries = append(queries, domain.ExecutableQuery{
			ID:           mapping.ID,
			DataCategory: dataCategoryForQuery(command.DataCategory, mapping),
			Source:       domain.SourceCassandra,
			Statement:    statement,
			Arguments:    arguments,
			Parameters: map[string]any{
				"projection_columns": cassandraProjectionColumns(command.Columns),
			},
		})
	}

	return queries, nil
}

func dataCategoryForQuery(commandCategory domain.DataCategory, mapping domain.Mapping) domain.DataCategory {
	if mapping.DataCategory != "" {
		return mapping.DataCategory
	}
	return commandCategory
}

func (b *CassandraQueryBuilder) resolveTable(mapping domain.Mapping) (string, bool) {
	id := strings.ToLower(strings.TrimSpace(mapping.CassandraID))
	for key, table := range b.tableMappings {
		if strings.Contains(id, strings.ToLower(key)) && strings.TrimSpace(table) != "" {
			return b.qualifiedTable(table), true
		}
	}

	return "", false
}

func (b *CassandraQueryBuilder) qualifiedTable(table string) string {
	table = strings.TrimSpace(table)
	if strings.Contains(table, ".") || b.keyspace == "" {
		return table
	}
	return b.keyspace + "." + table
}

func (b *CassandraQueryBuilder) buildStatement(table string, mapping domain.Mapping, command domain.Command) (string, []any, bool, error) {
	if mapping.CassandraID == "" {
		return "", nil, false, apperr.New(apperr.Invalid, fmt.Sprintf("mapping %d has no Cassandra ID", mapping.ID))
	}

	quoteIndices := command.QuoteIndices
	forceNoRows := false
	if len(quoteIndices) == 0 {
		quoteIndices = []int{1}
		forceNoRows = true
	}

	arguments := []any{mapping.CassandraID, quoteIndices}
	where := []string{"ts_id = ?", "quote_index IN ?"}
	columns := "ts_id, qte_y, qte_m, qte_d, quote_index, publish_time, del_y, del_m, del_d, del_h, del_min, del_offset, value"

	deliveryCQL, deliveryArguments, noRows, err := buildDeliveryFilters(command.Filters.Nodes, cassandraTimeZone(mapping.ID), quoteIndices[0])
	if err != nil {
		return "", nil, false, err
	}
	if forceNoRows {
		where = append(where, "(del_y, del_m, del_d, del_h) = (?, ?, ?, ?)")
		arguments = append(arguments, int16(1), int8(1), int8(1), int8(0))
	} else if noRows {
		return "", nil, true, nil
	} else if !noRows && deliveryCQL != "" {
		where = append(where, deliveryCQL)
		arguments = append(arguments, deliveryArguments...)
	}

	statement := fmt.Sprintf("SELECT %s FROM %s WHERE %s", columns, table, strings.Join(where, " AND "))
	return statement, arguments, false, nil
}

func cassandraProjectionColumns(columns []string) []string {
	if len(columns) > 0 {
		return append([]string(nil), columns...)
	}
	return []string{
		"Identifier",
		"ReferenceTime",
		"DeliveryStart",
		"DeliveryEnd",
		"RelativeDeliveryPeriod",
		"Value",
		"LegacyDeliveryBucketNumber",
	}
}

type localHourWindow struct {
	lower          *time.Time
	lowerInclusive bool
	upper          *time.Time
	upperInclusive bool
}

func (w localHourWindow) empty() bool {
	if w.lower == nil || w.upper == nil {
		return false
	}
	if w.lower.After(*w.upper) {
		return true
	}
	return w.lower.Equal(*w.upper) && !(w.lowerInclusive && w.upperInclusive)
}

func (w localHourWindow) point() bool {
	return w.lower != nil && w.upper != nil && w.lower.Equal(*w.upper) && w.lowerInclusive && w.upperInclusive
}

func buildDeliveryFilters(nodes []domain.FilterNode, timezone string, quoteIndex int) (string, []any, bool, error) {
	if len(nodes) == 0 {
		return "", nil, false, nil
	}

	location, err := loadCassandraLocation(timezone)
	if err != nil {
		return "", nil, false, apperr.Wrap(apperr.Invalid, fmt.Sprintf("invalid Cassandra timezone %q", timezone), err)
	}

	deliveryWindow, hasDelivery, err := buildDeliveryWindow(nodes, location)
	if err != nil {
		return "", nil, false, err
	}
	rdpWindow, hasRDP, err := buildRDPWindow(nodes, quoteIndex)
	if err != nil {
		return "", nil, false, err
	}
	if !hasDelivery && !hasRDP {
		return "", nil, false, nil
	}

	var window localHourWindow
	switch {
	case hasDelivery && hasRDP:
		window = intersectWindows(deliveryWindow, rdpWindow)
	case hasDelivery:
		window = deliveryWindow
	default:
		window = rdpWindow
	}
	if window.empty() {
		return "", nil, true, nil
	}

	cql, args := emitTupleClause(window)
	return cql, args, false, nil
}

func buildDeliveryWindow(nodes []domain.FilterNode, location *time.Location) (localHourWindow, bool, error) {
	var window localHourWindow
	found := false
	for _, node := range nodes {
		filter, ok := node.(domain.ComparisonFilter)
		if !ok || !isDeliveryField(filter.Field) {
			continue
		}
		found = true
		if strings.EqualFold(filter.Operator, "in") && filter.Value.Kind == domain.FilterValueTimeInterval {
			start, end, ok, err := intervalBounds(filter.Value)
			if err != nil {
				return localHourWindow{}, false, err
			}
			if !ok {
				continue
			}
			start = deliveryLocalHour(start, location, filter.Field)
			end = deliveryLocalHour(end, location, filter.Field)
			tightenLower(&window, start, true)
			tightenUpper(&window, end, true)
			continue
		}

		point, ok, err := pointTime(filter.Value)
		if err != nil {
			return localHourWindow{}, false, err
		}
		if !ok {
			continue
		}
		local := deliveryLocalHour(point, location, filter.Field)
		applyWindowComparison(&window, filter.Operator, local)
	}

	if !found {
		return localHourWindow{}, false, nil
	}
	if window.lower == nil && window.upper != nil {
		lower := window.upper.AddDate(-defaultYearsAdjustment, 0, 0)
		window.lower = &lower
		window.lowerInclusive = true
	} else if window.upper == nil && window.lower != nil {
		upper := window.lower.AddDate(defaultYearsAdjustment, 0, 0)
		window.upper = &upper
		window.upperInclusive = true
	}
	return window, true, nil
}

func pointTime(value domain.FilterValue) (time.Time, bool, error) {
	switch value.Kind {
	case domain.FilterValuePointInTime:
		point, err := transactional.ParsePointTime(value.Raw, nil)
		if err != nil {
			return time.Time{}, false, apperr.Wrap(apperr.Invalid, fmt.Sprintf("invalid point-in-time value %q", value.Raw), err)
		}
		return point, true, nil
	case domain.FilterValueTimeIntervalPointTime:
		point, ok, err := intervalPointTime(value.Raw)
		if err != nil {
			return time.Time{}, false, err
		}
		return point, ok, nil
	default:
		return time.Time{}, false, nil
	}
}

func intervalBounds(value domain.FilterValue) (time.Time, time.Time, bool, error) {
	if value.Kind != domain.FilterValueTimeInterval {
		return time.Time{}, time.Time{}, false, nil
	}
	if value.Start != "" && value.End != "" {
		start, err := transactional.ParsePointTime(value.Start, nil)
		if err != nil {
			return time.Time{}, time.Time{}, false, apperr.Wrap(apperr.Invalid, fmt.Sprintf("invalid interval start %q", value.Start), err)
		}
		end, err := transactional.ParsePointTime(value.End, nil)
		if err != nil {
			return time.Time{}, time.Time{}, false, apperr.Wrap(apperr.Invalid, fmt.Sprintf("invalid interval end %q", value.End), err)
		}
		return start, end, true, nil
	}

	start, end, ok, err := intervalFunctionBounds(value.Raw)
	if err != nil || !ok {
		return start, end, ok, err
	}
	return start, end, true, nil
}

func intervalPointTime(raw string) (time.Time, bool, error) {
	name, args, ok := functionCall(raw)
	if !ok {
		return time.Time{}, false, nil
	}
	start, end, ok, err := intervalFunctionBounds(args)
	if err != nil || !ok {
		return time.Time{}, ok, err
	}
	switch strings.ToLower(name) {
	case "begin":
		return start, true, nil
	case "end":
		return end, true, nil
	default:
		return time.Time{}, false, nil
	}
}

func intervalFunctionBounds(raw string) (time.Time, time.Time, bool, error) {
	name, args, ok := functionCall(raw)
	if !ok {
		return time.Time{}, time.Time{}, false, nil
	}

	parts := splitArguments(args)
	if len(parts) == 0 {
		return time.Time{}, time.Time{}, false, nil
	}

	name = strings.ToLower(name)
	if name == "ti" {
		if len(parts) != 2 {
			return time.Time{}, time.Time{}, false, nil
		}
		start, err := transactional.ParsePointTime(parts[0], nil)
		if err != nil {
			return time.Time{}, time.Time{}, false, err
		}
		end, err := transactional.ParsePointTime(parts[1], nil)
		if err != nil {
			return time.Time{}, time.Time{}, false, err
		}
		return start, end, true, nil
	}

	start, err := transactional.ParsePointTime(parts[0], nil)
	if err != nil {
		return time.Time{}, time.Time{}, false, err
	}
	switch name {
	case "tiday", "gasdayeurope":
		return start, start, true, nil
	case "tiweek", "gasweekeurope":
		return start, start.AddDate(0, 0, 6), true, nil
	case "timonth", "gasmontheurope":
		return start, start.AddDate(0, 1, -1), true, nil
	case "tiquarter", "gasquartereurope":
		return start, start.AddDate(0, 3, -1), true, nil
	case "tiyear", "gasyeareurope":
		return start, start.AddDate(1, 0, -1), true, nil
	default:
		return time.Time{}, time.Time{}, false, nil
	}
}

func functionCall(raw string) (name, args string, ok bool) {
	raw = strings.TrimSpace(raw)
	open := strings.Index(raw, "(")
	if open <= 0 || !strings.HasSuffix(raw, ")") {
		return "", "", false
	}
	return raw[:open], raw[open+1 : len(raw)-1], true
}

func splitArguments(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func buildRDPWindow(nodes []domain.FilterNode, quoteIndex int) (localHourWindow, bool, error) {
	var window localHourWindow
	found := false
	reference := time.Date(quoteIndex/10000, time.Month((quoteIndex/100)%100), quoteIndex%100, 0, 0, 0, 0, time.UTC)
	for _, node := range nodes {
		filter, ok := node.(domain.ComparisonFilter)
		if !ok || !strings.EqualFold(filter.Field, "RelativeDeliveryPeriod") {
			continue
		}
		found = true
		hours, err := strconv.Atoi(strings.TrimSpace(filter.Value.Raw))
		if err != nil {
			return localHourWindow{}, false, apperr.Wrap(apperr.Invalid, fmt.Sprintf("invalid RelativeDeliveryPeriod value %q", filter.Value.Raw), err)
		}
		local := reference.Add(time.Duration(hours) * time.Hour)
		applyWindowComparison(&window, filter.Operator, local)
	}
	return window, found, nil
}

func isDeliveryField(field string) bool {
	return strings.EqualFold(field, "DeliveryStart") || strings.EqualFold(field, "DeliveryEnd")
}

func deliveryLocalHour(value time.Time, location *time.Location, field string) time.Time {
	local := value.UTC().In(location)
	local = time.Date(local.Year(), local.Month(), local.Day(), local.Hour(), 0, 0, 0, time.UTC)
	if strings.EqualFold(field, "DeliveryEnd") {
		local = local.Add(-time.Hour)
	}
	return local
}

func applyWindowComparison(window *localHourWindow, operator string, value time.Time) {
	switch strings.ToLower(operator) {
	case "=":
		tightenLower(window, value, true)
		tightenUpper(window, value, true)
	case ">=":
		tightenLower(window, value, true)
	case ">":
		tightenLower(window, value, false)
	case "<=":
		tightenUpper(window, value, true)
	case "<":
		tightenUpper(window, value, false)
	}
}

func tightenLower(window *localHourWindow, value time.Time, inclusive bool) {
	if window.lower == nil || value.After(*window.lower) {
		window.lower = &value
		window.lowerInclusive = inclusive
		return
	}
	if value.Equal(*window.lower) {
		window.lowerInclusive = window.lowerInclusive && inclusive
	}
}

func tightenUpper(window *localHourWindow, value time.Time, inclusive bool) {
	if window.upper == nil || value.Before(*window.upper) {
		window.upper = &value
		window.upperInclusive = inclusive
		return
	}
	if value.Equal(*window.upper) {
		window.upperInclusive = window.upperInclusive && inclusive
	}
}

func intersectWindows(left, right localHourWindow) localHourWindow {
	var window localHourWindow
	if left.lower != nil {
		tightenLower(&window, *left.lower, left.lowerInclusive)
	}
	if right.lower != nil {
		tightenLower(&window, *right.lower, right.lowerInclusive)
	}
	if left.upper != nil {
		tightenUpper(&window, *left.upper, left.upperInclusive)
	}
	if right.upper != nil {
		tightenUpper(&window, *right.upper, right.upperInclusive)
	}
	return window
}

func emitTupleClause(window localHourWindow) (string, []any) {
	if window.point() {
		return "(del_y, del_m, del_d, del_h) = (?, ?, ?, ?)", tupleArgs(*window.lower)
	}

	parts := make([]string, 0, 2)
	args := make([]any, 0, 8)
	if window.lower != nil {
		operator := ">="
		if !window.lowerInclusive {
			operator = ">"
		}
		parts = append(parts, fmt.Sprintf("(del_y, del_m, del_d, del_h) %s (?, ?, ?, ?)", operator))
		args = append(args, tupleArgs(*window.lower)...)
	}
	if window.upper != nil {
		operator := "<="
		if !window.upperInclusive {
			operator = "<"
		}
		parts = append(parts, fmt.Sprintf("(del_y, del_m, del_d, del_h) %s (?, ?, ?, ?)", operator))
		args = append(args, tupleArgs(*window.upper)...)
	}
	return strings.Join(parts, " AND "), args
}

func tupleArgs(value time.Time) []any {
	return []any{int16(value.Year()), int8(value.Month()), int8(value.Day()), int8(value.Hour())}
}

func cassandraTimeZone(id domain.Identifier) string {
	switch id {
	case 536958751, 536959001, 536959251, 536959501:
		return "Australia/Sydney"
	case 536960251:
		return "Asia/Singapore"
	case 536959751, 536960001:
		return "Asia/Tokyo"
	case 537085751, 537119501:
		return "Pacific/Auckland"
	default:
		return "Europe/Zurich"
	}
}

func loadCassandraLocation(name string) (*time.Location, error) {
	return timeexpr.LoadLocation(name)
}
