# Local Access MCP Server

A Model Context Protocol (MCP) server that provides secure access to the local filesystem, allowing clients to perform basic file operations.

## Features

- **List Files**: Explore directory contents.
- **Read/Write Files**: Access and modify file content.
- **Create Directories**: Organize files into folders.
- **Delete Files/Directories**: Clean up or remove local data.

## Prerequisites

- [Go](https://go.dev/doc/install) (version 1.21 or higher recommended)

## Setup and Installation

1. Navigate to the `local_access` directory:
   ```bash
   cd local_access
   ```

2. Initialize the Go module and install dependencies:
   ```bash
   go mod init local_access
   go get github.com/modelcontextprotocol/go-sdk/mcp
   ```

3. Build the server:
   ```bash
   go build -o localaccess .
   ```

## Configuration

To use this server with an MCP client (like Claude Desktop), add it to your `mcp_config.json`:

```json
{
  "mcpServers": {
    "local-access": {
      "command": "/Users/tao/Documents/mcp-servers/localaccess"
    }
  }
}
```

## Tools

### `list_files`
- **Description**: List files and directories in a specified path.
- **Arguments**:
  - `path` (string): The directory path to list.

### `read_file`
- **Description**: Read the content of a local file.
- **Arguments**:
  - `path` (string): The path to the file to be read.

### `write_file`
- **Description**: Write content to a local file (appends if the file exists).
- **Arguments**:
  - `path` (string): The path to the file.
  - `content` (string): The text content to write.

### `create_directory`
- **Description**: Create a new directory (includes parent directories if necessary).
- **Arguments**:
  - `path` (string): The directory path to create.

### `delete_file_or_directory`
- **Description**: Delete a file or an empty directory.
- **Arguments**:
  - `path` (string): The path to delete.
