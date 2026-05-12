# User Management

An MCP server that provides user management capabilities backed by a PostgreSQL database.

## Tools

| Tool | Description |
|------|-------------|
| `list_users` | List all users in the database |
| `get_user` | Get a user by ID |
| `add_user` | Add a new user |
| `delete_user` | Delete a user by ID |

## Prerequisites

- [Go](https://go.dev/doc/install) 1.21+
- [Docker](https://www.docker.com/) & Docker Compose

## Getting Started

1. Create a `.env` file with your database credentials:
   ```env
   DB_HOST=db
   DB_PORT=5432
   DB_USER=<db_user>
   DB_PASSWORD=<db_password>
   DB_NAME=<db_name>
   ```
2. run `docker-compose up` to start the application (including the database), its to test the app directly, not the mcp server
3. to use the mcp server, first, run `docker-compose up db` to start the database, as it is the backend of the mcp server that will modify the database data. after that, build the server: `go build -o usermanagement main.go`
4. then, add the mcp server to your `mcp_config.json` of your mcp client: be careful with DB_HOST: localhost
   ```json
   {
     "mcpServers": {
       "user-management": {
         "command": "/path/to/usermanagement",
         "env": {
           "DB_HOST": "localhost",
           "DB_PORT": "5432",
           "DB_USER": "mcp_user",
           "DB_PASSWORD": "mcp_password_db",
           "DB_NAME": "app_management"
         }
       }
     }
   }
   ```
6. then, you can use the mcp server in Claude, Antigravity or Reboost