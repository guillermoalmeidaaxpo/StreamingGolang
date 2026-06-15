package transactional

import (
	"context"
	"encoding/json"
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

func TestPlannerRejectsRankOverForTimeseries(t *testing.T) {
	planner := NewPlanner(WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{rankOverMapping(10, domain.TimeSeries)}}))

	_, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.TimeSeries,
		EndpointKind: EndpointTransactional,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{rankOverRequest(10)})

	if err == nil || !strings.Contains(err.Error(), "rankover filters are only available for curves and surfaces") {
		t.Fatalf("err = %v, want rankover data-category validation error", err)
	}
}

func TestPlannerRejectsRankOverForCassandraHostedID(t *testing.T) {
	mapping := rankOverMapping(10, domain.Curves)
	mapping.CassandraID = "curve:10"
	planner := NewPlanner(WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{mapping}}))

	_, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		EndpointKind: EndpointTransactional,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{rankOverRequest(10)})

	if err == nil || !strings.Contains(err.Error(), "rankover filters are only supported for CMDP-hosted IDs") {
		t.Fatalf("err = %v, want rankover CMDP-hosted validation error", err)
	}
}

func TestPlannerRejectsRankOverPartitionColumnThatIsNotKey(t *testing.T) {
	mapping := rankOverMapping(10, domain.Curves)
	mapping.Columns = append(mapping.Columns, domain.ColumnMapping{
		MDSName:       "Bucket",
		SourceName:    "Bucket",
		IsKey:         false,
		IsProjectable: true,
	})
	planner := NewPlanner(WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{mapping}}))

	request := rankOverRequest(10)
	request.Filters.Parsed.Nodes = []domain.FilterNode{
		domain.RankOverFilter{
			PartitionBy: []string{"Bucket"},
			OrderBy:     []domain.SortExpression{{Field: "ReferenceTime", Direction: "desc"}},
		},
	}

	_, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		EndpointKind: EndpointTransactional,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{request})

	if err == nil || !strings.Contains(err.Error(), "invalid rankover partition column") {
		t.Fatalf("err = %v, want invalid rankover partition column error", err)
	}
}

func TestShapeNormalizerMatchesCSharpTokens(t *testing.T) {
	shape, err := normalizeShape(json.RawMessage(`{
		"months":["Mar","Jan"],
		"days":["Sun","Mon"],
		"time":[{"start":"T08:00:00","end":"T10:30:00"}],
		"holidayCalendar": 42
	}`))
	if err != nil {
		t.Fatalf("normalize shape failed: %v", err)
	}
	if got := shape.Months; len(got) != 2 || got[0] != 1 || got[1] != 3 {
		t.Fatalf("months = %#v, want [1 3]", got)
	}
	if got := shape.Days; len(got) != 2 || got[0] != 1 || got[1] != 7 {
		t.Fatalf("days = %#v, want [1 7]", got)
	}
	if got := shape.TimeSpans[0]; got.StartSeconds != 8*3600 || got.EndSeconds != 10*3600+30*60 {
		t.Fatalf("time span = %#v", got)
	}
	if shape.HolidayCalendar == nil || *shape.HolidayCalendar != 42 {
		t.Fatalf("holiday calendar = %#v", shape.HolidayCalendar)
	}
}

func TestValidatorRejectsOverlappingShapeTimeRanges(t *testing.T) {
	validator := NewValidator()
	err := validator.Validate(context.Background(), []Request{{
		IDs: []domain.Identifier{10},
		Filters: &Filters{Shape: json.RawMessage(`{
			"time":[
				{"start":"T08:00:00","end":"T10:00:00"},
				{"start":"T09:00:00","end":"T11:00:00"}
			]
		}`)},
	}})
	if err == nil || !strings.Contains(err.Error(), "overlapping shape time ranges") {
		t.Fatalf("err = %v, want overlapping shape time ranges", err)
	}
}

func TestPlannerCarriesNormalizedShape(t *testing.T) {
	planner := NewPlanner(
		WithMappingResolver(fixedMappingResolver{mappings: []domain.Mapping{rankOverMapping(10, domain.Curves)}}),
		WithQueryBuilder(PlaceholderQueryBuilder{}),
	)

	plan, err := planner.BuildPlan(context.Background(), RequestContext{
		DataCategory: domain.Curves,
		EndpointKind: EndpointTransactional,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{{
		IDs:     []domain.Identifier{10},
		Filters: &Filters{Shape: json.RawMessage(`{"months":["Jan"],"days":["Mon"],"time":[{"start":"T00:00:00","end":"T00:00:00"}]}`)},
	}})
	if err != nil {
		t.Fatalf("build plan failed: %v", err)
	}
	shape := plan.Steps[0].Command.Shape
	if shape == nil {
		t.Fatal("expected normalized shape")
	}
	if len(shape.Months) != 1 || shape.Months[0] != 1 {
		t.Fatalf("months = %#v", shape.Months)
	}
	if len(shape.Days) != 1 || shape.Days[0] != 1 {
		t.Fatalf("days = %#v", shape.Days)
	}
	if len(shape.TimeSpans) != 1 || shape.TimeSpans[0].EndSeconds != 24*3600 {
		t.Fatalf("time spans = %#v", shape.TimeSpans)
	}
}

func rankOverRequest(id domain.Identifier) Request {
	return Request{
		IDs: []domain.Identifier{id},
		Filters: &Filters{Parsed: FilterSet{Nodes: []domain.FilterNode{
			domain.RankOverFilter{
				PartitionBy: []string{"DeliveryStart"},
				OrderBy:     []domain.SortExpression{{Field: "ReferenceTime", Direction: "desc"}},
			},
		}}},
	}
}

func rankOverMapping(id domain.Identifier, category domain.DataCategory) domain.Mapping {
	return domain.Mapping{
		ID:           id,
		DataCategory: category,
		Source:       domain.SourceCMDP,
		ViewName:     "ACCESS.Data_RankOver",
		Columns: []domain.ColumnMapping{
			{MDSName: "ReferenceTime", SourceName: "QuoteTime", IsKey: true, IsProjectable: true},
			{MDSName: "DeliveryStart", SourceName: "AdjustedDeliveryStartDate", IsKey: true, IsProjectable: true},
			{MDSName: "Value", SourceName: "Value", IsProjectable: true},
		},
	}
}
