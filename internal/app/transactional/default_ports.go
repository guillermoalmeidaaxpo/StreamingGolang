package transactional

import (
	"context"
	"fmt"
	"time"

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

func (StaticMappingResolver) GetWatermark(_ context.Context, _ []domain.Mapping) (time.Time, error) {
	return time.Now().UTC(), nil
}

type PlaceholderQueryBuilder struct{}

func (PlaceholderQueryBuilder) BuildQueries(_ context.Context, command Command) ([]domain.ExecutableQuery, error) {
	mappingByID := make(map[domain.Identifier]domain.Mapping, len(command.Mappings))
	for _, mapping := range command.Mappings {
		mappingByID[mapping.ID] = mapping
	}

	queries := make([]domain.ExecutableQuery, 0, len(command.IDs))
	for _, id := range command.IDs {
		mapping := mappingByID[id]
		queries = append(queries, domain.ExecutableQuery{
			ID:           id,
			DataCategory: dataCategoryForQuery(command.DataCategory, mapping),
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

func dataCategoryForQuery(commandCategory domain.DataCategory, mapping domain.Mapping) domain.DataCategory {
	if mapping.DataCategory != "" {
		return mapping.DataCategory
	}
	return commandCategory
}
