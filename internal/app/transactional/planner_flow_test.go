package transactional

import (
	"context"
	"testing"

	"streaming-golang/internal/domain"
)

func TestPlannerBuildsExecutableQueriesThroughFullFlow(t *testing.T) {
	planner := NewPlanner(WithQueryStrategy(SplitQueryStrategy{
		QueriesCount:           1,
		ReferenceTimeSplitDays: 2,
	}))

	plan, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{{
		IDs: []domain.Identifier{10},
		Filters: &Filters{
			Expressions: []string{"ReferenceTime in ti(2023-01-01T00:00:00,2023-01-03T00:00:00)"},
			Parsed: FilterSet{Nodes: []domain.FilterNode{
				referenceTimeInterval("2023-01-01T00:00:00", "2023-01-03T00:00:00"),
			}},
		},
	}})
	if err != nil {
		t.Fatalf("build plan failed: %v", err)
	}

	if len(plan.Steps) != 3 {
		t.Fatalf("expected 3 split steps, got %d", len(plan.Steps))
	}
	if len(plan.Steps[0].Queries) != 1 {
		t.Fatalf("expected 1 query per step, got %d", len(plan.Steps[0].Queries))
	}
	query := plan.Steps[0].Queries[0]
	if query.DataCategory != domain.Curves {
		t.Fatalf("expected curves category, got %q", query.DataCategory)
	}
	if query.Source != domain.SourceCMDP {
		t.Fatalf("expected cmdp source, got %q", query.Source)
	}
	if query.IndexRange == nil || query.IndexRange.Start != 20221231 || query.IndexRange.End != 20230101 {
		t.Fatalf("unexpected first index range: %#v", query.IndexRange)
	}
}

func TestGenericPlannerUsesMappingCategories(t *testing.T) {
	planner := NewPlanner(
		WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{
			mappingWithColumns(10, domain.Curves),
			mappingWithColumns(20, domain.TimeSeries),
		}}),
		WithQueryBuilder(PlaceholderQueryBuilder{}),
	)

	plan, err := planner.BuildPlan(context.Background(), RequestContext{
		EndpointKind: EndpointGeneric,
		Stage:        "development",
		Mode:         ModeCSV,
	}, []Request{{
		IDs: []domain.Identifier{10, 20},
	}})
	if err != nil {
		t.Fatalf("build plan failed: %v", err)
	}

	if len(plan.Steps) != 2 {
		t.Fatalf("steps = %d, want one per mapping category", len(plan.Steps))
	}
	if got := plan.Steps[0].Command.DataCategory; got != domain.Curves {
		t.Fatalf("first command category = %q, want curves", got)
	}
	if got := plan.Steps[1].Command.DataCategory; got != domain.TimeSeries {
		t.Fatalf("second command category = %q, want timeseries", got)
	}
	if got := plan.Steps[0].Queries[0].DataCategory; got != domain.Curves {
		t.Fatalf("first query category = %q, want curves", got)
	}
	if got := plan.Steps[1].Queries[0].DataCategory; got != domain.TimeSeries {
		t.Fatalf("second query category = %q, want timeseries", got)
	}
}
