package transactional

import (
	"context"
	"testing"
	"time"

	"streaming-golang/internal/domain"
)

type watermarkResolver struct {
	fixedMappingResolver
	watermark time.Time
}

func (r watermarkResolver) GetWatermark(context.Context, []domain.Mapping) (time.Time, error) {
	return r.watermark, nil
}

func TestPlannerSplitsHybridCommand(t *testing.T) {
	watermark := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	resolver := watermarkResolver{
		fixedMappingResolver: fixedMappingResolver{mappings: []domain.Mapping{{
			ID:           536013751,
			DataCategory: domain.Curves,
			Source:       domain.SourceCassandra,
			CassandraID:  "test:1",
			ViewName:     "TestView",
			IndexField:   "QuoteDateIndex",
		}}},
		watermark: watermark,
	}

	planner := NewPlanner(
		WithMappingResolver(resolver),
		WithQueryBuilder(PlaceholderQueryBuilder{}),
	)

	// Request spanning across the watermark: 2023-12-30 to 2024-01-02
	plan, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{{
		IDs: []domain.Identifier{536013751},
		Filters: &Filters{
			Parsed: FilterSet{Nodes: []domain.FilterNode{
				referenceTimeInterval("2023-12-30T00:00:00", "2024-01-02T00:00:00"),
			}},
		},
	}})

	if err != nil {
		t.Fatalf("build plan failed: %v", err)
	}

	// Should produce 2 steps: one for Cassandra (< watermark) and one for CMDP (>= watermark)
	if len(plan.Steps) != 2 {
		t.Fatalf("expected 2 steps for hybrid split, got %d", len(plan.Steps))
	}

	// Step 1: Cassandra
	if plan.Steps[0].Command.Source != domain.SourceCassandra {
		t.Fatalf("expected step 0 source to be Cassandra, got %q", plan.Steps[0].Command.Source)
	}

	// Step 2: CMDP
	if plan.Steps[1].Command.Source != domain.SourceCMDP {
		t.Fatalf("expected step 1 source to be CMDP, got %q", plan.Steps[1].Command.Source)
	}
}
