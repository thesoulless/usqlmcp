package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/thesoulless/usqlmcp/api"
	_ "github.com/thesoulless/usqlmcp/internal"
	"github.com/xo/dburl"
	"github.com/xo/usql/drivers"
)

func main() {
	dsnFlag := flag.String("dsn", "", "Database connection string")
	flag.Parse()

	dsn := *dsnFlag
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "Error: DSN is required. Provide it using --dsn flag or DB_DSN environment variable.")
		os.Exit(100)
	}

	ctx := context.Background()

	u, err := dburl.Parse(dsn)
	if err != nil {
		log.Fatalf("Failed to parse DSN: %v", err)
		fmt.Fprintln(os.Stderr, "Error: Invalid DSN format.")
		os.Exit(101)
	}

	openCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	db, err := drivers.Open(openCtx, u, func() io.Writer { return os.Stdout }, func() io.Writer { return os.Stderr })
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		fmt.Fprintln(os.Stderr, "Error: Failed to open database.")
		os.Exit(102)
	}
	defer db.Close()

	s := server.NewMCPServer(
		"USQL MCP Server",
		"0.1.0",
		server.WithLogging(),
		server.WithRecovery(),
	)

	s.AddTool(mcp.NewTool(
		"read_query",
		mcp.WithDescription("Execute a SELECT query and return the results."),
		mcp.WithString("query", mcp.Required(), mcp.Description("The SELECT query to execute.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, ok := request.Params.Arguments["query"].(string)
		if !ok {
			return nil, errors.New("query must be a string")
		}

		results, err := api.ReadQuery(db, query)
		if err != nil {
			return nil, fmt.Errorf("failed to execute read query: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("%v", results)), nil
	})

	s.AddTool(mcp.NewTool(
		"write_query",
		mcp.WithDescription("Execute an INSERT, UPDATE, DELETE, or ALTER query and return the number of affected rows."),
		mcp.WithString("query", mcp.Required(), mcp.Description("The query to execute.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, ok := request.Params.Arguments["query"].(string)
		if !ok {
			return nil, errors.New("query must be a string")
		}

		affectedRows, err := api.WriteQuery(db, query)
		if err != nil {
			return nil, fmt.Errorf("failed to execute write query: %w", err)
		}

		if query[:6] == "ALTER " && affectedRows == 0 {
			return mcp.NewToolResultText("ALTER query executed successfully, but no rows were affected."), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("%d rows affected", affectedRows)), nil
	})

	s.AddTool(mcp.NewTool(
		"create_table",
		mcp.WithDescription("Execute a CREATE TABLE query."),
		mcp.WithString("query", mcp.Required(), mcp.Description("The CREATE TABLE query to execute.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, ok := request.Params.Arguments["query"].(string)
		if !ok {
			return nil, errors.New("query must be a string")
		}

		message, err := api.CreateTable(db, query)
		if err != nil {
			return nil, fmt.Errorf("failed to execute create table query: %w", err)
		}

		return mcp.NewToolResultText(message), nil
	})

	s.AddTool(mcp.NewTool(
		"list_tables",
		mcp.WithDescription("Retrieve a list of all table names in the database."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tables, err := api.ListTables(db)
		if err != nil {
			return nil, fmt.Errorf("failed to list tables: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("%v", tables)), nil
	})

	s.AddTool(mcp.NewTool(
		"describe_table",
		mcp.WithDescription("Retrieve schema information for a specific table."),
		mcp.WithString("table_name", mcp.Required(), mcp.Description("The name of the table to describe.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tableName, ok := request.Params.Arguments["table_name"].(string)
		if !ok {
			return nil, errors.New("table_name must be a string")
		}

		schema, err := api.DescribeTable(db, tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to describe table: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("%v", schema)), nil
	})

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
