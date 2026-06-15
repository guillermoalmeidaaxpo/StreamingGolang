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

		statement, parameters, err := buildCMDPStatement(mapping, command.Filters, command.IndexRange, command.Columns, command.Aggregations, command.TargetTimeZone, command.FilterTimeZone, command.Shape)
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

		statement, parameters, err := buildHyperscaleStatement(mapping, command.Filters, command.Columns, command.VersionAsOf, command.IncludeDeleted, command.IncludeIdentifier, command.LatestReferenceTime, command.Aggregations, command.TargetTimeZone)
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

func buildCMDPStatement(mapping domain.Mapping, filters domain.FilterSet, indexRange *domain.IndexRange, requestedColumns []string, aggregations *domain.Aggregations, targetTimeZone string, filterTimeZone string, shape *domain.NormalizedShape) (string, map[string]any, error) {
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
	shapePredicates, err := builder.shapePredicates(shape, filterTimeZone)
	if err != nil {
		return "", nil, err
	}
	where = append(where, shapePredicates...)
	if indexRange != nil && strings.TrimSpace(mapping.IndexField) != "" {
		builder.addParameter("indexStart", indexRange.Start)
		builder.addParameter("indexEnd", indexRange.End)
		where = append(where,
			fmt.Sprintf("%s >= @indexStart", qualify(mapping.IndexField)),
			fmt.Sprintf("%s <= @indexEnd", qualify(mapping.IndexField)),
		)
	}

	if rankOver, ok := singleRankOverFilter(filters.Nodes); ok {
		return builder.rankOverStatement(mapping, rankOver, requestedColumns, where)
	}
	if aggregations != nil {
		return builder.aggregationStatement(mapping, aggregations, requestedColumns, targetTimeZone, quoteTable(mapping.ViewName), where, false)
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

func buildHyperscaleStatement(mapping domain.Mapping, filters domain.FilterSet, requestedColumns []string, versionAsOf *time.Time, includeDeleted bool, includeIdentifier bool, latestReferenceTime bool, aggregations *domain.Aggregations, targetTimeZone string) (string, map[string]any, error) {
	if hasRankOverFilters(filters.Nodes) {
		return "", nil, apperr.New(apperr.Invalid, "rankover filters are only supported for CMDP-hosted curves and surfaces")
	}

	viewName, err := hyperscaleViewName(mapping, requestedColumns, versionAsOf, latestReferenceTime)
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
	latestFilter, hasLatestFilter := singleLatestFilter(filters.Nodes)

	from := quoteTable(viewName)
	where := make([]string, 0)
	if versionAsOf != nil {
		builder.addParameter("MdoId", int64(mapping.ID))
		builder.addParameter("CreatedOn", *versionAsOf)
		builder.addParameter("IncludeDeleted", includeDeleted)
		from += "(@MdoId, @CreatedOn, @IncludeDeleted)"
	} else {
		builder.addParameter("id", int64(mapping.ID))
		where = append(where, fmt.Sprintf("%s = @id", qualify(hyperscaleIdentifierColumn)))
		if !includeDeleted {
			where = append(where, fmt.Sprintf("%s = 0", qualify(hyperscaleDeletedColumn)))
		}
	}

	filterPredicates, err := builder.hyperscaleFilterPredicates(filters.Nodes)
	if err != nil {
		return "", nil, err
	}
	where = append(where, filterPredicates...)
	cte := ""
	if hasLatestFilter {
		var latestPredicate string
		cte, latestPredicate, err = builder.latestReferenceCTE(mapping, latestFilter, versionAsOf)
		if err != nil {
			return "", nil, err
		}
		where = append(where, latestPredicate)
	}

	if aggregations != nil {
		statement, parameters, err := builder.aggregationStatement(mapping, aggregations, requestedColumns, targetTimeZone, from, where, true)
		if err != nil {
			return "", nil, err
		}
		return cte + statement, parameters, nil
	}

	statement := fmt.Sprintf("SELECT %s FROM %s AS [d]",
		strings.Join(hyperscaleSelectColumns(mapping, requestedColumns, valueColumn, includeIdentifier, latestReferenceTime), ", "),
		from,
	)
	if len(where) > 0 {
		statement += " WHERE " + strings.Join(where, " AND ")
	}

	if order := hyperscaleOrderColumns(mapping.Columns, valueColumn, includeIdentifier, latestReferenceTime); len(order) > 0 {
		statement += " ORDER BY " + strings.Join(order, ", ")
	}
	return cte + statement, builder.parameters, nil
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
			continue
		}
	}
	return predicates, nil
}

func (b *sqlBuilder) hyperscaleFilterPredicates(nodes []domain.FilterNode) ([]string, error) {
	predicates := make([]string, 0, len(nodes))
	for _, node := range nodes {
		switch filter := node.(type) {
		case domain.ComparisonFilter:
			if filter.Value.Kind == domain.FilterValueLatest {
				continue
			}
			predicate, err := b.comparisonPredicate(filter)
			if err != nil {
				return nil, err
			}
			if predicate != "" {
				predicates = append(predicates, predicate)
			}
		case domain.RankOverFilter:
			continue
		}
	}
	return predicates, nil
}

