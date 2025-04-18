# usqlmcp

A universal SQL MCP (Model Context Protocol).

## Features

- **Tools**
  - `read_query`: Execute a `SELECT` query and return the results.
  - `write_query`: Execute an `INSERT`, `UPDATE`, `DELETE`, or `ALTER` query and return the number of affected rows.
  - `create_table`: Execute a `CREATE TABLE` query to define new tables in the database.
  - `list_tables`: Retrieve a list of all table names in the database.
  - `describe_table`: Retrieve schema information for a specific table.

# Setup

## MCP Integration in Cursor

Add the following configuration to your ~/.cursor/mcp.json file or configure via the settings menu in Cursor.

```json
{
    "mcpServers": {
        "usqlmcp": {
            "command": "usqlmcp",
            "args": ["--dsn", "sqlite3:///your/db/dsn/file.db"]
        }
    }
}
```

## Other tools

```json
{
    "servers": {
        "usqlmcp": {
            "type": "stdio",
            "command": "usqlmcp",
            "args": ["--dsn", "sqlite3:///your/db/dsn/file.db"]
        }
    }
}
```