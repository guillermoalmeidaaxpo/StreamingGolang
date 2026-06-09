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

type fixedMappingResolver struct {
	mappings []domain.Mapping
}

func (r fixedMappingResolver) ResolveMappings(context.Context, []domain.Identifier, domain.DataCategory, string) ([]domain.Mapping, error) {
	return r.mappings, nil
}

func mappingWithColumns(id domain.Identifier, category domain.DataCategory) domain.Mapping {
	return domain.Mapping{
		ID:           id,
		DataCategory: category,
		Source:       domain.SourceCMDP,
		ViewName:     "CurveView",
		IndexField:   "QuoteDateIndex_FID",
		Columns: []domain.ColumnMapping{
			{MDSName: "ReferenceTime", SourceName: "ReferenceTime"},
			{MDSName: "Value", SourceName: "Value"},
		},
	}
}
