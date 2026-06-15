package transactional

import (
	"context"
	"testing"
	"time"

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

func TestTransformationProcessorCalculatesCassandraRelativeDeliveryPeriodByDefault(t *testing.T) {
	processor := NewTransformationProcessor()
	item := DataItem{
		ID: 536013751,
		Fields: map[string]any{
			"ReferenceTime":              time.Date(2024, 4, 26, 0, 0, 0, 0, time.UTC),
			"DeliveryStart":              time.Date(2024, 4, 26, 3, 0, 0, 0, time.UTC),
			"RelativeDeliveryPeriod":     nil,
			"LegacyDeliveryBucketNumber": nil,
			"Value":                      115.9,
		},
	}
	command := Command{
		Source: domain.SourceCassandra,
		Mappings: []domain.Mapping{{
			Resolution: "PT1H",
		}},
		Columns: []string{"Value"},
	}

	processed := processor.Process(context.Background(), []DataItem{item}, command)

	if len(processed) != 1 {
		t.Fatalf("items = %d, want 1", len(processed))
	}
	if got := processed[0].Fields["RelativeDeliveryPeriod"]; got != int64(3) {
		t.Fatalf("RelativeDeliveryPeriod = %#v, want 3", got)
	}
}
