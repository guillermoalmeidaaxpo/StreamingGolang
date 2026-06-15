package cassandra

import (
	"testing"
	"time"

	"streaming-golang/internal/domain"
)

func TestMapCassandraRowKeepsCSharpContractColumnsWhenProjectionIsRequested(t *testing.T) {
	row := map[string]any{
		"qte_y":      2024,
		"qte_m":      4,
		"qte_d":      26,
		"del_y":      2025,
		"del_m":      11,
		"del_d":      10,
		"del_h":      13,
		"del_min":    0,
		"del_offset": 1,
		"value":      115.9137420654,
	}
	query := domain.ExecutableQuery{
		ID: 536013751,
		Parameters: map[string]any{
			"projection_columns": []string{"Value"},
		},
	}

	fields := mapCassandraRow(row, query)

	for _, column := range []string{
		"Identifier",
		"ReferenceTime",
		"DeliveryStart",
		"DeliveryEnd",
		"LegacyDeliveryBucketNumber",
		"RelativeDeliveryPeriod",
		"Value",
	} {
		if _, ok := fields[column]; !ok {
			t.Fatalf("field %q missing from %#v", column, fields)
		}
	}
	if fields["Identifier"] != int64(536013751) {
		t.Fatalf("Identifier = %#v, want 536013751", fields["Identifier"])
	}
	if got, ok := fields["ReferenceTime"].(time.Time); !ok || !got.Equal(time.Date(2024, 4, 26, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("ReferenceTime = %#v, want 2024-04-26T00:00:00Z", fields["ReferenceTime"])
	}
}
