package antlrparser

import (
	"context"
	"testing"

	"streaming-golang/internal/domain"
)

func TestParserAcceptsCurrentFilterSyntax(t *testing.T) {
	parser := New()

	filters, err := parser.Parse(context.Background(), []string{
		"ReferenceTime in ti(2023-05-21T00:00:00,2023-05-21T23:59:59)",
		"ReferenceTime = latest(DeliveryStart > 2023-01-01T00:00:00)",
		"rankover([DeliveryStart],[ReferenceTime desc],[1,last])",
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if got, want := len(filters.Nodes), 3; got != want {
		t.Fatalf("node count = %d, want %d", got, want)
	}
}

func TestParserBuildsComparisonAST(t *testing.T) {
	parser := New()

	filters, err := parser.Parse(context.Background(), []string{
		"ReferenceTime in ti(2023-05-21T00:00:00,2023-05-21T23:59:59)",
		"ReferenceTime = latest(DeliveryStart > 2023-01-01T00:00:00, QuoteIndex >= 1)",
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	interval, ok := filters.Nodes[0].(domain.ComparisonFilter)
	if !ok {
		t.Fatalf("node 0 type = %T, want domain.ComparisonFilter", filters.Nodes[0])
	}
	if interval.Field != "ReferenceTime" || interval.Operator != "in" {
		t.Fatalf("interval comparison = %#v", interval)
	}
	if interval.Value.Kind != domain.FilterValueTimeInterval {
		t.Fatalf("interval kind = %q, want %q", interval.Value.Kind, domain.FilterValueTimeInterval)
	}
	if interval.Value.Start != "2023-05-21T00:00:00" || interval.Value.End != "2023-05-21T23:59:59" {
		t.Fatalf("interval bounds = %q/%q", interval.Value.Start, interval.Value.End)
	}

	latest, ok := filters.Nodes[1].(domain.ComparisonFilter)
	if !ok {
		t.Fatalf("node 1 type = %T, want domain.ComparisonFilter", filters.Nodes[1])
	}
	if latest.Value.Kind != domain.FilterValueLatest {
		t.Fatalf("latest kind = %q, want %q", latest.Value.Kind, domain.FilterValueLatest)
	}
	if got, want := len(latest.Value.Arguments), 2; got != want {
		t.Fatalf("latest argument count = %d, want %d", got, want)
	}
	if latest.Value.Arguments[0].Field != "DeliveryStart" {
		t.Fatalf("latest argument 0 = %#v", latest.Value.Arguments[0])
	}
	if latest.Value.Arguments[1].Value.Kind != domain.FilterValueNumber {
		t.Fatalf("latest argument 1 kind = %q, want %q", latest.Value.Arguments[1].Value.Kind, domain.FilterValueNumber)
	}
}

func TestParserBuildsRankOverAST(t *testing.T) {
	parser := New()

	filters, err := parser.Parse(context.Background(), []string{
		"rankover([DeliveryStart,QuoteIndex],[ReferenceTime desc,DeliveryStart asc],[1,last])",
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	rankOver, ok := filters.Nodes[0].(domain.RankOverFilter)
	if !ok {
		t.Fatalf("node 0 type = %T, want domain.RankOverFilter", filters.Nodes[0])
	}
	if got, want := len(rankOver.PartitionBy), 2; got != want {
		t.Fatalf("partition count = %d, want %d", got, want)
	}
	if got, want := len(rankOver.OrderBy), 2; got != want {
		t.Fatalf("order count = %d, want %d", got, want)
	}
	if rankOver.OrderBy[0].Field != "ReferenceTime" || rankOver.OrderBy[0].Direction != "desc" {
		t.Fatalf("order 0 = %#v", rankOver.OrderBy[0])
	}
	if got, want := len(rankOver.Bounds), 1; got != want {
		t.Fatalf("bound count = %d, want %d", got, want)
	}
	if rankOver.Bounds[0].Start != "1" || rankOver.Bounds[0].End != "last" {
		t.Fatalf("bound = %#v", rankOver.Bounds[0])
	}
}

func TestParserRejectsInvalidFilterSyntax(t *testing.T) {
	parser := New()

	_, err := parser.Parse(context.Background(), []string{"ReferenceTime between 2023-01-01T00:00:00"})
	if err == nil {
		t.Fatal("expected parse error")
	}
}
