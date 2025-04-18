package api

import (
	"database/sql"
	"fmt"

	_ "github.com/xo/usql/drivers"
)

// CreateTable executes a CREATE TABLE SQL statement and returns a confirmation message.
func CreateTable(db *sql.DB, query string) (string, error) {
	_, err := db.Exec(query)
	if err != nil {
		return "", fmt.Errorf("failed to execute create table query: %w", err)
	}

	return "Table created successfully", nil
}
