package api

import (
	"database/sql"
	"fmt"

	_ "github.com/xo/usql/drivers"
)

type Row map[string]interface{}

// ReadQuery executes a SELECT query and returns the results as a slice of Row.
func ReadQuery(db *sql.DB, query string) ([]Row, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	numColumns := len(columns)

	results := []Row{}
	for rows.Next() {
		values := make([]interface{}, numColumns)

		// Create a slice of pointers to the elements in the 'values' slice.
		// Rows.Scan requires pointers to store the scanned values.
		scanArgs := make([]interface{}, numColumns)
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Scan the row's column values into the pointers.
		// The 'values' slice will be populated with the actual data.
		err := rows.Scan(scanArgs...)
		if err != nil {
			// Consider logging the error here too
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rowMap := make(Row, numColumns) // Pre-allocate map size for efficiency
		for i, colName := range columns {
			val := values[i]

			rowMap[colName] = val
		}

		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}
