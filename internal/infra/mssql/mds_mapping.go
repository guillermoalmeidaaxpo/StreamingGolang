package mssql

import (
	"database/sql"
	"fmt"
	"strings"

	"streaming-golang/internal/domain"
)

func mdsMappingQuery(ids []domain.Identifier) (string, []any) {
	placeholders := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for index, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("@p%d", index+1))
		args = append(args, int64(id))
	}

	return fmt.Sprintf(`
SELECT
	MdoId,
	CategoryName,
	ResolutionISO,
	ColumnName,
	DataType,
	OrderPriority,
	KeyColumnOrdering,
	ValueColumnOrdering,
	LatestVersionView,
	LatestReferenceTimeView,
	LatestVersionView + 'WithCreatedOn' AS LatestVersionWithCreatedOnView,
	LatestReferenceTimeView + 'WithCreatedOn' AS LatestReferenceTimeWithCreatedOnView,
	GetByCreatedOnView,
	GetByCreatedOnLatestReferenceTimeView,
	TimeZone
FROM [Api].[VI_MdsMappingDetails]
WHERE MdoId IN (%s)`, strings.Join(placeholders, ",")), args
}

func scanMDSMappingRow(scanner rowScanner) (mappingRow, error) {
	if scanner == nil {
		return mappingRow{}, errNilScanner
	}

	var (
		row                 mappingRow
		orderPriority       sql.NullInt64
		keyColumnOrdering   sql.NullInt64
		valueColumnOrdering sql.NullInt64
		timezone            sql.NullString
	)

	err := scanner.Scan(
		&row.TimeSeriesID,
		&row.MDSDataCategory,
		&row.Resolution,
		&row.CMDPColumnName,
		&row.DataType,
		&orderPriority,
		&keyColumnOrdering,
		&valueColumnOrdering,
		&row.LatestVersionView,
		&row.LatestReferenceTimeView,
		&row.LatestVersionWithCreatedOnView,
		&row.LatestReferenceTimeWithCreatedOnView,
		&row.GetByCreatedOnView,
		&row.GetByCreatedOnLatestReferenceTimeView,
		&timezone,
	)
	if err != nil {
		return mappingRow{}, err
	}

	isKey := keyColumnOrdering.Valid
	row.MDSColumnName = row.CMDPColumnName
	row.IsKey = sql.NullBool{Bool: isKey, Valid: true}
	row.IsProjectable = sql.NullBool{Bool: !isKey, Valid: true}
	row.OrderPriority = orderPriority
	row.KeyColumnOrdering = keyColumnOrdering
	row.ValueColumnOrdering = valueColumnOrdering
	row.HyperscaleID = sql.NullInt64{Int64: row.TimeSeriesID, Valid: true}
	row.SplitQuery = sql.NullBool{Bool: true, Valid: true}
	row.Timezone = timezone

	return row, nil
}
