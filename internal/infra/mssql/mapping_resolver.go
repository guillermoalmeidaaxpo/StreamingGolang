package mssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/domain"
)

type MappingResolver struct {
	cmdpMappingDB *sql.DB
	mdsDB         *sql.DB
	cmdpSQLDB     *sql.DB
	logger        *slog.Logger
	filterLimits  *filterLimitsMemoryCache
}

func NewMappingResolver(cmdpMappingDB, mdsDB, cmdpSQLDB *sql.DB, logger *slog.Logger) *MappingResolver {
	if logger == nil {
		logger = slog.Default()
	}
	return &MappingResolver{
		cmdpMappingDB: cmdpMappingDB,
		mdsDB:         mdsDB,
		cmdpSQLDB:     cmdpSQLDB,
		logger:        logger,
		filterLimits:  newFilterLimitsMemoryCache(time.Hour, 10*time.Minute),
	}
}

func (r *MappingResolver) GetWatermark(ctx context.Context, mappings []domain.Mapping) (time.Time, error) {
	if r.cmdpSQLDB == nil {
		return time.Now().UTC(), nil
	}

	if len(mappings) == 0 || strings.TrimSpace(mappings[0].ViewName) == "" {
		return time.Now().UTC(), nil
	}

	mapping := mappings[0]
	referenceTimeColumn := referenceTimeSourceColumn(mapping)
	if strings.TrimSpace(referenceTimeColumn) == "" || strings.TrimSpace(mapping.IndexField) == "" {
		return time.Now().UTC(), nil
	}
	cache := r.getFilterLimitsCache()
	cacheKey := filterLimitsCacheKey(mapping.ID, true, false, false)
	if limits, ok := cache.Get(cacheKey); ok {
		watermark := time.Now().UTC()
		if limits.MaxReferenceTime.Valid {
			watermark = limits.MaxReferenceTime.Time.UTC()
		}
		r.logger.InfoContext(ctx, "CMDP filter limits cache hit",
			slog.Int64("identifier", int64(mapping.ID)),
			slog.String("cache_key", cacheKey),
			slog.String("view", mapping.ViewName),
			slog.String("reference_time_column", referenceTimeColumn),
			slog.Time("watermark", watermark),
			slog.Bool("watermark_found", limits.MaxReferenceTime.Valid),
		)
		return watermark, nil
	}

	start := time.Now()
	var minReferenceTime sql.NullTime
	var maxReferenceTime sql.NullTime
	var minDeliveryStart sql.NullTime
	var maxDeliveryStart sql.NullTime

	commandText := calculateMinMaxReferenceTimeDeliveryStartCommandText()
	deliveryStartColumn := deliveryStartSourceColumn(mapping)

	r.logger.InfoContext(ctx, "reading CMDP filter limits",
		slog.Int64("identifier", int64(mapping.ID)),
		slog.String("view", mapping.ViewName),
		slog.String("reference_time_index_column", mapping.IndexField),
		slog.String("reference_time_column", referenceTimeColumn),
		slog.String("delivery_start_column", deliveryStartColumn),
		slog.Bool("get_min_reference_time", false),
		slog.Bool("get_max_reference_time", true),
		slog.Bool("get_min_max_delivery_start", false),
		slog.String("command", compactSQL(commandText)),
	)

	_, err := r.cmdpSQLDB.ExecContext(ctx, commandText,
		sql.Named("Id", int64(mapping.ID)),
		sql.Named("referenceTimeIndexedFieldName", mapping.IndexField),
		sql.Named("referenceTimeFieldName", referenceTimeColumn),
		sql.Named("deliveryStartFieldName", deliveryStartColumn),
		sql.Named("getMinReferenceTime", false),
		sql.Named("getMaxReferenceTime", true),
		sql.Named("getMinMaxDeliveryStart", false),
		sql.Named("schemaQualifiedViewName", mapping.ViewName),
		sql.Named("minReferenceTime", sql.Out{Dest: &minReferenceTime}),
		sql.Named("maxReferenceTime", sql.Out{Dest: &maxReferenceTime}),
		sql.Named("minDeliveryStart", sql.Out{Dest: &minDeliveryStart}),
		sql.Named("maxDeliveryStart", sql.Out{Dest: &maxDeliveryStart}),
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "read CMDP filter limits failed",
			slog.Int64("identifier", int64(mapping.ID)),
			slog.String("view", mapping.ViewName),
			slog.String("reference_time_index_column", mapping.IndexField),
			slog.String("reference_time_column", referenceTimeColumn),
			slog.String("delivery_start_column", deliveryStartColumn),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.Any("error", err),
		)
		return time.Time{}, apperr.Wrap(apperr.Unavailable, "read CMDP filter limits", err)
	}

	res := maxReferenceTime.Time
	if !maxReferenceTime.Valid {
		res = time.Now().UTC()
	}
	cache.Set(cacheKey, filterLimits{
		MinReferenceTime: minReferenceTime,
		MaxReferenceTime: maxReferenceTime,
		MinDeliveryStart: minDeliveryStart,
		MaxDeliveryStart: maxDeliveryStart,
	})

	r.logger.InfoContext(ctx, "CMDP filter limits read",
		slog.Int64("identifier", int64(mapping.ID)),
		slog.String("view", mapping.ViewName),
		slog.String("reference_time_column", referenceTimeColumn),
		slog.Time("watermark", res),
		slog.Bool("watermark_found", maxReferenceTime.Valid),
		slog.Duration("duration", time.Since(start)),
		slog.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return res.UTC(), nil
}

