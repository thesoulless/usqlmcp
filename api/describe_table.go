package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/xo/dburl"
	_ "github.com/xo/usql/drivers"
)

// DescribeTable retrieves schema information for a specific table.
func DescribeTable(db *sql.DB, tableName string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf(`PRAGMA table_info(%s);`, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to describe table: %w", err)
	}
	defer rows.Close()

	schema := []map[string]interface{}{}
	for rows.Next() {
		var (
			cid        int
			name       string
			typeInfo   string
			notnull    int
			defaultVal sql.NullString
			pk         int
		)
		if err := rows.Scan(&cid, &name, &typeInfo, &notnull, &defaultVal, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan table schema: %w", err)
		}

		column := map[string]interface{}{
			"cid":         cid,
			"name":        name,
			"type":        typeInfo,
			"notnull":     notnull,
			"default":     defaultVal.String,
			"primary_key": pk,
		}
		schema = append(schema, column)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return schema, nil
}

// TableColumn represents a standardized table column schema
type TableColumn struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Nullable     bool        `json:"nullable"`
	Default      interface{} `json:"default"`
	IsPrimaryKey bool        `json:"is_primary_key"`
}

// DescribeTableUniversal retrieves schema information for a specific table across different database types
func DescribeTableUniversal(db *sql.DB, tableName string, dsn string) ([]TableColumn, error) {
	u, err := dburl.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	driverName := strings.ToLower(u.Driver)

	switch driverName {
	case "sqlite", "sqlite3", "moderncsqlite":
		return describeSQLiteTable(db, tableName)
	case "postgres", "pgx":
		return describePostgresTable(db, tableName)
	case "mysql", "mymysql":
		return describeMySQLTable(db, tableName)
	case "sqlserver":
		return describeSQLServerTable(db, tableName)
	case "oracle", "godror":
		return describeOracleTable(db, tableName)
	case "clickhouse":
		return describeClickHouseTable(db, tableName)
	case "duckdb":
		return describeDuckDBTable(db, tableName)
	case "snowflake":
		return describeSnowflakeTable(db, tableName)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driverName)
	}
}

// describeSQLiteTable handles SQLite table schema
func describeSQLiteTable(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := fmt.Sprintf(`PRAGMA table_info(%s);`, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to describe SQLite table: %w", err)
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var (
			cid        int
			name       string
			typeInfo   string
			notnull    int
			defaultVal sql.NullString
			pk         int
		)
		if err := rows.Scan(&cid, &name, &typeInfo, &notnull, &defaultVal, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan SQLite table schema: %w", err)
		}

		var defaultValue interface{}
		if defaultVal.Valid {
			defaultValue = defaultVal.String
		}

		column := TableColumn{
			Name:         name,
			Type:         typeInfo,
			Nullable:     notnull == 0,
			Default:      defaultValue,
			IsPrimaryKey: pk == 1,
		}
		columns = append(columns, column)
	}

	return columns, rows.Err()
}

