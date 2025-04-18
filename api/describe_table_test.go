package api_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/thesoulless/usqlmcp/api"
)

func TestDescribeTable(t *testing.T) {
	// Create a temporary SQLite database file
	dbFile := "test_describe_table.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Create a test table
	createTableQuery := `CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT, age INTEGER);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Execute the DescribeTable function
	schema, err := api.DescribeTable(db, "test_table")
	if err != nil {
		t.Fatalf("DescribeTable failed: %v", err)
	}

	// Validate the schema
	expected := []map[string]interface{}{
		{"cid": 0, "name": "id", "type": "INTEGER", "notnull": 0, "default": "", "primary_key": 1},
		{"cid": 1, "name": "name", "type": "TEXT", "notnull": 0, "default": "", "primary_key": 0},
		{"cid": 2, "name": "age", "type": "INTEGER", "notnull": 0, "default": "", "primary_key": 0},
	}
	assert.Equal(t, expected, schema, "schema does not match expected output")
}
