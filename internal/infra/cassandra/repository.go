package cassandra

import (
	"context"
	"log/slog"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/gocql/gocql"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
	"streaming-golang/internal/domain/timeexpr"
)

type repository struct {
	session *gocql.Session
	logger  *slog.Logger
}

func NewRepository(session *gocql.Session, loggers ...*slog.Logger) transactional.Repository {
	logger := slog.Default()
	if len(loggers) > 0 && loggers[0] != nil {
		logger = loggers[0]
	}
	return &repository{session: session, logger: logger}
}

func (r *repository) Execute(ctx context.Context, query domain.ExecutableQuery) ([]transactional.DataItem, error) {
	start := time.Now()
	r.logger.InfoContext(ctx, "executing cassandra query",
		slog.Int64("identifier", int64(query.ID)),
		slog.String("source", string(query.Source)),
		slog.String("data_category", string(query.DataCategory)),
		slog.String("query", query.Statement),
		slog.Any("arguments", cassandraArguments(query.Arguments)),
	)

	iter := r.session.Query(query.Statement).Bind(query.Arguments...).WithContext(ctx).Iter()

	items := make([]transactional.DataItem, 0)
	for {
		row := make(map[string]any)
		if !iter.MapScan(row) {
			break
		}
		fields := mapCassandraRow(row, query)
		items = append(items, transactional.DataItem{
			ID:     query.ID,
			Fields: fields,
		})
	}

	if err := iter.Close(); err != nil {
		r.logger.ErrorContext(ctx, "execute cassandra query failed",
			slog.Int64("identifier", int64(query.ID)),
			slog.String("source", string(query.Source)),
			slog.String("data_category", string(query.DataCategory)),
			slog.String("query", query.Statement),
			slog.Any("arguments", cassandraArguments(query.Arguments)),
			slog.Int("row_count", len(items)),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.Any("error", err),
		)
		return nil, apperr.Wrap(apperr.Unavailable, "execute cassandra query", err)
	}

	r.logger.InfoContext(ctx, "cassandra query executed",
		slog.Int64("identifier", int64(query.ID)),
		slog.String("source", string(query.Source)),
		slog.String("data_category", string(query.DataCategory)),
		slog.Int("row_count", len(items)),
		slog.Duration("duration", time.Since(start)),
		slog.Int64("duration_ms", time.Since(start).Milliseconds()),
	)

	return items, nil
}

func (r *repository) Stream(ctx context.Context, query domain.ExecutableQuery) (transactional.Stream, error) {
	start := time.Now()
	r.logger.InfoContext(ctx, "streaming cassandra query",
		slog.Int64("identifier", int64(query.ID)),
		slog.String("source", string(query.Source)),
		slog.String("data_category", string(query.DataCategory)),
		slog.String("query", query.Statement),
		slog.Any("arguments", cassandraArguments(query.Arguments)),
	)

	iter := r.session.Query(query.Statement).Bind(query.Arguments...).WithContext(ctx).Iter()

	return &cassandraStream{
		ctx:    ctx,
		iter:   iter,
		id:     query.ID,
		query:  query,
		logger: r.logger,
		start:  start,
	}, nil
}

type cassandraStream struct {
	ctx       context.Context
	iter      *gocql.Iter
	id        domain.Identifier
	query     domain.ExecutableQuery
	logger    *slog.Logger
	start     time.Time
	firstRow  time.Time
	rowCount  int
	item      transactional.DataItem
	err       error
	closeOnce sync.Once
	closeErr  error
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
		Fields: mapCassandraRow(row, s.query),
	}
	s.rowCount++
	if s.rowCount == 1 {
		s.firstRow = time.Now()
	}
	return true
}

func (s *cassandraStream) Item() transactional.DataItem {
	return s.item
}

func (s *cassandraStream) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.close()
}

func (s *cassandraStream) Close() error {
	return s.close()
}

