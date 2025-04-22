package main

import (
	"context"
	"encoding/json"
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

	s.AddResourceTemplate(mcp.NewResourceTemplate(
		"schema://{table}",
		"Get schema information for a specific table",
		mcp.WithTemplateDescription("Returns schema information for a specific table"),
		mcp.WithTemplateMIMEType("application/json"),
	), func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract table name from URI
		tableName := request.Params.URI[len("schema://"):]
		if tableName == "" {
			return nil, errors.New("table name is required in schema URI")
		}

		schema, err := api.GetTableSchema(ctx, db, tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get table schema: %w", err)
		}

		jsonData, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal schema to JSON: %w", err)
		}

		return []mcp.ResourceContents{&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		}}, nil
	})

	s.AddResource(mcp.NewResource(
		"schema://all",
		"Get schema information for all tables",
	), func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		metadata, err := api.GetAllSchema(ctx, db)
		if err != nil {
			return nil, fmt.Errorf("failed to get all schema: %w", err)
		}

		jsonData, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata to JSON: %w", err)
		}

		return []mcp.ResourceContents{&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		}}, nil
	})

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