func (b *sqlBuilder) latestReferenceCTE(mapping domain.Mapping, filter domain.ComparisonFilter, versionAsOf *time.Time) (string, string, error) {
	if len(filter.Value.Arguments) != 1 {
		return "", "", apperr.New(apperr.Invalid, "latest filter requires one reference-time argument")
	}
	argument := filter.Value.Arguments[0]
	if !strings.EqualFold(argument.Field, "ReferenceTime") {
		return "", "", apperr.New(apperr.Invalid, "latest filter only supports ReferenceTime arguments")
	}
	if !isComparisonOperator(argument.Operator) {
		return "", "", apperr.New(apperr.Invalid, fmt.Sprintf("unsupported latest filter operator %q", argument.Operator))
	}
	value, err := sqlScalarValue(argument.Value)
	if err != nil {
		return "", "", err
	}
	b.addParameter("latestReferenceTime", value)

	predicate := fmt.Sprintf("%s = (SELECT [MaxReferenceTimeBefore] FROM [LatestReference])", qualify("ReferenceTime"))
	if versionAsOf != nil {
		tableName, err := hyperscaleVersionTableName(mapping)
		if err != nil {
			return "", "", err
		}
		cte := fmt.Sprintf("WITH [LatestReference] AS (SELECT TOP 1 [ReferenceTime] AS [MaxReferenceTimeBefore] FROM (SELECT [ReferenceTime], MAX([CreatedOn]) AS [CreatedOn] FROM %s WHERE [CreatedOn] <= @CreatedOn AND [MdoId] = @MdoId AND [ReferenceTime] %s @latestReferenceTime GROUP BY [ReferenceTime]) AS [VersionedRefs] ORDER BY [ReferenceTime] DESC) ",
			quoteTable("Core."+tableName+"Version"),
			argument.Operator,
		)
		return cte, predicate, nil
	}

	viewName := firstNonEmpty(mapping.Views.LatestVersion, defaultHyperscaleLatestVersionView(mapping.DataCategory))
	if strings.TrimSpace(viewName) == "" {
		return "", "", apperr.New(apperr.Invalid, fmt.Sprintf("mapping %d has no hyperscale latest-version view for latest filter", mapping.ID))
	}
	cte := fmt.Sprintf("WITH [LatestReference] AS (SELECT MAX([ReferenceTime]) AS [MaxReferenceTimeBefore] FROM %s WHERE [ReferenceTime] %s @latestReferenceTime AND [MdoId] = @id) ",
		quoteTable(viewName),
		argument.Operator,
	)
	return cte, predicate, nil
}

func (b *sqlBuilder) shapePredicates(shape *domain.NormalizedShape, filterTimeZone string) ([]string, error) {
	if shape == nil {
		return nil, nil
	}
	deliveryStart, ok := b.columnByMDSName("DeliveryStart")
	if !ok {
		return nil, apperr.New(apperr.Invalid, "shape filter requires a mapped DeliveryStart column")
	}
	deliveryStartExpression := b.shapeDeliveryStartExpression(deliveryStart, filterTimeZone)
	predicates := make([]string, 0, 3)

	if len(shape.Months) > 0 && len(shape.Months) < 12 {
		params := make([]string, 0, len(shape.Months))
		for i, month := range shape.Months {
			name := fmt.Sprintf("month%d", i)
			b.addParameter(name, month)
			params = append(params, "@"+name)
		}
		predicates = append(predicates, fmt.Sprintf("DATEPART(MONTH, %s) IN (%s)", deliveryStartExpression, strings.Join(params, ", ")))
	}

	if len(shape.Days) > 0 && len(shape.Days) < 7 {
		params := make([]string, 0, len(shape.Days))
		for i, day := range shape.Days {
			name := fmt.Sprintf("day%d", i)
			b.addParameter(name, day)
			params = append(params, "@"+name)
		}
		predicates = append(predicates, fmt.Sprintf("((DATEDIFF(DAY, '19000101', %s) %% 7) + 1) IN (%s)", deliveryStartExpression, strings.Join(params, ", ")))
	}

	if len(shape.TimeSpans) > 0 && !shapeCoversFullDay(shape.TimeSpans) {
		fragments := make([]string, 0, len(shape.TimeSpans))
		for i, span := range shape.TimeSpans {
			startName := fmt.Sprintf("t%d_start", i)
			b.addParameter(startName, sqlShapeTime(span.StartSeconds))
			if span.EndSeconds == 24*3600 {
				fragments = append(fragments, fmt.Sprintf("(CAST(%s AS time) >= @%s)", deliveryStartExpression, startName))
				continue
			}
			endName := fmt.Sprintf("t%d_end", i)
			b.addParameter(endName, sqlShapeTime(span.EndSeconds))
			fragments = append(fragments, fmt.Sprintf("(CAST(%s AS time) >= @%s AND CAST(%s AS time) < @%s)", deliveryStartExpression, startName, deliveryStartExpression, endName))
		}
		predicates = append(predicates, "("+strings.Join(fragments, " OR ")+")")
	}

	return predicates, nil
}

func (b *sqlBuilder) shapeDeliveryStartExpression(column domain.ColumnMapping, filterTimeZone string) string {
	expression := b.columnExpression(column)
	if strings.TrimSpace(filterTimeZone) == "" || isUTCTimeZone(filterTimeZone) {
		return expression
	}
	return fmt.Sprintf("%s AT TIME ZONE '%s'", expression, sqlServerTimeZone(filterTimeZone))
}

