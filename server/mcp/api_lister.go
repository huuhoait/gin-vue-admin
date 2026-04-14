package mcpTool

import (
	"context"

	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
)

func init() {
	RegisterTool(&ApiLister{})
}

type ApiInfo struct {
	ID          uint   `json:"id,omitempty"`
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
	ApiGroup    string `json:"apiGroup,omitempty"`
	Method      string `json:"method"`
	Source      string `json:"source"`
}

type ApiListResponse struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message"`
	DatabaseApis []ApiInfo `json:"databaseApis"`
	GinApis      []ApiInfo `json:"ginApis"`
	TotalCount   int       `json:"totalCount"`
}

type mcpRoutesResponse struct {
	Routes gin.RoutesInfo `json:"routes"`
}

type ApiLister struct{}

func (a *ApiLister) New() mcp.Tool {
	return mcp.NewTool("list_all_apis",
		mcp.WithDescription(`getSystemInAllofAPIAPI, PartForTwoGroup:

**Description:**
- ReturnDatabaseInAlreadyRegisterofAPIList
- ReturnginFrameworkInActualRegisterofRouteAPIList
- HelphelpFrontendJudgeYesUseAppearHaveAPIStillYesNeedCreateNewofAPI,IfapiAtFrontendNotUseAndNeedFrontendInvokeofWhenTime, PleaseToapiFileclipDownCorrespondingModuleofjsInAddmethodAndexposeToCurrentBusinessInvoke

**ReturnDataStructure:**
- databaseApis: DatabaseInAPIRecord(PackageIncludeID, description, Groupetc.CompleteArrangeInformation)
- ginApis: ginRouteInAPI(OnlyPackageIncludepathAndmethod), NeedAIAccording topathSelfRowguesspathofBusinessIncludemeaning, e.g. /api/user/:id TableShowAccording toUser IDget user info`),
		mcp.WithString("_placeholder",
			mcp.Description("placeholder to avoid json schema validation failure"),
		),
	)
}

func (a *ApiLister) Handle(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiResp, err := postUpstream[systemRes.SysAPIListResponse](ctx, "/api/getAllApis", map[string]any{})
	if err != nil {
		return nil, err
	}

	routeResp, err := postUpstream[mcpRoutesResponse](ctx, "/autoCode/mcpRoutes", map[string]any{})
	if err != nil {
		return nil, err
	}

	databaseApis := make([]ApiInfo, 0, len(apiResp.Data.Apis))
	for _, api := range apiResp.Data.Apis {
		databaseApis = append(databaseApis, ApiInfo{
			ID:          api.ID,
			Path:        api.Path,
			Description: api.Description,
			ApiGroup:    api.ApiGroup,
			Method:      api.Method,
			Source:      "database",
		})
	}

	ginApis := make([]ApiInfo, 0, len(routeResp.Data.Routes))
	for _, route := range routeResp.Data.Routes {
		ginApis = append(ginApis, ApiInfo{
			Path:   route.Path,
			Method: route.Method,
			Source: "gin",
		})
	}

	return textResultWithJSON("", ApiListResponse{
		Success:      true,
		Message:      "getAPIListsucceeded",
		DatabaseApis: databaseApis,
		GinApis:      ginApis,
		TotalCount:   len(databaseApis) + len(ginApis),
	})
}