func (s *cassandraStream) close() error {
	s.closeOnce.Do(func() {
		s.closeErr = s.iter.Close()
		if s.closeErr != nil && s.logger != nil {
			s.logger.ErrorContext(s.ctx, "stream cassandra query failed",
				slog.Int64("identifier", int64(s.query.ID)),
				slog.String("source", string(s.query.Source)),
				slog.String("data_category", string(s.query.DataCategory)),
				slog.String("query", s.query.Statement),
				slog.Any("arguments", cassandraArguments(s.query.Arguments)),
				slog.Duration("duration", time.Since(s.start)),
				slog.Int64("duration_ms", time.Since(s.start).Milliseconds()),
				slog.Any("error", s.closeErr),
			)
			return
		}
		if s.logger != nil {
			attrs := []slog.Attr{
				slog.Int64("identifier", int64(s.query.ID)),
				slog.String("source", string(s.query.Source)),
				slog.String("data_category", string(s.query.DataCategory)),
				slog.Int("row_count", s.rowCount),
				slog.Duration("duration", time.Since(s.start)),
				slog.Int64("duration_ms", time.Since(s.start).Milliseconds()),
			}
			if !s.firstRow.IsZero() {
				firstRowDuration := s.firstRow.Sub(s.start)
				attrs = append(attrs,
					slog.Duration("first_row_duration", firstRowDuration),
					slog.Int64("first_row_duration_ms", firstRowDuration.Milliseconds()),
				)
			}
			s.logger.LogAttrs(s.ctx, slog.LevelInfo, "cassandra stream completed", attrs...)
		}
	})
	return s.closeErr
}

func mapCassandraRow(row map[string]any, query domain.ExecutableQuery) map[string]any {
	location := cassandraResponseLocation(query.Parameters)
	deliveryStart := deliveryStartFromCassandra(row)
	fields := map[string]any{
		"Identifier":                 int64(query.ID),
		"ReferenceTime":              referenceTimeFromCassandra(row, location),
		"DeliveryStart":              deliveryStart,
		"DeliveryEnd":                deliveryEndFromCassandra(deliveryStart, location),
		"RelativeDeliveryPeriod":     nil,
		"Value":                      adaptiveRound(asFloat(row["value"])),
		"LegacyDeliveryBucketNumber": nil,
	}
	return projectCassandraFields(fields, query.Parameters)
}

func projectCassandraFields(fields map[string]any, parameters map[string]any) map[string]any {
	columns, ok := parameters["projection_columns"].([]string)
	if !ok || len(columns) == 0 {
		return fields
	}
	columns = ensureCassandraContractColumns(columns)
	projected := make(map[string]any, len(columns))
	for _, column := range columns {
		if value, exists := fields[column]; exists {
			projected[column] = value
		}
	}
	return projected
}

func ensureCassandraContractColumns(columns []string) []string {
	result := append([]string(nil), columns...)
	for _, column := range []string{
		"Identifier",
		"ReferenceTime",
		"DeliveryStart",
		"DeliveryEnd",
		"LegacyDeliveryBucketNumber",
		"RelativeDeliveryPeriod",
		"Value",
	} {
		if !hasColumn(result, column) {
			result = append(result, column)
		}
	}
	return result
}

func hasColumn(columns []string, name string) bool {
	for _, column := range columns {
		if strings.EqualFold(column, name) {
			return true
		}
	}
	return false
}

func referenceTimeFromCassandra(row map[string]any, location *time.Location) time.Time {
	return time.Date(
		asInt(row["qte_y"]),
		time.Month(asInt(row["qte_m"])),
		asInt(row["qte_d"]),
		0,
		0,
		0,
		0,
		location,
	)
}

func cassandraResponseLocation(parameters map[string]any) *time.Location {
	timezone := "Europe/Zurich"
	if value, ok := parameters["cassandra_timezone"].(string); ok && strings.TrimSpace(value) != "" {
		timezone = value
	}
	location, err := timeexpr.LoadLocation(timezone)
	if err != nil {
		return time.UTC
	}
	return location
}

func deliveryStartFromCassandra(row map[string]any) time.Time {
	offset := time.FixedZone("", asInt(row["del_offset"])*int(time.Hour/time.Second))
	return time.Date(
		asInt(row["del_y"]),
		time.Month(asInt(row["del_m"])),
		asInt(row["del_d"]),
		asInt(row["del_h"]),
		asInt(row["del_min"]),
		0,
		0,
		offset,
	)
}

func deliveryEndFromCassandra(deliveryStart time.Time, location *time.Location) time.Time {
	local := deliveryStart.In(location)
	return time.Date(
		local.Year(),
		local.Month(),
		local.Day(),
		local.Hour()+1,
		local.Minute(),
		local.Second(),
		local.Nanosecond(),
		location,
	)
}

func asInt(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	default:
		return 0
	}
}

func asFloat(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

func adaptiveRound(value float64) float64 {
	return math.Round(value*1e10) / 1e10
}

func cassandraArguments(arguments []any) []map[string]any {
	if len(arguments) == 0 {
		return nil
	}
	values := make([]map[string]any, 0, len(arguments))
	for i, value := range arguments {
		values = append(values, map[string]any{
			"index": i,
			"value": value,
		})
	}
	return values
}
