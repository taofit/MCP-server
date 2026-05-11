# Weather MCP Server

A Model Context Protocol (MCP) server that provides weather information and alerts for locations within the United States using the National Weather Service (NWS) API.

## Features

- **Get Forecast**: Retrieve a 5-day weather forecast for any US location using latitude and longitude.
- **Get Alerts**: Check for active weather alerts and warnings in any US state.

## Prerequisites

- [Go](https://go.dev/doc/install) (version 1.21 or higher recommended)
- Internet connection (to access `api.weather.gov`)

## Setup and Installation

1. Clone the repository and navigate to the `weather` directory:
   ```bash
   cd weather
   ```

2. Initialize the Go module and install dependencies:
   ```bash
   go mod init weather
   go get github.com/modelcontextprotocol/go-sdk/mcp
   ```

3. Build the server:
   ```bash
   go build -o weather .
   ```

## Configuration

To use this server with an MCP client (like Claude Desktop), add it to your `mcp_config.json`:

```json
{
  "mcpServers": {
    "weather": {
      "command": "/Users/tao/Documents/mcp-servers/weather"
    }
  }
}
```

## Tools

### `get_forecast`
- **Description**: Get weather forecast for a location in the United States.
- **Arguments**:
  - `latitude` (number): Latitude of the location.
  - `longitude` (number): Longitude of the location.

### `get_alerts`
- **Description**: Get active weather alerts for a US state.
- **Arguments**:
  - `state` (string): Two-letter US state code (e.g., "CA", "NY").