func shapeCoversFullDay(spans []domain.ShapeTimeSpan) bool {
	if len(spans) != 1 {
		return false
	}
	return spans[0].StartSeconds == 0 && spans[0].EndSeconds == 24*3600
}

func sqlShapeTime(seconds int) string {
	if seconds == 24*3600 {
		return "00:00:00"
	}
	hour := seconds / 3600
	minute := (seconds % 3600) / 60
	second := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}

func isUTCTimeZone(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "utc", "z", "etc/utc":
		return true
	default:
		return false
	}
}

func (b *sqlBuilder) rankOverStatement(mapping domain.Mapping, filter domain.RankOverFilter, requestedColumns []string, where []string) (string, map[string]any, error) {
	partitionBy, err := b.rankPartitionByClause(filter)
	if err != nil {
		return "", nil, err
	}
	orderBy, err := b.rankOrderByClause(filter)
	if err != nil {
		return "", nil, err
	}
	rankPredicate, err := rankOverPredicate(filter)
	if err != nil {
		return "", nil, err
	}

	columns := cmdpSelectColumnSpecs(mapping.Columns, requestedColumns)
	innerSelect := make([]string, 0, len(columns))
	outerSelect := make([]string, 0, len(columns))
	for _, column := range columns {
		innerSelect = append(innerSelect, column.Inner)
		outerSelect = append(outerSelect, column.Outer)
	}
	if len(innerSelect) == 0 {
		innerSelect = []string{"[d].*"}
		outerSelect = []string{"[d].*"}
	}

	statement := fmt.Sprintf("SELECT %s FROM (SELECT %s, RANK() OVER (PARTITION BY %s ORDER BY %s) AS [rank] FROM %s AS [d] WHERE %s) AS [d] WHERE %s",
		strings.Join(outerSelect, ", "),
		strings.Join(innerSelect, ", "),
		partitionBy,
		orderBy,
		quoteTable(mapping.ViewName),
		strings.Join(where, " AND "),
		rankPredicate,
	)
	if order := outerOrderColumns(mapping.Columns); len(order) > 0 {
		statement += " ORDER BY " + strings.Join(order, ", ")
	}
	return statement, b.parameters, nil
}

func (b *sqlBuilder) aggregationStatement(mapping domain.Mapping, aggregations *domain.Aggregations, requestedColumns []string, targetTimeZone string, from string, where []string, isHyperscale bool) (string, map[string]any, error) {
	selectColumns, err := b.aggregationSelectColumns(mapping, aggregations, requestedColumns, targetTimeZone, isHyperscale)
	if err != nil {
		return "", nil, err
	}
	groupBy, err := b.aggregationGroupByColumns(aggregations, targetTimeZone)
	if err != nil {
		return "", nil, err
	}
	orderBy, err := b.aggregationOrderByColumns(aggregations, targetTimeZone)
	if err != nil {
		return "", nil, err
	}

	statement := fmt.Sprintf("SELECT %s FROM %s AS [d]", strings.Join(selectColumns, ", "), from)
	if len(where) > 0 {
		statement += " WHERE " + strings.Join(where, " AND ")
	}
	if len(groupBy) > 0 {
		statement += " GROUP BY " + strings.Join(groupBy, ", ")
	}
	if len(orderBy) > 0 {
		statement += " ORDER BY " + strings.Join(orderBy, ", ")
	}
	return statement, b.parameters, nil
}

func (b *sqlBuilder) aggregationSelectColumns(mapping domain.Mapping, aggregations *domain.Aggregations, requestedColumns []string, targetTimeZone string, isHyperscale bool) ([]string, error) {
	requested := requestedColumnSet(requestedColumns)
	columns := make([]string, 0)
	if includeAggregationColumn(requested, "Identifier") {
		columns = append(columns, fmt.Sprintf("%d AS [Identifier]", mapping.ID))
	}

	referenceTimeGrouped := false
	for _, group := range aggregations.GroupBy {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(group.Expression)), "aggregate(") && strings.EqualFold(aggregateBucketColumn(group.Expression), "ReferenceTime") {
			expression, err := b.aggregationGroupExpression(group.Expression, targetTimeZone)
			if err != nil {
				return nil, err
			}
			if includeAggregationColumn(requested, "ReferenceTime") {
				columns = append(columns, expression+" AS [ReferenceTime]")
			}
			referenceTimeGrouped = true
			break
		}
	}
	if !referenceTimeGrouped && includeAggregationColumn(requested, "ReferenceTime") {
		columns = append(columns, fmt.Sprintf("MIN(%s) AS [ReferenceTime]", b.requiredColumnExpression("ReferenceTime")))
	}
	if includeAggregationColumn(requested, "DeliveryStart") {
		columns = append(columns, fmt.Sprintf("MIN(%s) AS [DeliveryStart]", b.requiredColumnExpression("DeliveryStart")))
	}
	if includeAggregationColumn(requested, "DeliveryEnd") {
		columns = append(columns, fmt.Sprintf("MAX(%s) AS [DeliveryEnd]", b.requiredColumnExpression("DeliveryEnd")))
	}
	if includeAggregationColumn(requested, "RelativeDeliveryPeriod") {
		columns = append(columns, "NULL AS [RelativeDeliveryPeriod]")
	}
	if !isHyperscale && includeAggregationColumn(requested, "LegacyDeliveryBucketNumber") {
		columns = append(columns, "NULL AS [LegacyDeliveryBucketNumber]")
	}

	for _, group := range aggregations.GroupBy {
		trimmed := strings.TrimSpace(group.Expression)
		if strings.HasPrefix(strings.ToLower(trimmed), "aggregate(") {
			if !includeAggregationColumn(requested, group.Alias) || strings.EqualFold(group.Alias, "ReferenceTime") {
				continue
			}
			expression, err := b.aggregationGroupExpression(trimmed, targetTimeZone)
			if err != nil {
				return nil, err
			}
			columns = append(columns, expression+" AS "+quoteIdentifier(group.Alias))
			continue
		}
		if includeAggregationColumn(requested, group.Alias) {
			column, ok := b.columnByMDSName(trimmed)
			if !ok {
				return nil, apperr.New(apperr.Invalid, fmt.Sprintf("aggregation group column %q is not mapped", trimmed))
			}
			columns = append(columns, b.columnExpression(column)+" AS "+quoteIdentifier(group.Alias))
		}
	}

	for _, aggregation := range aggregations.Expressions {
		if !includeAggregationColumn(requested, aggregation.Alias) {
			continue
		}
		expression, err := b.aggregationValueExpression(aggregation.Expression)
		if err != nil {
			return nil, err
		}
		columns = append(columns, expression+" AS "+quoteIdentifier(aggregation.Alias))
	}

	if len(columns) == 0 {
		return nil, apperr.New(apperr.Invalid, "no aggregation columns selected")
	}
	return columns, nil
}

