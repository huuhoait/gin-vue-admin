package mcpTool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/mark3labs/mcp-go/mcp"
)

func init() {
	RegisterTool(&MenuCreator{})
}

type MenuCreateRequest struct {
	ParentId    uint                   `json:"parentId"`
	Path        string                 `json:"path"`
	Name        string                 `json:"name"`
	Hidden      bool                   `json:"hidden"`
	Component   string                 `json:"component"`
	Sort        int                    `json:"sort"`
	Title       string                 `json:"title"`
	Icon        string                 `json:"icon"`
	KeepAlive   bool                   `json:"keepAlive"`
	DefaultMenu bool                   `json:"defaultMenu"`
	CloseTab    bool                   `json:"closeTab"`
	ActiveName  string                 `json:"activeName"`
	Parameters  []MenuParameterRequest `json:"parameters"`
	MenuBtn     []MenuButtonRequest    `json:"menuBtn"`
}

type MenuParameterRequest struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MenuButtonRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type MenuCreateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	MenuID  uint   `json:"menuId"`
	Name    string `json:"name"`
	Path    string `json:"path"`
}

type MenuCreator struct{}

func (m *MenuCreator) New() mcp.Tool {
	return mcp.NewTool("create_menu",
		mcp.WithDescription(`CreateFrontendMenurecord used forAIeditDeviceAutomaticAddFrontendPageWhenauto-createCorrespondingMenuItem. 

**Important limits:**
- when usinggva_auto_generate tool and needCreatedModules=truewhen true, module creation auto-generatesMenuItem, NotShouldInvokeThisUtility
- only use when:1) create standaloneMenu(not involving module creation);2) AIeditDeviceAutomaticAddFrontendPageWhen`),
		mcp.WithNumber("parentId",
			mcp.Description("Parentmenu ID, 0TableShowRootMenu"),
			mcp.DefaultNumber(0),
		),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("route path, Such As:userList"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("route name, Used forVue Router, Such As:userList"),
		),
		mcp.WithBoolean("hidden",
			mcp.Description("YesNoAtMenuListInhidden"),
		),
		mcp.WithString("component",
			mcp.Required(),
			mcp.Description("CorrespondingFrontendVueComponentpath, Such As:view/user/list.vue"),
		),
		mcp.WithNumber("sort",
			mcp.Description("Menusortorder number; smaller sorts first"),
			mcp.DefaultNumber(1),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("MenuShowTitle"),
		),
		mcp.WithString("icon",
			mcp.Description("MenuIconname"),
			mcp.DefaultString("menu"),
		),
		mcp.WithBoolean("keepAlive",
			mcp.Description("YesNoKeep-alive"),
		),
		mcp.WithBoolean("defaultMenu",
			mcp.Description("YesNoYesBasicRoute"),
		),
		mcp.WithBoolean("closeTab",
			mcp.Description("YesNoAutomaticDisabletab"),
		),
		mcp.WithString("activeName",
			mcp.Description("HighBrightmenu titleName"),
		),
		mcp.WithString("parameters",
			mcp.Description("Route ParamsJSONstring, format: [{\"type\":\"params\",\"key\":\"id\",\"value\":\"1\"}]"),
		),
		mcp.WithString("menuBtn",
			mcp.Description("menu buttonsJSONstring, format: [{\"name\":\"add\",\"desc\":\"create\"}]"),
		),
	)
}

func (m *MenuCreator) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	path, ok := args["path"].(string)
	if !ok || path == "" {
		return nil, errors.New("path parameter is required")
	}
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, errors.New("name ParameterYesrequired")
	}
	component, ok := args["component"].(string)
	if !ok || component == "" {
		return nil, errors.New("component ParameterYesrequired")
	}
	title, ok := args["title"].(string)
	if !ok || title == "" {
		return nil, errors.New("title ParameterYesrequired")
	}

	parentID := uint(0)
	if value, ok := args["parentId"].(float64); ok {
		parentID = uint(value)
	}
	hidden, _ := args["hidden"].(bool)
	sort := 1
	if value, ok := args["sort"].(float64); ok {
		sort = int(value)
	}
	icon := "menu"
	if value, ok := args["icon"].(string); ok && value != "" {
		icon = value
	}
	keepAlive, _ := args["keepAlive"].(bool)
	defaultMenu, _ := args["defaultMenu"].(bool)
	closeTab, _ := args["closeTab"].(bool)
	activeName, _ := args["activeName"].(string)

	parameters := make([]system.SysBaseMenuParameter, 0)
	if parametersStr, ok := args["parameters"].(string); ok && parametersStr != "" {
		var paramReqs []MenuParameterRequest
		if err := json.Unmarshal([]byte(parametersStr), &paramReqs); err != nil {
			return nil, fmt.Errorf("parameters invalid parameter format: %v", err)
		}
		for _, param := range paramReqs {
			parameters = append(parameters, system.SysBaseMenuParameter{
				Type:  param.Type,
				Key:   param.Key,
				Value: param.Value,
			})
		}
	}

	menuBtns := make([]system.SysBaseMenuBtn, 0)
	if menuBtnStr, ok := args["menuBtn"].(string); ok && menuBtnStr != "" {
		var buttonReqs []MenuButtonRequest
		if err := json.Unmarshal([]byte(menuBtnStr), &buttonReqs); err != nil {
			return nil, fmt.Errorf("menuBtn invalid parameter format: %v", err)
		}
		for _, button := range buttonReqs {
			menuBtns = append(menuBtns, system.SysBaseMenuBtn{
				Name: button.Name,
				Desc: button.Desc,
			})
		}
	}

	menu := system.SysBaseMenu{
		ParentId:  parentID,
		Path:      path,
		Name:      name,
		Hidden:    hidden,
		Component: component,
		Sort:      sort,
		Meta: system.Meta{
			Title:       title,
			Icon:        icon,
			KeepAlive:   keepAlive,
			DefaultMenu: defaultMenu,
			CloseTab:    closeTab,
			ActiveName:  activeName,
		},
		Parameters: parameters,
		MenuBtn:    menuBtns,
	}

	if _, err := postUpstream[map[string]any](ctx, "/menu/addBaseMenu", menu); err != nil {
		return nil, fmt.Errorf("CreateMenufailed: %v", err)
	}

	menuID := uint(0)
	if menuListResp, err := postUpstream[[]system.SysBaseMenu](ctx, "/menu/getMenuList", map[string]any{}); err == nil {
		menuID = findMenuID(menuListResp.Data, name, path)
	}

	return textResultWithJSON("Menucreation result:", &MenuCreateResponse{
		Success: true,
		Message: fmt.Sprintf("succeededCreateMenu %s", title),
		MenuID:  menuID,
		Name:    name,
		Path:    path,
	})
}

func findMenuID(menus []system.SysBaseMenu, name, path string) uint {
	for _, menu := range menus {
		if menu.Name == name && menu.Path == path {
			return menu.ID
		}
		if id := findMenuID(menu.Children, name, path); id != 0 {
			return id
		}
	}
	return 0
}
