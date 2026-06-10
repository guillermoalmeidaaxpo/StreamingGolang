package cassandra

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"streaming-golang/internal/domain"
)

func TestCassandraQueryBuilderMatchesCSharpShape(t *testing.T) {
	builder := NewCassandraQueryBuilder(map[string]string{"power": "hpfc"}, "ts")
	command := domain.Command{
		DataCategory: domain.Curves,
		QuoteIndices: []int{
			20250706,
			20250707,
		},
		Mappings: []domain.Mapping{{
			ID:           312091001,
			DataCategory: domain.Curves,
			Source:       domain.SourceCassandra,
			CassandraID:  "power:312091001",
		}},
	}

	queries, err := builder.BuildQueries(context.Background(), command)
	if err != nil {
		t.Fatalf("BuildQueries returned error: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("len(queries) = %d, want 1", len(queries))
	}

	wantStatement := "SELECT ts_id, qte_y, qte_m, qte_d, quote_index, publish_time, del_y, del_m, del_d, del_h, del_min, del_offset, value FROM ts.hpfc WHERE ts_id = ? AND quote_index IN ?"
	if queries[0].Statement != wantStatement {
		t.Fatalf("statement = %q, want %q", queries[0].Statement, wantStatement)
	}

	wantArguments := []any{"power:312091001", []int{20250706, 20250707}}
	if !reflect.DeepEqual(queries[0].Arguments, wantArguments) {
		t.Fatalf("arguments = %#v, want %#v", queries[0].Arguments, wantArguments)
	}
}

func TestCassandraQueryBuilderReturnsNoRowsWhenQuoteIndicesAreEmpty(t *testing.T) {
	builder := NewCassandraQueryBuilder(map[string]string{"power": "hpfc"}, "ts")
	command := domain.Command{
		DataCategory: domain.Curves,
		Mappings: []domain.Mapping{{
			ID:           312091001,
			DataCategory: domain.Curves,
			Source:       domain.SourceCassandra,
			CassandraID:  "power:312091001",
		}},
	}

	queries, err := builder.BuildQueries(context.Background(), command)
	if err != nil {
		t.Fatalf("BuildQueries returned error: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("len(queries) = %d, want 1", len(queries))
	}
	if !strings.Contains(queries[0].Statement, "(del_y, del_m, del_d, del_h) = (?, ?, ?, ?)") {
		t.Fatalf("statement missing no-row guard: %s", queries[0].Statement)
	}
	wantArguments := []any{"power:312091001", []int{1}, int16(1), int8(1), int8(1), int8(0)}
	if !reflect.DeepEqual(queries[0].Arguments, wantArguments) {
		t.Fatalf("arguments = %#v, want %#v", queries[0].Arguments, wantArguments)
	}
}

func TestCassandraQueryBuilderAddsDeliveryTupleFilters(t *testing.T) {
	builder := NewCassandraQueryBuilder(map[string]string{"power": "hpfc"}, "ts")
	command := domain.Command{
		DataCategory: domain.Curves,
		QuoteIndices: []int{20250707},
		Filters: domain.FilterSet{Nodes: []domain.FilterNode{
			domain.ComparisonFilter{
				Field:    "DeliveryStart",
				Operator: ">=",
				Value:    domain.FilterValue{Kind: domain.FilterValuePointInTime, Raw: "2025-07-07T10:00:00"},
			},
			domain.ComparisonFilter{
				Field:    "DeliveryEnd",
				Operator: "<=",
				Value:    domain.FilterValue{Kind: domain.FilterValuePointInTime, Raw: "2025-07-07T13:00:00"},
			},
		}},
		Mappings: []domain.Mapping{{
			ID:           312091001,
			DataCategory: domain.Curves,
			Source:       domain.SourceCassandra,
			CassandraID:  "power:312091001",
		}},
	}

	queries, err := builder.BuildQueries(context.Background(), command)
	if err != nil {
		t.Fatalf("BuildQueries returned error: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("len(queries) = %d, want 1", len(queries))
	}
	if !strings.Contains(queries[0].Statement, "(del_y, del_m, del_d, del_h) >= (?, ?, ?, ?)") {
		t.Fatalf("statement missing lower delivery filter: %s", queries[0].Statement)
	}
	if !strings.Contains(queries[0].Statement, "(del_y, del_m, del_d, del_h) <= (?, ?, ?, ?)") {
		t.Fatalf("statement missing upper delivery filter: %s", queries[0].Statement)
	}

	wantSuffix := []any{int16(2025), int8(7), int8(7), int8(12), int16(2025), int8(7), int8(7), int8(14)}
	gotSuffix := queries[0].Arguments[len(queries[0].Arguments)-len(wantSuffix):]
	if !reflect.DeepEqual(gotSuffix, wantSuffix) {
		t.Fatalf("delivery args = %#v, want %#v", gotSuffix, wantSuffix)
	}
}

func TestCassandraQueryBuilderAddsRelativeDeliveryPeriodFilters(t *testing.T) {
	builder := NewCassandraQueryBuilder(map[string]string{"power": "hpfc"}, "ts")
	command := domain.Command{
		DataCategory: domain.Curves,
		QuoteIndices: []int{20250707},
		Filters: domain.FilterSet{Nodes: []domain.FilterNode{
			domain.ComparisonFilter{
				Field:    "RelativeDeliveryPeriod",
				Operator: ">=",
				Value:    domain.FilterValue{Kind: domain.FilterValueNumber, Raw: "1"},
			},
		}},
		Mappings: []domain.Mapping{{
			ID:           312091001,
			DataCategory: domain.Curves,
			Source:       domain.SourceCassandra,
			CassandraID:  "power:312091001",
		}},
	}

	queries, err := builder.BuildQueries(context.Background(), command)
	if err != nil {
		t.Fatalf("BuildQueries returned error: %v", err)
	}
	if !strings.Contains(queries[0].Statement, "(del_y, del_m, del_d, del_h) >= (?, ?, ?, ?)") {
		t.Fatalf("statement missing RDP lower filter: %s", queries[0].Statement)
	}
	wantSuffix := []any{int16(2025), int8(7), int8(7), int8(1)}
	gotSuffix := queries[0].Arguments[len(queries[0].Arguments)-len(wantSuffix):]
	if !reflect.DeepEqual(gotSuffix, wantSuffix) {
		t.Fatalf("RDP args = %#v, want %#v", gotSuffix, wantSuffix)
	}
}