func (b *sqlBuilder) aggregationGroupByColumns(aggregations *domain.Aggregations, targetTimeZone string) ([]string, error) {
	columns := make([]string, 0)
	if referenceTime, ok := b.columnByMDSName("ReferenceTime"); ok {
		columns = append(columns, b.columnExpression(referenceTime))
	}
	for _, group := range aggregations.GroupBy {
		trimmed := strings.TrimSpace(group.Expression)
		if strings.HasPrefix(strings.ToLower(trimmed), "aggregate(") {
			expression, err := b.aggregationGroupExpression(trimmed, targetTimeZone)
			if err != nil {
				return nil, err
			}
			columns = append(columns, expression)
			continue
		}
		column, ok := b.columnByMDSName(trimmed)
		if !ok {
			return nil, apperr.New(apperr.Invalid, fmt.Sprintf("aggregation group column %q is not mapped", trimmed))
		}
		columns = append(columns, b.columnExpression(column))
	}
	return uniqueSQLParts(columns), nil
}

func (b *sqlBuilder) aggregationOrderByColumns(aggregations *domain.Aggregations, targetTimeZone string) ([]string, error) {
	columns := make([]string, 0)
	if referenceTime, ok := b.columnByMDSName("ReferenceTime"); ok {
		columns = append(columns, b.columnExpression(referenceTime))
	}
	for _, group := range aggregations.GroupBy {
		trimmed := strings.TrimSpace(group.Expression)
		if strings.HasPrefix(strings.ToLower(trimmed), "aggregate(") {
			expression, err := b.aggregationGroupExpression(trimmed, targetTimeZone)
			if err != nil {
				return nil, err
			}
			columns = append(columns, expression)
			continue
		}
		column, ok := b.columnByMDSName(trimmed)
		if !ok {
			return nil, apperr.New(apperr.Invalid, fmt.Sprintf("aggregation group column %q is not mapped", trimmed))
		}
		columns = append(columns, b.columnExpression(column))
	}
	return uniqueSQLParts(columns), nil
}

func (b *sqlBuilder) aggregationGroupExpression(expression string, targetTimeZone string) (string, error) {
	columnName, period, ok := parseAggregateBucket(expression)
	if !ok {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("invalid aggregation key expression %q", expression))
	}
	column, ok := b.columnByMDSName(columnName)
	if !ok {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("aggregation key column %q is not mapped", columnName))
	}
	datePart, interval, anchor, err := sqlPeriodBucket(period)
	if err != nil {
		return "", err
	}
	columnExpression := b.columnExpression(column)
	if targetTimeZone = strings.TrimSpace(targetTimeZone); targetTimeZone != "" && !strings.EqualFold(targetTimeZone, "UTC") {
		columnExpression = fmt.Sprintf("CAST(%s AT TIME ZONE '%s' AS datetimeoffset)", columnExpression, sqlServerTimeZone(targetTimeZone))
	}
	return fmt.Sprintf("CAST(DATEADD(%s, (DATEDIFF(%s, '%s', %s) / %d) * %d, '%s') AS datetimeoffset)", datePart, datePart, anchor, columnExpression, interval, interval, anchor), nil
}

func (b *sqlBuilder) aggregationValueExpression(expression string) (string, error) {
	functionName, columnName, ok := parseAggregationFunction(expression)
	if !ok {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("invalid aggregation expression %q", expression))
	}
	if strings.EqualFold(columnName, "*") {
		if !strings.EqualFold(functionName, "COUNT") {
			return "", apperr.New(apperr.Invalid, fmt.Sprintf("invalid aggregation expression %q", expression))
		}
		return "COUNT(*)", nil
	}
	column, ok := b.columnByMDSName(columnName)
	if !ok {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("aggregation value column %q is not mapped", columnName))
	}
	return fmt.Sprintf("%s(%s)", strings.ToUpper(functionName), b.columnExpression(column)), nil
}

