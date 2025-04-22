package api_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/thesoulless/usqlmcp/api"
)

func TestListTables(t *testing.T) {
	dbFile := "test_list_tables.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("failed to open SQLite database: %v", err)
	}
	defer db.Close()

	createTableQueries := []string{
		`CREATE TABLE table1 (id INTEGER PRIMARY KEY);`,
		`CREATE TABLE table2 (name TEXT);`,
	}
	for _, query := range createTableQueries {
		_, err = db.Exec(query)
		if err != nil {
			t.Fatalf("failed to create table: %v", err)
		}
	}

	tables, err := api.ListTables(db)
	if err != nil {
		t.Fatalf("ListTables failed: %v", err)
	}

	expected := []string{"table1", "table2"}
	assert.ElementsMatch(t, expected, tables, "table list does not match expected output")
}
