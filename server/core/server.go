package core

import (
	"fmt"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/initialize"
	mcpTool "github.com/huuhoait/gin-vue-admin/server/mcp"
	"github.com/huuhoait/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
)

func RunServer() {
	if global.GVA_CONFIG.System.UseRedis {
		initialize.Redis()
		if global.GVA_CONFIG.System.UseMultipoint {
			initialize.RedisList()
		}
	}

	if global.GVA_CONFIG.System.UseMongo {
		if err := initialize.Mongo.Initialization(); err != nil {
			zap.L().Error(fmt.Sprintf("%+v", err))
		}
	}

	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	mcpBaseURL := mcpTool.ResolveMCPServiceURL()

	fmt.Printf(`
	WelcomeUse gin-vue-admin
	CurrentVersion:%s
	ProjectAddress:https://github.com/huuhoait/gin-vue-admin
	Plugin Market:https://plugin.gin-vue-admin.com
	defaultAutomatictransformDocumentationAddress:http://127.0.0.1%s/swagger/index.html
	MCP standalone servicePleaseManualStart: go run ./cmd/mcp -config ./cmd/mcp/config.yaml
	defaultMCP StreamHTTPAddress:%s
	defaultFrontendFileRunAddress:http://127.0.0.1:8080
`, global.Version, address, mcpBaseURL)

	initServer(address, Router, 10*time.Minute, 10*time.Minute)
}
