package api_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thesoulless/usqlmcp/api"
)

func TestDescribeTable(t *testing.T) {
	dbFile := "test_describe_table.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("failed to open SQLite database: %v", err)
	}
	defer db.Close()

	createTableQuery := `CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT, age INTEGER);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	schema, err := api.DescribeTable(db, "test_table")
	if err != nil {
		t.Fatalf("DescribeTable failed: %v", err)
	}

	expected := []map[string]interface{}{
		{"cid": 0, "name": "id", "type": "INTEGER", "notnull": 0, "default": "", "primary_key": 1},
		{"cid": 1, "name": "name", "type": "TEXT", "notnull": 0, "default": "", "primary_key": 0},
		{"cid": 2, "name": "age", "type": "INTEGER", "notnull": 0, "default": "", "primary_key": 0},
	}
	assert.Equal(t, expected, schema, "schema does not match expected output")
}

func TestDescribeTableUniversal(t *testing.T) {
	dbFile := "test_describe_table_universal.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	require.NoError(t, err, "failed to open SQLite database")
	defer db.Close()

	createTableQuery := `CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT NOT NULL, age INTEGER, email TEXT DEFAULT 'no-email');`
	_, err = db.Exec(createTableQuery)
	require.NoError(t, err, "failed to create table")

	dsn := "sqlite3://" + dbFile
	schema, err := api.DescribeTableUniversal(db, "test_table", dsn)
	require.NoError(t, err, "DescribeTableUniversal failed")

	require.Len(t, schema, 4, "expected 4 columns")

	// Check first column (id)
	assert.Equal(t, "id", schema[0].Name)
	assert.Equal(t, "INTEGER", schema[0].Type)
	assert.True(t, schema[0].Nullable)
	assert.True(t, schema[0].IsPrimaryKey)

	// Check second column (name)
	assert.Equal(t, "name", schema[1].Name)
	assert.Equal(t, "TEXT", schema[1].Type)
	assert.False(t, schema[1].Nullable)
	assert.False(t, schema[1].IsPrimaryKey)

	// Check third column (age)
	assert.Equal(t, "age", schema[2].Name)
	assert.Equal(t, "INTEGER", schema[2].Type)
	assert.True(t, schema[2].Nullable)
	assert.False(t, schema[2].IsPrimaryKey)

	// Check fourth column (email)
	assert.Equal(t, "email", schema[3].Name)
	assert.Equal(t, "TEXT", schema[3].Type)
	assert.True(t, schema[3].Nullable)
	assert.Equal(t, "'no-email'", schema[3].Default)
	assert.False(t, schema[3].IsPrimaryKey)
}

func TestDescribeTableUniversalUnsupportedDriver(t *testing.T) {
	dbFile := "test_unsupported.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	require.NoError(t, err, "failed to open SQLite database")
	defer db.Close()

	// Test with invalid DSN first
	dsn := "unsupported://localhost/test"
	_, err = api.DescribeTableUniversal(db, "test_table", dsn)
	assert.Error(t, err, "expected error for invalid DSN")
	assert.Contains(t, err.Error(), "failed to parse DSN")

	// Test with valid DSN format but unsupported driver (using a driver that dburl recognizes but we don't handle)
	dsn = "adodb://localhost/test"
	_, err = api.DescribeTableUniversal(db, "test_table", dsn)
	assert.Error(t, err, "expected error for unsupported driver")
	assert.Contains(t, err.Error(), "unsupported database driver")
}

func TestListTables(t *testing.T) {
	dbFile := "test_list_tables.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	require.NoError(t, err, "failed to open SQLite database")
	defer db.Close()

	// Create some test tables
	_, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);`)
	require.NoError(t, err, "failed to create users table")

	_, err = db.Exec(`CREATE TABLE products (id INTEGER PRIMARY KEY, name TEXT, price REAL);`)
	require.NoError(t, err, "failed to create products table")

	dsn := "sqlite3://" + dbFile
	tables, err := api.ListTables(db, dsn)
	require.NoError(t, err, "ListTables failed")

	require.Len(t, tables, 2, "expected 2 tables")
	assert.Contains(t, tables, "users")
	assert.Contains(t, tables, "products")
}
