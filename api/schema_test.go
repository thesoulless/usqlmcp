package api

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetTableSchema(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT);`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	ctx := context.Background()

	_, err = GetTableSchema(ctx, db, "non_existent_table")
	if err == nil {
		t.Error("Expected error for non-existent table, got nil")
	}

	schema, err := GetTableSchema(ctx, db, "test_table")
	if err != nil {
		t.Errorf("Expected no error for existing table, got %v", err)
	}
	if schema == nil {
		t.Fatal("Expected schema for existing table, got nil")
	}
	if schema.Name != "test_table" {
		t.Errorf("Expected table name 'test_table', got '%s'", schema.Name)
	}
	if len(schema.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(schema.Columns))
	}
}
