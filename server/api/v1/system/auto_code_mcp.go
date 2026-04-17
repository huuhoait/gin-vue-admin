package system

import (
	"strings"

	"github.com/huuhoait/gin-vue-admin/server/global"
	mcpTool "github.com/huuhoait/gin-vue-admin/server/mcp"
	"github.com/huuhoait/gin-vue-admin/server/mcp/client"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
)

func (a *AutoCodeTemplateApi) MCP(c *gin.Context) {
	var info request.AutoMcpTool
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	toolFilePath, err := autoCodeTemplateService.CreateMcp(c.Request.Context(), info)
	if err != nil {
		response.FailWithMessage("Creation failed", c)
		global.GVA_LOG.Error(err.Error())
		return
	}
	response.OkWithMessage("Created successfully, MCP Tool path: "+toolFilePath, c)
}

func (a *AutoCodeTemplateApi) MCPStatus(c *gin.Context) {
	response.OkWithData(gin.H{
		"status":          mcpTool.GetManagedStandaloneStatus(c.Request.Context()),
		"mcpServerConfig": buildMCPServerConfig(),
	}, c)
}

func (a *AutoCodeTemplateApi) MCPStart(c *gin.Context) {
	status, err := mcpTool.StartManagedStandalone(c.Request.Context())
	if err != nil {
		response.FailWithDetailed(gin.H{
			"status":          status,
			"mcpServerConfig": buildMCPServerConfig(),
		}, err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"status":          status,
		"mcpServerConfig": buildMCPServerConfig(),
	}, "MCP standalone service started", c)
}

func (a *AutoCodeTemplateApi) MCPStop(c *gin.Context) {
	status, err := mcpTool.StopManagedStandalone(c.Request.Context())
	if err != nil {
		response.FailWithDetailed(gin.H{
			"status":          status,
			"mcpServerConfig": buildMCPServerConfig(),
		}, err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"status":          status,
		"mcpServerConfig": buildMCPServerConfig(),
	}, "MCP standalone service stopped", c)
}

func (a *AutoCodeTemplateApi) MCPList(c *gin.Context) {
	baseURL := mcpTool.ResolveMCPServiceURL()
	testClient, err := client.NewClient(baseURL, "testClient", "v1.0.0", mcpServerName(), incomingMCPHeaders(c))
	if err != nil {
		response.FailWithDetailed(gin.H{
			"status":          mcpTool.GetManagedStandaloneStatus(c.Request.Context()),
			"mcpServerConfig": buildMCPServerConfig(),
		}, "Failed to connect to MCP service: "+err.Error(), c)
		return
	}
	defer testClient.Close()

	list, err := testClient.ListTools(c.Request.Context(), mcp.ListToolsRequest{})
	if err != nil {
		response.FailWithDetailed(gin.H{
			"status":          mcpTool.GetManagedStandaloneStatus(c.Request.Context()),
			"mcpServerConfig": buildMCPServerConfig(),
		}, "Failed to retrieve tool list: "+err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"status":          mcpTool.GetManagedStandaloneStatus(c.Request.Context()),
		"mcpServerConfig": buildMCPServerConfig(),
		"list":            list,
	}, c)
}

func (a *AutoCodeTemplateApi) MCPRoutes(c *gin.Context) {
	response.OkWithData(gin.H{
		"routes": global.GVA_ROUTERS,
	}, c)
}

func (a *AutoCodeTemplateApi) MCPTest(c *gin.Context) {
	var testRequest struct {
		Name      string                 `json:"name" binding:"required"`
		Arguments map[string]interface{} `json:"arguments" binding:"required"`
	}
	if err := c.ShouldBindJSON(&testRequest); err != nil {
		response.FailWithMessage("Parameter parsing failed: "+err.Error(), c)
		return
	}

	baseURL := mcpTool.ResolveMCPServiceURL()
	testClient, err := client.NewClient(baseURL, "testClient", "v1.0.0", mcpServerName(), incomingMCPHeaders(c))
	if err != nil {
		response.FailWithMessage("Failed to connect to MCP service: "+err.Error(), c)
		return
	}
	defer testClient.Close()

	callRequest := mcp.CallToolRequest{}
	callRequest.Params.Name = testRequest.Name
	callRequest.Params.Arguments = testRequest.Arguments

	result, err := testClient.CallTool(c.Request.Context(), callRequest)
	if err != nil {
		response.FailWithMessage("Tool invocation failed: "+err.Error(), c)
		return
	}
	if len(result.Content) == 0 {
		response.FailWithMessage("Tool returned no content", c)
		return
	}

	response.OkWithData(result.Content, c)
}

func incomingMCPHeaders(c *gin.Context) map[string]string {
	headerName := mcpTool.ConfiguredAuthHeader()
	headerValue := c.GetHeader(headerName)
	if headerValue == "" {
		return nil
	}

	return map[string]string{
		headerName: headerValue,
	}
}

func buildMCPServerConfig() map[string]interface{} {
	baseURL := mcpTool.ResolveMCPServiceURL()
	authHeader := mcpTool.ConfiguredAuthHeader()
	serverName := mcpServerName()

	return map[string]interface{}{
		"mcpServers": map[string]interface{}{
			serverName: map[string]interface{}{
				"url": baseURL,
				"headers": map[string]string{
					authHeader: "${YOUR_GVA_TOKEN}",
				},
			},
		},
	}
}

func mcpServerName() string {
	if name := strings.TrimSpace(global.GVA_CONFIG.MCP.Name); name != "" {
		return name
	}
	return "gva"
}
