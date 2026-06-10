package transactional

import (
	"context"
	"reflect"
	"testing"
	"time"

	"streaming-golang/internal/domain"
)

func TestCassandraQuoteIndexPlannerGeneratesIndicesInMDOTimeZone(t *testing.T) {
	indices, err := CassandraQuoteIndexPlanner{}.PlanQuoteIndices(context.Background(), Command{
		Filters: domain.FilterSet{Nodes: []domain.FilterNode{
			domain.ComparisonFilter{
				Field:    "ReferenceTime",
				Operator: "in",
				Value: domain.FilterValue{
					Kind:  domain.FilterValueTimeInterval,
					Start: "2025-08-22T22:00:00",
					End:   "2025-08-24T22:00:00",
				},
			},
		}},
		Mappings: []domain.Mapping{{ID: 536013751, Source: domain.SourceCassandra}},
	})
	if err != nil {
		t.Fatalf("PlanQuoteIndices returned error: %v", err)
	}

	want := []int{20250823, 20250824, 20250825}
	if !reflect.DeepEqual(indices, want) {
		t.Fatalf("indices = %#v, want %#v", indices, want)
	}
}

func TestCassandraQuoteIndexPlannerUsesDefaultLookbackForOpenEndedReferenceTime(t *testing.T) {
	indices, err := CassandraQuoteIndexPlanner{
		Now: func() time.Time {
			return time.Date(2026, 6, 10, 9, 0, 0, 0, time.UTC)
		},
	}.PlanQuoteIndices(context.Background(), Command{
		Filters: domain.FilterSet{Nodes: []domain.FilterNode{
			domain.ComparisonFilter{
				Field:    "ReferenceTime",
				Operator: "<",
				Value: domain.FilterValue{
					Kind: domain.FilterValuePointInTime,
					Raw:  "2026-06-10T22:00:00",
				},
			},
		}},
		Mappings: []domain.Mapping{{ID: 536013751, Source: domain.SourceCassandra}},
	})
	if err != nil {
		t.Fatalf("PlanQuoteIndices returned error: %v", err)
	}
	if len(indices) == 0 {
		t.Fatal("indices is empty, want default Cassandra lookback range")
	}
	if indices[0] != 20230610 {
		t.Fatalf("first index = %d, want 20230610", indices[0])
	}
	if indices[len(indices)-1] != 20260610 {
		t.Fatalf("last index = %d, want 20260610", indices[len(indices)-1])
	}
}

func TestCassandraQuoteIndexPlannerDoesNotInventDefaultWithoutReferenceTime(t *testing.T) {
	indices, err := CassandraQuoteIndexPlanner{
		Now: func() time.Time {
			return time.Date(2026, 6, 10, 9, 0, 0, 0, time.UTC)
		},
	}.PlanQuoteIndices(context.Background(), Command{
		Mappings: []domain.Mapping{{ID: 536013751, Source: domain.SourceCassandra}},
	})
	if err != nil {
		t.Fatalf("PlanQuoteIndices returned error: %v", err)
	}
	if len(indices) != 0 {
		t.Fatalf("indices = %#v, want empty until filter provider supplies default ReferenceTime", indices)
	}
}
