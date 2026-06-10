package mssql

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/domain"
)

const (
	cmdpIdentifierColumn       = "TimeSeries_FID"
	hyperscaleIdentifierColumn = "MdoId"
	hyperscaleDeletedColumn    = "Deleted"
)

type CMDPQueryBuilder struct{}

func NewCMDPQueryBuilder() CMDPQueryBuilder {
	return CMDPQueryBuilder{}
}

type HyperscaleQueryBuilder struct{}

func NewHyperscaleQueryBuilder() HyperscaleQueryBuilder {
	return HyperscaleQueryBuilder{}
}

func (CMDPQueryBuilder) BuildQueries(_ context.Context, command domain.Command) ([]domain.ExecutableQuery, error) {
	mappings := command.Mappings
	if len(mappings) == 0 {
		return nil, apperr.New(apperr.Invalid, "cannot build CMDP query without mappings")
	}

	queries := make([]domain.ExecutableQuery, 0, len(mappings))
	for _, mapping := range mappings {
		if mapping.Source != "" && mapping.Source != domain.SourceCMDP {
			continue
		}

		statement, parameters, err := buildCMDPStatement(mapping, command.Filters, command.IndexRange, command.Columns)
		if err != nil {
			return nil, err
		}
		queries = append(queries, domain.ExecutableQuery{
			ID:           mapping.ID,
			DataCategory: dataCategoryForQuery(command.DataCategory, mapping),
			Source:       domain.SourceCMDP,
			Filters:      command.Filters,
			IndexRange:   command.IndexRange,
			Statement:    statement,
			Parameters:   parameters,
		})
	}

	if len(queries) == 0 {
		return nil, nil
	}
	return queries, nil
}

func (HyperscaleQueryBuilder) BuildQueries(_ context.Context, command domain.Command) ([]domain.ExecutableQuery, error) {
	mappings := command.Mappings
	if len(mappings) == 0 {
		return nil, apperr.New(apperr.Invalid, "cannot build hyperscale query without mappings")
	}

	queries := make([]domain.ExecutableQuery, 0, len(mappings))
	for _, mapping := range mappings {
		if mapping.Source != domain.SourceHyperscale {
			continue
		}

		statement, parameters, err := buildHyperscaleStatement(mapping, command.Filters, command.Columns)
		if err != nil {
			return nil, err
		}
		queries = append(queries, domain.ExecutableQuery{
			ID:           mapping.ID,
			DataCategory: dataCategoryForQuery(command.DataCategory, mapping),
			Source:       domain.SourceHyperscale,
			Filters:      command.Filters,
			IndexRange:   command.IndexRange,
			Statement:    statement,
			Parameters:   parameters,
		})
	}

	if len(queries) == 0 {
		return nil, nil
	}
	return queries, nil
}

func dataCategoryForQuery(commandCategory domain.DataCategory, mapping domain.Mapping) domain.DataCategory {
	if mapping.DataCategory != "" {
		return mapping.DataCategory
	}
	return commandCategory
}

func buildCMDPStatement(mapping domain.Mapping, filters domain.FilterSet, indexRange *domain.IndexRange, requestedColumns []string) (string, map[string]any, error) {
	if strings.TrimSpace(mapping.ViewName) == "" {
		return "", nil, apperr.New(apperr.Invalid, fmt.Sprintf("mapping %d has no CMDP view name", mapping.ID))
	}

	builder := sqlBuilder{
		mapping:    mapping,
		parameters: make(map[string]any),
	}
	builder.addParameter("id", int64(mapping.ID))

	where := []string{fmt.Sprintf("%s = @id", qualify(cmdpIdentifierColumn))}
	filterPredicates, err := builder.filterPredicates(filters.Nodes)
	if err != nil {
		return "", nil, err
	}
	where = append(where, filterPredicates...)
	if indexRange != nil && strings.TrimSpace(mapping.IndexField) != "" {
		builder.addParameter("indexStart", indexRange.Start)
		builder.addParameter("indexEnd", indexRange.End)
		where = append(where,
			fmt.Sprintf("%s >= @indexStart", qualify(mapping.IndexField)),
			fmt.Sprintf("%s <= @indexEnd", qualify(mapping.IndexField)),
		)
	}

	statement := fmt.Sprintf("SELECT %s FROM %s AS [d] WHERE %s",
		strings.Join(selectColumns(mapping.Columns, requestedColumns), ", "),
		quoteTable(mapping.ViewName),
		strings.Join(where, " AND "),
	)

	if order := orderColumns(mapping.Columns); len(order) > 0 {
		statement += " ORDER BY " + strings.Join(order, ", ")
	}
	return statement, builder.parameters, nil
}

