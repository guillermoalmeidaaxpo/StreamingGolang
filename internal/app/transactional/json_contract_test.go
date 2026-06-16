package transactional

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDataItemMarshalJSONFormatsTimesLikeCSharpTimestamp(t *testing.T) {
	zurich := time.FixedZone("CEST", 2*60*60)
	item := DataItem{
		ID: 536013751,
		Fields: map[string]any{
			"ReferenceTime": time.Date(2024, 4, 26, 0, 0, 0, 0, zurich),
		},
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marshal item: %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("decode item: %v", err)
	}
	referenceTime, ok := body["ReferenceTime"].([]any)
	if !ok || len(referenceTime) != 1 || referenceTime[0] != "2024-04-26T00:00:00.000" {
		t.Fatalf("ReferenceTime = %#v, want C# local timestamp with milliseconds", body["ReferenceTime"])
	}
}
