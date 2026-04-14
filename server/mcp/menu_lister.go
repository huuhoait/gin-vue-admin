package mcpTool

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/mark3labs/mcp-go/mcp"
)

func init() {
	RegisterTool(&MenuLister{})
}

type MenuListResponse struct {
	Success     bool                 `json:"success"`
	Message     string               `json:"message"`
	Menus       []system.SysBaseMenu `json:"menus"`
	TotalCount  int                  `json:"totalCount"`
	Description string               `json:"description"`
}

type MenuLister struct{}

func (m *MenuLister) New() mcp.Tool {
	return mcp.NewTool("list_all_menus",
		mcp.WithDescription(`getSystemInAllMenuInformation, PackageIncludeMenuTreeStructure, route information, Componentpathetc., Used forFrontendeditWritevue-routerWhenCorrectJump

**Description:**
- ReturnCompleteArrangeofMenutree structure
- PackageIncludeRouteconfigurationInformation(path, name, component)
- PackageIncludeMenuYuanData(title, icon, keepAliveetc.)
- PackageIncludeMenuParameterAndButtonconfiguration
- SupportParentchild menusRelationshipDisplay

**Use cases:**
- FrontendRouteconfiguration:getAllMenuInformationUsed forconfigurationvue-router
- MenuPermissionmanagement:DonesolveSystemInAllCanUseofMenuItem
- GuidenavigateComponentDevelopment:BuildDynamicGuidenavigateMenu
- SystemarchitecturePartAnalyze:DonesolveSystemofMenuStructureAndPageOrganization`),
		mcp.WithString("_placeholder",
			mcp.Description("placeholder to avoid json schema validation failure"),
		),
	)
}

func (m *MenuLister) Handle(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resp, err := postUpstream[[]system.SysBaseMenu](ctx, "/menu/getMenuList", map[string]any{})
	if err != nil {
		return nil, err
	}

	return textResultWithJSON("", MenuListResponse{
		Success:     true,
		Message:     "getMenuListsucceeded",
		Menus:       resp.Data,
		TotalCount:  len(resp.Data),
		Description: "SystemInAllMenuInformationofTagReadyList, PackageIncludeRouteconfigurationAndComponentInformation",
	})
}