func buildHyperscaleStatement(mapping domain.Mapping, filters domain.FilterSet, requestedColumns []string) (string, map[string]any, error) {
	viewName, err := hyperscaleViewName(mapping)
	if err != nil {
		return "", nil, err
	}
	valueColumn, err := hyperscaleValueColumn(mapping.DataCategory)
	if err != nil {
		return "", nil, err
	}

	builder := sqlBuilder{
		mapping:         mapping,
		parameters:      make(map[string]any),
		jsonValueColumn: valueColumn,
	}
	builder.addParameter("id", int64(mapping.ID))

	where := []string{fmt.Sprintf("%s = @id", qualify(hyperscaleIdentifierColumn))}
	filterPredicates, err := builder.filterPredicates(filters.Nodes)
	if err != nil {
		return "", nil, err
	}
	where = append(where, filterPredicates...)
	where = append(where, fmt.Sprintf("%s = 0", qualify(hyperscaleDeletedColumn)))

	statement := fmt.Sprintf("SELECT %s FROM %s AS [d] WHERE %s",
		strings.Join(hyperscaleSelectColumns(mapping.Columns, requestedColumns, valueColumn), ", "),
		quoteTable(viewName),
		strings.Join(where, " AND "),
	)

	if order := hyperscaleOrderColumns(mapping.Columns, valueColumn); len(order) > 0 {
		statement += " ORDER BY " + strings.Join(order, ", ")
	}
	return statement, builder.parameters, nil
}

type sqlBuilder struct {
	mapping         domain.Mapping
	parameters      map[string]any
	nextParam       int
	jsonValueColumn string
}

func (b *sqlBuilder) filterPredicates(nodes []domain.FilterNode) ([]string, error) {
	predicates := make([]string, 0, len(nodes))
	for _, node := range nodes {
		switch filter := node.(type) {
		case domain.ComparisonFilter:
			predicate, err := b.comparisonPredicate(filter)
			if err != nil {
				return nil, err
			}
			if predicate != "" {
				predicates = append(predicates, predicate)
			}
		case domain.RankOverFilter:
			return nil, apperr.New(apperr.Invalid, "rankover filters are not supported by the CMDP SQL builder yet")
		}
	}
	return predicates, nil
}

func (b *sqlBuilder) comparisonPredicate(filter domain.ComparisonFilter) (string, error) {
	column, ok := b.columnByMDSName(filter.Field)
	if !ok {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("filter field %q is not mapped for CMDP view %q", filter.Field, b.mapping.ViewName))
	}

	switch {
	case strings.EqualFold(filter.Operator, "in"):
		return b.intervalPredicate(column, filter.Value)
	case isComparisonOperator(filter.Operator):
		return b.scalarPredicate(column, filter.Operator, filter.Value)
	default:
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("unsupported filter operator %q", filter.Operator))
	}
}

func (b *sqlBuilder) intervalPredicate(column domain.ColumnMapping, value domain.FilterValue) (string, error) {
	if value.Kind != domain.FilterValueTimeInterval {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("field %q uses IN with unsupported value %q", column.MDSName, value.Raw))
	}

	start, end, ok, err := sqlIntervalBounds(value)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("cannot convert interval %q into CMDP SQL bounds", value.Raw))
	}

	startParam := b.nextParameter(start)
	endParam := b.nextParameter(end)
	columnExpression := b.columnExpression(column)
	return fmt.Sprintf("(%s >= @%s AND %s <= @%s)", columnExpression, startParam, columnExpression, endParam), nil
}

func (b *sqlBuilder) scalarPredicate(column domain.ColumnMapping, operator string, value domain.FilterValue) (string, error) {
	if value.Kind == domain.FilterValueLatest || value.Kind == domain.FilterValueLatestGlobal {
		return "", apperr.New(apperr.Invalid, "latest filters are not supported by the CMDP SQL builder yet")
	}
	if value.Kind == domain.FilterValueTimeIntervalPointTime {
		point, ok, err := sqlIntervalPointTime(value.Raw)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", apperr.New(apperr.Invalid, fmt.Sprintf("cannot convert interval point %q into CMDP SQL", value.Raw))
		}
		param := b.nextParameter(point)
		return fmt.Sprintf("%s %s @%s", b.columnExpression(column), operator, param), nil
	}

	paramValue, err := sqlScalarValue(value)
	if err != nil {
		return "", err
	}
	param := b.nextParameter(paramValue)
	return fmt.Sprintf("%s %s @%s", b.columnExpression(column), operator, param), nil
}

