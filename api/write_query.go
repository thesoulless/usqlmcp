package api

import (
	"database/sql"
	"fmt"

	_ "github.com/xo/usql/drivers"
)

// WriteQuery executes an INSERT, UPDATE, DELETE, or ALTER query and returns the number of affected rows.
func WriteQuery(db *sql.DB, query string) (int64, error) {
	result, err := db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	return affectedRows, nil
}
