package transactional

import (
	"streaming-golang/internal/domain"
)

type QueryStrategy interface {
	Plan(command Command) []Command
}

type SingleQueryStrategy struct{}

func (SingleQueryStrategy) Plan(command Command) []Command {
	return []Command{command}
}

type SplitQueryStrategy struct {
	QueriesCount           int
	ReferenceTimeSplitDays int
}

func (s SplitQueryStrategy) Plan(command Command) []Command {
	if len(command.QuoteIndices) == 0 {
		return []Command{command}
	}

	batches := s.batchQuoteIndices(command.DataCategory, command.QuoteIndices)
	commands := make([]Command, 0, len(batches))
	for _, batch := range batches {
		if len(batch) == 0 {
			continue
		}

		split := command
		split.QuoteIndices = append([]int(nil), batch...)
		split.IndexRange = &domain.IndexRange{
			Start: batch[0],
			End:   batch[len(batch)-1],
		}
		commands = append(commands, split)
	}

	if len(commands) == 0 {
		return []Command{command}
	}
	return commands
}

func (s SplitQueryStrategy) batchQuoteIndices(category domain.DataCategory, quoteIndices []int) [][]int {
	queryCount := s.QueriesCount
	if queryCount <= 0 {
		queryCount = 1
	}

	batchSize := s.ReferenceTimeSplitDays
	if category == domain.TimeSeries {
		batchSize = len(quoteIndices) / queryCount
	}
	if batchSize <= 0 {
		batchSize = 1
	}

	batches := make([][]int, 0, (len(quoteIndices)+batchSize-1)/batchSize)
	for start := 0; start < len(quoteIndices); start += batchSize {
		end := start + batchSize
		if end > len(quoteIndices) {
			end = len(quoteIndices)
		}
		batches = append(batches, quoteIndices[start:end])
	}
	return batches
}
