package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityMenuApi struct{}

// GetMenu
// @Tags      AuthorityMenu
// @Summary   Get user dynamic routes
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.Empty                                                  true  "Empty"
// @Success   200   {object}  response.Response{data=systemRes.SysMenusResponse,msg=string}  "Get user dynamic routes, returns system menu detail list"
// @Router    /menu/getMenu [post]
func (a *AuthorityMenuApi) GetMenu(c *gin.Context) {
	menus, err := menuService.GetMenuTree(utils.GetUserAuthorityId(c))
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	if menus == nil {
		menus = []system.SysMenu{}
	}
	response.OkWithDetailed(systemRes.SysMenusResponse{Menus: menus}, "Retrieved successfully", c)
}

// GetBaseMenuTree
// @Tags      AuthorityMenu
// @Summary   Get base menu tree
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.Empty                                                      true  "Empty"
// @Success   200   {object}  response.Response{data=systemRes.SysBaseMenusResponse,msg=string}  "Get base menu tree, returns system menu list"
// @Router    /menu/getBaseMenuTree [post]
func (a *AuthorityMenuApi) GetBaseMenuTree(c *gin.Context) {
	authority := utils.GetUserAuthorityId(c)
	menus, err := menuService.GetBaseMenuTree(authority)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(systemRes.SysBaseMenusResponse{Menus: menus}, "Retrieved successfully", c)
}

