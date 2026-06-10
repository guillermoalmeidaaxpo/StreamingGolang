package transactional

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/domain"
)

const referenceTimeField = "ReferenceTime"

var isoPeriodPattern = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)

type FilterQuoteIndexPlanner struct{}

func (FilterQuoteIndexPlanner) PlanQuoteIndices(_ context.Context, command Command) ([]int, error) {
	if len(command.Filters.Nodes) == 0 || !hasQuoteIndexField(command.Mappings) {
		return nil, nil
	}

	window, found, err := referenceTimeWindow(command.Filters.Nodes)
	if err != nil || !found {
		return nil, err
	}
	if window.start == nil || window.end == nil {
		return nil, nil
	}
	if window.end.Before(*window.start) {
		return nil, apperr.New(apperr.Invalid, "ReferenceTime filters produce an empty quote index window")
	}

	return quoteIndicesBetween(*window.start, *window.end), nil
}

type quoteIndexWindow struct {
	start *time.Time
	end   *time.Time
}

func referenceTimeWindow(nodes []domain.FilterNode) (quoteIndexWindow, bool, error) {
	var window quoteIndexWindow
	found := false

	for _, node := range nodes {
		comparison, ok := node.(domain.ComparisonFilter)
		if !ok || !strings.EqualFold(comparison.Field, referenceTimeField) {
			continue
		}
		found = true

		filterWindow, finite, err := referenceTimeComparisonWindow(comparison)
		if err != nil {
			return quoteIndexWindow{}, false, err
		}
		if !finite {
			continue
		}
		window = intersectQuoteIndexWindow(window, filterWindow)
	}

	return window, found, nil
}

func referenceTimeComparisonWindow(filter domain.ComparisonFilter) (quoteIndexWindow, bool, error) {
	switch {
	case strings.EqualFold(filter.Operator, "in"):
		return intervalQuoteIndexWindow(filter.Value)
	case filter.Operator == "=":
		point, ok, err := pointTime(filter.Value)
		if err != nil || !ok {
			return quoteIndexWindow{}, ok, err
		}
		return paddedPointWindow(point), true, nil
	case filter.Operator == ">" || filter.Operator == ">=":
		point, ok, err := pointTime(filter.Value)
		if err != nil || !ok {
			return quoteIndexWindow{}, ok, err
		}
		start := quoteIndexDate(point).AddDate(0, 0, -1)
		return quoteIndexWindow{start: &start}, true, nil
	case filter.Operator == "<" || filter.Operator == "<=":
		point, ok, err := pointTime(filter.Value)
		if err != nil || !ok {
			return quoteIndexWindow{}, ok, err
		}
		end := quoteIndexDate(point).AddDate(0, 0, 1)
		return quoteIndexWindow{end: &end}, true, nil
	default:
		return quoteIndexWindow{}, false, nil
	}
}

func intervalQuoteIndexWindow(value domain.FilterValue) (quoteIndexWindow, bool, error) {
	if value.Kind != domain.FilterValueTimeInterval {
		return quoteIndexWindow{}, false, nil
	}

	start, end, ok, err := intervalBounds(value)
	if err != nil || !ok {
		return quoteIndexWindow{}, ok, err
	}

	windowStart := quoteIndexDate(start).AddDate(0, 0, -1)
	windowEnd := quoteIndexDate(end).AddDate(0, 0, 1)
	return quoteIndexWindow{start: &windowStart, end: &windowEnd}, true, nil
}

func intervalBounds(value domain.FilterValue) (time.Time, time.Time, bool, error) {
	if value.Start != "" && value.End != "" {
		start, err := parsePointTime(value.Start)
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidReferenceTime(value.Start, err)
		}
		end, err := parsePointTime(value.End)
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidReferenceTime(value.End, err)
		}
		start, end, err = applyTimeArithmetic(start, end, value.Arithmetic)
		if err != nil {
			return time.Time{}, time.Time{}, false, err
		}
		return start, end, true, nil
	}

	start, end, ok, err := intervalFunctionBounds(stripTrailingTimeArithmetic(value.Raw, value.Arithmetic))
	if err != nil || !ok {
		return start, end, ok, err
	}
	start, end, err = applyTimeArithmetic(start, end, value.Arithmetic)
	if err != nil {
		return time.Time{}, time.Time{}, false, err
	}
	return start, end, true, nil
}

