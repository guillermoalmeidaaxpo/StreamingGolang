package mssql

import (
	"database/sql"
)

type rowScanner interface {
	Scan(dest ...any) error
}

type mappingRow struct {
	TimeSeriesID        int64
	CMDPViewName        sql.NullString
	MDSDataCategory     string
	MDSDataStructure    sql.NullString
	Resolution          sql.NullString
	CMDPColumnName      string
	MDSColumnName       string
	DataType            string
	IsProjectable       sql.NullBool
	IsKey               sql.NullBool
	OrderPriority       sql.NullInt64
	KeyColumnOrdering   sql.NullInt64
	ValueColumnOrdering sql.NullInt64
	CassandraID         sql.NullString
	HyperscaleID        sql.NullInt64
	SplitQuery          sql.NullBool
	IndexField          sql.NullString
	SwitchOver          sql.NullString
	Timezone            sql.NullString
}

func scanMappingRow(scanner rowScanner) (mappingRow, error) {
	if scanner == nil {
		return mappingRow{}, errNilScanner
	}

	var row mappingRow
	err := scanner.Scan(
		&row.TimeSeriesID,
		&row.CMDPViewName,
		&row.MDSDataCategory,
		&row.MDSDataStructure,
		&row.Resolution,
		&row.CMDPColumnName,
		&row.MDSColumnName,
		&row.DataType,
		&row.IsProjectable,
		&row.IsKey,
		&row.OrderPriority,
		&row.KeyColumnOrdering,
		&row.ValueColumnOrdering,
		&row.CassandraID,
		&row.HyperscaleID,
		&row.SplitQuery,
		&row.IndexField,
		&row.SwitchOver,
	)
	return row, err
}
