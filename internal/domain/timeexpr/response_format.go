package timeexpr

import "time"

const ResponseDateTimeOffsetLayout = "2006-01-02T15:04:05.000-07:00"

func FormatResponseTime(value time.Time) string {
	return value.Format(ResponseDateTimeOffsetLayout)
}

func FormatResponseValue(value any) any {
	switch typed := value.(type) {
	case time.Time:
		return FormatResponseTime(typed)
	case []time.Time:
		result := make([]any, len(typed))
		for i, value := range typed {
			result[i] = FormatResponseTime(value)
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for i, value := range typed {
			result[i] = FormatResponseValue(value)
		}
		return result
	default:
		return value
	}
}
