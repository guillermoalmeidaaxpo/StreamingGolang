package transactional

import (
	"context"
	"testing"

	"streaming-golang/internal/domain"
)

func TestValidatorRejectsDuplicateIDsAcrossRequests(t *testing.T) {
	validator := NewValidator()

	err := validator.Validate(context.Background(), []Request{
		{IDs: []domain.Identifier{10}},
		{IDs: []domain.Identifier{10}},
	})

	if err == nil {
		t.Fatal("expected duplicate ID validation error")
	}
}

func TestValidatorRejectsInvalidProjectionColumnName(t *testing.T) {
	validator := NewValidator()

	err := validator.Validate(context.Background(), []Request{{
		IDs:     []domain.Identifier{10},
		Columns: []string{"Reference-Time"},
	}})

	if err == nil {
		t.Fatal("expected invalid column validation error")
	}
}

func TestValidatorRejectsInvalidTimeZone(t *testing.T) {
	validator := NewValidator()

	err := validator.Validate(context.Background(), []Request{{
		IDs: []domain.Identifier{10},
		Transformations: &Transformations{
			TargetTimeZone: "Not/AZone",
		},
	}})

	if err == nil {
		t.Fatal("expected invalid timezone validation error")
	}
}

func TestValidatorRequiresAggregationKeysAndValuesTogether(t *testing.T) {
	validator := NewValidator()

	err := validator.Validate(context.Background(), []Request{{
		IDs: []domain.Identifier{10},
		Transformations: &Transformations{
			Keys: []string{"Aggregate(Delivery, PT1H)"},
		},
	}})

	if err == nil {
		t.Fatal("expected aggregation keys/values validation error")
	}
}
