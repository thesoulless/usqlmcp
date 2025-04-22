package api_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/thesoulless/usqlmcp/api"
)

func TestCreateTable(t *testing.T) {
	dbFile := "test_create_table.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("failed to open SQLite database: %v", err)
	}
	defer db.Close()

	createTableQuery := `CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT);`

	message, err := api.CreateTable(db, createTableQuery)
	if err != nil {
		t.Fatalf("CreateTable failed: %v", err)
	}

	assert.Equal(t, "Table created successfully", message, "unexpected success message")

	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='test_table';`
	row := db.QueryRow(query)
	var tableName string
	if err := row.Scan(&tableName); err != nil {
		t.Fatalf("failed to verify table existence: %v", err)
	}
	assert.Equal(t, "test_table", tableName, "table name does not match")
}
