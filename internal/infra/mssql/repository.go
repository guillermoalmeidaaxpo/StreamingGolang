package mssql

import (
	"context"
	"database/sql"
	"log/slog"
	"sort"
	"time"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
)

type repository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewRepository(db *sql.DB, loggers ...*slog.Logger) transactional.Repository {
	logger := slog.Default()
	if len(loggers) > 0 && loggers[0] != nil {
		logger = loggers[0]
	}
	return &repository{db: db, logger: logger}
}

func (r *repository) Execute(ctx context.Context, query domain.ExecutableQuery) ([]transactional.DataItem, error) {
	args := r.namedArgs(query.Parameters)
	start := time.Now()
	r.logger.InfoContext(ctx, "executing mssql query",
		slog.Int64("identifier", int64(query.ID)),
		slog.String("source", string(query.Source)),
		slog.String("data_category", string(query.DataCategory)),
		slog.String("query", compactSQL(query.Statement)),
		slog.Any("parameters", sortedParameters(query.Parameters)),
	)

	rows, err := r.db.QueryContext(ctx, query.Statement, args...)
	if err != nil {
		r.logger.ErrorContext(ctx, "execute mssql query failed",
			slog.Int64("identifier", int64(query.ID)),
			slog.String("source", string(query.Source)),
			slog.String("data_category", string(query.DataCategory)),
			slog.String("query", compactSQL(query.Statement)),
			slog.Any("parameters", sortedParameters(query.Parameters)),
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		return nil, apperr.Wrap(apperr.Unavailable, "execute mssql query", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, apperr.Wrap(apperr.Internal, "get mssql columns", err)
	}

	items := make([]transactional.DataItem, 0)
	for rows.Next() {
		fields, err := r.scanRow(rows, cols)
		if err != nil {
			return nil, err
		}
		items = append(items, transactional.DataItem{
			ID:     query.ID,
			Fields: fields,
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "iterate mssql rows failed",
			slog.Int64("identifier", int64(query.ID)),
			slog.String("source", string(query.Source)),
			slog.String("data_category", string(query.DataCategory)),
			slog.String("query", compactSQL(query.Statement)),
			slog.Any("parameters", sortedParameters(query.Parameters)),
			slog.Int("row_count", len(items)),
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		return nil, apperr.Wrap(apperr.Unavailable, "iterate mssql rows", err)
	}

	r.logger.InfoContext(ctx, "mssql query executed",
		slog.Int64("identifier", int64(query.ID)),
		slog.String("source", string(query.Source)),
		slog.String("data_category", string(query.DataCategory)),
		slog.Int("row_count", len(items)),
		slog.Duration("duration", time.Since(start)),
	)

	return items, nil
}

func (r *repository) Stream(ctx context.Context, query domain.ExecutableQuery) (transactional.Stream, error) {
	args := r.namedArgs(query.Parameters)
	start := time.Now()
	r.logger.InfoContext(ctx, "streaming mssql query",
		slog.Int64("identifier", int64(query.ID)),
		slog.String("source", string(query.Source)),
		slog.String("data_category", string(query.DataCategory)),
		slog.String("query", compactSQL(query.Statement)),
		slog.Any("parameters", sortedParameters(query.Parameters)),
	)

	rows, err := r.db.QueryContext(ctx, query.Statement, args...)
	if err != nil {
		r.logger.ErrorContext(ctx, "stream mssql query failed",
			slog.Int64("identifier", int64(query.ID)),
			slog.String("source", string(query.Source)),
			slog.String("data_category", string(query.DataCategory)),
			slog.String("query", compactSQL(query.Statement)),
			slog.Any("parameters", sortedParameters(query.Parameters)),
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)
		return nil, apperr.Wrap(apperr.Unavailable, "stream mssql query", err)
	}

	cols, err := rows.Columns()
	if err != nil {
		rows.Close()
		return nil, apperr.Wrap(apperr.Internal, "get mssql columns", err)
	}

	return &mssqlStream{
		ctx:  ctx,
		rows: rows,
		cols: cols,
		id:   query.ID,
		repo: r,
	}, nil
}

func (r *repository) namedArgs(parameters map[string]any) []any {
	args := make([]any, 0, len(parameters))
	for name, value := range parameters {
		args = append(args, sql.Named(name, value))
	}
	return args
}

func sortedParameters(parameters map[string]any) []any {
	if len(parameters) == 0 {
		return nil
	}
	names := make([]string, 0, len(parameters))
	for name := range parameters {
		names = append(names, name)
	}
	sort.Strings(names)

	values := make([]any, 0, len(names))
	for _, name := range names {
		values = append(values, map[string]any{
			"name":  name,
			"value": parameters[name],
		})
	}
	return values
}

func (r *repository) scanRow(rows *sql.Rows, cols []string) (map[string]any, error) {
	values := make([]any, len(cols))
	pointers := make([]any, len(cols))
	for i := range values {
		pointers[i] = &values[i]
	}

	if err := rows.Scan(pointers...); err != nil {
		return nil, apperr.Wrap(apperr.Internal, "scan mssql row", err)
	}

	fields := make(map[string]any, len(cols))
	for i, col := range cols {
		val := values[i]
		if b, ok := val.([]byte); ok {
			val = string(b)
		}
		fields[col] = val
	}
	return fields, nil
}

type mssqlStream struct {
	ctx  context.Context
	rows *sql.Rows
	cols []string
	id   domain.Identifier
	item transactional.DataItem
	err  error
	repo *repository
}

func (s *mssqlStream) Next(ctx context.Context) bool {
	if s.err != nil || s.ctx.Err() != nil || ctx.Err() != nil {
		return false
	}

	if !s.rows.Next() {
		return false
	}

	fields, err := s.repo.scanRow(s.rows, s.cols)
	if err != nil {
		s.err = err
		return false
	}

	s.item = transactional.DataItem{
		ID:     s.id,
		Fields: fields,
	}
	return true
}

func (s *mssqlStream) Item() transactional.DataItem {
	return s.item
}

func (s *mssqlStream) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.rows.Err()
}

func (s *mssqlStream) Close() error {
	return s.rows.Close()
}
