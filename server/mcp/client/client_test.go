package client

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"testing"
)

// test MCP CustomerEndConnection
func TestMcpClientConnection(t *testing.T) {
	t.Skip("requires running MCP server at localhost:8888")
	c, err := NewClient("http://localhost:8888/sse", "test-client", "1.0.0", "gin-vue-admin MCPService")
	defer c.Close()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
}

func TestTools(t *testing.T) {
	t.Skip("requires running MCP server at localhost:8888")
	t.Run("currentTime", func(t *testing.T) {
		c, err := NewClient("http://localhost:8888/sse", "test-client", "1.0.0", "gin-vue-admin MCPService")
		defer c.Close()
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}
		ctx := context.Background()

		request := mcp.CallToolRequest{}
		request.Params.Name = "currentTime"
		request.Params.Arguments = map[string]interface{}{
			"timezone": "UTC+8",
		}

		result, err := c.CallTool(ctx, request)
		if err != nil {
			t.Fatalf("methodInvokeError: %v", err)
		}

		if len(result.Content) != 1 {
			t.Errorf("ShouldThisHaveAnd OnlyReturn1RowInformation, ButYesAppearAtHave %d", len(result.Content))
		}
		if content, ok := result.Content[0].(mcp.TextContent); ok {
			t.Logf("succeededReturnInformation%s", content.Text)
		} else {
			t.Logf("ReturnForstoptypeInformation%+v", content)
		}
	})

	t.Run("getNickname", func(t *testing.T) {

		c, err := NewClient("http://localhost:8888/sse", "test-client", "1.0.0", "gin-vue-admin MCPService")
		defer c.Close()
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}
		ctx := context.Background()

		// Initialize
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "test-client",
			Version: "1.0.0",
		}

		_, err = c.Initialize(ctx, initRequest)
		if err != nil {
			t.Fatalf("initializefailed: %v", err)
		}

		request := mcp.CallToolRequest{}
		request.Params.Name = "getNickname"
		request.Params.Arguments = map[string]interface{}{
			"username": "admin",
		}

		result, err := c.CallTool(ctx, request)
		if err != nil {
			t.Fatalf("methodInvokeError: %v", err)
		}

		if len(result.Content) != 1 {
			t.Errorf("ShouldThisHaveAnd OnlyReturn1RowInformation, ButYesAppearAtHave %d", len(result.Content))
		}
		if content, ok := result.Content[0].(mcp.TextContent); ok {
			t.Logf("succeededReturnInformation%s", content.Text)
		} else {
			t.Logf("ReturnForstoptypeInformation%+v", content)
		}
	})
}

func TestGetTools(t *testing.T) {
	t.Skip("requires running MCP server at localhost:8888")
	c, err := NewClient("http://localhost:8888/sse", "test-client", "1.0.0", "gin-vue-admin MCPService")
	defer c.Close()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	ctx := context.Background()

	toolsRequest := mcp.ListToolsRequest{}

	toolListResult, err := c.ListTools(ctx, toolsRequest)
	if err != nil {
		t.Fatalf("getUtilityListfailed: %v", err)
	}
	for i := range toolListResult.Tools {
		tool := toolListResult.Tools[i]
		fmt.Printf("Utilityname: %s\n", tool.Name)
		fmt.Printf("Utilitydescription: %s\n", tool.Description)

		// PrintParameterInformation
		if tool.InputSchema.Properties != nil {
			fmt.Println("Refercount list:")
			for paramName, prop := range tool.InputSchema.Properties {
				required := "No"
				// CheckParameterYesNoAtRequiredListIn
				for _, reqField := range tool.InputSchema.Required {
					if reqField == paramName {
						required = "Yes"
						break
					}
				}
				fmt.Printf("  - %s (type: %s, description: %s, Required: %s)\n",
					paramName, prop.(map[string]any)["type"], prop.(map[string]any)["description"], required)
			}
		} else {
			fmt.Println("ThisUtilityNoneParameter")
		}
		fmt.Println("-------------------")
	}
}
