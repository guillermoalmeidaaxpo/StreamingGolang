package transactional

import (
	"context"
	"strings"
	"testing"

	"streaming-golang/internal/domain"
)

func TestPlannerRejectsProjectionColumnOutsideMapping(t *testing.T) {
	planner := NewPlanner(WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{mappingWithColumns(10, domain.Curves)}}))

	_, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		EndpointKind: EndpointTransactional,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{{
		IDs:     []domain.Identifier{10},
		Columns: []string{"Unknown"},
	}})

	if err == nil || !strings.Contains(err.Error(), "unmapped request projection column") {
		t.Fatalf("err = %v, want unmapped projection column error", err)
	}
}

func TestPlannerRejectsFilterColumnOutsideMapping(t *testing.T) {
	planner := NewPlanner(WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{mappingWithColumns(10, domain.Curves)}}))

	_, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		EndpointKind: EndpointTransactional,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{{
		IDs: []domain.Identifier{10},
		Filters: &Filters{
			Parsed: FilterSet{Nodes: []domain.FilterNode{
				domain.ComparisonFilter{Field: "Unknown", Operator: "=", Value: domain.FilterValue{Kind: domain.FilterValueText, Raw: "x"}},
			}},
		},
	}})

	if err == nil || !strings.Contains(err.Error(), "filter field") {
		t.Fatalf("err = %v, want unmapped filter column error", err)
	}
}

func TestPlannerRejectsTransactionalOffsetTransformation(t *testing.T) {
	planner := NewPlanner(WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{mappingWithColumns(10, domain.Curves)}}))
	offset := true

	_, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		EndpointKind: EndpointTransactional,
		Stage:        "development",
		Mode:         ModeCSV,
	}, []Request{{
		IDs: []domain.Identifier{10},
		Transformations: &Transformations{
			Offset: &offset,
		},
	}})

	if err == nil || !strings.Contains(err.Error(), "offset transformation") {
		t.Fatalf("err = %v, want offset endpoint validation error", err)
	}
}

func TestPlannerAllowsGenericOffsetTransformation(t *testing.T) {
	planner := NewPlanner(WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{mappingWithColumns(10, domain.Curves)}}))
	offset := true

	_, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		EndpointKind: EndpointGeneric,
		Stage:        "development",
		Mode:         ModeCSV,
	}, []Request{{
		IDs: []domain.Identifier{10},
		Transformations: &Transformations{
			Offset: &offset,
		},
	}})

	if err != nil {
		t.Fatalf("build plan failed: %v", err)
	}
}
