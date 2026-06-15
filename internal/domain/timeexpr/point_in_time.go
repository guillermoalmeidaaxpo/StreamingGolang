package timeexpr

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
)

func ParsePointInTime(raw string, loc *time.Location) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if loc == nil {
		loc = time.UTC
	}

	if hasExplicitOffset(raw) {
		parsed, err := time.Parse(time.RFC3339Nano, raw)
		if err != nil {
			return time.Time{}, err
		}
		return parsed.UTC(), nil
	}

	return ParsePointInTimeToken(raw, loc)
}

func ParsePointInTimeToken(raw string, loc *time.Location) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if loc == nil {
		loc = time.UTC
	}

	datePart, timePart, ok := strings.Cut(raw, "T")
	if !ok {
		return time.Time{}, fmt.Errorf("invalid point-in-time token %q", raw)
	}

	year, month, day, err := parseDate(datePart)
	if err != nil {
		return time.Time{}, err
	}
	hour, minute, second, nanosecond, err := parseTimeOfDay(timePart)
	if err != nil {
		return time.Time{}, err
	}

	parsed := time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, loc)
	if parsed.Year() != year || int(parsed.Month()) != month || parsed.Day() != day ||
		parsed.Hour() != hour || parsed.Minute() != minute || parsed.Second() != second {
		return time.Time{}, fmt.Errorf("invalid point-in-time token %q", raw)
	}
	return parsed.UTC(), nil
}

func FormatUTC(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

func LoadLocation(name string) (*time.Location, error) {
	switch strings.ToUpper(strings.TrimSpace(name)) {
	case "", "UTC":
		return time.UTC, nil
	case "CET":
		name = "Europe/Zurich"
	}
	return time.LoadLocation(name)
}

func hasExplicitOffset(raw string) bool {
	if strings.HasSuffix(raw, "Z") {
		return true
	}
	tIndex := strings.Index(raw, "T")
	if tIndex < 0 {
		return false
	}
	return strings.LastIndex(raw[tIndex+1:], "+") >= 0 || strings.LastIndex(raw[tIndex+1:], "-") >= 0
}

func parseDate(raw string) (int, int, int, error) {
	parts := strings.Split(raw, "-")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid date token %q", raw)
	}
	year, err := atoiPart(parts[0], "year")
	if err != nil {
		return 0, 0, 0, err
	}
	month, err := atoiPart(parts[1], "month")
	if err != nil {
		return 0, 0, 0, err
	}
	day, err := atoiPart(parts[2], "day")
	if err != nil {
		return 0, 0, 0, err
	}
	return year, month, day, nil
}

func parseTimeOfDay(raw string) (int, int, int, int, error) {
	parts := strings.Split(raw, ":")
	if len(parts) != 3 {
		return 0, 0, 0, 0, fmt.Errorf("invalid time token %q", raw)
	}
	hour, err := atoiPart(parts[0], "hour")
	if err != nil {
		return 0, 0, 0, 0, err
	}
	minute, err := atoiPart(parts[1], "minute")
	if err != nil {
		return 0, 0, 0, 0, err
	}

	secondPart := parts[2]
	millisecond := 0
	if secondText, fraction, ok := strings.Cut(secondPart, "."); ok {
		secondPart = secondText
		millisecond, err = atoiPart(fraction, "millisecond")
		if err != nil {
			return 0, 0, 0, 0, err
		}
	}
	second, err := atoiPart(secondPart, "second")
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return hour, minute, second, millisecond * int(time.Millisecond), nil
}

func atoiPart(raw string, name string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s token %q", name, raw)
	}
	return value, nil
}
