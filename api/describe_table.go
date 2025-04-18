package api

import (
	"database/sql"
	"fmt"

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
