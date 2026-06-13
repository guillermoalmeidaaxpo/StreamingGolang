package transactional

import (
	"fmt"
	"strings"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/domain"
)

func validateAgainstMappings(requestContext RequestContext, request Request, command Command, mappings []domain.Mapping) error {
	if err := validateEndpointRules(requestContext, request); err != nil {
		return err
	}
	if len(mappings) == 0 {
		return nil
	}
	if err := validateMappingCategory(requestContext, mappings); err != nil {
		return err
	}
	if err := validateShape(command, mappings); err != nil {
		return err
	}
	if err := validateAggregations(command, mappings); err != nil {
		return err
	}
	if err := validateProjectionColumns(request, mappings); err != nil {
		return err
	}
	if err := validateFilterColumns(request, command, mappings); err != nil {
		return err
	}
	return validateMappingContent(mappings)
}

func validateEndpointRules(requestContext RequestContext, request Request) error {
	if request.Transformations == nil {
		return nil
	}

	if request.Transformations.Offset != nil && endpointKind(requestContext) == EndpointTransactional {
		return apperr.New(apperr.Invalid, "offset transformation is only available on generic requests")
	}

	if request.Transformations.Nested != "" {
		if endpointKind(requestContext) == EndpointGeneric || requestContext.DataCategory == domain.TimeSeries {
			return apperr.New(apperr.Invalid, "nested transformation is only available for curves and surfaces transactional requests")
		}
	}

	return nil
}

func endpointKind(requestContext RequestContext) EndpointKind {
	if requestContext.EndpointKind == "" {
		return EndpointTransactional
	}
	return requestContext.EndpointKind
}

func validateMappingCategory(requestContext RequestContext, mappings []domain.Mapping) error {
	if endpointKind(requestContext) == EndpointGeneric && requestContext.DataCategory == "" {
		return nil
	}
	if requestContext.DataCategory == "" {
		return nil
	}

	invalid := make([]string, 0)
	for _, mapping := range mappings {
		if mapping.DataCategory != "" && mapping.DataCategory != requestContext.DataCategory {
			invalid = append(invalid, fmt.Sprint(mapping.ID))
		}
	}
	if len(invalid) == 0 {
		return nil
	}
	return apperr.New(apperr.Invalid,
		fmt.Sprintf("request contains IDs that are not of data category %q: %s", requestContext.DataCategory, strings.Join(invalid, ", ")))
}

func validateShape(command Command, mappings []domain.Mapping) error {
	if !command.HasShape {
		return nil
	}
	if command.DataCategory != domain.Curves {
		return apperr.New(apperr.Invalid, "shape filters are only available for curves")
	}
	for _, mapping := range mappings {
		if mapping.HyperscaleID != nil {
			return apperr.New(apperr.Invalid, "shape filters are not supported for hyperscale IDs")
		}
	}
	return nil
}

func validateProjectionColumns(request Request, mappings []domain.Mapping) error {
	if request.Transformations != nil && (len(request.Transformations.Keys) > 0 || len(request.Transformations.Values) > 0) {
		return nil
	}
	if len(request.Columns) == 0 || !hasColumnMetadata(mappings) {
		return nil
	}

	for _, requested := range request.Columns {
		if strings.EqualFold(requested, "CreatedOn") && anyHyperscale(mappings) {
			continue
		}
		if !mappedColumnExists(mappings, requested) {
			return apperr.New(apperr.Invalid, fmt.Sprintf("unmapped request projection column %q", requested))
		}
	}
	return nil
}

func validateAggregations(command Command, mappings []domain.Mapping) error {
	if !command.HasAggregations || command.Aggregations == nil {
		return nil
	}
	if command.DataCategory != domain.Curves {
		return apperr.New(apperr.Invalid, "Aggregation endpoint is only supported for data category 'Curve'.")
	}
	for _, group := range command.Aggregations.GroupBy {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(group.Expression)), "aggregate(") {
			continue
		}
		if !mappedColumnExists(mappings, group.Expression) {
			return apperr.New(apperr.Invalid, fmt.Sprintf("aggregation key column %q is not mapped", group.Expression))
		}
	}
	aliases := make(map[string]struct{})
	for _, expression := range append([]domain.AggregationColumn{}, command.Aggregations.GroupBy...) {
		alias := strings.ToLower(strings.TrimSpace(expression.Alias))
		if alias == "" {
			continue
		}
		if _, exists := aliases[alias]; exists {
			return apperr.New(apperr.Invalid, fmt.Sprintf("duplicate aggregation alias %q", expression.Alias))
		}
		aliases[alias] = struct{}{}
	}
	for _, expression := range command.Aggregations.Expressions {
		alias := strings.ToLower(strings.TrimSpace(expression.Alias))
		if alias != "" {
			if _, exists := aliases[alias]; exists {
				return apperr.New(apperr.Invalid, fmt.Sprintf("duplicate aggregation alias %q", expression.Alias))
			}
			aliases[alias] = struct{}{}
		}
		columnName := aggregationValueColumn(expression.Expression)
		if strings.EqualFold(columnName, "*") {
			continue
		}
		column, ok := mappedColumn(mappings, columnName)
		if !ok {
			return apperr.New(apperr.Invalid, fmt.Sprintf("aggregation value column %q is not mapped", columnName))
		}
		if column.IsKey {
			return apperr.New(apperr.Invalid, fmt.Sprintf("aggregation value column %q cannot be a key column", columnName))
		}
	}
	return nil
}

