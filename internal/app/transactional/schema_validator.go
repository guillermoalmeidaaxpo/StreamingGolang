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
			if command.HasAggregations {
				return apperr.New(apperr.Invalid, "rankover filters cannot be combined with aggregations")
			}
			for _, field := range filter.PartitionBy {
				if !mappedColumnExists(mappings, field) {
					return apperr.New(apperr.Invalid, fmt.Sprintf("rankover partition column %q is not mapped", field))
				}
			}
			for _, order := range filter.OrderBy {
				if !mappedColumnExists(mappings, order.Field) {
					return apperr.New(apperr.Invalid, fmt.Sprintf("rankover order column %q is not mapped", order.Field))
				}
			}
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
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	for _, mapping := range mappings {
		for _, column := range mapping.Columns {
			if strings.EqualFold(column.MDSName, name) || strings.EqualFold(column.SourceName, name) {
				return true
			}
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
