# usqlmcp

A universal SQL MCP (Model Context Protocol).

[releases]: https://github.com/thesoulless/usqlmcp/releases "Releases"

## Features

- **Tools**
  - `read_query`: Execute a `SELECT` query and return the results.
  - `write_query`: Execute an `INSERT`, `UPDATE`, `DELETE`, or `ALTER` query and return the number of affected rows.
  - `create_table`: Execute a `CREATE TABLE` query to define new tables in the database.
  - `list_tables`: Retrieve a list of all table names in the database.
  - `describe_table`: Retrieve schema information for a specific table.

- **Resources**
  - `schema://all`: Get schema information for all tables
  - `schema://{table}`: Get schema information for a specific table

## Installing
`usqlmcp` is available [via Release][]

[via Release]: #installing-via-release

### Installing via Release

1. [Download a release for your platform][releases]
2. Extract the `usqlmcp` or `usqlcmp.exe` file from the `.tar.bz2` or `.zip` file
3. Move the extracted executable to somewhere on your `$PATH` (Linux/macOS) or
   `%PATH%` (Windows)

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

## Docker Usage

You can also run usqlmcp using Docker:

```bash
docker run -i --rm ghcr.io/thesoulless/usqlmcp:latest --dsn postgres://username:password@host.docker.internal:5432/dbname?sslmode=disable
```

Note: When connecting to a database on your host machine, use `host.docker.internal` instead of `localhost` or `127.0.0.1`.

## Acknowledgments

This project depends on [usql](https://github.com/xo/usql), a universal command-line interface for SQL databases.
