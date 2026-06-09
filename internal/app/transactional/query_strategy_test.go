package transactional

import (
	"testing"

	"streaming-golang/internal/domain"
)

func TestSingleQueryStrategyKeepsOneCommand(t *testing.T) {
	command := Command{
		IDs:          []domain.Identifier{1},
		QuoteIndices: []int{20230101, 20230102},
	}

	commands := SingleQueryStrategy{}.Plan(command)
	if len(commands) != 1 {
		t.Fatalf("expected one command, got %d", len(commands))
	}
	if commands[0].IndexRange != nil {
		t.Fatal("single query strategy should not add an index range")
	}
}

func TestSplitQueryStrategySplitsCurvesByReferenceTimeDays(t *testing.T) {
	command := Command{
		IDs:          []domain.Identifier{1},
		DataCategory: domain.Curves,
		QuoteIndices: []int{1, 2, 3, 4, 5},
	}

	commands := SplitQueryStrategy{ReferenceTimeSplitDays: 2}.Plan(command)
	if len(commands) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(commands))
	}
	if commands[0].IndexRange.Start != 1 || commands[0].IndexRange.End != 2 {
		t.Fatalf("unexpected first range: %#v", commands[0].IndexRange)
	}
	if commands[2].IndexRange.Start != 5 || commands[2].IndexRange.End != 5 {
		t.Fatalf("unexpected last range: %#v", commands[2].IndexRange)
	}
}

func TestSplitQueryStrategySplitsTimeseriesByQueryCount(t *testing.T) {
	command := Command{
		IDs:          []domain.Identifier{1},
		DataCategory: domain.TimeSeries,
		QuoteIndices: []int{1, 2, 3, 4, 5, 6},
	}

	commands := SplitQueryStrategy{QueriesCount: 3, ReferenceTimeSplitDays: 10}.Plan(command)
	if len(commands) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(commands))
	}
	if commands[1].IndexRange.Start != 3 || commands[1].IndexRange.End != 4 {
		t.Fatalf("unexpected second range: %#v", commands[1].IndexRange)
	}
}