func (b *sqlBuilder) requiredColumnExpression(name string) string {
	if column, ok := b.columnByMDSName(name); ok {
		return b.columnExpression(column)
	}
	return qualify(name)
}

func includeAggregationColumn(requested map[string]struct{}, name string) bool {
	if len(requested) == 0 {
		return true
	}
	_, ok := requested[strings.ToLower(strings.TrimSpace(name))]
	return ok
}

func parseAggregateBucket(expression string) (column string, period string, ok bool) {
	name, args, ok := sqlFunctionCall(strings.TrimSpace(expression))
	if !ok || !strings.EqualFold(name, "Aggregate") {
		return "", "", false
	}
	parts := splitSQLArguments(args)
	if len(parts) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
}

func aggregateBucketColumn(expression string) string {
	column, _, ok := parseAggregateBucket(expression)
	if !ok {
		return ""
	}
	return column
}

func parseAggregationFunction(expression string) (functionName string, column string, ok bool) {
	name, args, ok := sqlFunctionCall(strings.TrimSpace(expression))
	if !ok {
		return "", "", false
	}
	switch strings.ToUpper(strings.TrimSpace(name)) {
	case "AVG", "SUM", "MIN", "MAX", "COUNT":
		return strings.ToUpper(strings.TrimSpace(name)), strings.TrimSpace(args), true
	default:
		return "", "", false
	}
}

func sqlPeriodBucket(period string) (datePart string, interval int, anchor string, err error) {
	period = strings.ToUpper(strings.TrimSpace(period))
	anchor = "20000101"
	switch {
	case strings.HasPrefix(period, "PT") && strings.HasSuffix(period, "H"):
		interval, err = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(period, "PT"), "H"))
		return "HOUR", interval, anchor, invalidPeriodIfNeeded(period, interval, err)
	case strings.HasPrefix(period, "PT") && strings.HasSuffix(period, "M"):
		interval, err = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(period, "PT"), "M"))
		return "MINUTE", interval, anchor, invalidPeriodIfNeeded(period, interval, err)
	case strings.HasPrefix(period, "P") && strings.HasSuffix(period, "D"):
		interval, err = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(period, "P"), "D"))
		return "DAY", interval, anchor, invalidPeriodIfNeeded(period, interval, err)
	case strings.HasPrefix(period, "P") && strings.HasSuffix(period, "W"):
		interval, err = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(period, "P"), "W"))
		return "WEEK", interval, "20000103", invalidPeriodIfNeeded(period, interval, err)
	case strings.HasPrefix(period, "P") && strings.HasSuffix(period, "M"):
		interval, err = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(period, "P"), "M"))
		return "MONTH", interval, anchor, invalidPeriodIfNeeded(period, interval, err)
	case strings.HasPrefix(period, "P") && strings.HasSuffix(period, "Y"):
		interval, err = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(period, "P"), "Y"))
		return "YEAR", interval, anchor, invalidPeriodIfNeeded(period, interval, err)
	default:
		return "", 0, "", apperr.New(apperr.Invalid, fmt.Sprintf("unsupported aggregation period %q", period))
	}
}

func invalidPeriodIfNeeded(period string, interval int, err error) error {
	if err != nil || interval <= 0 {
		return apperr.New(apperr.Invalid, fmt.Sprintf("unsupported aggregation period %q", period))
	}
	return nil
}

func sqlServerTimeZone(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "cet", "europe/zurich":
		return "W. Europe Standard Time"
	default:
		return strings.ReplaceAll(name, "'", "''")
	}
}

func uniqueSQLParts(parts []string) []string {
	seen := make(map[string]struct{}, len(parts))
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		key := strings.ToLower(strings.TrimSpace(part))
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, part)
	}
	return result
}

func (b *sqlBuilder) rankPartitionByClause(filter domain.RankOverFilter) (string, error) {
	parts := make([]string, 0, len(filter.PartitionBy))
	for _, field := range filter.PartitionBy {
		column, ok := b.columnByMDSName(field)
		if !ok {
			return "", apperr.New(apperr.Invalid, fmt.Sprintf("rankover partition column %q is not mapped", field))
		}
		parts = append(parts, b.columnExpression(column))
	}
	if len(parts) == 0 {
		return "", apperr.New(apperr.Invalid, "rankover requires at least one partition column")
	}
	return strings.Join(parts, " , "), nil
}

func (b *sqlBuilder) rankOrderByClause(filter domain.RankOverFilter) (string, error) {
	parts := make([]string, 0, len(filter.OrderBy))
	for _, order := range filter.OrderBy {
		column, ok := b.columnByMDSName(order.Field)
		if !ok {
			return "", apperr.New(apperr.Invalid, fmt.Sprintf("rankover order column %q is not mapped", order.Field))
		}
		direction := strings.ToUpper(strings.TrimSpace(order.Direction))
		switch direction {
		case "":
			parts = append(parts, b.columnExpression(column))
		case "ASC", "DESC":
			parts = append(parts, b.columnExpression(column)+" "+direction)
		default:
			return "", apperr.New(apperr.Invalid, fmt.Sprintf("rankover order direction %q is not supported", order.Direction))
		}
	}
	if len(parts) == 0 {
		return "", apperr.New(apperr.Invalid, "rankover requires at least one order column")
	}
	return strings.Join(parts, " , "), nil
}

