package cassandra

import (
	"context"

	"github.com/gocql/gocql"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
)

type repository struct {
	session *gocql.Session
}

func NewRepository(session *gocql.Session) transactional.Repository {
	return &repository{session: session}
}

func (r *repository) Execute(ctx context.Context, query domain.ExecutableQuery) ([]transactional.DataItem, error) {
	iter := r.session.Query(query.Statement).Bind(query.Arguments...).WithContext(ctx).Iter()
	
	items := make([]transactional.DataItem, 0)
	for {
		row := make(map[string]any)
		if !iter.MapScan(row) {
			break
		}
		items = append(items, transactional.DataItem{
			ID:     query.ID,
			Fields: row,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, apperr.Wrap(apperr.Unavailable, "execute cassandra query", err)
	}

	return items, nil
}

func (r *repository) Stream(ctx context.Context, query domain.ExecutableQuery) (transactional.Stream, error) {
	iter := r.session.Query(query.Statement).Bind(query.Arguments...).WithContext(ctx).Iter()
	
	return &cassandraStream{
		ctx:  ctx,
		iter: iter,
		id:   query.ID,
	}, nil
}

type cassandraStream struct {
	ctx  context.Context
	iter *gocql.Iter
	id   domain.Identifier
	item transactional.DataItem
	err  error
}

func (s *cassandraStream) Next(ctx context.Context) bool {
	if s.err != nil || s.ctx.Err() != nil || ctx.Err() != nil {
		return false
	}

	row := make(map[string]any)
	if !s.iter.MapScan(row) {
		return false
	}

	s.item = transactional.DataItem{
		ID:     s.id,
		Fields: row,
	}
	return true
}

func (s *cassandraStream) Item() transactional.DataItem {
	return s.item
}

func (s *cassandraStream) Err() error {
	return s.iter.Close()
}

func (s *cassandraStream) Close() error {
	return s.iter.Close()
}