func (b *sqlBuilder) columnExpression(column domain.ColumnMapping) string {
	if b.jsonValueColumn != "" && !column.IsKey {
		return hyperscaleJSONValueExpression(b.jsonValueColumn, firstNonEmpty(column.MDSName, column.SourceName), column.DataType)
	}
	return qualify(column.SourceName)
}

func (b *sqlBuilder) columnByMDSName(name string) (domain.ColumnMapping, bool) {
	for _, column := range b.mapping.Columns {
		if strings.EqualFold(column.MDSName, name) || strings.EqualFold(column.SourceName, name) {
			if strings.TrimSpace(column.SourceName) == "" {
				column.SourceName = column.MDSName
			}
			return column, true
		}
	}
	return domain.ColumnMapping{}, false
}

func (b *sqlBuilder) nextParameter(value any) string {
	name := fmt.Sprintf("p%d", b.nextParam)
	b.nextParam++
	b.addParameter(name, value)
	return name
}

func (b *sqlBuilder) addParameter(name string, value any) {
	b.parameters[name] = value
}

func selectColumns(columns []domain.ColumnMapping, requestedColumns []string) []string {
	requested := requestedColumnSet(requestedColumns)
	selected := make([]domain.ColumnMapping, 0, len(columns))
	for _, column := range columns {
		if strings.TrimSpace(column.SourceName) == "" {
			continue
		}
		if len(requested) > 0 && !column.IsKey && !isRequestedColumn(column, requested) {
			continue
		}
		if column.IsKey || column.IsProjectable {
			selected = append(selected, column)
		}
	}
	if len(selected) == 0 {
		return []string{"[d].*"}
	}

	sort.SliceStable(selected, func(i, j int) bool {
		left := columnSortValue(selected[i])
		right := columnSortValue(selected[j])
		if left == right {
			return selected[i].MDSName < selected[j].MDSName
		}
		return left < right
	})

	seen := make(map[string]struct{}, len(selected))
	expressions := make([]string, 0, len(selected))
	for _, column := range selected {
		key := strings.ToLower(column.SourceName)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		expression := qualify(column.SourceName)
		if column.MDSName != "" && !strings.EqualFold(column.MDSName, column.SourceName) {
			expression += " AS " + quoteIdentifier(column.MDSName)
		}
		expressions = append(expressions, expression)
	}
	return expressions
}

func hyperscaleSelectColumns(columns []domain.ColumnMapping, requestedColumns []string, valueColumn string) []string {
	requested := requestedColumnSet(requestedColumns)
	selected := make([]domain.ColumnMapping, 0, len(columns))
	for _, column := range columns {
		if strings.TrimSpace(column.SourceName) == "" && strings.TrimSpace(column.MDSName) == "" {
			continue
		}
		if len(requested) > 0 && !column.IsKey && !isRequestedColumn(column, requested) {
			continue
		}
		if column.IsKey || column.IsProjectable {
			selected = append(selected, column)
		}
	}
	if len(selected) == 0 {
		return []string{"[d].*"}
	}

	sort.SliceStable(selected, func(i, j int) bool {
		left := columnSortValue(selected[i])
		right := columnSortValue(selected[j])
		if left == right {
			return selected[i].MDSName < selected[j].MDSName
		}
		return left < right
	})

	seen := make(map[string]struct{}, len(selected))
	expressions := make([]string, 0, len(selected))
	for _, column := range selected {
		outputName := firstNonEmpty(column.MDSName, column.SourceName)
		key := strings.ToLower(outputName)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		var expression string
		if column.IsKey {
			sourceName := firstNonEmpty(column.SourceName, column.MDSName)
			expression = qualify(sourceName)
			if outputName != "" && !strings.EqualFold(outputName, sourceName) {
				expression += " AS " + quoteIdentifier(outputName)
			}
		} else {
			expression = hyperscaleJSONValueExpression(valueColumn, outputName, column.DataType)
			if outputName != "" {
				expression += " AS " + quoteIdentifier(outputName)
			}
		}
		expressions = append(expressions, expression)
	}
	return expressions
}

func requestedColumnSet(columns []string) map[string]struct{} {
	if len(columns) == 0 {
		return nil
	}
	requested := make(map[string]struct{}, len(columns))
	for _, column := range columns {
		column = strings.ToLower(strings.TrimSpace(column))
		if column != "" {
			requested[column] = struct{}{}
		}
	}
	return requested
}

