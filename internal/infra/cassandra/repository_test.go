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
			"cassandra_timezone": "Europe/Zurich",
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
	got, ok := fields["ReferenceTime"].(time.Time)
	if !ok {
		t.Fatalf("ReferenceTime = %#v, want time.Time", fields["ReferenceTime"])
	}
	if got.Format(time.RFC3339) != "2024-04-26T00:00:00+02:00" {
		t.Fatalf("ReferenceTime = %s, want 2024-04-26T00:00:00+02:00", got.Format(time.RFC3339))
	}
	deliveryStart, ok := fields["DeliveryStart"].(time.Time)
	if !ok {
		t.Fatalf("DeliveryStart = %#v, want time.Time", fields["DeliveryStart"])
	}
	if deliveryStart.Format(time.RFC3339) != "2025-11-10T13:00:00+01:00" {
		t.Fatalf("DeliveryStart = %s, want 2025-11-10T13:00:00+01:00", deliveryStart.Format(time.RFC3339))
	}
	deliveryEnd, ok := fields["DeliveryEnd"].(time.Time)
	if !ok {
		t.Fatalf("DeliveryEnd = %#v, want time.Time", fields["DeliveryEnd"])
	}
	if deliveryEnd.Format(time.RFC3339) != "2025-11-10T14:00:00+01:00" {
		t.Fatalf("DeliveryEnd = %s, want 2025-11-10T14:00:00+01:00", deliveryEnd.Format(time.RFC3339))
	}
}

func TestMapCassandraRowCalculatesDeliveryEndLikeCSharpAcrossDST(t *testing.T) {
	row := map[string]any{
		"qte_y":      2024,
		"qte_m":      3,
		"qte_d":      31,
		"del_y":      2024,
		"del_m":      3,
		"del_d":      31,
		"del_h":      1,
		"del_min":    0,
		"del_offset": 1,
		"value":      64.34,
	}
	query := domain.ExecutableQuery{
		ID: 536013751,
		Parameters: map[string]any{
			"projection_columns": []string{"DeliveryStart", "DeliveryEnd"},
			"cassandra_timezone": "Europe/Zurich",
		},
	}

	fields := mapCassandraRow(row, query)

	deliveryStart := fields["DeliveryStart"].(time.Time)
	if deliveryStart.Format(time.RFC3339) != "2024-03-31T01:00:00+01:00" {
		t.Fatalf("DeliveryStart = %s, want 2024-03-31T01:00:00+01:00", deliveryStart.Format(time.RFC3339))
	}
	deliveryEnd := fields["DeliveryEnd"].(time.Time)
	if deliveryEnd.Format(time.RFC3339) != "2024-03-31T03:00:00+02:00" {
		t.Fatalf("DeliveryEnd = %s, want C# wall-clock DST-adjusted 2024-03-31T03:00:00+02:00", deliveryEnd.Format(time.RFC3339))
	}
}