func aggregationValueColumn(expression string) string {
	open := strings.Index(expression, "(")
	close := strings.LastIndex(expression, ")")
	if open < 0 || close <= open {
		return ""
	}
	return strings.TrimSpace(expression[open+1 : close])
}

func validateFilterColumns(request Request, command Command, mappings []domain.Mapping) error {
	if request.Filters == nil || len(command.Filters.Nodes) == 0 || !hasColumnMetadata(mappings) {
		return nil
	}

	for _, node := range command.Filters.Nodes {
		switch filter := node.(type) {
		case domain.ComparisonFilter:
			if err := validateComparisonFilter(filter, mappings); err != nil {
				return err
			}
		case domain.RankOverFilter:
			return validateRankOverFilter(command, mappings, filter)
		}
	}
	return nil
}

func validateRankOverFilter(command Command, mappings []domain.Mapping, filter domain.RankOverFilter) error {
	if command.HasAggregations {
		return apperr.New(apperr.Invalid, "rankover filters cannot be combined with aggregations")
	}
	if command.DataCategory == domain.TimeSeries {
		return apperr.New(apperr.Invalid, "rankover filters are only available for curves and surfaces")
	}
	for _, mapping := range mappings {
		if strings.TrimSpace(mapping.CassandraID) != "" || mapping.HyperscaleID != nil || mapping.Source == domain.SourceHyperscale || mapping.Source == domain.SourceCassandra {
			return apperr.New(apperr.Invalid, "rankover filters are only supported for CMDP-hosted IDs")
		}
	}
	for _, field := range filter.PartitionBy {
		column, ok := mappedColumn(mappings, field)
		if !ok {
			return apperr.New(apperr.Invalid, fmt.Sprintf("rankover partition column %q is not mapped", field))
		}
		if !column.IsKey || isRankOverExcludedColumn(column, field) {
			return apperr.New(apperr.Invalid, fmt.Sprintf("invalid rankover partition column %q", field))
		}
	}
	for _, order := range filter.OrderBy {
		column, ok := mappedColumn(mappings, order.Field)
		if !ok {
			return apperr.New(apperr.Invalid, fmt.Sprintf("rankover order column %q is not mapped", order.Field))
		}
		if isRankOverExcludedColumn(column, order.Field) {
			return apperr.New(apperr.Invalid, fmt.Sprintf("invalid rankover order column %q", order.Field))
		}
	}
	return nil
}

func validateComparisonFilter(filter domain.ComparisonFilter, mappings []domain.Mapping) error {
	if !mappedColumnExists(mappings, filter.Field) {
		return apperr.New(apperr.Invalid, fmt.Sprintf("filter field %q is not mapped", filter.Field))
	}
	if filter.Value.Kind == domain.FilterValueLatest && anyCassandra(mappings) {
		return apperr.New(apperr.Invalid, "latest filters are not supported for Cassandra IDs")
	}
	for _, latest := range filter.Value.Arguments {
		if !mappedColumnExists(mappings, latest.Field) {
			return apperr.New(apperr.Invalid, fmt.Sprintf("latest filter field %q is not mapped", latest.Field))
		}
	}
	return nil
}

func validateMappingContent(mappings []domain.Mapping) error {
	for _, mapping := range mappings {
		if mapping.Source == "" {
			return apperr.New(apperr.Invalid, fmt.Sprintf("mapping %d has no source", mapping.ID))
		}
		if mapping.Source == domain.SourceCMDP && strings.TrimSpace(mapping.ViewName) == "" {
			return apperr.New(apperr.Invalid, fmt.Sprintf("mapping %d has no CMDP view name", mapping.ID))
		}
	}
	return nil
}

func hasColumnMetadata(mappings []domain.Mapping) bool {
	for _, mapping := range mappings {
		if len(mapping.Columns) > 0 {
			return true
		}
	}
	return false
}

func mappedColumnExists(mappings []domain.Mapping, name string) bool {
	_, ok := mappedColumn(mappings, name)
	return ok
}

func mappedColumn(mappings []domain.Mapping, name string) (domain.ColumnMapping, bool) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.ColumnMapping{}, false
	}
	for _, mapping := range mappings {
		for _, column := range mapping.Columns {
			if strings.EqualFold(column.MDSName, name) || strings.EqualFold(column.SourceName, name) {
				return column, true
			}
		}
	}
	return domain.ColumnMapping{}, false
}

func isRankOverExcludedColumn(column domain.ColumnMapping, requested string) bool {
	for _, name := range []string{requested, column.MDSName, column.SourceName} {
		if strings.EqualFold(strings.TrimSpace(name), "Identifier") ||
			strings.EqualFold(strings.TrimSpace(name), "MdoId") ||
			strings.EqualFold(strings.TrimSpace(name), "RelativeDeliveryPeriod") {
			return true
		}
	}
	return false
}

func anyHyperscale(mappings []domain.Mapping) bool {
	for _, mapping := range mappings {
		if mapping.HyperscaleID != nil {
			return true
		}
	}
	return false
}

func anyCassandra(mappings []domain.Mapping) bool {
	for _, mapping := range mappings {
		if strings.TrimSpace(mapping.CassandraID) != "" || mapping.Source == domain.SourceCassandra {
			return true
		}
	}
	return false
}
