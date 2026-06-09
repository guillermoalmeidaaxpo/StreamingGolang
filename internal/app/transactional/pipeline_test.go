package transactional

import (
	"context"
	"testing"

	"streaming-golang/internal/domain"
)

func TestPipelineExecutePlansRequestedIDs(t *testing.T) {
	pipeline := NewPipeline(NewValidator(), testFilterParser{}, NewPlanner(), NewExecutor(map[domain.SourceKind]Repository{
		domain.SourceCMDP: testRepository{},
	}, 0))

	response, err := pipeline.Execute(context.Background(), RequestContext{
		DataCategory: domain.TimeSeries,
		Stage:        "development",
		Mode:         ModeJSON,
	}, []Request{{IDs: []domain.Identifier{10, 20}}})
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if len(response.TransactionalData) != 2 {
		t.Fatalf("expected 2 data items, got %d", len(response.TransactionalData))
	}

	if response.ReferenceData[0] != 10 || response.ReferenceData[1] != 20 {
		t.Fatalf("unexpected requested IDs: %#v", response.ReferenceData)
	}
}

type testFilterParser struct{}

func (testFilterParser) Parse(_ context.Context, expressions []string) (FilterSet, error) {
	return FilterSet{Expressions: expressions}, nil
}

type testRepository struct{}

func (testRepository) Execute(_ context.Context, query domain.ExecutableQuery) ([]DataItem, error) {
	return []DataItem{{
		ID: query.ID,
		Fields: map[string]any{
			"status":         "planned",
			"source":         query.Source,
			"dataCategory":   query.DataCategory,
			"statement":      query.Statement,
			"parameterCount": len(query.Parameters),
		},
	}}, nil
}

func (r testRepository) Stream(ctx context.Context, query domain.ExecutableQuery) (Stream, error) {
	items, err := r.Execute(ctx, query)
	if err != nil {
		return nil, err
	}
	return &sliceStream{items: items}, nil
}