func (r *MappingResolver) getFilterLimitsCache() *filterLimitsMemoryCache {
	if r.filterLimits == nil {
		r.filterLimits = newFilterLimitsMemoryCache(time.Hour, 10*time.Minute)
	}
	return r.filterLimits
}

type filterLimits struct {
	MinReferenceTime sql.NullTime
	MaxReferenceTime sql.NullTime
	MinDeliveryStart sql.NullTime
	MaxDeliveryStart sql.NullTime
}

type filterLimitsCacheEntry struct {
	value    filterLimits
	created  time.Time
	accessed time.Time
}

type filterLimitsMemoryCache struct {
	mu       sync.Mutex
	entries  map[string]filterLimitsCacheEntry
	absolute time.Duration
	sliding  time.Duration
	now      func() time.Time
}

func newFilterLimitsMemoryCache(absolute, sliding time.Duration) *filterLimitsMemoryCache {
	return &filterLimitsMemoryCache{
		entries:  make(map[string]filterLimitsCacheEntry),
		absolute: absolute,
		sliding:  sliding,
		now:      time.Now,
	}
}

func (c *filterLimitsMemoryCache) Get(key string) (filterLimits, bool) {
	if c == nil {
		return filterLimits{}, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return filterLimits{}, false
	}
	now := c.now()
	if c.expired(entry, now) {
		delete(c.entries, key)
		return filterLimits{}, false
	}
	entry.accessed = now
	c.entries[key] = entry
	return entry.value, true
}

func (c *filterLimitsMemoryCache) Set(key string, value filterLimits) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	c.entries[key] = filterLimitsCacheEntry{
		value:    value,
		created:  now,
		accessed: now,
	}
}

func (c *filterLimitsMemoryCache) expired(entry filterLimitsCacheEntry, now time.Time) bool {
	if c.absolute > 0 && now.Sub(entry.created) >= c.absolute {
		return true
	}
	return c.sliding > 0 && now.Sub(entry.accessed) >= c.sliding
}

func filterLimitsCacheKey(id domain.Identifier, getMaxReferenceTime, getMinReferenceTime, getMinMaxDeliveryStart bool) string {
	return fmt.Sprintf("FilterLimits_%d_%s_%s_%s",
		id,
		csharpBool(getMaxReferenceTime),
		csharpBool(getMinReferenceTime),
		csharpBool(getMinMaxDeliveryStart),
	)
}

func csharpBool(value bool) string {
	if value {
		return "True"
	}
	return "False"
}

func calculateMinMaxReferenceTimeDeliveryStartCommandText() string {
	return `
EXEC [MDS].[CalculateMinMaxReferenceTimeDeliveryStart]
	@Id = @Id,
	@referenceTimeIndexedFieldName = @referenceTimeIndexedFieldName,
	@referenceTimeFieldName = @referenceTimeFieldName,
	@deliveryStartFieldName = @deliveryStartFieldName,
	@getMinReferenceTime = @getMinReferenceTime,
	@getMaxReferenceTime = @getMaxReferenceTime,
	@getMinMaxDeliveryStart = @getMinMaxDeliveryStart,
	@schemaQualifiedViewName = @schemaQualifiedViewName,
	@minReferenceTime = @minReferenceTime OUTPUT,
	@maxReferenceTime = @maxReferenceTime OUTPUT,
	@minDeliveryStart = @minDeliveryStart OUTPUT,
	@maxDeliveryStart = @maxDeliveryStart OUTPUT`
}

