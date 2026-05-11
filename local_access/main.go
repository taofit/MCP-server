package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)


type CreateDirectoryArgs struct {
	Path string `json:"path"`
}

type ListFilesArgs struct {
	Path string `json:"path"`
}

type ReadFileArgs struct {
	Path string `json:"path"`
}

type WriteFileArgs struct {
	Path string `json:"path"`
	Content string `json:"content"`
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "local-access",
		Version: "0.0.1",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_files",
		Description: "list files in a directory",
	}, listFiles)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_file",
		Description: "read a file",
	}, readFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "write_file",
		Description: "write to a file",
	}, writeFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_directory",
		Description: "create a directory",
	}, createDirectory)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

func createDirectory(ctx context.Context, req *mcp.CallToolRequest, args *CreateDirectoryArgs) (*mcp.CallToolResult, any, error) {
	err := os.MkdirAll(args.Path, 0755)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error creating directory: %v", err)},
			},
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Successfully created directory: %s", args.Path)},
		},
	}, nil, nil
}

func writeFile(ctx context.Context, req *mcp.CallToolRequest, args *WriteFileArgs) (*mcp.CallToolResult, any, error){
	file, err := os.OpenFile(args.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error opening file: %v", err)},
			},
		}, nil, nil
	}
	defer file.Close()

	_, err = file.WriteString(args.Content)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error writing to file: %v", err)},
			},
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Successfully wrote to file: %s", args.Path)},
		},
	}, nil, nil
}

func readFile(ctx context.Context, req *mcp.CallToolRequest, args *ReadFileArgs) (*mcp.CallToolResult, any, error) {
	data, err := os.ReadFile(args.Path)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error reading file: %v", err)},
			},
		}, nil, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}


func listFiles(ctx context.Context, req *mcp.CallToolRequest, args *ListFilesArgs) (*mcp.CallToolResult, any, error) {
	// Read the directory contents
	entries, err := os.ReadDir(args.Path)
	if err != nil {
		log.Printf("Error reading directory %s: %v", args.Path, err)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error reading directory: %v", err),
				},
			},
		}, nil, nil
	}

	// Build list of files and directories
	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			files = append(files, fmt.Sprintf("[DIR] %s", entry.Name()))
		} else {
			files = append(files, entry.Name())
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Files in %s:\n%s", args.Path, strings.Join(files, "\n")),
			},
		},
	}, nil, nil
}
