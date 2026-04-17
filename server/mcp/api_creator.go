package mcpTool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	commonReq "github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/mark3labs/mcp-go/mcp"
)

func init() {
	RegisterTool(&ApiCreator{})
}

type ApiCreateRequest struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	ApiGroup    string `json:"apiGroup"`
	Method      string `json:"method"`
}

type ApiCreateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ApiID   uint   `json:"apiId"`
	Path    string `json:"path"`
	Method  string `json:"method"`
}

type ApiCreator struct{}

func (a *ApiCreator) New() mcp.Tool {
	return mcp.NewTool("create_api",
		mcp.WithDescription(`CreateAfterEndAPIrecord used forAIeditDeviceAutomaticAddAPIAPIWhenauto-createCorrespondingAPIPermissionRecord. 

**Important limits:**
- when usinggva_auto_generate tool and needCreatedModules=truewhen true, module creation auto-generatesAPIPermission, NotShouldInvokeThisUtility
- only use when:1) create standaloneAPI(not involving module creation);2) AIeditDeviceAutomaticAddAPI;3) routerDownofFileGeneratepathVariationWhen`),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("APIpath, Such As:/user/create"),
		),
		mcp.WithString("description",
			mcp.Required(),
			mcp.Description("APIInTextdescription, Such As:CreateUser"),
		),
		mcp.WithString("apiGroup",
			mcp.Required(),
			mcp.Description("APIGroupname, Used forPartClassmanagement, Such As:User Management"),
		),
		mcp.WithString("method",
			mcp.Description("HTTPmethod"),
			mcp.DefaultString("POST"),
		),
		mcp.WithString("apis",
			mcp.Description("BatchCreateAPIofJSONstring, format: [{\"path\":\"/user/create\",\"description\":\"CreateUser\",\"apiGroup\":\"User Management\",\"method\":\"POST\"}]"),
		),
	)
}

func (a *ApiCreator) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	var apis []ApiCreateRequest
	if apisStr, ok := args["apis"].(string); ok && apisStr != "" {
		if err := json.Unmarshal([]byte(apisStr), &apis); err != nil {
			return nil, fmt.Errorf("apis invalid parameter format: %w", err)
		}
	} else {
		path, ok := args["path"].(string)
		if !ok || path == "" {
			return nil, errors.New("path parameter is required")
		}
		description, ok := args["description"].(string)
		if !ok || description == "" {
			return nil, errors.New("description ParameterYesrequired")
		}
		apiGroup, ok := args["apiGroup"].(string)
		if !ok || apiGroup == "" {
			return nil, errors.New("apiGroup ParameterYesrequired")
		}

		method := "POST"
		if value, ok := args["method"].(string); ok && value != "" {
			method = value
		}

		apis = append(apis, ApiCreateRequest{
			Path:        path,
			Description: description,
			ApiGroup:    apiGroup,
			Method:      method,
		})
	}

	if len(apis) == 0 {
		return nil, errors.New("NoneRequireCreateofAPI")
	}

	responses := make([]ApiCreateResponse, 0, len(apis))
	successCount := 0

	for _, apiReq := range apis {
		_, err := postUpstream[map[string]any](ctx, "/api/createApi", system.SysApi{
			Path:        apiReq.Path,
			Description: apiReq.Description,
			ApiGroup:    apiReq.ApiGroup,
			Method:      apiReq.Method,
		})
		if err != nil {
			responses = append(responses, ApiCreateResponse{
				Success: false,
				Message: fmt.Sprintf("CreateAPIfailed: %v", err),
				Path:    apiReq.Path,
				Method:  apiReq.Method,
			})
			continue
		}

		lookupResp, lookupErr := postUpstream[pageResultData[[]system.SysApi]](ctx, "/api/getApiList", systemReq.SearchApiParams{
			SysApi: system.SysApi{
				Path:   apiReq.Path,
				Method: apiReq.Method,
			},
			PageInfo: commonReq.PageInfo{
				Page:     1,
				PageSize: 1,
			},
		})

		var apiID uint
		if lookupErr == nil && len(lookupResp.Data.List) > 0 {
			apiID = lookupResp.Data.List[0].ID
		}

		responses = append(responses, ApiCreateResponse{
			Success: true,
			Message: fmt.Sprintf("succeededCreateAPI %s %s", apiReq.Method, apiReq.Path),
			ApiID:   apiID,
			Path:    apiReq.Path,
			Method:  apiReq.Method,
		})
		successCount++
	}

	result := map[string]any{
		"success":      successCount > 0,
		"totalCount":   len(apis),
		"successCount": successCount,
		"failedCount":  len(apis) - successCount,
		"details":      responses,
	}

	return textResultWithJSON("APIcreation result:", result)
}
