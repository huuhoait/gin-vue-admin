package system

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MenuRouter struct{}

func (s *MenuRouter) InitMenuRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	menuRouter := Router.Group("menu").Use(middleware.OperationRecord())
	menuRouterWithoutRecord := Router.Group("menu")
	{
		menuRouter.POST("addBaseMenu", authorityMenuApi.AddBaseMenu)           // create menu
		menuRouter.POST("addMenuAuthority", authorityMenuApi.AddMenuAuthority) //	add menu-role association
		menuRouter.POST("deleteBaseMenu", authorityMenuApi.DeleteBaseMenu)     // delete menu
		menuRouter.POST("updateBaseMenu", authorityMenuApi.UpdateBaseMenu)     // update menu
		menuRouter.POST("setMenuRoles", authorityMenuApi.SetMenuRoles)         // full overwriteMenuAssociationRole
	}
	{
		menuRouterWithoutRecord.POST("getMenu", authorityMenuApi.GetMenu)                   // get menu tree
		menuRouterWithoutRecord.POST("getMenuList", authorityMenuApi.GetMenuList)           // PaginationgetBasicmenuList
		menuRouterWithoutRecord.POST("getBaseMenuTree", authorityMenuApi.GetBaseMenuTree)   // getuser dynamic routes
		menuRouterWithoutRecord.POST("getMenuAuthority", authorityMenuApi.GetMenuAuthority) // get menus by role
		menuRouterWithoutRecord.POST("getBaseMenuById", authorityMenuApi.GetBaseMenuById)   // According toidgetMenu
		menuRouterWithoutRecord.GET("getMenuRoles", authorityMenuApi.GetMenuRoles)          // getMenuAssociationrole IDList
	}
	return menuRouter
}
