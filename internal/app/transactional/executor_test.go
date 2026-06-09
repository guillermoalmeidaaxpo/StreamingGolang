package transactional

import (
	"context"
	"testing"

	"streaming-golang/internal/domain"
)

func TestExecutorReadsRowsFromRepository(t *testing.T) {
	executor := NewExecutor(map[domain.SourceKind]Repository{
		domain.SourceCMDP: testRepository{},
	}, 2)

	response, err := executor.Execute(context.Background(), Plan{Steps: []PlanStep{{
		Command: Command{IDs: []domain.Identifier{10}},
		Queries: []ExecutableQuery{{
			ID:           10,
			Source:       domain.SourceCMDP,
			DataCategory: domain.Curves,
			Statement:    "select Value from CurveView",
			Parameters:   map[string]any{"id": int64(10)},
		}},
	}}})
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if len(response.TransactionalData) != 1 {
		t.Fatalf("items = %d, want 1", len(response.TransactionalData))
	}
	if response.TransactionalData[0].Fields["statement"] != "select Value from CurveView" {
		t.Fatalf("fields = %#v", response.TransactionalData[0].Fields)
	}
	if len(response.ReferenceData) != 1 || response.ReferenceData[0] != 10 {
		t.Fatalf("referenceData = %#v, want [10]", response.ReferenceData)
	}
}

func TestExecutorRejectsMissingRepository(t *testing.T) {
	executor := NewExecutor(nil, 1)

	_, err := executor.Execute(context.Background(), Plan{Steps: []PlanStep{{
		Command: Command{IDs: []domain.Identifier{10}},
		Queries: []ExecutableQuery{{
			ID:     10,
			Source: domain.SourceCMDP,
		}},
	}}})

	if err == nil {
		t.Fatal("expected missing repository error")
	}
}

func TestExecutorReferenceDataDeduplicatesSplitSteps(t *testing.T) {
	executor := NewExecutor(map[domain.SourceKind]Repository{
		domain.SourceCMDP: testRepository{},
	}, 2)

	response, err := executor.Execute(context.Background(), Plan{Steps: []PlanStep{
		{
			Command: Command{IDs: []domain.Identifier{10}},
			Queries: []ExecutableQuery{{
				ID:     10,
				Source: domain.SourceCMDP,
			}},
		},
		{
			Command: Command{IDs: []domain.Identifier{10}},
			Queries: []ExecutableQuery{{
				ID:     10,
				Source: domain.SourceCMDP,
			}},
		},
	}})
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if len(response.ReferenceData) != 1 || response.ReferenceData[0] != 10 {
		t.Fatalf("referenceData = %#v, want deduplicated [10]", response.ReferenceData)
	}
}