func rankOverPredicate(filter domain.RankOverFilter) (string, error) {
	if len(filter.Bounds) == 0 {
		return "[d].[rank] = 1", nil
	}
	bound := filter.Bounds[0]
	start := strings.TrimSpace(bound.Start)
	if start == "" {
		return "", apperr.New(apperr.Invalid, "rankover lower bound is required")
	}
	if _, err := strconv.Atoi(start); err != nil {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("rankover lower bound %q is invalid", start))
	}

	end := strings.TrimSpace(bound.End)
	if end == "" {
		return fmt.Sprintf("[d].[rank] = %s", start), nil
	}
	if strings.EqualFold(end, "last") {
		return fmt.Sprintf("[d].[rank] >= %s", start), nil
	}
	if _, err := strconv.Atoi(end); err != nil {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("rankover upper bound %q is invalid", end))
	}
	return fmt.Sprintf("[d].[rank] >= %s AND [d].[rank] <= %s", start, end), nil
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
	if b.jsonValueColumn != "" {
		if !column.IsKey {
			return hyperscaleJSONValueExpression(b.jsonValueColumn, firstNonEmpty(column.MDSName, column.SourceName), column.DataType)
		}
		return qualify(hyperscalePhysicalColumnName(column))
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
	specs := cmdpSelectColumnSpecs(columns, requestedColumns)
	if len(specs) == 0 {
		return []string{"[d].*"}
	}
	expressions := make([]string, 0, len(specs))
	for _, spec := range specs {
		expressions = append(expressions, spec.Inner)
	}
	return expressions
}

type selectColumnSpec struct {
	Inner string
	Outer string
}

func cmdpSelectColumnSpecs(columns []domain.ColumnMapping, requestedColumns []string) []selectColumnSpec {
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
		return nil
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
	specs := make([]selectColumnSpec, 0, len(selected))
	for _, column := range selected {
		key := strings.ToLower(column.SourceName)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		outputName := firstNonEmpty(column.MDSName, column.SourceName)
		expression := qualify(column.SourceName)
		if outputName != "" && !strings.EqualFold(outputName, column.SourceName) {
			expression += " AS " + quoteIdentifier(outputName)
		}
		specs = append(specs, selectColumnSpec{
			Inner: expression,
			Outer: "[d]." + quoteIdentifier(outputName),
		})
	}
	return specs
}