// describePostgresTable handles PostgreSQL table schema
func describePostgresTable(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable = 'YES' AS nullable,
			column_default,
			CASE WHEN column_name IN (
				SELECT a.attname
				FROM pg_index i
				JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
				WHERE i.indrelid = $2::regclass AND i.indisprimary
			) THEN true ELSE false END AS is_primary_key
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position;`

	rows, err := db.Query(query, tableName, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to describe PostgreSQL table: %w", err)
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var (
			name       string
			dataType   string
			nullable   bool
			defaultVal sql.NullString
			isPK       bool
		)
		if err := rows.Scan(&name, &dataType, &nullable, &defaultVal, &isPK); err != nil {
			return nil, fmt.Errorf("failed to scan PostgreSQL table schema: %w", err)
		}

		var defaultValue interface{}
		if defaultVal.Valid {
			defaultValue = defaultVal.String
		}

		column := TableColumn{
			Name:         name,
			Type:         dataType,
			Nullable:     nullable,
			Default:      defaultValue,
			IsPrimaryKey: isPK,
		}
		columns = append(columns, column)
	}

	return columns, rows.Err()
}

// describeMySQLTable handles MySQL table schema
func describeMySQLTable(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable = 'YES' AS nullable,
			column_default,
			column_key = 'PRI' AS is_primary_key
		FROM information_schema.columns
		WHERE table_name = ?
		ORDER BY ordinal_position;`

	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to describe MySQL table: %w", err)
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var (
			name       string
			dataType   string
			nullable   bool
			defaultVal sql.NullString
			isPK       bool
		)
		if err := rows.Scan(&name, &dataType, &nullable, &defaultVal, &isPK); err != nil {
			return nil, fmt.Errorf("failed to scan MySQL table schema: %w", err)
		}

		var defaultValue interface{}
		if defaultVal.Valid {
			defaultValue = defaultVal.String
		}

		column := TableColumn{
			Name:         name,
			Type:         dataType,
			Nullable:     nullable,
			Default:      defaultValue,
			IsPrimaryKey: isPK,
		}
		columns = append(columns, column)
	}

	return columns, rows.Err()
}

// describeSQLServerTable handles SQL Server table schema
func describeSQLServerTable(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := `
		SELECT
			c.COLUMN_NAME,
			c.DATA_TYPE,
			CASE WHEN c.IS_NULLABLE = 'YES' THEN 1 ELSE 0 END AS nullable,
			c.COLUMN_DEFAULT,
			CASE WHEN pk.COLUMN_NAME IS NOT NULL THEN 1 ELSE 0 END AS is_primary_key
		FROM INFORMATION_SCHEMA.COLUMNS c
		LEFT JOIN (
			SELECT kcu.COLUMN_NAME
			FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
			WHERE tc.TABLE_NAME = ? AND tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
		) pk ON c.COLUMN_NAME = pk.COLUMN_NAME
		WHERE c.TABLE_NAME = ?
		ORDER BY c.ORDINAL_POSITION;`

	rows, err := db.Query(query, tableName, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to describe SQL Server table: %w", err)
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var (
			name       string
			dataType   string
			nullable   int
			defaultVal sql.NullString
			isPK       int
		)
		if err := rows.Scan(&name, &dataType, &nullable, &defaultVal, &isPK); err != nil {
			return nil, fmt.Errorf("failed to scan SQL Server table schema: %w", err)
		}

		var defaultValue interface{}
		if defaultVal.Valid {
			defaultValue = defaultVal.String
		}

		column := TableColumn{
			Name:         name,
			Type:         dataType,
			Nullable:     nullable == 1,
			Default:      defaultValue,
			IsPrimaryKey: isPK == 1,
		}
		columns = append(columns, column)
	}

	return columns, rows.Err()
}

// describeOracleTable handles Oracle table schema
func describeOracleTable(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := `
		SELECT
			c.COLUMN_NAME,
			c.DATA_TYPE,
			CASE WHEN c.NULLABLE = 'Y' THEN 1 ELSE 0 END AS nullable,
			c.DATA_DEFAULT,
			CASE WHEN pk.COLUMN_NAME IS NOT NULL THEN 1 ELSE 0 END AS is_primary_key
		FROM ALL_TAB_COLUMNS c
		LEFT JOIN (
			SELECT acc.COLUMN_NAME
			FROM ALL_CONSTRAINTS ac
			JOIN ALL_CONS_COLUMNS acc ON ac.CONSTRAINT_NAME = acc.CONSTRAINT_NAME
			WHERE ac.TABLE_NAME = UPPER(?) AND ac.CONSTRAINT_TYPE = 'P'
		) pk ON c.COLUMN_NAME = pk.COLUMN_NAME
		WHERE c.TABLE_NAME = UPPER(?)
		ORDER BY c.COLUMN_ID;`

	rows, err := db.Query(query, tableName, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to describe Oracle table: %w", err)
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var (
			name       string
			dataType   string
			nullable   int
			defaultVal sql.NullString
			isPK       int
		)
		if err := rows.Scan(&name, &dataType, &nullable, &defaultVal, &isPK); err != nil {
			return nil, fmt.Errorf("failed to scan Oracle table schema: %w", err)
		}

		var defaultValue interface{}
		if defaultVal.Valid {
			defaultValue = defaultVal.String
		}

		column := TableColumn{
			Name:         name,
			Type:         dataType,
			Nullable:     nullable == 1,
			Default:      defaultValue,
			IsPrimaryKey: isPK == 1,
		}
		columns = append(columns, column)
	}

	return columns, rows.Err()
}

