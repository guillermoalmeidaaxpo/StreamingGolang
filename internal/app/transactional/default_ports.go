package transactional

import (
	"context"
	"fmt"

	"streaming-golang/internal/domain"
)

type StaticMappingResolver struct{}

func (StaticMappingResolver) ResolveMappings(_ context.Context, ids []domain.Identifier, category domain.DataCategory, _ string) ([]domain.Mapping, error) {
	mappings := make([]domain.Mapping, 0, len(ids))
	for _, id := range ids {
		mappings = append(mappings, domain.Mapping{
			ID:           id,
			DataCategory: category,
			Source:       domain.SourceCMDP,
			ViewName:     fmt.Sprintf("pending_%s_view", category),
			IndexField:   "QuoteDateIndex_FID",
		})
	}
	return mappings, nil
}

type PlaceholderQueryBuilder struct{}

func (PlaceholderQueryBuilder) BuildQueries(_ context.Context, command Command) ([]domain.ExecutableQuery, error) {
	queries := make([]domain.ExecutableQuery, 0, len(command.IDs))
	for _, id := range command.IDs {
		queries = append(queries, domain.ExecutableQuery{
			ID:           id,
			DataCategory: command.DataCategory,
			Source:       command.Source,
			Filters:      command.Filters,
			IndexRange:   command.IndexRange,
			Statement:    "pending_query_generation",
			Parameters: map[string]any{
				"id": id,
			},
		})
	}
	return queries, nil
}
