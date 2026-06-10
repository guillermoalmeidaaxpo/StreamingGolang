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
	"streaming-golang/internal/domain/timeexpr"
)

const referenceTimeField = "ReferenceTime"
const cmdpDefaultLookBackYears = -3

var isoPeriodPattern = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)

type FilterQuoteIndexPlanner struct{}

func (FilterQuoteIndexPlanner) PlanQuoteIndices(_ context.Context, command Command) ([]int, error) {
	if len(command.Filters.Nodes) == 0 || !hasQuoteIndexField(command.Mappings) {
		return nil, nil
	}

	loc, _ := loadLocation(command.FilterTimeZone)
	window, found, err := referenceTimeWindow(command.Filters.Nodes, loc)
	if err != nil || !found {
		return nil, err
	}
	if window.start == nil || window.end == nil {
		now := quoteIndexDate(time.Now().UTC())
		if window.start == nil && window.end != nil {
			start := window.end.AddDate(cmdpDefaultLookBackYears, 0, 0)
			window.start = &start
		} else if window.start != nil && window.end == nil {
			end := now
			if window.start.After(end) {
				end = *window.start
			}
			end = end.AddDate(0, 0, 2)
			window.end = &end
		}
	}
	if window.start == nil || window.end == nil {
		return nil, nil
	}
	if window.end.Before(*window.start) {
		return nil, nil
	}

	return quoteIndicesBetween(*window.start, *window.end), nil
}

type quoteIndexWindow struct {
	start *time.Time
	end   *time.Time
}

func referenceTimeWindow(nodes []domain.FilterNode, loc *time.Location) (quoteIndexWindow, bool, error) {
	var window quoteIndexWindow
	found := false

	for _, node := range nodes {
		comparison, ok := node.(domain.ComparisonFilter)
		if !ok || !strings.EqualFold(comparison.Field, referenceTimeField) {
			continue
		}
		found = true

		filterWindow, finite, err := referenceTimeComparisonWindow(comparison, loc)
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

func referenceTimeComparisonWindow(filter domain.ComparisonFilter, loc *time.Location) (quoteIndexWindow, bool, error) {
	switch {
	case strings.EqualFold(filter.Operator, "in"):
		return intervalQuoteIndexWindow(filter.Value, loc)
	case filter.Operator == "=":
		point, ok, err := pointTime(filter.Value, loc)
		if err != nil || !ok {
			return quoteIndexWindow{}, ok, err
		}
		return cmdpEqualityQuoteIndexWindow(point), true, nil
	case filter.Operator == ">" || filter.Operator == ">=":
		point, ok, err := pointTime(filter.Value, loc)
		if err != nil || !ok {
			return quoteIndexWindow{}, ok, err
		}
		adjustment := -2
		if filter.Operator == ">" {
			adjustment = -3
		}
		start := quoteIndexDate(point).AddDate(0, 0, adjustment)
		return quoteIndexWindow{start: &start}, true, nil
	case filter.Operator == "<" || filter.Operator == "<=":
		point, ok, err := pointTime(filter.Value, loc)
		if err != nil || !ok {
			return quoteIndexWindow{}, ok, err
		}
		adjustment := 2
		if filter.Operator == "<" {
			adjustment = 3
		}
		end := quoteIndexDate(point).AddDate(0, 0, adjustment)
		return quoteIndexWindow{end: &end}, true, nil
	default:
		return quoteIndexWindow{}, false, nil
	}
}

func intervalQuoteIndexWindow(value domain.FilterValue, loc *time.Location) (quoteIndexWindow, bool, error) {
	if value.Kind != domain.FilterValueTimeInterval {
		return quoteIndexWindow{}, false, nil
	}

	start, end, ok, err := intervalBounds(value, loc)
	if err != nil || !ok {
		return quoteIndexWindow{}, ok, err
	}

	windowStart := quoteIndexDate(start).AddDate(0, 0, -2)
	windowEnd := quoteIndexDate(end).AddDate(0, 0, 2)
	return quoteIndexWindow{start: &windowStart, end: &windowEnd}, true, nil
}

func intervalBounds(value domain.FilterValue, loc *time.Location) (time.Time, time.Time, bool, error) {
	if value.Start != "" && value.End != "" {
		start, err := ParsePointTime(value.Start, loc)
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidReferenceTime(value.Start, err)
		}
		end, err := ParsePointTime(value.End, loc)
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidReferenceTime(value.End, err)
		}
		start, end, err = applyTimeArithmetic(start, end, value.Arithmetic)
		if err != nil {
			return time.Time{}, time.Time{}, false, err
		}
		return start, end, true, nil
	}

	start, end, ok, err := intervalFunctionBounds(stripTrailingTimeArithmetic(value.Raw, value.Arithmetic), loc)
	if err != nil || !ok {
		return start, end, ok, err
	}
	start, end, err = applyTimeArithmetic(start, end, value.Arithmetic)
	if err != nil {
		return time.Time{}, time.Time{}, false, err
	}
	return start, end, true, nil
}

