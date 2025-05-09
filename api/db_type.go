package api

import (
	"fmt"
	"strings"

	"github.com/xo/dburl"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GetDBType returns the database type based on the DSN.
func GetDBType(dsn string) (string, error) {
	u, err := dburl.Parse(dsn)
	if err != nil {
		return "", fmt.Errorf("failed to parse DSN: %w", err)
	}

	driverName := u.Driver
	switch strings.ToLower(driverName) {
	case "postgres", "pgx":
		return "PostgreSQL", nil
	case "mysql", "mymysql":
		return "MySQL", nil
	case "sqlite", "sqlite3", "moderncsqlite":
		return "SQLite", nil
	case "sqlserver":
		return "SQL Server", nil
	case "oracle", "godror":
		return "Oracle", nil
	case "clickhouse":
		return "ClickHouse", nil
	case "cassandra":
		return "Cassandra", nil
	case "couchbase":
		return "Couchbase", nil
	case "dynamodb":
		return "DynamoDB", nil
	case "duckdb":
		return "DuckDB", nil
	case "firebird":
		return "Firebird", nil
	case "h2":
		return "H2", nil
	case "hive":
		return "Hive", nil
	case "mssql":
		return "Microsoft SQL Server", nil
	case "odbc":
		return "ODBC", nil
	case "presto":
		return "Presto", nil
	case "snowflake":
		return "Snowflake", nil
	case "trino":
		return "Trino", nil
	case "vertica":
		return "Vertica", nil
	default:
		caser := cases.Title(language.English)
		return caser.String(driverName), nil
	}
}
