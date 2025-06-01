package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/xo/dburl"
	_ "github.com/xo/usql/drivers"
)

// ListTables returns a list of all tables in the database
func ListTables(db *sql.DB, dsn string) ([]string, error) {
	u, err := dburl.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	driverName := strings.ToLower(u.Driver)

	switch driverName {
	case "sqlite", "sqlite3", "moderncsqlite":
		return listSQLiteTables(db)
	case "postgres", "pgx":
		return listPostgresTables(db)
	case "mysql", "mymysql":
		return listMySQLTables(db)
	case "sqlserver":
		return listSQLServerTables(db)
	case "oracle", "godror":
		return listOracleTables(db)
	case "clickhouse":
		return listClickHouseTables(db)
	case "duckdb":
		return listDuckDBTables(db)
	case "snowflake":
		return listSnowflakeTables(db)
	default:
		return nil, fmt.Errorf("unsupported database driver for table listing: %s", driverName)
	}
}

func listSQLiteTables(db *sql.DB) ([]string, error) {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list SQLite tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func listPostgresTables(db *sql.DB) ([]string, error) {
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE';`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list PostgreSQL tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func listMySQLTables(db *sql.DB) ([]string, error) {
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE';`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list MySQL tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func listSQLServerTables(db *sql.DB) ([]string, error) {
	query := `SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE';`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list SQL Server tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func listOracleTables(db *sql.DB) ([]string, error) {
	query := `SELECT table_name FROM user_tables;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list Oracle tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func listClickHouseTables(db *sql.DB) ([]string, error) {
	query := `SELECT name FROM system.tables WHERE database = currentDatabase() AND engine NOT LIKE '%View';`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list ClickHouse tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func listDuckDBTables(db *sql.DB) ([]string, error) {
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'main' AND table_type = 'BASE TABLE';`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list DuckDB tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func listSnowflakeTables(db *sql.DB) ([]string, error) {
	query := `SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE';`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list Snowflake tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}
