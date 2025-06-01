# usqlmcp

A universal SQL MCP (Model Context Protocol).

[releases]: https://github.com/thesoulless/usqlmcp/releases "Releases"

## Features

- **Tools**
  - `read_query`: Execute a `SELECT` query and return the results.
  - `write_query`: Execute an `INSERT`, `UPDATE`, `DELETE`, or `ALTER` query and return the number of affected rows.
  - `create_table`: Execute a `CREATE TABLE` query to define new tables in the database.
  - `describe_table_schema`: Get the JSON schema for a given table, including column names and data types, for all supported databases.

- **Resources**
  - `usqlmcp://<table>/schema`: Access table schema as JSON resource for any table in the database.
  - Individual table schema resources are automatically discovered and registered for each table.

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

You can configure usqlmcp to run via Docker in Cursor by specifying the appropriate command and arguments in your MCP JSON configuration.

### Example: Postgres (Docker)

Add the following to your `~/.cursor/mcp.json`:

```json
{
    "mcpServers": {
        "usqlmcp": {
            "command": "docker",
            "args": [
                "run", "-i", "--rm",
                "ghcr.io/thesoulless/usqlmcp:latest",
                "--dsn", "postgres://username:password@host.docker.internal:5432/dbname?sslmode=disable"
            ]
        }
    }
}
```

Note: When connecting to a database on your host machine, use `host.docker.internal` instead of `localhost` or `127.0.0.1`.

### Example: SQLite with Volume Mounting (Docker)

To use a SQLite database file from your local machine, add the following to your `~/.cursor/mcp.json`:

```json
{
    "mcpServers": {
        "usqlmcp": {
            "command": "docker",
            "args": [
                "run", "-i", "--rm",
                "-v", "/path/to/local/mydatabase.db:/data/mydatabase.db",
                "ghcr.io/thesoulless/usqlmcp:latest",
                "--dsn", "sqlite3:///data/mydatabase.db"
            ]
        }
    }
}
```

This mounts your local SQLite database file directly into the container, providing access to only what's needed.

## Acknowledgments

This project depends on [usql](https://github.com/xo/usql), a universal command-line interface for SQL databases.
