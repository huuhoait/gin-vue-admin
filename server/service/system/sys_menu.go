package system

import (
	"errors"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
	"strconv"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getMenuTreeMap
//@description: get the full route tree map
//@param: authorityId string
//@return: treeMap map[string][]system.SysMenu, err error

type MenuService struct{}

var MenuServiceApp = new(MenuService)

func (menuService *MenuService) getMenuTreeMap(authorityId uint) (treeMap map[uint][]system.SysMenu, err error) {
	var allMenus []system.SysMenu
	var baseMenu []system.SysBaseMenu
	var btns []system.SysAuthorityBtn
	treeMap = make(map[uint][]system.SysMenu)

	var SysAuthorityMenus []system.SysAuthorityMenu
	err = global.GVA_DB.Where("sys_authority_authority_id = ?", authorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}

	var MenuIds []string

	for i := range SysAuthorityMenus {
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}

	err = global.GVA_DB.Where("id in (?)", MenuIds).Order("sort").Preload("Parameters").Find(&baseMenu).Error
	if err != nil {
		return
	}

	for i := range baseMenu {
		allMenus = append(allMenus, system.SysMenu{
			SysBaseMenu: baseMenu[i],
			AuthorityId: authorityId,
			MenuId:      baseMenu[i].ID,
			Parameters:  baseMenu[i].Parameters,
		})
	}

	err = global.GVA_DB.Where("authority_id = ?", authorityId).Preload("SysBaseMenuBtn").Find(&btns).Error
	if err != nil {
		return
	}
	var btnMap = make(map[uint]map[string]uint)
	for _, v := range btns {
		if btnMap[v.SysMenuID] == nil {
			btnMap[v.SysMenuID] = make(map[string]uint)
		}
		btnMap[v.SysMenuID][v.SysBaseMenuBtn.Name] = authorityId
	}
	for _, v := range allMenus {
		v.Btns = btnMap[v.SysBaseMenu.ID]
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetMenuTree
//@description: get dynamic menu tree
//@param: authorityId string
//@return: menus []system.SysMenu, err error

func (menuService *MenuService) GetMenuTree(authorityId uint) (menus []system.SysMenu, err error) {
	menuTree, err := menuService.getMenuTreeMap(authorityId)
	menus = menuTree[0]
	for i := 0; i < len(menus); i++ {
		err = menuService.getChildrenList(&menus[i], menuTree)
	}
	return menus, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getChildrenList
//@description: get child menus
//@param: menu *model.SysMenu, treeMap map[string][]model.SysMenu
//@return: err error

func (menuService *MenuService) getChildrenList(menu *system.SysMenu, treeMap map[uint][]system.SysMenu) (err error) {
	menu.Children = treeMap[menu.MenuId]
	for i := 0; i < len(menu.Children); i++ {
		err = menuService.getChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetInfoList
//@description: get paginated route list
//@return: list interface{}, total int64,err error

func (menuService *MenuService) GetInfoList(authorityID uint) (list interface{}, err error) {
	var menuList []system.SysBaseMenu
	treeMap, err := menuService.getBaseMenuTreeMap(authorityID)
	menuList = treeMap[0]
	for i := 0; i < len(menuList); i++ {
		err = menuService.getBaseChildrenList(&menuList[i], treeMap)
	}
	return menuList, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getBaseChildrenList
//@description: get child menus of a menu
//@param: menu *model.SysBaseMenu, treeMap map[string][]model.SysBaseMenu
//@return: err error

func (menuService *MenuService) getBaseChildrenList(menu *system.SysBaseMenu, treeMap map[uint][]system.SysBaseMenu) (err error) {
	menu.Children = treeMap[menu.ID]
	for i := 0; i < len(menu.Children); i++ {
		err = menuService.getBaseChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: AddBaseMenu
//@description: add a base route
//@param: menu model.SysBaseMenu
//@return: error

func (menuService *MenuService) AddBaseMenu(menu system.SysBaseMenu) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// check if name is duplicated
		if !errors.Is(tx.Where("name = ?", menu.Name).First(&system.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("duplicate name exists, please modify the name")
		}

		if menu.ParentId != 0 {
			// check if parent menu exists
			var parentMenu system.SysBaseMenu
			if err := tx.First(&parentMenu, menu.ParentId).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("parent menu does not exist")
				}
				return err
			}

			// check the number of existing child menus under the parent menu
			var existingChildrenCount int64
			err := tx.Model(&system.SysBaseMenu{}).Where("parent_id = ?", menu.ParentId).Count(&existingChildrenCount).Error
			if err != nil {
				return err
			}

			// if the parent menu was a leaf menu (no children) and is now becoming a branch menu, clear its permission assignments
			if existingChildrenCount == 0 {
				// check if the parent menu is set as the home page by other roles
				var defaultRouterCount int64
				err := tx.Model(&system.SysAuthority{}).Where("default_router = ?", parentMenu.Name).Count(&defaultRouterCount).Error
				if err != nil {
					return err
				}
				if defaultRouterCount > 0 {
					return errors.New("parent menu is already used as the home page by other roles, please release the home page permission first")
				}

				// clear all permission assignments of the parent menu
				err = tx.Where("sys_base_menu_id = ?", menu.ParentId).Delete(&system.SysAuthorityMenu{}).Error
				if err != nil {
					return err
				}
			}
		}

		// create menu
		return tx.Create(&menu).Error
	})
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getBaseMenuTreeMap
//@description: get the full route tree map
//@return: treeMap map[string][]system.SysBaseMenu, err error

func (menuService *MenuService) getBaseMenuTreeMap(authorityID uint) (treeMap map[uint][]system.SysBaseMenu, err error) {
	parentAuthorityID, err := AuthorityServiceApp.GetParentAuthorityID(authorityID)
	if err != nil {
		return nil, err
	}

	var allMenus []system.SysBaseMenu
	treeMap = make(map[uint][]system.SysBaseMenu)
	db := global.GVA_DB.Order("sort").Preload("MenuBtn").Preload("Parameters")

	// when strict tree-based roles are enabled and parent authority is not 0, menu filtering is required
	if global.GVA_CONFIG.System.UseStrictAuth && parentAuthorityID != 0 {
		var authorityMenus []system.SysAuthorityMenu
		err = global.GVA_DB.Where("sys_authority_authority_id = ?", authorityID).Find(&authorityMenus).Error
		if err != nil {
			return nil, err
		}
		var menuIds []string
		for i := range authorityMenus {
			menuIds = append(menuIds, authorityMenus[i].MenuId)
		}
		db = db.Where("id in (?)", menuIds)
	}

	err = db.Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetBaseMenuTree
//@description: get the base route tree
//@return: menus []system.SysBaseMenu, err error

func (menuService *MenuService) GetBaseMenuTree(authorityID uint) (menus []system.SysBaseMenu, err error) {
	treeMap, err := menuService.getBaseMenuTreeMap(authorityID)
	menus = treeMap[0]
	for i := 0; i < len(menus); i++ {
		err = menuService.getBaseChildrenList(&menus[i], treeMap)
	}
	return menus, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: AddMenuAuthority
//@description: add menu tree for a role
//@param: menus []model.SysBaseMenu, authorityId string
//@return: err error

func (menuService *MenuService) AddMenuAuthority(menus []system.SysBaseMenu, adminAuthorityID, authorityId uint) (err error) {
	var auth system.SysAuthority
	auth.AuthorityId = authorityId
	auth.SysBaseMenus = menus

	err = AuthorityServiceApp.CheckAuthorityIDAuth(adminAuthorityID, authorityId)
	if err != nil {
		return err
	}

	var authority system.SysAuthority
	_ = global.GVA_DB.First(&authority, "authority_id = ?", adminAuthorityID).Error
	var menuIds []string

	// when strict tree-based roles are enabled and parent authority is not 0, menu filtering is required
	if global.GVA_CONFIG.System.UseStrictAuth && *authority.ParentId != 0 {
		var authorityMenus []system.SysAuthorityMenu
		err = global.GVA_DB.Where("sys_authority_authority_id = ?", adminAuthorityID).Find(&authorityMenus).Error
		if err != nil {
			return err
		}
		for i := range authorityMenus {
			menuIds = append(menuIds, authorityMenus[i].MenuId)
		}

		for i := range menus {
			hasMenu := false
			for j := range menuIds {
				idStr := strconv.Itoa(int(menus[i].ID))
				if idStr == menuIds[j] {
					hasMenu = true
				}
			}
			if !hasMenu {
				return errors.New("add failed, cross-level operation is not allowed")
			}
		}
	}

	err = AuthorityServiceApp.SetMenuAuthority(&auth)
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetMenuAuthority
//@description: view the current role's menu tree
//@param: info *request.GetAuthorityId
//@return: menus []system.SysMenu, err error

func (menuService *MenuService) GetMenuAuthority(info *request.GetAuthorityId) (menus []system.SysMenu, err error) {
	var baseMenu []system.SysBaseMenu
	var SysAuthorityMenus []system.SysAuthorityMenu
	err = global.GVA_DB.Where("sys_authority_authority_id = ?", info.AuthorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}

	var MenuIds []string

	for i := range SysAuthorityMenus {
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}

	err = global.GVA_DB.Where("id in (?) ", MenuIds).Order("sort").Find(&baseMenu).Error

	for i := range baseMenu {
		menus = append(menus, system.SysMenu{
			SysBaseMenu: baseMenu[i],
			AuthorityId: info.AuthorityId,
			MenuId:      baseMenu[i].ID,
			Parameters:  baseMenu[i].Parameters,
		})
	}
	return menus, err
}

// GetAuthoritiesByMenuId retrieves all authority IDs that have the specified menu
func (menuService *MenuService) GetAuthoritiesByMenuId(menuId uint) (authorityIds []uint, err error) {
	var records []system.SysAuthorityMenu
	err = global.GVA_DB.Where("sys_base_menu_id = ?", menuId).Find(&records).Error
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		id, e := strconv.Atoi(r.AuthorityId)
		if e == nil {
			authorityIds = append(authorityIds, uint(id))
		}
	}
	return authorityIds, nil
}

// GetDefaultRouterAuthorityIds retrieves the list of authority IDs that set the specified menu as the home page
func (menuService *MenuService) GetDefaultRouterAuthorityIds(menuId uint) (authorityIds []uint, err error) {
	var menu system.SysBaseMenu
	err = global.GVA_DB.First(&menu, menuId).Error
	if err != nil {
		return nil, err
	}
	var authorities []system.SysAuthority
	err = global.GVA_DB.Where("default_router = ?", menu.Name).Find(&authorities).Error
	if err != nil {
		return nil, err
	}
	for _, auth := range authorities {
		authorityIds = append(authorityIds, auth.AuthorityId)
	}
	return authorityIds, nil
}

// SetMenuAuthorities fully replaces the role list associated with a menu
func (menuService *MenuService) SetMenuAuthorities(menuId uint, authorityIds []uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 1. delete all existing role associations for this menu
		if err := tx.Where("sys_base_menu_id = ?", menuId).Delete(&system.SysAuthorityMenu{}).Error; err != nil {
			return err
		}
		// 2. batch insert new association records
		if len(authorityIds) > 0 {
			menuIdStr := strconv.Itoa(int(menuId))
			newRecords := make([]system.SysAuthorityMenu, 0, len(authorityIds))
			for _, authorityId := range authorityIds {
				newRecords = append(newRecords, system.SysAuthorityMenu{
					MenuId:      menuIdStr,
					AuthorityId: strconv.Itoa(int(authorityId)),
				})
			}
			if err := tx.Create(&newRecords).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UserAuthorityDefaultRouter checks the default router for a user's authority
//
//	Author [SliverHorn](https://github.com/SliverHorn)
func (menuService *MenuService) UserAuthorityDefaultRouter(user *system.SysUser) {
	var menuIds []string
	err := global.GVA_DB.Model(&system.SysAuthorityMenu{}).Where("sys_authority_authority_id = ?", user.AuthorityId).Pluck("sys_base_menu_id", &menuIds).Error
	if err != nil {
		return
	}
	var am system.SysBaseMenu
	err = global.GVA_DB.First(&am, "name = ? and id in (?)", user.Authority.DefaultRouter, menuIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user.Authority.DefaultRouter = "404"
	}
}