func isRequestedColumn(column domain.ColumnMapping, requested map[string]struct{}) bool {
	for _, name := range []string{column.MDSName, column.SourceName} {
		if _, ok := requested[strings.ToLower(strings.TrimSpace(name))]; ok {
			return true
		}
	}
	return false
}

func columnSortValue(column domain.ColumnMapping) int {
	switch {
	case column.KeyColumnOrdering != nil:
		return *column.KeyColumnOrdering
	case column.OrderPriority != nil:
		return 1000 + *column.OrderPriority
	case column.ValueColumnOrdering != nil:
		return 2000 + *column.ValueColumnOrdering
	default:
		return 3000
	}
}

func orderColumns(columns []domain.ColumnMapping) []string {
	ordered := make([]domain.ColumnMapping, 0)
	for _, column := range columns {
		if strings.TrimSpace(column.SourceName) == "" {
			continue
		}
		if column.OrderPriority != nil || column.KeyColumnOrdering != nil {
			ordered = append(ordered, column)
		}
	}

	if len(ordered) == 0 {
		for _, name := range []string{"ReferenceTime", "DeliveryStart"} {
			for _, column := range columns {
				if strings.EqualFold(column.MDSName, name) && strings.TrimSpace(column.SourceName) != "" {
					ordered = append(ordered, column)
				}
			}
		}
	}

	sort.SliceStable(ordered, func(i, j int) bool {
		return columnSortValue(ordered[i]) < columnSortValue(ordered[j])
	})

	order := make([]string, 0, len(ordered))
	seen := make(map[string]struct{}, len(ordered))
	for _, column := range ordered {
		key := strings.ToLower(column.SourceName)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		order = append(order, qualify(column.SourceName))
	}
	return order
}

func hyperscaleOrderColumns(columns []domain.ColumnMapping, valueColumn string) []string {
	ordered := make([]domain.ColumnMapping, 0)
	for _, column := range columns {
		if strings.TrimSpace(column.SourceName) == "" && strings.TrimSpace(column.MDSName) == "" {
			continue
		}
		if column.OrderPriority != nil || column.KeyColumnOrdering != nil {
			ordered = append(ordered, column)
		}
	}

	if len(ordered) == 0 {
		for _, name := range []string{"ReferenceTime", "DeliveryStart"} {
			for _, column := range columns {
				if strings.EqualFold(column.MDSName, name) && strings.TrimSpace(firstNonEmpty(column.SourceName, column.MDSName)) != "" {
					ordered = append(ordered, column)
				}
			}
		}
	}

	sort.SliceStable(ordered, func(i, j int) bool {
		return columnSortValue(ordered[i]) < columnSortValue(ordered[j])
	})

	order := make([]string, 0, len(ordered))
	seen := make(map[string]struct{}, len(ordered))
	for _, column := range ordered {
		outputName := firstNonEmpty(column.MDSName, column.SourceName)
		key := strings.ToLower(outputName)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		if column.IsKey {
			order = append(order, qualify(firstNonEmpty(column.SourceName, column.MDSName)))
			continue
		}
		order = append(order, hyperscaleJSONValueExpression(valueColumn, outputName, column.DataType))
	}
	return order
}

func hyperscaleViewName(mapping domain.Mapping) (string, error) {
	if strings.TrimSpace(mapping.ViewName) != "" {
		return mapping.ViewName, nil
	}
	switch mapping.DataCategory {
	case domain.Curves:
		return "Api.VI_Curve", nil
	case domain.Surfaces:
		return "Api.VI_Surface", nil
	case domain.TimeSeries:
		return "Api.VI_Timeseries", nil
	default:
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("mapping %d has no hyperscale view for data category %q", mapping.ID, mapping.DataCategory))
	}
}

func hyperscaleValueColumn(category domain.DataCategory) (string, error) {
	switch category {
	case domain.Curves:
		return "CurveValue", nil
	case domain.Surfaces:
		return "SurfaceValue", nil
	case domain.TimeSeries:
		return "TimeseriesValue", nil
	default:
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("no hyperscale value column for data category %q", category))
	}
}

