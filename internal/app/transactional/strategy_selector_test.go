package transactional

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"streaming-golang/internal/domain"
)

func TestPlannerSelectsCassandraForNormalZurichCassandraMapping(t *testing.T) {
	plan := buildStrategyPlan(t, domain.Mapping{
		ID:           536013751,
		DataCategory: domain.Curves,
		Source:       domain.SourceCassandra,
		CassandraID:  "price_modelled:536013751",
		SplitQuery:   true,
	})

	assertStrategySource(t, plan, domain.SourceCassandra)
}

func TestPlannerSelectsCMDPForSpecialHPFCIDs(t *testing.T) {
	plan := buildStrategyPlan(t, domain.Mapping{
		ID:           536346251,
		DataCategory: domain.Curves,
		Source:       domain.SourceCassandra,
		CassandraID:  "hpfc:536346251",
		ViewName:     "ACCESS.Data_HPFC",
		IndexField:   "QuoteDateIndex_FID",
		SplitQuery:   true,
	})

	assertStrategySource(t, plan, domain.SourceCMDP)
}

func TestPlannerDoesNotHybridSplitSpecialHPFCIDs(t *testing.T) {
	resolver := &watermarkResolver{
		fixedMappingResolver: fixedMappingResolver{mappings: []domain.Mapping{{
			ID:           536346251,
			DataCategory: domain.Curves,
			Source:       domain.SourceCassandra,
			CassandraID:  "hpfc:536346251",
			ViewName:     "ACCESS.Data_HPFC",
			IndexField:   "QuoteDateIndex_FID",
			SplitQuery:   true,
		}}},
		watermark: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	planner := NewPlanner(
		WithMappingResolver(resolver),
		WithQueryBuilder(PlaceholderQueryBuilder{}),
	)

	plan, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{{
		IDs: []domain.Identifier{536346251},
		Filters: &Filters{
			FilterTimeZone: "Europe/Zurich",
			Parsed: FilterSet{Nodes: []domain.FilterNode{
				referenceTimePoint("=", "2024-04-26T00:00:00"),
			}},
		},
	}})
	if err != nil {
		t.Fatalf("build plan failed: %v", err)
	}

	assertStrategySource(t, plan, domain.SourceCMDP)
	if resolver.watermarkCalled {
		t.Fatal("special HPFC ID should stay on CMDP and not call hybrid watermark")
	}
}

func TestPlannerSelectsCMDPForAggregations(t *testing.T) {
	plan := buildStrategyPlan(t, domain.Mapping{
		ID:           536013751,
		DataCategory: domain.Curves,
		Source:       domain.SourceCassandra,
		CassandraID:  "price_modelled:536013751",
		ViewName:     "ACCESS.Data_PriceModelled",
		IndexField:   "QuoteDateIndex_FID",
		SplitQuery:   true,
	}, func(request *Request) {
		request.Transformations = &Transformations{Keys: []string{"ReferenceTime"}}
	})

	assertStrategySource(t, plan, domain.SourceCMDP)
}

func TestPlannerSelectsCMDPForShapeFilters(t *testing.T) {
	plan := buildStrategyPlan(t, domain.Mapping{
		ID:           536013751,
		DataCategory: domain.Curves,
		Source:       domain.SourceCassandra,
		CassandraID:  "price_modelled:536013751",
		ViewName:     "ACCESS.Data_PriceModelled",
		IndexField:   "QuoteDateIndex_FID",
		SplitQuery:   true,
	}, func(request *Request) {
		request.Filters = &Filters{Shape: json.RawMessage(`{"kind":"period"}`)}
	})

	assertStrategySource(t, plan, domain.SourceCMDP)
}

func TestPlannerKeepsHyperscaleAheadOfCMDPRules(t *testing.T) {
	hyperscaleID := domain.Identifier(1000000001)
	plan := buildStrategyPlan(t, domain.Mapping{
		ID:           1000000001,
		DataCategory: domain.Surfaces,
		Source:       domain.SourceHyperscale,
		HyperscaleID: &hyperscaleID,
	})

	assertStrategySource(t, plan, domain.SourceHyperscale)
}

func TestPlannerSelectsCMDPForNonZurichCassandraTimeZone(t *testing.T) {
	plan := buildStrategyPlan(t, domain.Mapping{
		ID:           536960251,
		DataCategory: domain.Curves,
		Source:       domain.SourceCassandra,
		CassandraID:  "curve:536960251",
		ViewName:     "ACCESS.Data_Singapore",
		IndexField:   "QuoteDateIndex_FID",
		SplitQuery:   true,
	})

	assertStrategySource(t, plan, domain.SourceCMDP)
}

func buildStrategyPlan(t *testing.T, mapping domain.Mapping, mutateRequest ...func(*Request)) Plan {
	t.Helper()

	planner := NewPlanner(
		WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{mapping}}),
		WithQueryBuilder(PlaceholderQueryBuilder{}),
	)

	request := Request{IDs: []domain.Identifier{mapping.ID}}
	for _, mutate := range mutateRequest {
		mutate(&request)
	}

	plan, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: mapping.DataCategory,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{request})
	if err != nil {
		t.Fatalf("build plan failed: %v", err)
	}
	return plan
}

func assertStrategySource(t *testing.T, plan Plan, want domain.SourceKind) {
	t.Helper()

	if len(plan.Steps) != 1 {
		t.Fatalf("expected 1 plan step, got %d", len(plan.Steps))
	}
	command := plan.Steps[0].Command
	if command.Source != want {
		t.Fatalf("command source = %q, want %q", command.Source, want)
	}
	if command.Mappings[0].Source != want {
		t.Fatalf("mapping source = %q, want %q", command.Mappings[0].Source, want)
	}
	if plan.Steps[0].Queries[0].Source != want {
		t.Fatalf("query source = %q, want %q", plan.Steps[0].Queries[0].Source, want)
	}
}
