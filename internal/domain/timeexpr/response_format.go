package timeexpr

import "time"

const ResponseDateTimeLayout = "2006-01-02T15:04:05.000"

func FormatResponseTime(value time.Time) string {
	return value.Format(ResponseDateTimeLayout)
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