func intervalFunctionBounds(raw string) (time.Time, time.Time, bool, error) {
	name, args, ok := functionCall(raw)
	if !ok {
		return time.Time{}, time.Time{}, false, nil
	}

	name = strings.ToLower(name)
	if name == "ti" {
		parts := splitArguments(args)
		if len(parts) != 2 {
			return time.Time{}, time.Time{}, false, nil
		}
		start, err := parsePointTime(parts[0])
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidReferenceTime(parts[0], err)
		}
		end, err := parsePointTime(parts[1])
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidReferenceTime(parts[1], err)
		}
		return start, end, true, nil
	}

	parts := splitArguments(args)
	if len(parts) == 0 {
		return time.Time{}, time.Time{}, false, nil
	}
	start, err := parsePointTime(parts[0])
	if err != nil {
		return time.Time{}, time.Time{}, false, invalidReferenceTime(parts[0], err)
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

func pointTime(value domain.FilterValue) (time.Time, bool, error) {
	switch value.Kind {
	case domain.FilterValuePointInTime:
		point, err := parsePointTime(value.Raw)
		if err != nil {
			return time.Time{}, false, invalidReferenceTime(value.Raw, err)
		}
		return point, true, nil
	case domain.FilterValueTimeIntervalPointTime:
		point, ok, err := intervalPointTime(value.Raw)
		if err != nil || !ok {
			return time.Time{}, ok, err
		}
		return point, true, nil
	default:
		return time.Time{}, false, nil
	}
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

func paddedPointWindow(point time.Time) quoteIndexWindow {
	day := quoteIndexDate(point)
	start := day.AddDate(0, 0, -1)
	end := day.AddDate(0, 0, 1)
	return quoteIndexWindow{start: &start, end: &end}
}

func intersectQuoteIndexWindow(left, right quoteIndexWindow) quoteIndexWindow {
	if right.start != nil && (left.start == nil || right.start.After(*left.start)) {
		left.start = right.start
	}
	if right.end != nil && (left.end == nil || right.end.Before(*left.end)) {
		left.end = right.end
	}
	return left
}

func quoteIndicesBetween(start, end time.Time) []int {
	start = quoteIndexDate(start)
	end = quoteIndexDate(end)

	indices := make([]int, 0, int(end.Sub(start).Hours()/24)+1)
	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {
		indices = append(indices, quoteIndex(day))
	}
	return indices
}

func quoteIndexDate(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

func quoteIndex(value time.Time) int {
	day := quoteIndexDate(value)
	return day.Year()*10000 + int(day.Month())*100 + day.Day()
}

func parsePointTime(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	base, arithmeticOperator, period := splitPointTimeArithmetic(raw)

	var parsed time.Time
	var err error

	if strings.EqualFold(base, "now()") {
		parsed = time.Now().UTC()
		if arithmeticOperator == "" {
			return parsed, nil
		}
		return applyPeriod(parsed, arithmeticOperator, period)
	}

	for _, layout := range []string{"2006-01-02T15:04:05.000", "2006-01-02T15:04:05"} {
		parsed, err = time.ParseInLocation(layout, base, time.UTC)
		if err == nil {
			if arithmeticOperator == "" {
				return parsed, nil
			}
			return applyPeriod(parsed, arithmeticOperator, period)
		}
	}
	return time.Time{}, err
}

func splitPointTimeArithmetic(raw string) (base, operator, period string) {
	for _, marker := range []string{"+P", "-P"} {
		if index := strings.Index(raw, marker); index > 0 {
			return raw[:index], raw[index : index+1], raw[index+1:]
		}
	}
	return raw, "", ""
}

func applyPeriod(value time.Time, operator, rawPeriod string) (time.Time, error) {
	parts := isoPeriodPattern.FindStringSubmatch(rawPeriod)
	if parts == nil {
		return time.Time{}, fmt.Errorf("invalid ISO period %q", rawPeriod)
	}

	sign := 1
	if operator == "-" {
		sign = -1
	}

	years := atoi(parts[1]) * sign
	months := atoi(parts[2]) * sign
	weeks := atoi(parts[3])
	days := (atoi(parts[4]) + weeks*7) * sign
	hours := atoi(parts[5]) * sign
	minutes := atoi(parts[6]) * sign
	seconds := atoi(parts[7]) * sign

	value = value.AddDate(years, months, days)
	return value.Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second), nil
}

func applyTimeArithmetic(start, end time.Time, arithmetic *domain.TimeArithmetic) (time.Time, time.Time, error) {
	if arithmetic == nil {
		return start, end, nil
	}
	shiftedStart, err := applyPeriod(start, arithmetic.Operator, arithmetic.Period)
	if err != nil {
		return time.Time{}, time.Time{}, invalidReferenceTime(arithmetic.Period, err)
	}
	shiftedEnd, err := applyPeriod(end, arithmetic.Operator, arithmetic.Period)
	if err != nil {
		return time.Time{}, time.Time{}, invalidReferenceTime(arithmetic.Period, err)
	}
	return shiftedStart, shiftedEnd, nil
}

func stripTrailingTimeArithmetic(raw string, arithmetic *domain.TimeArithmetic) string {
	if arithmetic == nil {
		return raw
	}
	suffix := arithmetic.Operator + arithmetic.Period
	return strings.TrimSuffix(raw, suffix)
}

func atoi(raw string) int {
	if raw == "" {
		return 0
	}
	value, _ := strconv.Atoi(raw)
	return value
}

func functionCall(raw string) (name, args string, ok bool) {
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

func hasQuoteIndexField(mappings []domain.Mapping) bool {
	if len(mappings) == 0 {
		return true
	}
	for _, mapping := range mappings {
		if mapping.IndexField != "" {
			return true
		}
	}
	return false
}

func invalidReferenceTime(raw string, err error) error {
	return apperr.Wrap(apperr.Invalid, fmt.Sprintf("invalid ReferenceTime value %q", raw), err)
}
