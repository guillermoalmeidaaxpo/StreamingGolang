package mssql

import (
	"database/sql"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/microsoft/go-mssqldb/azuread"
)

func OpenSQLServer(dsn string) (*sql.DB, error) {
	if !IsConfiguredDSN(dsn) {
		return nil, nil
	}
	driverName := DriverNameForDSN(dsn)
	dsn = NormalizeDSNForDriver(dsn)
	return sql.Open(driverName, dsn)
}

func IsConfiguredDSN(dsn string) bool {
	dsn = strings.TrimSpace(dsn)
	return dsn != "" && !strings.EqualFold(dsn, "NOT SET")
}

func DriverNameForDSN(dsn string) string {
	if usesAzureADAuth(dsn) {
		return azuread.DriverName
	}
	return "sqlserver"
}

func NormalizeDSNForDriver(dsn string) string {
	if usesAzureADAuth(dsn) {
		return normalizeAzureADDSN(dsn)
	}
	return dsn
}

func AuthModeForDSN(dsn string) string {
	if !IsConfiguredDSN(dsn) {
		return "not-configured"
	}

	for _, part := range strings.Split(dsn, ";") {
		key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}

		key = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(key), " ", ""))
		value = strings.TrimSpace(value)
		switch key {
		case "fedauth":
			return "fedauth:" + value
		case "authentication":
			return "authentication:" + value
		case "integratedsecurity", "trusted_connection":
			return "integrated-security:" + value
		}
	}

	return "sql-auth-or-default"
}

func usesAzureADAuth(dsn string) bool {
	dsn = strings.ToLower(dsn)
	return strings.Contains(dsn, "fedauth=activedirectory") ||
		strings.Contains(dsn, "authentication=active directory")
}

func normalizeAzureADDSN(dsn string) string {
	parts := strings.Split(dsn, ";")
	normalized := make([]string, 0, len(parts)+1)
	hasFedAuth := false
	mappedFedAuth := ""

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		key, value, ok := strings.Cut(part, "=")
		if !ok {
			normalized = append(normalized, part)
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		normalizedKey := strings.ToLower(strings.ReplaceAll(key, " ", ""))

		switch normalizedKey {
		case "fedauth":
			hasFedAuth = true
			normalized = append(normalized, key+"="+value)
		case "authentication":
			if mapped := mapAzureADAuthentication(value); mapped != "" {
				mappedFedAuth = mapped
			} else {
				normalized = append(normalized, part)
			}
		default:
			normalized = append(normalized, key+"="+value)
		}
	}

	if !hasFedAuth && mappedFedAuth != "" {
		normalized = append(normalized, "fedauth="+mappedFedAuth)
	}

	return strings.Join(normalized, ";")
}

func mapAzureADAuthentication(value string) string {
	switch strings.ToLower(strings.Join(strings.Fields(value), " ")) {
	case "active directory default":
		return "ActiveDirectoryDefault"
	case "active directory interactive":
		return "ActiveDirectoryDefault"
	case "active directory managed identity", "active directory msi":
		return "ActiveDirectoryManagedIdentity"
	case "active directory service principal", "active directory application":
		return "ActiveDirectoryServicePrincipal"
	case "active directory password":
		return "ActiveDirectoryPassword"
	case "active directory integrated":
		return "ActiveDirectoryIntegrated"
	default:
		return ""
	}
}
