# Create a new directory for our project
mkdir weather
cd weather

# Initialize Go module
go mod init weather

# Install dependencies
go get github.com/modelcontextprotocol/go-sdk/mcp

# Create our server file
touch main.go

# build and run
go build -o weather .

The compiled binary("./weather") can be placed in any directory and run as an MCP server. For example, if you place it in `/Users/tao/Documents/reboost/repos/go/MCP-server/weather`, you can run it by executing `/Users/tao/Documents/reboost/repos/go/MCP-server/weather`.

In the mcp_config.json, add the following:

```json
{
  "mcpServers": {
    "weather": {
      "command": "/Users/tao/Documents/reboost/repos/go/MCP-server/weather/weather"
    }
  }
}