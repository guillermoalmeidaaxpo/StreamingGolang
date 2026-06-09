package cassandra

import (
	"context"
	"fmt"
	"strings"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/domain"
)

type CassandraQueryBuilder struct {
	tableMappings map[string]string
}

func NewCassandraQueryBuilder(tableMappings map[string]string) *CassandraQueryBuilder {
	return &CassandraQueryBuilder{tableMappings: tableMappings}
}

func (b *CassandraQueryBuilder) BuildQueries(_ context.Context, command domain.Command) ([]domain.ExecutableQuery, error) {
	if len(command.Mappings) == 0 {
		return nil, apperr.New(apperr.Invalid, "cannot build Cassandra query without mappings")
	}

	queries := make([]domain.ExecutableQuery, 0)
	for _, mapping := range command.Mappings {
		if mapping.Source != domain.SourceCassandra && mapping.CassandraID == "" {
			continue
		}

		table, ok := b.resolveTable(mapping)
		if !ok {
			return nil, apperr.New(apperr.Invalid, fmt.Sprintf("no Cassandra table mapping for %q", mapping.DataCategory))
		}

		statement, arguments, err := b.buildStatement(table, mapping, command)
		if err != nil {
			return nil, err
		}

		queries = append(queries, domain.ExecutableQuery{
			ID:           mapping.ID,
			DataCategory: command.DataCategory,
			Source:       domain.SourceCassandra,
			Statement:    statement,
			Arguments:    arguments,
		})
	}

	return queries, nil
}

func (b *CassandraQueryBuilder) resolveTable(mapping domain.Mapping) (string, bool) {
	// Usually mappings like 'power' -> 'hpfc'
	// We might need more complex logic here if category isn't enough
	table, ok := b.tableMappings[string(mapping.DataCategory)]
	if ok {
		return table, true
	}
	// Fallback to category name
	return string(mapping.DataCategory), true
}

func (b *CassandraQueryBuilder) buildStatement(table string, mapping domain.Mapping, command domain.Command) (string, []any, error) {
	if mapping.CassandraID == "" {
		return "", nil, apperr.New(apperr.Invalid, fmt.Sprintf("mapping %d has no Cassandra ID", mapping.ID))
	}

	arguments := []any{mapping.CassandraID}
	where := []string{"id = ?"}

	if len(command.QuoteIndices) > 0 {
		arguments = append(arguments, command.QuoteIndices)
		where = append(where, "quote_date_index IN ?")
	}

	// Dynamic columns
	columns := "*"
	if len(command.Columns) > 0 {
		columns = strings.Join(command.Columns, ", ")
	}

	statement := fmt.Sprintf("SELECT %s FROM %s WHERE %s", columns, table, strings.Join(where, " AND "))
	return statement, arguments, nil
}
