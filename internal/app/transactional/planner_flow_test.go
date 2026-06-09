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
