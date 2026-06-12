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
		"ReferenceTime = latest(ReferenceTime > 2023-01-01T00:00:00)",
		"rankover([DeliveryStart],[ReferenceTime desc],[1,last])",
	}, "Europe/Zurich")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if got, want := len(filters.Nodes), 4; got != want {
		t.Fatalf("node count = %d, want %d", got, want)
	}
}

func TestParserBuildsComparisonAST(t *testing.T) {
	parser := New()

	filters, err := parser.Parse(context.Background(), []string{
		"ReferenceTime in ti(2023-05-21T00:00:00,2023-05-21T23:59:59)",
		"ReferenceTime = latest(ReferenceTime > 2023-01-01T00:00:00)",
	}, "Europe/Zurich")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	interval, ok := filters.Nodes[0].(domain.ComparisonFilter)
	if !ok {
		t.Fatalf("node 0 type = %T, want domain.ComparisonFilter", filters.Nodes[0])
	}
	if interval.Field != "ReferenceTime" || interval.Operator != ">=" {
		t.Fatalf("interval lower comparison = %#v", interval)
	}
	if interval.Value.Kind != domain.FilterValuePointInTime {
		t.Fatalf("interval lower kind = %q, want %q", interval.Value.Kind, domain.FilterValuePointInTime)
	}
	if interval.Value.Raw != "2023-05-20T22:00:00Z" {
		t.Fatalf("interval lower raw = %q", interval.Value.Raw)
	}

	upper, ok := filters.Nodes[1].(domain.ComparisonFilter)
	if !ok {
		t.Fatalf("node 1 type = %T, want domain.ComparisonFilter", filters.Nodes[1])
	}
	if upper.Field != "ReferenceTime" || upper.Operator != "<" || upper.Value.Raw != "2023-05-21T21:59:59Z" {
		t.Fatalf("interval upper comparison = %#v", upper)
	}

	latest, ok := filters.Nodes[2].(domain.ComparisonFilter)
	if !ok {
		t.Fatalf("node 2 type = %T, want domain.ComparisonFilter", filters.Nodes[2])
	}
	if latest.Value.Kind != domain.FilterValueLatest {
		t.Fatalf("latest kind = %q, want %q", latest.Value.Kind, domain.FilterValueLatest)
	}
	if got, want := len(latest.Value.Arguments), 1; got != want {
		t.Fatalf("latest argument count = %d, want %d", got, want)
	}
	if latest.Value.Arguments[0].Field != "ReferenceTime" {
		t.Fatalf("latest argument 0 = %#v", latest.Value.Arguments[0])
	}
	if latest.Value.Arguments[0].Value.Raw != "2022-12-31T23:00:00Z" {
		t.Fatalf("latest argument 0 raw = %q", latest.Value.Arguments[0].Value.Raw)
	}
}

func TestParserBuildsRankOverAST(t *testing.T) {
	parser := New()

	filters, err := parser.Parse(context.Background(), []string{
		"rankover([DeliveryStart,QuoteIndex],[ReferenceTime desc,DeliveryStart asc],[1,last])",
	}, "")
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

	_, err := parser.Parse(context.Background(), []string{"ReferenceTime between 2023-01-01T00:00:00"}, "")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParserRejectsLatestWithMultipleParameters(t *testing.T) {
	parser := New()

	_, err := parser.Parse(context.Background(), []string{
		"ReferenceTime = latest(ReferenceTime > 2023-01-01T00:00:00, QuoteIndex >= 1)",
	}, "Europe/Zurich")
	if err == nil {
		t.Fatal("expected latest parameter count error")
	}
}

func TestParserRejectsLatestOnNonReferenceTime(t *testing.T) {
	parser := New()

	_, err := parser.Parse(context.Background(), []string{
		"DeliveryStart = latest(ReferenceTime > 2023-01-01T00:00:00)",
	}, "Europe/Zurich")
	if err == nil {
		t.Fatal("expected latest target column error")
	}
}

func TestParserNormalizesPointInTimeInFilterTimeZone(t *testing.T) {
	parser := New()

	filters, err := parser.Parse(context.Background(), []string{
		"ReferenceTime = 2024-04-26T00:00:00",
	}, "Europe/Zurich")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	comparison, ok := filters.Nodes[0].(domain.ComparisonFilter)
	if !ok {
		t.Fatalf("node 0 type = %T, want domain.ComparisonFilter", filters.Nodes[0])
	}
	if comparison.Value.Raw != "2024-04-25T22:00:00Z" {
		t.Fatalf("raw point-in-time = %q", comparison.Value.Raw)
	}
}

func TestParserAppliesPointInTimeArithmetic(t *testing.T) {
	parser := New()

	filters, err := parser.Parse(context.Background(), []string{
		"ReferenceTime = 2024-04-26T00:00:00+P1D",
	}, "Europe/Zurich")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	comparison, ok := filters.Nodes[0].(domain.ComparisonFilter)
	if !ok {
		t.Fatalf("node 0 type = %T, want domain.ComparisonFilter", filters.Nodes[0])
	}
	if comparison.Value.Raw != "2024-04-26T22:00:00Z" {
		t.Fatalf("raw point-in-time = %q", comparison.Value.Raw)
	}
}

func TestParserExpandsTimeIntervalFunction(t *testing.T) {
	parser := New()

	filters, err := parser.Parse(context.Background(), []string{
		"ReferenceTime in tiDay(2024-04-26T00:00:00)",
	}, "Europe/Zurich")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if got, want := len(filters.Nodes), 2; got != want {
		t.Fatalf("node count = %d, want %d", got, want)
	}
	lower := filters.Nodes[0].(domain.ComparisonFilter)
	upper := filters.Nodes[1].(domain.ComparisonFilter)
	if lower.Operator != ">=" || lower.Value.Raw != "2024-04-25T22:00:00Z" {
		t.Fatalf("lower = %#v", lower)
	}
	if upper.Operator != "<" || upper.Value.Raw != "2024-04-26T22:00:00Z" {
		t.Fatalf("upper = %#v", upper)
	}
}
