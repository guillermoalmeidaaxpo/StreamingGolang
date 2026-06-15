package transactional

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/domain"
)

type requestValidator struct{}

var (
	columnNamePattern            = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	aggregateKeyPattern          = regexp.MustCompile(`^Aggregate\(([a-zA-Z0-9]+),\s*([A-Z0-9]+)\)$`)
	aggregationExpressionPattern = regexp.MustCompile(`(?i)^(count\(\*\)|(avg|sum|min|max|count)\([a-zA-Z0-9]+\))$`)
	aggregationPeriodPattern     = regexp.MustCompile(`^P(?:\d+W|(?:\d+Y)?(?:\d+M)?(?:\d+D)?(?:T(?:\d+H)?(?:\d+M)?(?:\d+(?:\.\d+)?S)?)?)$`)
)

func NewValidator() Validator {
	return requestValidator{}
}

func (requestValidator) Validate(_ context.Context, requests []Request) error {
	if len(requests) == 0 {
		return apperr.New(apperr.Invalid, "Invalid Request Body: Empty request")
	}

	allIDs := make(map[domain.Identifier]struct{})
	for index, request := range requests {
		if len(request.IDs) == 0 {
			return apperr.New(apperr.Invalid, fmt.Sprintf("request %d has no identifier specified", index))
		}

		seen := make(map[domain.Identifier]struct{}, len(request.IDs))
		for _, id := range request.IDs {
			if id <= 0 {
				return apperr.New(apperr.Invalid, "No identifier specified. Please provide a valid identifier")
			}
			if _, exists := seen[id]; exists {
				return apperr.New(apperr.Invalid, "Duplicated Id. Please provide unique ids")
			}
			seen[id] = struct{}{}
			if _, exists := allIDs[id]; exists {
				return apperr.New(apperr.Invalid, "Duplicated Id. Please provide unique ids")
			}
			allIDs[id] = struct{}{}
		}

		if err := validateColumns(request.Columns); err != nil {
			return err
		}
		if err := validateFilters(request.Filters); err != nil {
			return err
		}
		if err := validateShapePayload(request.Filters); err != nil {
			return err
		}
		if err := validateTransformations(request.Transformations); err != nil {
			return err
		}
	}

	return nil
}

func validateShapePayload(filters *Filters) error {
	if filters == nil || len(filters.Shape) == 0 {
		return nil
	}
	_, err := normalizeShape(filters.Shape)
	return err
}

func validateColumns(columns []string) error {
	for _, column := range columns {
		if strings.TrimSpace(column) == "" {
			return apperr.New(apperr.Invalid, "columns cannot contain empty values")
		}
		if !columnNamePattern.MatchString(column) {
			return apperr.New(apperr.Invalid, fmt.Sprintf("column name %q contains invalid characters", column))
		}
	}
	return nil
}

func validateFilters(filters *Filters) error {
	if filters == nil {
		return nil
	}
	if filters.FilterTimeZone != "" {
		if err := validateTimeZone(filters.FilterTimeZone); err != nil {
			return err
		}
	}
	for _, expression := range filters.Expressions {
		if strings.TrimSpace(expression) == "" {
			return apperr.New(apperr.Invalid, "filter expression cannot be empty")
		}
	}
	return nil
}

func validateTransformations(transformations *Transformations) error {
	if transformations == nil {
		return nil
	}
	if transformations.Timezone != "" {
		if err := validateTimeZone(transformations.Timezone); err != nil {
			return err
		}
	}
	if transformations.TargetTimeZone != "" {
		if err := validateTimeZone(transformations.TargetTimeZone); err != nil {
			return err
		}
	}
	if transformations.Nested != "" && transformations.Nested != "yes" && transformations.Nested != "no" {
		return apperr.New(apperr.Invalid, "Invalid expression. Nested parameter possible values are only 'yes' or 'no'.")
	}

	hasKeys := len(transformations.Keys) > 0
	hasValues := len(transformations.Values) > 0
	if hasKeys != hasValues {
		return apperr.New(apperr.Invalid, "Both Keys and Values must be provided for aggregation.")
	}
	if err := validateAggregationKeys(transformations.Keys); err != nil {
		return err
	}
	if err := validateAggregationValues(transformations.Values); err != nil {
		return err
	}
	return nil
}

func validateAggregationKeys(keys []string) error {
	if len(keys) > 1 {
		return apperr.New(apperr.Invalid, "Keys cannot contain more than 1 element.")
	}
	for index, key := range keys {
		if strings.TrimSpace(key) == "" {
			return apperr.New(apperr.Invalid, fmt.Sprintf("Keys[%d] contains an empty or whitespace value.", index))
		}
		expression := strings.SplitN(key, "=", 2)[0]
		match := aggregateKeyPattern.FindStringSubmatch(expression)
		if len(match) == 0 {
			return apperr.New(apperr.Invalid, fmt.Sprintf("Keys[%d]: Invalid key expression %q.", index, expression))
		}
		if !strings.EqualFold(match[1], "Delivery") {
			return apperr.New(apperr.Invalid, fmt.Sprintf("Keys[%d]: Invalid value %q in Aggregate expression. 'Delivery' period is the only one available.", index, match[1]))
		}
		if !aggregationPeriodPattern.MatchString(match[2]) {
			return apperr.New(apperr.Invalid, fmt.Sprintf("Keys[%d]: Invalid ISO 8601 period %q in Aggregate expression.", index, match[2]))
		}
	}
	return nil
}

func validateAggregationValues(values [][]string) error {
	for index, group := range values {
		if len(group) != 2 {
			return apperr.New(apperr.Invalid, fmt.Sprintf("Values[%d] must contain exactly 2 elements: [expression, columnName].", index))
		}
		if strings.TrimSpace(group[0]) == "" || !aggregationExpressionPattern.MatchString(group[0]) {
			return apperr.New(apperr.Invalid, fmt.Sprintf("Values[%d][0]: Invalid aggregation expression %q.", index, group[0]))
		}
		if strings.TrimSpace(group[1]) == "" || !columnNamePattern.MatchString(group[1]) {
			return apperr.New(apperr.Invalid, fmt.Sprintf("Values[%d][1]: Invalid column name %q.", index, group[1]))
		}
	}
	return nil
}

func validateTimeZone(name string) error {
	if strings.EqualFold(name, "UTC") {
		return nil
	}
	if _, err := time.LoadLocation(name); err != nil {
		return apperr.New(apperr.Invalid, fmt.Sprintf("invalid time zone %q", name))
	}
	return nil
}
