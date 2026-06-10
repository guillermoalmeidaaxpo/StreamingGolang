package transactional

import (
	"context"
	"time"

	"streaming-golang/internal/domain"
)

type MappingResolver interface {
	ResolveMappings(context.Context, []domain.Identifier, domain.DataCategory, string) ([]domain.Mapping, error)
	GetWatermark(ctx context.Context, mappings []domain.Mapping) (time.Time, error)
}

type QuoteIndexPlanner interface {
	PlanQuoteIndices(context.Context, Command) ([]int, error)
}

type QueryBuilder interface {
	BuildQueries(context.Context, Command) ([]domain.ExecutableQuery, error)
}

type Repository interface {
	Execute(ctx context.Context, query domain.ExecutableQuery) ([]DataItem, error)
	Stream(ctx context.Context, query domain.ExecutableQuery) (Stream, error)
}

type CompositeQueryBuilder struct {
	builders []QueryBuilder
}

func NewCompositeQueryBuilder(builders ...QueryBuilder) QueryBuilder {
	return &CompositeQueryBuilder{builders: builders}
}

func (c *CompositeQueryBuilder) BuildQueries(ctx context.Context, command Command) ([]domain.ExecutableQuery, error) {
	allQueries := make([]domain.ExecutableQuery, 0)
	for _, builder := range c.builders {
		queries, err := builder.BuildQueries(ctx, command)
		if err != nil {
			return nil, err
		}
		allQueries = append(allQueries, queries...)
	}
	return allQueries, nil
}