func referenceTimeSourceColumn(mapping domain.Mapping) string {
	for _, column := range mapping.Columns {
		if strings.EqualFold(column.MDSName, "ReferenceTime") && strings.TrimSpace(column.SourceName) != "" {
			return column.SourceName
		}
	}
	return "ReferenceTime"
}

func deliveryStartSourceColumn(mapping domain.Mapping) string {
	for _, column := range mapping.Columns {
		if isDeliveryStartColumn(column.MDSName) && strings.TrimSpace(column.SourceName) != "" {
			return column.SourceName
		}
	}
	return ""
}

func isDeliveryStartColumn(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "deliverystart", "underlyingstart", "optionstart":
		return true
	default:
		return false
	}
}

func (r *MappingResolver) ResolveMappings(ctx context.Context, ids []domain.Identifier, category domain.DataCategory, stage string) ([]domain.Mapping, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	if r == nil || r.cmdpMappingDB == nil {
		return nil, apperr.New(apperr.Unavailable, "mapping database is not configured")
	}

	logger := r.resolverLogger()
	logger.InfoContext(ctx, "resolving data mappings",
		slog.Any("identifiers", ids),
		slog.String("data_category", string(category)),
		slog.String("stage", stage),
	)

	if usesMDSMappings(stage) {
		return r.readMDSDomainMappings(ctx, ids, category)
	}

	rows, err := r.readCMDPMappings(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		logger.WarnContext(ctx, "no CMDP mappings found",
			slog.Any("identifiers", ids),
			slog.String("data_category", string(category)),
			slog.String("stage", stage),
		)
		return nil, apperr.New(apperr.NotFound, "requested identifiers do not have mappings")
	}

	cmdpMappings := buildDomainMappings(rows, category)
	if usesMigrationMappings(stage) {
		return r.processMigrationMappings(ctx, cmdpMappings, category)
	}
	return r.processRegularMappings(ctx, cmdpMappings, category)
}

