# Go MCP Servers Collection

This repository contains a collection of Model Context Protocol (MCP) servers implemented in Go. These servers allow AI models (like Claude) to interact with local tools and external APIs.

## Included Servers

### 1. [Hello (Greeter)](./hello)
A simple introductory server that provides greeting functionality and tracks uptime.
- **Tools**: `greet`, `age`

### 2. [Local Access](./local_access)
Provides secure access to the local filesystem for file and directory operations.
- **Tools**: `list_files`, `read_file`, `write_file`, `create_directory`, `delete_file_or_directory`

### 3. [Weather](./weather)
Retrieves weather forecasts and active alerts for US locations using the National Weather Service API.
- **Tools**: `get_forecast`, `get_alerts`

## Prerequisites

- [Go](https://go.dev/doc/install) (version 1.21+)
- Git

## Getting Started

Each server is located in its own directory. To build a specific server:

1. Navigate to the server directory:
   ```bash
   cd <server_name>
   ```
2. Build the binary:
   ```bash
   go build -o <binary_name> .
   ```

## Configuration

To use these servers with Claude Desktop or any other MCP client, add them to your configuration file (usually `mcp_config.json` or `claude_desktop_config.json`).

Example configuration:

```json
{
  "mcpServers": {
    "hello": {
      "command": "/Users/tao/Documents/mcp-servers/hello"
    },
    "local-access": {
      "command": "/Users/tao/Documents/mcp-servers/localaccess"
    },
    "weather": {
      "command": "/Users/tao/Documents/mcp-servers/weather"
    }
  }
}
```

> [!NOTE]
> Ensure the `command` paths point to the actual locations of the compiled binaries on your system.

## Development

Each server directory contains its own `go.mod` and `main.go`. You can modify the tools or add new ones by following the MCP Go SDK patterns demonstrated in the source code.
