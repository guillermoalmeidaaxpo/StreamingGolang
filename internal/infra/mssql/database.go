package mssql

import (
	"database/sql"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
)

func OpenSQLServer(dsn string) (*sql.DB, error) {
	if !IsConfiguredDSN(dsn) {
		return nil, nil
	}
	return sql.Open("sqlserver", dsn)
}

func IsConfiguredDSN(dsn string) bool {
	dsn = strings.TrimSpace(dsn)
	return dsn != "" && !strings.EqualFold(dsn, "NOT SET")
}