// describeClickHouseTable handles ClickHouse table schema
func describeClickHouseTable(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := fmt.Sprintf(`DESCRIBE TABLE %s;`, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to describe ClickHouse table: %w", err)
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var (
			name        string
			typeInfo    string
			defaultType sql.NullString
			defaultExpr sql.NullString
			comment     sql.NullString
			codecExpr   sql.NullString
			ttlExpr     sql.NullString
		)
		if err := rows.Scan(&name, &typeInfo, &defaultType, &defaultExpr, &comment, &codecExpr, &ttlExpr); err != nil {
			return nil, fmt.Errorf("failed to scan ClickHouse table schema: %w", err)
		}

		var defaultValue interface{}
		if defaultExpr.Valid && defaultExpr.String != "" {
			defaultValue = defaultExpr.String
		}

		// ClickHouse has nullable types indicated by Nullable() wrapper
		nullable := strings.Contains(typeInfo, "Nullable(")

		column := TableColumn{
			Name:         name,
			Type:         typeInfo,
			Nullable:     nullable,
			Default:      defaultValue,
			IsPrimaryKey: false, // ClickHouse doesn't have traditional primary keys
		}
		columns = append(columns, column)
	}

	return columns, rows.Err()
}

// describeDuckDBTable handles DuckDB table schema
func describeDuckDBTable(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := fmt.Sprintf(`PRAGMA table_info('%s');`, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to describe DuckDB table: %w", err)
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var (
			cid        int
			name       string
			typeInfo   string
			notnull    bool
			defaultVal sql.NullString
			pk         bool
		)
		if err := rows.Scan(&cid, &name, &typeInfo, &notnull, &defaultVal, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan DuckDB table schema: %w", err)
		}

		var defaultValue interface{}
		if defaultVal.Valid {
			defaultValue = defaultVal.String
		}

		column := TableColumn{
			Name:         name,
			Type:         typeInfo,
			Nullable:     !notnull,
			Default:      defaultValue,
			IsPrimaryKey: pk,
		}
		columns = append(columns, column)
	}

	return columns, rows.Err()
}

// describeSnowflakeTable handles Snowflake table schema
func describeSnowflakeTable(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable = 'YES' AS nullable,
			column_default,
			CASE WHEN column_name IN (
				SELECT column_name
				FROM information_schema.table_constraints tc
				JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
				WHERE tc.table_name = ? AND tc.constraint_type = 'PRIMARY KEY'
			) THEN true ELSE false END AS is_primary_key
		FROM information_schema.columns
		WHERE table_name = ?
		ORDER BY ordinal_position;`

	rows, err := db.Query(query, strings.ToUpper(tableName), strings.ToUpper(tableName))
	if err != nil {
		return nil, fmt.Errorf("failed to describe Snowflake table: %w", err)
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var (
			name       string
			dataType   string
			nullable   bool
			defaultVal sql.NullString
			isPK       bool
		)
		if err := rows.Scan(&name, &dataType, &nullable, &defaultVal, &isPK); err != nil {
			return nil, fmt.Errorf("failed to scan Snowflake table schema: %w", err)
		}

		var defaultValue interface{}
		if defaultVal.Valid {
			defaultValue = defaultVal.String
		}

		column := TableColumn{
			Name:         name,
			Type:         dataType,
			Nullable:     nullable,
			Default:      defaultValue,
			IsPrimaryKey: isPK,
		}
		columns = append(columns, column)
	}

	return columns, rows.Err()
}
