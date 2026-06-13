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

	if len(plan.Steps) != 4 {
		t.Fatalf("expected 4 split steps, got %d", len(plan.Steps))
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
	if query.IndexRange == nil || query.IndexRange.Start != 20221230 || query.IndexRange.End != 20221231 {
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

func TestPlannerBuildsAggregationCommandColumnsLikeCSharp(t *testing.T) {
	planner := NewPlanner(
		WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{{
			ID:           10,
			DataCategory: domain.Curves,
			Source:       domain.SourceCMDP,
			ViewName:     "CurveView",
			Columns: []domain.ColumnMapping{
				{MDSName: "ReferenceTime", SourceName: "QuoteTime", IsKey: true},
				{MDSName: "DeliveryStart", SourceName: "DeliveryStart", IsKey: true},
				{MDSName: "DeliveryEnd", SourceName: "DeliveryEnd", IsKey: true},
				{MDSName: "Value", SourceName: "Value", IsProjectable: true},
			},
		}}}),
		WithQueryBuilder(PlaceholderQueryBuilder{}),
	)

	plan, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		Stage:        "development",
		Mode:         ModeCSV,
	}, []Request{{
		IDs: []domain.Identifier{10},
		Transformations: &Transformations{
			Keys:   []string{"Aggregate(Delivery, PT1H)=DeliveryBucket"},
			Values: [][]string{{"AVG(Value)", "AverageValue"}},
		},
	}})
	if err != nil {
		t.Fatalf("build plan failed: %v", err)
	}
	if len(plan.Steps) != 1 {
		t.Fatalf("steps = %d, want 1", len(plan.Steps))
	}
	command := plan.Steps[0].Command
	if command.Aggregations == nil {
		t.Fatal("expected aggregation metadata")
	}
	if got := command.Aggregations.GroupBy[0].Expression; got != "Aggregate(DeliveryStart, PT1H)" {
		t.Fatalf("group expression = %q", got)
	}
	if got := command.TargetTimeZone; got != "UTC" {
		t.Fatalf("target timezone = %q, want UTC", got)
	}
	wantColumns := []string{"Identifier", "ReferenceTime", "DeliveryStart", "DeliveryEnd", "RelativeDeliveryPeriod", "LegacyDeliveryBucketNumber", "DeliveryBucket", "AverageValue"}
	if !equalStringSlices(command.Columns, wantColumns) {
		t.Fatalf("columns = %#v, want %#v", command.Columns, wantColumns)
	}
}

func equalStringSlices(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
