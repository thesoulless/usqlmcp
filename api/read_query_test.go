package api_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/thesoulless/usqlmcp/api"
)

func TestReadQueryWithSQLite(t *testing.T) {
	// Create a temporary SQLite database file
	dbFile := "test.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Create a test table
	createTableQuery := `CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Insert test data
	insertDataQuery := `INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 25);`
	_, err = db.Exec(insertDataQuery)
	if err != nil {
		t.Fatalf("failed to insert data: %v", err)
	}

	// Execute the ReadQuery function
	selectQuery := `SELECT * FROM users;`
	results, err := api.ReadQuery(db, selectQuery)
	if err != nil {
		t.Fatalf("ReadQuery failed: %v", err)
	}

	// Validate the results
	expected := []api.Row{
		{"id": int64(1), "name": "Alice", "age": int64(30)},
		{"id": int64(2), "name": "Bob", "age": int64(25)},
	}
	assert.Equal(t, expected, results, "results do not match expected output")
}
