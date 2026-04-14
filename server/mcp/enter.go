package mcpTool

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// McpTool defines the interface MCP tools must implement
type McpTool interface {
	// Handle returns tool invocation information
	Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	// New returns tool registration information
	New() mcp.Tool
}

// tool registry
var toolRegister = make(map[string]McpTool)

// RegisterTool is called from init to register the tool in the registry
func RegisterTool(tool McpTool) {
	mcpTool := tool.New()
	toolRegister[mcpTool.Name] = tool
}

// RegisterAllTools registers all tools with the MCP server
func RegisterAllTools(mcpServer *server.MCPServer) {
	for _, tool := range toolRegister {
		mcpServer.AddTool(tool.New(), tool.Handle)
	}
}
