package api

import (
	"context"
	"database/sql"
	"fmt"
)

type TableSchema struct {
	Name        string           `json:"name"`
	Columns     []ColumnInfo     `json:"columns"`
	Indexes     []IndexInfo      `json:"indexes"`
	Constraints []ConstraintInfo `json:"constraints"`
	Statistics  TableStatistics  `json:"statistics"`
}

type ColumnInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Default  string `json:"default,omitempty"`
}

type IndexInfo struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
	Primary bool     `json:"primary"`
}

type ConstraintInfo struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Columns    []string `json:"columns"`
	References string   `json:"references,omitempty"`
}

type TableStatistics struct {
	RowCount     int64  `json:"row_count"`
	SizeBytes    int64  `json:"size_bytes"`
	IndexSize    int64  `json:"index_size"`
	LastAnalyzed string `json:"last_analyzed"`
}

type DatabaseMetadata struct {
	Version    string        `json:"version"`
	SizeBytes  int64         `json:"size_bytes"`
	TableCount int           `json:"table_count"`
	Tables     []TableSchema `json:"tables"`
}

func GetTableSchema(ctx context.Context, db *sql.DB, tableName string) (*TableSchema, error) {
	columns, err := getColumnInfo(ctx, db, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get column info: %w", err)
	}

	indexes, err := getIndexInfo(ctx, db, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get index info: %w", err)
	}

	constraints, err := getConstraintInfo(ctx, db, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get constraint info: %w", err)
	}

	stats, err := getTableStatistics(ctx, db, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table statistics: %w", err)
	}

	return &TableSchema{
		Name:        tableName,
		Columns:     columns,
		Indexes:     indexes,
		Constraints: constraints,
		Statistics:  stats,
	}, nil
}

func GetAllSchema(ctx context.Context, db *sql.DB) (*DatabaseMetadata, error) {
	tables, err := ListTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}

	metadata := &DatabaseMetadata{
		Tables: make([]TableSchema, 0, len(tables)),
	}

	for _, tableName := range tables {
		schema, err := GetTableSchema(ctx, db, tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get schema for table %s: %w", tableName, err)
		}
		metadata.Tables = append(metadata.Tables, *schema)
	}

	version, size, err := getDatabaseInfo(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to get database info: %w", err)
	}

	metadata.Version = version
	metadata.SizeBytes = size
	metadata.TableCount = len(tables)

	return metadata, nil
}

func getColumnInfo(ctx context.Context, db *sql.DB, tableName string) ([]ColumnInfo, error) {
	var exists int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if table exists: %w", err)
	}
	if exists == 0 {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	rows, err := db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return nil, fmt.Errorf("failed to get column info: %w", err)
	}
	defer rows.Close()

	var columns []ColumnInfo
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
			return nil, fmt.Errorf("failed to scan column info: %w", err)
		}

		columns = append(columns, ColumnInfo{
			Name:     name,
			Type:     typeInfo,
			Nullable: notnull == 0,
			Default:  defaultVal.String,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating column info: %w", err)
	}

	return columns, nil
}

func getIndexInfo(ctx context.Context, db *sql.DB, tableName string) ([]IndexInfo, error) {
	// Implementation depends on the database driver
	return nil, nil
}

func getConstraintInfo(ctx context.Context, db *sql.DB, tableName string) ([]ConstraintInfo, error) {
	// Implementation depends on the database driver
	return nil, nil
}

func getTableStatistics(ctx context.Context, db *sql.DB, tableName string) (TableStatistics, error) {
	// Implementation depends on the database driver
	return TableStatistics{}, nil
}

func getDatabaseInfo(ctx context.Context, db *sql.DB) (string, int64, error) {
	// Implementation depends on the database driver
	return "", 0, nil
}
