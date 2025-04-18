package api_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/thesoulless/usqlmcp/api"
)

func TestWriteQuery(t *testing.T) {
	// Create a temporary SQLite database file
	dbFile := "test_write_query.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Create a test table
	createTableQuery := `CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Insert data using WriteQuery
	insertQuery := `INSERT INTO test_table (name) VALUES ('Alice'), ('Bob');`
	affectedRows, err := api.WriteQuery(db, insertQuery)
	if err != nil {
		t.Fatalf("WriteQuery failed: %v", err)
	}

	// Validate the number of affected rows
	assert.Equal(t, int64(2), affectedRows, "unexpected number of affected rows")

	// Verify the data was inserted
	query := `SELECT COUNT(*) FROM test_table;`
	row := db.QueryRow(query)
	var count int
	if err := row.Scan(&count); err != nil {
		t.Fatalf("failed to verify inserted data: %v", err)
	}
	assert.Equal(t, 2, count, "unexpected row count in table")
}