func hyperscaleSelectColumns(mapping domain.Mapping, requestedColumns []string, valueColumn string, includeIdentifier bool, latestReferenceTime bool) []string {
	requested := requestedColumnSet(requestedColumns)
	if latestReferenceTime && mapping.DataCategory == domain.TimeSeries {
		return hyperscaleLatestReferenceTimeTimeseriesColumns(mapping.Columns, requested, valueColumn, includeIdentifier)
	}

	columns := mapping.Columns
	selected := make([]domain.ColumnMapping, 0, len(columns))
	for _, column := range columns {
		if strings.TrimSpace(column.SourceName) == "" && strings.TrimSpace(column.MDSName) == "" {
			continue
		}
		if len(requested) > 0 && !column.IsKey && !isRequestedColumn(column, requested) {
			continue
		}
		if column.IsKey || column.IsProjectable {
			if isIdentifierColumn(column) && !includeIdentifier {
				continue
			}
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
			sourceName := hyperscalePhysicalColumnName(column)
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
	if hasRequestedColumnInSet(requested, "CreatedOn") {
		expressions = append(expressions, qualify("CreatedOn"))
	}
	return expressions
}

func hyperscaleLatestReferenceTimeTimeseriesColumns(columns []domain.ColumnMapping, requested map[string]struct{}, valueColumn string, includeIdentifier bool) []string {
	expressions := make([]string, 0)
	keyColumns := make([]domain.ColumnMapping, 0)
	for _, column := range columns {
		if !column.IsKey {
			continue
		}
		if isIdentifierColumn(column) && !includeIdentifier {
			continue
		}
		keyColumns = append(keyColumns, column)
	}
	sort.SliceStable(keyColumns, func(i, j int) bool {
		return columnSortValue(keyColumns[i]) < columnSortValue(keyColumns[j])
	})
	if len(keyColumns) == 0 {
		expressions = append(expressions, qualify("ReferenceTime"))
	} else {
		for _, column := range keyColumns {
			expressions = append(expressions, qualify(hyperscalePhysicalColumnName(column)))
		}
	}

	valueColumns := make([]domain.ColumnMapping, 0)
	for _, column := range columns {
		if column.IsKey || strings.EqualFold(column.MDSName, "Identifier") || strings.EqualFold(column.MDSName, "ReferenceTime") {
			continue
		}
		if !column.IsProjectable {
			continue
		}
		if len(requested) > 0 && !hasRequestedColumnInSet(requested, column.MDSName) && !hasRequestedColumnInSet(requested, column.SourceName) && !hasRequestedColumnInSet(requested, "CreatedOn") {
			continue
		}
		valueColumns = append(valueColumns, column)
	}
	if len(valueColumns) == 0 {
		valueColumns = append(valueColumns, domain.ColumnMapping{MDSName: "Value", SourceName: "Value"})
	}
	sort.SliceStable(valueColumns, func(i, j int) bool {
		return columnSortValue(valueColumns[i]) < columnSortValue(valueColumns[j])
	})
	for index, column := range valueColumns {
		name := firstNonEmpty(column.MDSName, column.SourceName, "Value")
		expressions = append(expressions, hyperscaleJSONValueExpression(valueColumn, name, column.DataType)+fmt.Sprintf(" AS [Property%d]", index))
	}
	if hasRequestedColumnInSet(requested, "CreatedOn") {
		expressions = append(expressions, qualify("CreatedOn"))
	}
	return expressions
}

func findColumn(columns []domain.ColumnMapping, name string) (domain.ColumnMapping, bool) {
	for _, column := range columns {
		if strings.EqualFold(column.MDSName, name) || strings.EqualFold(column.SourceName, name) {
			return column, true
		}
	}
	return domain.ColumnMapping{}, false
}

func hasRequestedColumnInSet(requested map[string]struct{}, name string) bool {
	if len(requested) == 0 {
		return false
	}
	_, ok := requested[strings.ToLower(strings.TrimSpace(name))]
	return ok
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

func outerOrderColumns(columns []domain.ColumnMapping) []string {
	ordered := make([]domain.ColumnMapping, 0)
	for _, column := range columns {
		if strings.TrimSpace(firstNonEmpty(column.MDSName, column.SourceName)) == "" {
			continue
		}
		if isIdentifierColumn(column) {
			continue
		}
		if column.OrderPriority != nil || column.KeyColumnOrdering != nil {
			ordered = append(ordered, column)
		}
	}

	if len(ordered) == 0 {
		for _, name := range []string{"ReferenceTime", "DeliveryStart"} {
			for _, column := range columns {
				if strings.EqualFold(column.MDSName, name) && strings.TrimSpace(firstNonEmpty(column.MDSName, column.SourceName)) != "" {
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
		name := firstNonEmpty(column.MDSName, column.SourceName)
		key := strings.ToLower(name)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		order = append(order, "[d]."+quoteIdentifier(name))
	}
	return order
}

func singleRankOverFilter(nodes []domain.FilterNode) (domain.RankOverFilter, bool) {
	for _, node := range nodes {
		if filter, ok := node.(domain.RankOverFilter); ok {
			return filter, true
		}
	}
	return domain.RankOverFilter{}, false
}

func singleLatestFilter(nodes []domain.FilterNode) (domain.ComparisonFilter, bool) {
	for _, node := range nodes {
		filter, ok := node.(domain.ComparisonFilter)
		if !ok {
			continue
		}
		if filter.Value.Kind == domain.FilterValueLatest {
			return filter, true
		}
	}
	return domain.ComparisonFilter{}, false
}

func hasRankOverFilters(nodes []domain.FilterNode) bool {
	_, ok := singleRankOverFilter(nodes)
	return ok
}

func orderByMDSColumns(columns []domain.ColumnMapping) []string {
	ordered := make([]domain.ColumnMapping, 0)
	for _, column := range columns {
		if column.OrderPriority != nil {
			ordered = append(ordered, column)
		}
	}
	if len(ordered) == 0 {
		return []string{qualify("ReferenceTime")}
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		if *ordered[i].OrderPriority == *ordered[j].OrderPriority {
			return ordered[i].MDSName < ordered[j].MDSName
		}
		return *ordered[i].OrderPriority < *ordered[j].OrderPriority
	})
	order := make([]string, 0, len(ordered))
	seen := make(map[string]struct{}, len(ordered))
	for _, column := range ordered {
		name := hyperscalePhysicalColumnName(column)
		key := strings.ToLower(name)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		order = append(order, qualify(name))
	}
	return order
}

func hyperscaleOrderColumns(columns []domain.ColumnMapping, valueColumn string, includeIdentifier bool, latestReferenceTime bool) []string {
	if latestReferenceTime {
		order := orderByMDSColumns(columns)
		if len(order) > 0 {
			return order
		}
		return []string{qualify("ReferenceTime")}
	}

	ordered := make([]domain.ColumnMapping, 0)
	for _, column := range columns {
		if strings.TrimSpace(column.SourceName) == "" && strings.TrimSpace(column.MDSName) == "" {
			continue
		}
		if isIdentifierColumn(column) && !includeIdentifier {
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
			order = append(order, qualify(hyperscalePhysicalColumnName(column)))
			continue
		}
		order = append(order, hyperscaleJSONValueExpression(valueColumn, outputName, column.DataType))
	}
	return order
}

func hyperscalePhysicalColumnName(column domain.ColumnMapping) string {
	if isIdentifierColumn(column) {
		return hyperscaleIdentifierColumn
	}
	return firstNonEmpty(column.SourceName, column.MDSName)
}

func isIdentifierColumn(column domain.ColumnMapping) bool {
	return strings.EqualFold(column.MDSName, "Identifier") ||
		strings.EqualFold(column.SourceName, "Identifier") ||
		strings.EqualFold(column.MDSName, "MdoId") ||
		strings.EqualFold(column.SourceName, "MdoId")
}

func hyperscaleViewName(mapping domain.Mapping, requestedColumns []string, versionAsOf *time.Time, latestReferenceTime bool) (string, error) {
	var viewName string
	if versionAsOf != nil {
		if latestReferenceTime {
			viewName = firstNonEmpty(mapping.Views.GetByCreatedOnLatestReferenceTime, defaultHyperscaleGetByCreatedOnLatestReferenceTimeView(mapping.DataCategory))
		} else {
			viewName = firstNonEmpty(mapping.Views.GetByCreatedOn, defaultHyperscaleGetByCreatedOnView(mapping.DataCategory))
		}
	} else if latestReferenceTime {
		viewName = firstNonEmpty(
			mapping.Views.LatestReferenceTime,
			mapping.Views.LatestReferenceTimeWithCreatedOn,
			defaultHyperscaleLatestReferenceTimeView(mapping.DataCategory),
			defaultHyperscaleLatestReferenceTimeWithCreatedOnView(mapping.DataCategory),
		)
	} else if hasRequestedColumn(requestedColumns, "CreatedOn") {
		viewName = firstNonEmpty(mapping.Views.LatestVersionWithCreatedOn, defaultHyperscaleLatestVersionWithCreatedOnView(mapping.DataCategory))
	} else {
		viewName = firstNonEmpty(mapping.Views.LatestVersion, defaultHyperscaleLatestVersionView(mapping.DataCategory))
	}
	if strings.TrimSpace(viewName) == "" {
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("mapping %d has no hyperscale view for data category %q", mapping.ID, mapping.DataCategory))
	}
	return viewName, nil
}

func defaultHyperscaleLatestReferenceTimeView(category domain.DataCategory) string {
	name, ok := hyperscaleCategoryName(category)
	if !ok {
		return ""
	}
	return fmt.Sprintf("Api.VI_%sLatestVersionLatestReferenceTime", name)
}

func defaultHyperscaleLatestReferenceTimeWithCreatedOnView(category domain.DataCategory) string {
	name, ok := hyperscaleCategoryName(category)
	if !ok {
		return ""
	}
	if category == domain.TimeSeries {
		return fmt.Sprintf("Api.VI_%sLatestVersionLatestReferenceTime", name)
	}
	return fmt.Sprintf("Api.VI_%sLatestVersionLatestReferenceTimeWithCreatedOn", name)
}

func hasRequestedColumn(columns []string, name string) bool {
	for _, column := range columns {
		if strings.EqualFold(strings.TrimSpace(column), name) {
			return true
		}
	}
	return false
}

func defaultHyperscaleLatestVersionView(category domain.DataCategory) string {
	name, ok := hyperscaleCategoryName(category)
	if !ok {
		return ""
	}
	return fmt.Sprintf("Api.VI_%sLatestVersion", name)
}

func defaultHyperscaleLatestVersionWithCreatedOnView(category domain.DataCategory) string {
	name, ok := hyperscaleCategoryName(category)
	if !ok {
		return ""
	}
	if category == domain.TimeSeries {
		return fmt.Sprintf("Api.VI_%sLatestVersion", name)
	}
	return fmt.Sprintf("Api.VI_%sLatestVersionWithCreatedOn", name)
}

func defaultHyperscaleGetByCreatedOnView(category domain.DataCategory) string {
	name, ok := hyperscaleCategoryName(category)
	if !ok {
		return ""
	}
	return fmt.Sprintf("Api.TVF_Get%sByCreatedOn", name)
}

func defaultHyperscaleGetByCreatedOnLatestReferenceTimeView(category domain.DataCategory) string {
	name, ok := hyperscaleCategoryName(category)
	if !ok {
		return ""
	}
	return fmt.Sprintf("Api.TVF_Get%sByCreatedOnLatestReferenceTime", name)
}

func hyperscaleCategoryName(category domain.DataCategory) (string, bool) {
	switch category {
	case domain.Curves:
		return "Curve", true
	case domain.Surfaces:
		return "Surface", true
	case domain.TimeSeries:
		return "Timeseries", true
	default:
		return "", false
	}
}

func hyperscaleVersionTableName(mapping domain.Mapping) (string, error) {
	viewName := firstNonEmpty(mapping.Views.LatestVersion, defaultHyperscaleLatestVersionView(mapping.DataCategory))
	for _, token := range strings.Split(viewName, ".") {
		token = strings.Trim(token, "[] ")
		if strings.HasPrefix(token, "VI_") {
			token = strings.TrimPrefix(token, "VI_")
		}
		if strings.HasSuffix(token, "LatestVersion") {
			token = strings.TrimSuffix(token, "LatestVersion")
			if token != "" {
				return token, nil
			}
		}
	}
	if name, ok := hyperscaleCategoryName(mapping.DataCategory); ok {
		return name, nil
	}
	return "", apperr.New(apperr.Invalid, fmt.Sprintf("cannot infer hyperscale version table for data category %q", mapping.DataCategory))
}

func hyperscaleValueColumn(category domain.DataCategory) (string, error) {
	switch category {
	case domain.Curves:
		return "CurveValue", nil
	case domain.Surfaces:
		return "SurfaceValue", nil
	case domain.TimeSeries:
		return "TimeSeriesValue", nil
	default:
		return "", apperr.New(apperr.Invalid, fmt.Sprintf("no hyperscale value column for data category %q", category))
	}
}

func hyperscaleJSONTextValueExpression(valueColumn, fieldName string) string {
	return fmt.Sprintf("JSON_VALUE(%s, '$.\"%s\"')", qualify(valueColumn), strings.ReplaceAll(fieldName, `"`, `\"`))
}

func hyperscaleJSONValueExpression(valueColumn, fieldName, dataType string) string {
	jsonValue := hyperscaleJSONTextValueExpression(valueColumn, fieldName)
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
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		parsed, err := time.Parse(layout, raw)
		if err == nil {
			return parsed, nil
		}
	}
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
