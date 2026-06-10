package transactional

import (
	"time"
)

type RDPCalculator struct{}

func (RDPCalculator) Calculate(referenceTime, deliveryStart time.Time, resolution, deliveryResolution string) *int64 {
	// Truncate to UTC as per C# .DateTime usage which usually implies ignoring timezone for simple diffs in this context
	ref := referenceTime.UTC()
	del := deliveryStart.UTC()

	var result int64
	switch resolution {
	case "P1Y":
		result = int64(getYearsAdjustedPeriod(ref, del, deliveryResolution))
	case "P1M":
		result = int64(getMonthsAdjustedPeriod(ref, del))
	case "P3M":
		result = int64(getQuarterAdjustedPeriod(ref, del))
	case "P6M":
		result = int64(getHalfYearAdjustedPeriod(ref, del, deliveryResolution))
	case "P1D":
		result = int64(del.Sub(ref).Hours() / 24)
	case "P1W":
		result = int64(getWeekAdjustedPeriod(ref, del))
	case "PT1H":
		result = int64(del.Truncate(time.Hour).Sub(ref.Truncate(time.Hour)).Hours())
	case "PT30M":
		result = int64(del.Truncate(30 * time.Minute).Sub(ref.Truncate(30 * time.Minute)).Minutes() / 30)
	case "PT15M":
		result = int64(del.Truncate(15 * time.Minute).Sub(ref.Truncate(15 * time.Minute)).Minutes() / 15)
	case "PT5M":
		result = int64(del.Truncate(5 * time.Minute).Sub(ref.Truncate(5 * time.Minute)).Minutes() / 5)
	case "PT1M":
		result = int64(del.Truncate(time.Minute).Sub(ref.Truncate(time.Minute)).Minutes())
	case "PT4S":
		result = int64(del.Truncate(4 * time.Second).Sub(ref.Truncate(4 * time.Second)).Seconds() / 4)
	default:
		return nil
	}

	return &result
}

func getStartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -(weekday - 1))
}

func getWeekAdjustedPeriod(ref, del time.Time) int {
	startRef := getStartOfWeek(ref)
	startDel := getStartOfWeek(del)
	return int(startDel.Sub(startRef).Hours() / (24 * 7))
}

func getQuarterAdjustedPeriod(ref, del time.Time) int {
	startQuarter := (int(ref.Month()) - 1) / 3
	endQuarter := (int(del.Month()) - 1) / 3
	yearDiff := del.Year() - ref.Year()
	return yearDiff*4 + (endQuarter - startQuarter)
}

func getHalfYearAdjustedPeriod(ref, del time.Time, deliveryResolution string) int {
	gasYearOffset := -9
	if deliveryResolution == "Season" {
		ref = ref.AddDate(0, gasYearOffset, 0)
		del = del.AddDate(0, gasYearOffset, 0)
	}
	refHalf := (int(ref.Month()) - 1) / 6
	delHalf := (int(del.Month()) - 1) / 6
	yearDiff := del.Year() - ref.Year()
	return yearDiff*2 + (delHalf - refHalf)
}

func getYearsAdjustedPeriod(ref, del time.Time, deliveryResolution string) int {
	if deliveryResolution == "Year" {
		ref = ref.AddDate(0, -9, 0)
		del = del.AddDate(0, -9, 0)
	}
	return del.Year() - ref.Year()
}

func getMonthsAdjustedPeriod(ref, del time.Time) int {
	return getYearsAdjustedPeriod(ref, del, "")*12 + int(del.Month()) - int(ref.Month())
}
