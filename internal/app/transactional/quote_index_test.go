package transactional

import (
	"context"
	"reflect"
	"testing"

	"streaming-golang/internal/domain"
)

func TestFilterQuoteIndexPlannerGeneratesIndicesFromReferenceTimeInterval(t *testing.T) {
	indices, err := FilterQuoteIndexPlanner{}.PlanQuoteIndices(context.Background(), Command{
		Filters: FilterSet{Nodes: []domain.FilterNode{
			referenceTimeInterval("2023-01-01T00:00:00", "2023-01-03T00:00:00"),
		}},
		Mappings: []domain.Mapping{{IndexField: "QuoteDateIndex_FID"}},
	})
	if err != nil {
		t.Fatalf("plan quote indices failed: %v", err)
	}

	want := []int{20221231, 20230101, 20230102, 20230103, 20230104}
	if !reflect.DeepEqual(indices, want) {
		t.Fatalf("indices = %#v, want %#v", indices, want)
	}
}

func TestFilterQuoteIndexPlannerPadsEqualityAcrossMonthBoundary(t *testing.T) {
	indices, err := FilterQuoteIndexPlanner{}.PlanQuoteIndices(context.Background(), Command{
		Filters: FilterSet{Nodes: []domain.FilterNode{
			referenceTimePoint("=", "2023-03-01T23:00:00"),
		}},
		Mappings: []domain.Mapping{{IndexField: "QuoteDateIndex_FID"}},
	})
	if err != nil {
		t.Fatalf("plan quote indices failed: %v", err)
	}

	want := []int{20230228, 20230301, 20230302}
	if !reflect.DeepEqual(indices, want) {
		t.Fatalf("indices = %#v, want %#v", indices, want)
	}
}

func TestFilterQuoteIndexPlannerIntersectsReferenceTimeBounds(t *testing.T) {
	indices, err := FilterQuoteIndexPlanner{}.PlanQuoteIndices(context.Background(), Command{
		Filters: FilterSet{Nodes: []domain.FilterNode{
			referenceTimePoint(">=", "2023-01-02T10:00:00"),
			referenceTimePoint("<=", "2023-01-04T12:00:00"),
		}},
		Mappings: []domain.Mapping{{IndexField: "QuoteDateIndex_FID"}},
	})
	if err != nil {
		t.Fatalf("plan quote indices failed: %v", err)
	}

	want := []int{20230101, 20230102, 20230103, 20230104, 20230105}
	if !reflect.DeepEqual(indices, want) {
		t.Fatalf("indices = %#v, want %#v", indices, want)
	}
}

func TestFilterQuoteIndexPlannerIgnoresNonReferenceTimeFilters(t *testing.T) {
	indices, err := FilterQuoteIndexPlanner{}.PlanQuoteIndices(context.Background(), Command{
		Filters: FilterSet{Nodes: []domain.FilterNode{
			domain.ComparisonFilter{
				Field:    "DeliveryStart",
				Operator: "=",
				Value: domain.FilterValue{
					Kind: domain.FilterValuePointInTime,
					Raw:  "2023-01-01T00:00:00",
				},
			},
		}},
		Mappings: []domain.Mapping{{IndexField: "QuoteDateIndex_FID"}},
	})
	if err != nil {
		t.Fatalf("plan quote indices failed: %v", err)
	}
	if indices != nil {
		t.Fatalf("indices = %#v, want nil", indices)
	}
}

func TestFilterQuoteIndexPlannerSkipsOpenEndedReferenceTimeWindow(t *testing.T) {
	indices, err := FilterQuoteIndexPlanner{}.PlanQuoteIndices(context.Background(), Command{
		Filters: FilterSet{Nodes: []domain.FilterNode{
			referenceTimePoint(">=", "2023-01-02T10:00:00"),
		}},
		Mappings: []domain.Mapping{{IndexField: "QuoteDateIndex_FID"}},
	})
	if err != nil {
		t.Fatalf("plan quote indices failed: %v", err)
	}
	if indices != nil {
		t.Fatalf("indices = %#v, want nil", indices)
	}
}

func referenceTimeInterval(start, end string) domain.ComparisonFilter {
	return domain.ComparisonFilter{
		Field:    "ReferenceTime",
		Operator: "in",
		Value: domain.FilterValue{
			Kind:  domain.FilterValueTimeInterval,
			Raw:   "ti(" + start + "," + end + ")",
			Start: start,
			End:   end,
		},
	}
}

func referenceTimePoint(operator, raw string) domain.ComparisonFilter {
	return domain.ComparisonFilter{
		Field:    "ReferenceTime",
		Operator: operator,
		Value: domain.FilterValue{
			Kind: domain.FilterValuePointInTime,
			Raw:  raw,
		},
	}
}
