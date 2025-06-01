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
	"strings"
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
		"0.3.0",
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	s.AddTool(mcp.NewTool(
		"db_type",
		mcp.WithDescription("Get the database type based on the DSN."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dbType, err := api.GetDBType(dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to get database type: %w", err)
		}

		return mcp.NewToolResultText(dbType), nil
	})

	s.AddTool(mcp.NewTool(
		"read_query",
		mcp.WithDescription("Execute a SELECT query and return the results."),
		mcp.WithString("query", mcp.Required(), mcp.Description("The SELECT query to execute.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid arguments format")
		}
		query, ok := args["query"].(string)
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
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid arguments format")
		}
		query, ok := args["query"].(string)
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
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid arguments format")
		}
		query, ok := args["query"].(string)
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
		"describe_table_schema",
		mcp.WithDescription("Get the JSON schema for a given table, including column names and data types, for all supported databases."),
		mcp.WithString("table", mcp.Required(), mcp.Description("The name of the table to describe.")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid arguments format")
		}
		tableName, ok := args["table"].(string)
		if !ok {
			return nil, errors.New("table must be a string")
		}

		schema, err := api.DescribeTableUniversal(db, tableName, dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to describe table schema: %w", err)
		}

		schemaJSON, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal schema to JSON: %w", err)
		}

		return mcp.NewToolResultText(string(schemaJSON)), nil
	})

	// Add resource template for usqlmcp://<table>/schema
	template := mcp.NewResourceTemplate(
		"usqlmcp://{table}/schema",
		"Table Schema",
		mcp.WithTemplateDescription("Returns the JSON schema for a given table, including column names and data types"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		uri := request.Params.URI

		if len(uri) < 10 || uri[:9] != "usqlmcp://" {
			return nil, fmt.Errorf("invalid URI scheme, expected usqlmcp://")
		}

		path := uri[9:] // Remove "usqlmcp://"
		parts := strings.Split(path, "/")
		if len(parts) != 2 || parts[1] != "schema" {
			return nil, fmt.Errorf("invalid URI format, expected usqlmcp://<table>/schema")
		}

		tableName := parts[0]
		if tableName == "" {
			return nil, fmt.Errorf("table name cannot be empty")
		}

		schema, err := api.DescribeTableUniversal(db, tableName, dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to describe table schema: %w", err)
		}

		schemaJSON, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal schema to JSON: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(schemaJSON),
			},
		}, nil
	})

	tables, err := api.ListTables(db, dsn)
	if err != nil {
		log.Printf("Warning: failed to list tables for resource registration: %v", err)
	} else {
		for _, tableName := range tables {
			resourceURI := fmt.Sprintf("usqlmcp://%s/schema", tableName)
			resource := mcp.NewResource(
				resourceURI,
				fmt.Sprintf("Schema for table %s", tableName),
			)

			tableNameCopy := tableName
			s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
				schema, err := api.DescribeTableUniversal(db, tableNameCopy, dsn)
				if err != nil {
					return nil, fmt.Errorf("failed to describe table schema: %w", err)
				}

				schemaJSON, err := json.MarshalIndent(schema, "", "  ")
				if err != nil {
					return nil, fmt.Errorf("failed to marshal schema to JSON: %w", err)
				}

				return []mcp.ResourceContents{
					mcp.TextResourceContents{
						URI:      request.Params.URI,
						MIMEType: "application/json",
						Text:     string(schemaJSON),
					},
				}, nil
			})
		}
		log.Printf("Registered %d table schema resources", len(tables))
	}

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
