package transactional

import (
	"context"
	"strings"
	"time"

	"streaming-golang/internal/domain"
)

const cassandraDefaultLookBackYears = -3

type CassandraQuoteIndexPlanner struct {
	Now func() time.Time
}

func (p CassandraQuoteIndexPlanner) PlanQuoteIndices(_ context.Context, command Command) ([]int, error) {
	now := time.Now
	if p.Now != nil {
		now = p.Now
	}
	return cassandraQuoteIndices(command.Filters.Nodes, cassandraTimeZone(command.Mappings), now())
}

type cassandraDateRange struct {
	start *time.Time
	end   *time.Time
}

func cassandraQuoteIndices(nodes []domain.FilterNode, timezone string, now time.Time) ([]int, error) {
	location, err := loadCassandraLocation(timezone)
	if err != nil {
		return nil, err
	}

	dateRange, err := cassandraReferenceTimeRange(nodes, location, now)
	if err != nil {
		return nil, err
	}
	return dateRange.cassandraQuoteIndices(now), nil
}

func cassandraReferenceTimeRange(nodes []domain.FilterNode, location *time.Location, now time.Time) (cassandraDateRange, error) {
	var dateRange cassandraDateRange
	foundReferenceTime := false

	for _, node := range nodes {
		filter, ok := node.(domain.ComparisonFilter)
		if !ok || !strings.EqualFold(filter.Field, referenceTimeField) {
			continue
		}
		foundReferenceTime = true

		if strings.EqualFold(filter.Operator, "in") && filter.Value.Kind == domain.FilterValueTimeInterval {
			start, end, ok, err := intervalBounds(filter.Value)
			if err != nil || !ok {
				return cassandraDateRange{}, err
			}
			tightenCassandraStart(&dateRange, cassandraStartDate(start, location, true))
			tightenCassandraEnd(&dateRange, cassandraEndDate(end, location, true))
			continue
		}

		point, ok, err := pointTime(filter.Value)
		if err != nil || !ok {
			return cassandraDateRange{}, err
		}

		switch filter.Operator {
		case ">=":
			tightenCassandraStart(&dateRange, cassandraStartDate(point, location, true))
		case ">":
			tightenCassandraStart(&dateRange, cassandraStartDate(point, location, false))
		case "<=":
			tightenCassandraEnd(&dateRange, cassandraEndDate(point, location, true))
		case "<":
			tightenCassandraEnd(&dateRange, cassandraEndDate(point, location, false))
		case "=":
			local := point.In(location)
			if !isLocalMidnight(local) {
				start := quoteIndexDate(now)
				end := start.AddDate(0, 0, -1)
				dateRange.start = &start
				dateRange.end = &end
				return dateRange, nil
			}
			day := localDate(local)
			dateRange.start = &day
			dateRange.end = &day
			return dateRange, nil
		}
	}

	if !foundReferenceTime {
		return dateRange, nil
	}
	if dateRange.start == nil {
		start := quoteIndexDate(now.UTC().AddDate(cassandraDefaultLookBackYears, 0, 0))
		dateRange.start = &start
	}
	return dateRange, nil
}

func (r cassandraDateRange) cassandraQuoteIndices(now time.Time) []int {
	if r.start == nil && r.end == nil {
		return nil
	}
	if r.start == nil {
		return nil
	}

	end := r.end
	if end == nil {
		defaultEnd := quoteIndexDate(now.UTC())
		if !r.start.Before(defaultEnd) {
			defaultEnd = *r.start
		}
		end = &defaultEnd
	}

	if r.start.After(*end) {
		return nil
	}
	return quoteIndicesBetween(*r.start, *end)
}

func cassandraStartDate(value time.Time, location *time.Location, inclusive bool) time.Time {
	local := value.UTC().In(location)
	day := localDate(local)
	if inclusive && isLocalMidnight(local) {
		return day
	}
	return day.AddDate(0, 0, 1)
}

func cassandraEndDate(value time.Time, location *time.Location, inclusive bool) time.Time {
	local := value.UTC().In(location)
	day := localDate(local)
	if !inclusive && isLocalMidnight(local) {
		return day.AddDate(0, 0, -1)
	}
	return day
}

func tightenCassandraStart(dateRange *cassandraDateRange, candidate time.Time) {
	if dateRange.start == nil || candidate.After(*dateRange.start) {
		dateRange.start = &candidate
	}
}

func tightenCassandraEnd(dateRange *cassandraDateRange, candidate time.Time) {
	if dateRange.end == nil || candidate.Before(*dateRange.end) {
		dateRange.end = &candidate
	}
}

func localDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func isLocalMidnight(value time.Time) bool {
	return value.Hour() == 0 && value.Minute() == 0 && value.Second() == 0 && value.Nanosecond() == 0
}

func cassandraTimeZone(mappings []domain.Mapping) string {
	if len(mappings) == 0 {
		return "Europe/Zurich"
	}
	switch mappings[0].ID {
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
	switch strings.ToUpper(strings.TrimSpace(name)) {
	case "", "CET":
		name = "Europe/Zurich"
	case "UTC":
		name = "UTC"
	}
	return time.LoadLocation(name)
}
