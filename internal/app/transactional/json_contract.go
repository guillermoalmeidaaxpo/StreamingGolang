package transactional

import (
	"encoding/json"
	"sort"

	"streaming-golang/internal/domain"
)

func (item DataItem) MarshalJSON() ([]byte, error) {
	value := make(map[string]any, len(item.Fields)+1)
	value["Identifier"] = int64(item.ID)

	keys := make([]string, 0, len(item.Fields))
	for key := range item.Fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value[key] = jsonArrayValue(item.Fields[key])
	}
	return json.Marshal(value)
}

func jsonArrayValue(value any) any {
	switch typed := value.(type) {
	case nil:
		return []any{nil}
	case []any:
		return typed
	case []string:
		result := make([]any, len(typed))
		for i, value := range typed {
			result[i] = value
		}
		return result
	case []int:
		result := make([]any, len(typed))
		for i, value := range typed {
			result[i] = value
		}
		return result
	case []int64:
		result := make([]any, len(typed))
		for i, value := range typed {
			result[i] = value
		}
		return result
	case []float64:
		result := make([]any, len(typed))
		for i, value := range typed {
			result[i] = value
		}
		return result
	case []domain.Identifier:
		result := make([]any, len(typed))
		for i, value := range typed {
			result[i] = int64(value)
		}
		return result
	default:
		return []any{value}
	}
}