// AddMenuAuthority
// @Tags      AuthorityMenu
// @Summary   Add menu-authority association
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.AddMenuAuthorityInfo  true  "Authority ID"
// @Success   200   {object}  response.Response{msg=string}   "Add menu-authority association"
// @Router    /menu/addMenuAuthority [post]
func (a *AuthorityMenuApi) AddMenuAuthority(c *gin.Context) {
	var authorityMenu systemReq.AddMenuAuthorityInfo
	err := c.ShouldBindJSON(&authorityMenu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(authorityMenu, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	adminAuthorityID := utils.GetUserAuthorityId(c)
	if err := menuService.AddMenuAuthority(authorityMenu.Menus, adminAuthorityID, authorityMenu.AuthorityId); err != nil {
		global.GVA_LOG.Error("Failed to add!", zap.Error(err))
		response.FailWithMessage("Failed to add", c)
	} else {
		response.OkWithMessage("Added successfully", c)
	}
}

// GetMenuAuthority
// @Tags      AuthorityMenu
// @Summary   Get menu for specified authority
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetAuthorityId                                     true  "Authority ID"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Get menu for specified authority"
// @Router    /menu/getMenuAuthority [post]
func (a *AuthorityMenuApi) GetMenuAuthority(c *gin.Context) {
	var param request.GetAuthorityId
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	menus, err := menuService.GetMenuAuthority(&param)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysMenusResponse{Menus: menus}, "Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"menus": menus}, "Retrieved successfully", c)
}

// AddBaseMenu
// @Tags      Menu
// @Summary   Add base menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysBaseMenu             true  "Route path, parent menu ID, route name, frontend file path, sort order"
// @Success   200   {object}  response.Response{msg=string}  "Add base menu"
// @Router    /menu/addBaseMenu [post]
func (a *AuthorityMenuApi) AddBaseMenu(c *gin.Context) {
	var menu system.SysBaseMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu, utils.MenuVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu.Meta, utils.MenuMetaVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = menuService.AddBaseMenu(menu)
	if err != nil {
		global.GVA_LOG.Error("Failed to add!", zap.Error(err))
		response.FailWithMessage("Failed to add: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Added successfully", c)
}

// DeleteBaseMenu
// @Tags      Menu
// @Summary   Delete menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "Menu ID"
// @Success   200   {object}  response.Response{msg=string}  "Delete menu"
// @Router    /menu/deleteBaseMenu [post]
func (a *AuthorityMenuApi) DeleteBaseMenu(c *gin.Context) {
	var menu request.GetById
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = baseMenuService.DeleteBaseMenu(menu.ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// UpdateBaseMenu
// @Tags      Menu
// @Summary   Update menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysBaseMenu             true  "Route path, parent menu ID, route name, frontend file path, sort order"
// @Success   200   {object}  response.Response{msg=string}  "Update menu"
// @Router    /menu/updateBaseMenu [post]
func (a *AuthorityMenuApi) UpdateBaseMenu(c *gin.Context) {
	var menu system.SysBaseMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu, utils.MenuVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(menu.Meta, utils.MenuMetaVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = baseMenuService.UpdateBaseMenu(menu)
	if err != nil {
		global.GVA_LOG.Error("Failed to update!", zap.Error(err))
		response.FailWithMessage("Update failed", c)
		return
	}
	response.OkWithMessage("Updated successfully", c)
}

// GetBaseMenuById
// @Tags      Menu
// @Summary   Get menu by ID
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                                                   true  "Menu ID"
// @Success   200   {object}  response.Response{data=systemRes.SysBaseMenuResponse,msg=string}  "Get menu by ID, returns system menu list"
// @Router    /menu/getBaseMenuById [post]
func (a *AuthorityMenuApi) GetBaseMenuById(c *gin.Context) {
	var idInfo request.GetById
	err := c.ShouldBindJSON(&idInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	menu, err := baseMenuService.GetBaseMenuById(idInfo.ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(systemRes.SysBaseMenuResponse{Menu: menu}, "Retrieved successfully", c)
}

// GetMenuRoles
// @Tags      AuthorityMenu
// @Summary   Get role IDs that have access to specified menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     menuId  query     uint                                                         true  "Menu ID"
// @Success   200     {object}  response.Response{data=map[string]interface{},msg=string}    "Retrieved successfully"
// @Router    /menu/getMenuRoles [get]
func (a *AuthorityMenuApi) GetMenuRoles(c *gin.Context) {
	var req systemReq.SetMenuAuthorities
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.MenuId == 0 {
		response.FailWithMessage("Menu ID cannot be empty", c)
		return
	}
	authorityIds, err := menuService.GetAuthoritiesByMenuId(req.MenuId)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve"+err.Error(), c)
		return
	}
	if authorityIds == nil {
		authorityIds = []uint{}
	}
	defaultRouterAuthorityIds, err := menuService.GetDefaultRouterAuthorityIds(req.MenuId)
	if err != nil {
		global.GVA_LOG.Error("Failed to get default router authority IDs!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve"+err.Error(), c)
		return
	}
	if defaultRouterAuthorityIds == nil {
		defaultRouterAuthorityIds = []uint{}
	}
	response.OkWithDetailed(gin.H{
		"authorityIds":              authorityIds,
		"defaultRouterAuthorityIds": defaultRouterAuthorityIds,
	}, "Retrieved successfully", c)
}

// SetMenuRoles
// @Tags      AuthorityMenu
// @Summary   Fully replace the role list associated with a menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetMenuAuthorities   true  "Menu ID and authority ID list"
// @Success   200   {object}  response.Response{msg=string}  "Set successfully"
// @Router    /menu/setMenuRoles [post]
func (a *AuthorityMenuApi) SetMenuRoles(c *gin.Context) {
	var req systemReq.SetMenuAuthorities
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.MenuId == 0 {
		response.FailWithMessage("Menu ID cannot be empty", c)
		return
	}
	if err := menuService.SetMenuAuthorities(req.MenuId, req.AuthorityIds); err != nil {
		global.GVA_LOG.Error("Failed to set!", zap.Error(err))
		response.FailWithMessage("Failed to set"+err.Error(), c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}

// GetMenuList
// @Tags      Menu
// @Summary   Get base menu list with pagination
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "Page number, page size"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Get paginated base menu list, returns list, total, page, page size"
// @Router    /menu/getMenuList [post]
func (a *AuthorityMenuApi) GetMenuList(c *gin.Context) {
	authorityID := utils.GetUserAuthorityId(c)
	menuList, err := menuService.GetInfoList(authorityID)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(menuList, "Retrieved successfully", c)
}
