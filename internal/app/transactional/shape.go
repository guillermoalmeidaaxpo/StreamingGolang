package transactional

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/domain"
)

type shapePayload struct {
	HolidayCalendar      *int             `json:"holidayCalendar"`
	Months               []string         `json:"months"`
	Days                 []string         `json:"days"`
	Time                 []shapeTimeRange `json:"time"`
	HolidayCalendarUpper *int             `json:"HolidayCalendar"`
	MonthsUpper          []string         `json:"Months"`
	DaysUpper            []string         `json:"Days"`
	TimeUpper            []shapeTimeRange `json:"Time"`
}

type shapeTimeRange struct {
	Start      string `json:"start"`
	End        string `json:"end"`
	StartUpper string `json:"Start"`
	EndUpper   string `json:"End"`
}

var shapeMonths = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

var shapeDays = map[string]int{
	"mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6, "sun": 7,
}

func normalizeShape(raw json.RawMessage) (*domain.NormalizedShape, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var payload shapePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, apperr.Wrap(apperr.Invalid, "invalid shape filter", err)
	}
	payload.mergeUpper()

	months, err := normalizeShapeTokens("month", payload.Months, shapeMonths)
	if err != nil {
		return nil, err
	}
	days, err := normalizeShapeTokens("day", payload.Days, shapeDays)
	if err != nil {
		return nil, err
	}
	timeSpans, err := normalizeShapeTimeSpans(payload.Time)
	if err != nil {
		return nil, err
	}

	return &domain.NormalizedShape{
		Months:          months,
		Days:            days,
		TimeSpans:       timeSpans,
		HolidayCalendar: payload.HolidayCalendar,
	}, nil
}

func (p *shapePayload) mergeUpper() {
	if p.HolidayCalendar == nil {
		p.HolidayCalendar = p.HolidayCalendarUpper
	}
	if len(p.Months) == 0 {
		p.Months = p.MonthsUpper
	}
	if len(p.Days) == 0 {
		p.Days = p.DaysUpper
	}
	if len(p.Time) == 0 {
		p.Time = p.TimeUpper
	}
}

func normalizeShapeTokens(kind string, values []string, allowed map[string]int) ([]int, error) {
	if len(values) == 0 {
		return nil, nil
	}
	seen := make(map[int]string, len(values))
	result := make([]int, 0, len(values))
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		number, ok := allowed[normalized]
		if !ok {
			return nil, apperr.New(apperr.Invalid, fmt.Sprintf("invalid shape %s %q", kind, value))
		}
		if previous, exists := seen[number]; exists {
			return nil, apperr.New(apperr.Invalid, fmt.Sprintf("duplicate shape %s %q", kind, previous))
		}
		seen[number] = value
		result = append(result, number)
	}
	sort.Ints(result)
	return result, nil
}

func normalizeShapeTimeSpans(ranges []shapeTimeRange) ([]domain.ShapeTimeSpan, error) {
	if len(ranges) == 0 {
		return nil, nil
	}
	spans := make([]domain.ShapeTimeSpan, 0, len(ranges))
	for _, item := range ranges {
		item.mergeUpper()
		start, err := parseShapeTime(item.Start, false)
		if err != nil {
			return nil, err
		}
		end, err := parseShapeTime(item.End, true)
		if err != nil {
			return nil, err
		}
		if start >= end {
			return nil, apperr.New(apperr.Invalid, fmt.Sprintf("invalid shape time range T%s-T%s", secondsToShapeTime(start), secondsToShapeTime(end)))
		}
		spans = append(spans, domain.ShapeTimeSpan{StartSeconds: start, EndSeconds: end})
	}
	sort.SliceStable(spans, func(i, j int) bool {
		if spans[i].StartSeconds == spans[j].StartSeconds {
			return spans[i].EndSeconds < spans[j].EndSeconds
		}
		return spans[i].StartSeconds < spans[j].StartSeconds
	})
	for i := 1; i < len(spans); i++ {
		if spans[i-1] == spans[i] {
			return nil, apperr.New(apperr.Invalid, "duplicate shape time range")
		}
		if spans[i].StartSeconds < spans[i-1].EndSeconds {
			return nil, apperr.New(apperr.Invalid, "overlapping shape time ranges")
		}
	}
	return spans, nil
}

func (r *shapeTimeRange) mergeUpper() {
	if r.Start == "" {
		r.Start = r.StartUpper
	}
	if r.End == "" {
		r.End = r.EndUpper
	}
}

func parseShapeTime(value string, isEnd bool) (int, error) {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "T")
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return 0, apperr.New(apperr.Invalid, fmt.Sprintf("invalid shape time %q", value))
	}
	var hour, minute, second int
	if _, err := fmt.Sscanf(value, "%02d:%02d:%02d", &hour, &minute, &second); err != nil {
		return 0, apperr.Wrap(apperr.Invalid, fmt.Sprintf("invalid shape time %q", value), err)
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 || second < 0 || second > 59 {
		return 0, apperr.New(apperr.Invalid, fmt.Sprintf("invalid shape time %q", value))
	}
	total := hour*3600 + minute*60 + second
	if isEnd && total == 0 {
		return 24 * 3600, nil
	}
	return total, nil
}

func secondsToShapeTime(seconds int) string {
	if seconds == 24*3600 {
		return "00:00:00"
	}
	hour := seconds / 3600
	minute := (seconds % 3600) / 60
	second := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}