func hyperscaleJSONValueExpression(valueColumn, fieldName, dataType string) string {
	jsonValue := fmt.Sprintf("JSON_VALUE(%s, '$.\"%s\"')", qualify(valueColumn), strings.ReplaceAll(fieldName, `"`, `\"`))
	switch strings.ToLower(strings.TrimSpace(dataType)) {
	case "int", "integer", "bigint", "long":
		return fmt.Sprintf("CAST(%s AS BIGINT)", jsonValue)
	case "number", "decimal", "float", "double", "real":
		return fmt.Sprintf("CAST(%s AS FLOAT)", jsonValue)
	default:
		return jsonValue
	}
}

func sqlScalarValue(value domain.FilterValue) (any, error) {
	switch value.Kind {
	case domain.FilterValueNumber:
		if strings.Contains(value.Raw, ".") {
			parsed, err := strconv.ParseFloat(value.Raw, 64)
			if err != nil {
				return nil, invalidSQLFilterValue(value.Raw, err)
			}
			return parsed, nil
		}
		parsed, err := strconv.ParseInt(value.Raw, 10, 64)
		if err != nil {
			return nil, invalidSQLFilterValue(value.Raw, err)
		}
		return parsed, nil
	case domain.FilterValuePointInTime:
		return parseSQLPointTime(value.Raw)
	case domain.FilterValueText, domain.FilterValueGeneric:
		return value.Raw, nil
	default:
		if value.Raw == "" {
			return nil, apperr.New(apperr.Invalid, "empty filter value")
		}
		return value.Raw, nil
	}
}

func sqlIntervalBounds(value domain.FilterValue) (time.Time, time.Time, bool, error) {
	if value.Start != "" && value.End != "" {
		start, err := parseSQLPointTime(value.Start)
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidSQLFilterValue(value.Start, err)
		}
		end, err := parseSQLPointTime(value.End)
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidSQLFilterValue(value.End, err)
		}
		return start, end, true, nil
	}

	return sqlIntervalFunctionBounds(value.Raw)
}

func sqlIntervalFunctionBounds(raw string) (time.Time, time.Time, bool, error) {
	name, args, ok := sqlFunctionCall(raw)
	if !ok {
		return time.Time{}, time.Time{}, false, nil
	}
	name = strings.ToLower(name)
	parts := splitSQLArguments(args)

	if name == "ti" {
		if len(parts) != 2 {
			return time.Time{}, time.Time{}, false, nil
		}
		start, err := parseSQLPointTime(parts[0])
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidSQLFilterValue(parts[0], err)
		}
		end, err := parseSQLPointTime(parts[1])
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidSQLFilterValue(parts[1], err)
		}
		return start, end, true, nil
	}

	if len(parts) == 0 {
		return time.Time{}, time.Time{}, false, nil
	}
	start, err := parseSQLPointTime(parts[0])
	if err != nil {
		return time.Time{}, time.Time{}, false, invalidSQLFilterValue(parts[0], err)
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

func sqlIntervalPointTime(raw string) (time.Time, bool, error) {
	name, args, ok := sqlFunctionCall(raw)
	if !ok {
		return time.Time{}, false, nil
	}
	start, end, ok, err := sqlIntervalFunctionBounds(args)
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

func parseSQLPointTime(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	for _, layout := range []string{"2006-01-02T15:04:05.000", "2006-01-02T15:04:05"} {
		parsed, err := time.ParseInLocation(layout, raw, time.UTC)
		if err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time %q", raw)
}

func sqlFunctionCall(raw string) (name, args string, ok bool) {
	open := strings.Index(raw, "(")
	if open <= 0 || !strings.HasSuffix(raw, ")") {
		return "", "", false
	}
	return raw[:open], raw[open+1 : len(raw)-1], true
}

func splitSQLArguments(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func isComparisonOperator(operator string) bool {
	switch operator {
	case "=", ">", ">=", "<", "<=":
		return true
	default:
		return false
	}
}

func qualify(identifier string) string {
	return "[d]." + quoteIdentifier(identifier)
}

func quoteTable(name string) string {
	name = strings.TrimSpace(name)
	if strings.Contains(name, "[") {
		return name
	}
	parts := strings.Split(name, ".")
	for i := range parts {
		parts[i] = quoteIdentifier(strings.TrimSpace(parts[i]))
	}
	return strings.Join(parts, ".")
}

func quoteIdentifier(identifier string) string {
	return "[" + strings.ReplaceAll(strings.TrimSpace(identifier), "]", "]]") + "]"
}

func invalidSQLFilterValue(raw string, err error) error {
	return apperr.Wrap(apperr.Invalid, fmt.Sprintf("invalid SQL filter value %q", raw), err)
}