func (r *MappingResolver) readCMDPMappings(ctx context.Context, ids []domain.Identifier) ([]mappingRow, error) {
	query, args := mappingQuery(ids)
	logger := r.resolverLogger()
	start := time.Now()
	logger.InfoContext(ctx, "reading CMDP mappings",
		slog.Any("identifiers", ids),
		slog.Int("parameter_count", len(args)),
		slog.String("query", compactSQL(query)),
	)

	rows, err := r.cmdpMappingDB.QueryContext(ctx, query, args...)
	if err != nil {
		logger.ErrorContext(ctx, "read CMDP mappings failed",
			slog.Any("identifiers", ids),
			slog.Int("parameter_count", len(args)),
			slog.String("query", compactSQL(query)),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.Any("error", err),
		)
		return nil, apperr.Wrap(apperr.Unavailable, "read CMDP mappings", err)
	}
	defer rows.Close()

	result := make([]mappingRow, 0)
	for rows.Next() {
		row, err := scanMappingRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		logger.ErrorContext(ctx, "iterate CMDP mappings failed",
			slog.Any("identifiers", ids),
			slog.Int("row_count", len(result)),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.Any("error", err),
		)
		return nil, apperr.Wrap(apperr.Unavailable, "iterate CMDP mappings", err)
	}

	logger.InfoContext(ctx, "CMDP mappings read",
		slog.Any("identifiers", ids),
		slog.Int("row_count", len(result)),
		slog.Duration("duration", time.Since(start)),
		slog.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return result, nil
}

func (r *MappingResolver) readMDSDomainMappings(ctx context.Context, ids []domain.Identifier, category domain.DataCategory) ([]domain.Mapping, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	if r.mdsDB == nil {
		return nil, apperr.New(apperr.Unavailable, "MDS mapping database is not configured")
	}

	query, args := mdsMappingQuery(ids)
	logger := r.resolverLogger()
	start := time.Now()
	logger.InfoContext(ctx, "reading MDS mappings",
		slog.Any("identifiers", ids),
		slog.String("data_category", string(category)),
		slog.Int("parameter_count", len(args)),
		slog.String("query", compactSQL(query)),
	)

	rows, err := r.mdsDB.QueryContext(ctx, query, args...)
	if err != nil {
		logger.ErrorContext(ctx, "read MDS mappings failed",
			slog.Any("identifiers", ids),
			slog.String("data_category", string(category)),
			slog.Int("parameter_count", len(args)),
			slog.String("query", compactSQL(query)),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.Any("error", err),
		)
		return nil, apperr.Wrap(apperr.Unavailable, "read MDS mappings", err)
	}
	defer rows.Close()

	result := make([]mappingRow, 0)
	for rows.Next() {
		row, err := scanMDSMappingRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		logger.ErrorContext(ctx, "iterate MDS mappings failed",
			slog.Any("identifiers", ids),
			slog.String("data_category", string(category)),
			slog.Int("row_count", len(result)),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.Any("error", err),
		)
		return nil, apperr.Wrap(apperr.Unavailable, "iterate MDS mappings", err)
	}
	if len(result) == 0 {
		logger.WarnContext(ctx, "no MDS mappings found",
			slog.Any("identifiers", ids),
			slog.String("data_category", string(category)),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil, apperr.New(apperr.NotFound, "requested identifiers do not have MDS mappings")
	}

	logger.InfoContext(ctx, "MDS mappings read",
		slog.Any("identifiers", ids),
		slog.String("data_category", string(category)),
		slog.Int("row_count", len(result)),
		slog.Duration("duration", time.Since(start)),
		slog.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return buildDomainMappings(result, category), nil
}

func (r *MappingResolver) processRegularMappings(ctx context.Context, cmdpMappings []domain.Mapping, category domain.DataCategory) ([]domain.Mapping, error) {
	groups := groupBySwitchover(cmdpMappings)
	mdsIDs := make([]domain.Identifier, 0)

	for _, mapping := range groups.MDSSwitchover {
		mdsIDs = append(mdsIDs, hyperscaleOrOwnID(mapping))
	}

	hyperscaleMDOIDs := make(map[domain.Identifier]struct{})
	for _, mapping := range groups.NoSwitchover {
		if mapping.HyperscaleID != nil {
			hyperscaleMDOIDs[mapping.ID] = struct{}{}
			mdsIDs = append(mdsIDs, mapping.ID)
		}
	}

	result := make([]domain.Mapping, 0, len(cmdpMappings))
	if len(mdsIDs) > 0 {
		mdsMappings, err := r.readMDSDomainMappings(ctx, distinctIdentifiers(mdsIDs), category)
		if err != nil {
			return nil, err
		}
		result = append(result, enrichMDSMappings(mdsMappings, groups.MDSSwitchover)...)
	}

	for _, mapping := range groups.CMDPSwitchover {
		result = append(result, forceCMDP(mapping))
	}
	for _, mapping := range groups.NoSwitchover {
		if _, isHyperscale := hyperscaleMDOIDs[mapping.ID]; !isHyperscale {
			result = append(result, mapping)
		}
	}

	return result, nil
}

func (r *MappingResolver) processMigrationMappings(ctx context.Context, cmdpMappings []domain.Mapping, category domain.DataCategory) ([]domain.Mapping, error) {
	groups := groupBySwitchover(cmdpMappings)
	mdsIDs := make([]domain.Identifier, 0)

	for _, mapping := range groups.CMDPSwitchover {
		mdsIDs = append(mdsIDs, hyperscaleOrOwnID(mapping))
	}
	for _, mapping := range groups.NoSwitchover {
		mdsIDs = append(mdsIDs, mapping.ID)
	}

	result := make([]domain.Mapping, 0, len(cmdpMappings))
	if len(mdsIDs) > 0 {
		mdsMappings, err := r.readMDSDomainMappings(ctx, distinctIdentifiers(mdsIDs), category)
		if err != nil {
			return nil, err
		}
		result = append(result, enrichMDSMappings(mdsMappings, groups.CMDPSwitchover)...)
	}

	for _, mapping := range groups.MDSSwitchover {
		result = append(result, forceCMDP(mapping))
	}

	return result, nil
}

func mappingQuery(ids []domain.Identifier) (string, []any) {
	placeholders := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for index, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("@p%d", index+1))
		args = append(args, int64(id))
	}

	return fmt.Sprintf(`
SELECT
	TIMESERIES_ID,
	CMDP_VIEW_NAME,
	MDS_DATA_CATEGORY,
	MDS_DATA_STRUCTURE,
	RESOLUTION,
	CMDP_COLUMN_NAME,
	MDS_COLUMN_NAME,
	DATA_TYPE,
	IS_PROJECTABLE,
	IS_KEY,
	ORDER_PRIORITY,
	KEY_COLUMN_ORDERING,
	VALUE_COLUMN_ORDERING,
	CASSANDRA_ID,
	HYPERSCALE_ID,
	SPLIT_QUERY,
	CMDP_COLUMN_INDEXED,
	SWITCHOVER
FROM CMDP_TO_MDS_MAPPING
WHERE TIMESERIES_ID IN (%s)`, strings.Join(placeholders, ",")), args
}

func buildDomainMappings(rows []mappingRow, fallbackCategory domain.DataCategory) []domain.Mapping {
	byID := make(map[domain.Identifier][]mappingRow)
	order := make([]domain.Identifier, 0)
	for _, row := range rows {
		id := domain.Identifier(row.TimeSeriesID)
		if _, exists := byID[id]; !exists {
			order = append(order, id)
		}
		byID[id] = append(byID[id], row)
	}

	mappings := make([]domain.Mapping, 0, len(byID))
	for _, id := range order {
		group := byID[id]
		first := group[0]
		category := parseDataCategory(first.MDSDataCategory, fallbackCategory)
		mapping := domain.Mapping{
			ID:           id,
			DataCategory: category,
			Source:       sourceKind(first),
			ViewName:     first.CMDPViewName.String,
			Views: domain.MappingViews{
				LatestVersion:                     first.LatestVersionView.String,
				LatestReferenceTime:               first.LatestReferenceTimeView.String,
				LatestVersionWithCreatedOn:        first.LatestVersionWithCreatedOnView.String,
				LatestReferenceTimeWithCreatedOn:  first.LatestReferenceTimeWithCreatedOnView.String,
				GetByCreatedOn:                    first.GetByCreatedOnView.String,
				GetByCreatedOnLatestReferenceTime: first.GetByCreatedOnLatestReferenceTimeView.String,
			},
			IndexField:  first.IndexField.String,
			Resolution:  first.Resolution.String,
			CassandraID: first.CassandraID.String,
			SwitchOver:  first.SwitchOver.String,
			SplitQuery:  boolValue(first.SplitQuery, true),
			Timezone:    first.Timezone.String,
			Columns:     make([]domain.ColumnMapping, 0, len(group)),
		}
		if first.HyperscaleID.Valid {
			hyperscaleID := domain.Identifier(first.HyperscaleID.Int64)
			mapping.HyperscaleID = &hyperscaleID
		}

		for _, row := range group {
			mapping.Columns = append(mapping.Columns, domain.ColumnMapping{
				MDSName:             row.MDSColumnName,
				SourceName:          row.CMDPColumnName,
				DataType:            row.DataType,
				IsKey:               row.IsKey.Bool,
				IsProjectable:       boolValue(row.IsProjectable, false),
				OrderPriority:       nullableInt(row.OrderPriority),
				KeyColumnOrdering:   nullableInt(row.KeyColumnOrdering),
				ValueColumnOrdering: nullableInt(row.ValueColumnOrdering),
			})
		}
		mappings = append(mappings, mapping)
	}

	return mappings
}

func parseDataCategory(value string, fallback domain.DataCategory) domain.DataCategory {
	switch strings.ToLower(value) {
	case "curve", "curves":
		return domain.Curves
	case "surface", "surfaces":
		return domain.Surfaces
	case "timeseries", "time_series", "time series":
		return domain.TimeSeries
	default:
		return fallback
	}
}

func sourceKind(row mappingRow) domain.SourceKind {
	if row.HyperscaleID.Valid {
		return domain.SourceHyperscale
	}
	if strings.HasPrefix(strings.ToLower(row.SwitchOver.String), "mds") {
		return domain.SourceHyperscale
	}
	if row.CassandraID.Valid && row.CassandraID.String != "" {
		return domain.SourceCassandra
	}
	return domain.SourceCMDP
}

func usesMDSMappings(stage string) bool {
	stage = strings.ToLower(stage)
	return strings.Contains(stage, "design") || strings.Contains(stage, "validation")
}

func usesMigrationMappings(stage string) bool {
	return strings.Contains(strings.ToLower(stage), "migration")
}

func boolValue(value sql.NullBool, fallback bool) bool {
	if !value.Valid {
		return fallback
	}
	return value.Bool
}

func nullableInt(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}
	result := int(value.Int64)
	return &result
}

var errNilScanner = errors.New("nil scanner")

func (r *MappingResolver) resolverLogger() *slog.Logger {
	if r == nil || r.logger == nil {
		return slog.Default()
	}
	return r.logger
}

func compactSQL(query string) string {
	return strings.Join(strings.Fields(query), " ")
}