func intervalFunctionBounds(raw string, loc *time.Location) (time.Time, time.Time, bool, error) {
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
		start, err := ParsePointTime(parts[0], loc)
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidReferenceTime(parts[0], err)
		}
		end, err := ParsePointTime(parts[1], loc)
		if err != nil {
			return time.Time{}, time.Time{}, false, invalidReferenceTime(parts[1], err)
		}
		return start, end, true, nil
	}

	parts := splitArguments(args)
	if len(parts) == 0 {
		return time.Time{}, time.Time{}, false, nil
	}
	start, err := ParsePointTime(parts[0], loc)
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

func cmdpEqualityQuoteIndexWindow(point time.Time) quoteIndexWindow {
	day := quoteIndexDate(point)
	start := day.AddDate(0, 0, -1)
	end := day.AddDate(0, 0, 1)
	return quoteIndexWindow{start: &start, end: &end}
}

func intersectQuoteIndexWindow(w1, w2 quoteIndexWindow) quoteIndexWindow {
	start := w1.start
	if w2.start != nil && (start == nil || w2.start.After(*start)) {
		start = w2.start
	}

	end := w1.end
	if w2.end != nil && (end == nil || w2.end.Before(*end)) {
		end = w2.end
	}

	return quoteIndexWindow{start: start, end: end}
}

func quoteIndicesBetween(start, end time.Time) []int {
	indices := make([]int, 0)
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

func ParsePointTime(raw string, loc *time.Location) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	base, arithmeticOperator, period := splitPointTimeArithmetic(raw)

	if loc == nil {
		loc = time.UTC
	}

	if strings.EqualFold(base, "now()") {
		parsed := time.Now().In(loc)
		if arithmeticOperator == "" {
			return parsed.UTC(), nil
		}
		res, err := applyPeriod(parsed, arithmeticOperator, period)
		if err != nil {
			return time.Time{}, err
		}
		return res.UTC(), nil
	}

	parsed, err := timeexpr.ParsePointInTime(base, loc)
	if err != nil {
		return time.Time{}, err
	}
	if arithmeticOperator == "" {
		return parsed.UTC(), nil
	}
	res, err := applyPeriod(parsed, arithmeticOperator, period)
	if err != nil {
		return time.Time{}, err
	}
	return res.UTC(), nil
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

func pointTime(value domain.FilterValue, loc *time.Location) (time.Time, bool, error) {
	switch value.Kind {
	case domain.FilterValuePointInTime:
		point, err := ParsePointTime(value.Raw, loc)
		if err != nil {
			return time.Time{}, false, invalidReferenceTime(value.Raw, err)
		}
		return point, true, nil
	case domain.FilterValueTimeIntervalPointTime:
		point, ok, err := intervalPointTime(value.Raw, loc)
		if err != nil {
			return time.Time{}, false, err
		}
		return point, ok, nil
	default:
		return time.Time{}, false, nil
	}
}

func intervalPointTime(raw string, loc *time.Location) (time.Time, bool, error) {
	_, args, ok := functionCall(raw)
	if !ok {
		return time.Time{}, false, nil
	}
	parts := splitArguments(args)
	if len(parts) == 0 {
		return time.Time{}, false, nil
	}
	point, err := ParsePointTime(parts[0], loc)
	return point, err == nil, err
}

func loadLocation(name string) (*time.Location, error) {
	return timeexpr.LoadLocation(name)
}

func invalidReferenceTime(raw string, err error) error {
	return apperr.Wrap(apperr.Invalid, fmt.Sprintf("invalid ReferenceTime value %q", raw), err)
}
