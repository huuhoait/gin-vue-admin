package utils

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
)

var (
	ApiMap  = make(map[string][]system.SysApi)
	MenuMap = make(map[string][]system.SysBaseMenu)
	DictMap = make(map[string][]system.SysDictionary)
	rw      sync.Mutex
)

func getPluginName() string {
	_, file, _, ok := runtime.Caller(2)
	pluginName := ""
	if ok {
		file = filepath.ToSlash(file)
		const key = "server/plugin/"
		if idx := strings.Index(file, key); idx != -1 {
			remain := file[idx+len(key):]
			parts := strings.Split(remain, "/")
			if len(parts) > 0 {
				pluginName = parts[0]
			}
		}
	}
	return pluginName
}

func RegisterApis(apis ...system.SysApi) {
	name := getPluginName()
	if name != "" {
		rw.Lock()
		ApiMap[name] = apis
		rw.Unlock()
	}

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, api := range apis {
			err := tx.Model(system.SysApi{}).Where("path = ? AND method = ? AND api_group = ? ", api.Path, api.Method, api.ApiGroup).FirstOrCreate(&api).Error
			if err != nil {
				zap.L().Error("Failed to register API", zap.Error(err), zap.String("api", api.Path), zap.String("method", api.Method), zap.String("apiGroup", api.ApiGroup))
				return err
			}
		}
		return nil
	})
	if err != nil {
		zap.L().Error("Failed to register API", zap.Error(err))
	}
	ensureDefaultSuperAdminCasbin(apis)
}

// apiGroupsGrantedToDefaultSuperAdmin gates which plugin ApiGroups receive
// auto-Casbin grants on boot. Whitelisting matches the menu-grant pattern in
// menusGrantedToDefaultSuperAdmin — first-party plugins (announcement, gva)
// stay opt-out so this change does not silently widen their permission set.
var apiGroupsGrantedToDefaultSuperAdmin = map[string]struct{}{
	"OnlineUsers":      {},
	"SysMonitor":       {},
	"OAuth2Client":     {},
	"OAuth2":           {},
	"Tenant":           {},
	"TenantMembership": {},
	"TenantPackage":    {},
}

// ensureDefaultSuperAdminCasbin grants the default super-admin authority
// (888) Casbin access to a plugin's APIs. Idempotent — checks for an
// existing rule before inserting. Skips silently if the gorm-adapter table
// is unavailable.
func ensureDefaultSuperAdminCasbin(apis []system.SysApi) {
	if global.GVA_DB == nil || len(apis) == 0 {
		return
	}
	const superAdminAuthority = "888"
	for _, api := range apis {
		if _, ok := apiGroupsGrantedToDefaultSuperAdmin[api.ApiGroup]; !ok {
			continue
		}
		var n int64
		if err := global.GVA_DB.Model(&adapter.CasbinRule{}).
			Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
				"p", superAdminAuthority, api.Path, api.Method).
			Count(&n).Error; err != nil {
			continue
		}
		if n > 0 {
			continue
		}
		if err := global.GVA_DB.Create(&adapter.CasbinRule{
			Ptype: "p",
			V0:    superAdminAuthority,
			V1:    api.Path,
			V2:    api.Method,
		}).Error; err != nil {
			zap.L().Warn("link api to default super admin failed",
				zap.String("api", api.Path),
				zap.String("method", api.Method),
				zap.Error(err))
		}
	}
}

// resolveLegacyPluginParentMenus maps ParentId 9 (upstream gin-vue-admin
// convention) onto the real "plugin" menu row ID for this fork's menu seed.
func resolveLegacyPluginParentMenus(menus []system.SysBaseMenu) {
	if global.GVA_DB == nil {
		return
	}
	var pluginID uint
	if err := global.GVA_DB.Model(&system.SysBaseMenu{}).Select("id").Where("name = ?", "plugin").Limit(1).Scan(&pluginID).Error; err != nil || pluginID == 0 {
		return
	}
	for i := range menus {
		if menus[i].ParentId == 9 {
			menus[i].ParentId = pluginID
		}
	}
}

// menusGrantedToDefaultSuperAdmin are linked to authority_id 888 on boot so
// the sidebar shows them without running SQL seeds manually.
var menusGrantedToDefaultSuperAdmin = map[string]struct{}{
	"security": {}, "oauth2Clients": {}, "onlineUsers": {},
	"sysmonitor": {}, "tenants": {}, "tenantList": {}, "tenantMembers": {}, "tenantPackages": {},
}

func ensureDefaultSuperAdminMenus(menuNames []string) {
	if global.GVA_DB == nil || len(menuNames) == 0 {
		return
	}
	const superAdminAuthority uint = 888
	aid := strconv.FormatUint(uint64(superAdminAuthority), 10)
	for _, name := range menuNames {
		if _, ok := menusGrantedToDefaultSuperAdmin[name]; !ok {
			continue
		}
		var mid uint
		if err := global.GVA_DB.Model(&system.SysBaseMenu{}).Select("id").Where("name = ?", name).Limit(1).Scan(&mid).Error; err != nil || mid == 0 {
			continue
		}
		midStr := strconv.FormatUint(uint64(mid), 10)
		var n int64
		if err := global.GVA_DB.Model(&system.SysAuthorityMenu{}).
			Where("sys_authority_authority_id = ? AND sys_base_menu_id = ?", aid, midStr).
			Count(&n).Error; err != nil {
			continue
		}
		if n > 0 {
			continue
		}
		if err := global.GVA_DB.Create(&system.SysAuthorityMenu{
			AuthorityId: aid,
			MenuId:      midStr,
		}).Error; err != nil {
			zap.L().Warn("link menu to default super admin failed", zap.String("menu", name), zap.Error(err))
		}
	}
}

func RegisterMenus(menus ...system.SysBaseMenu) {
	if len(menus) == 0 {
		return
	}
	name := getPluginName()
	if name != "" {
		rw.Lock()
		MenuMap[name] = menus
		rw.Unlock()
	}

	resolveLegacyPluginParentMenus(menus)

	parentMenu := menus[0]
	otherMenus := menus[1:]
	// FirstOrCreate overwrites the in-memory struct with the existing DB row
	// (when one matches), so capture the seed's display fields up front to
	// re-apply them afterwards.
	parentSeedTitle := parentMenu.Meta.Title
	parentSeedIcon := parentMenu.Meta.Icon
	childSeedTitles := make([]string, len(otherMenus))
	childSeedIcons := make([]string, len(otherMenus))
	for i := range otherMenus {
		childSeedTitles[i] = otherMenus[i].Meta.Title
		childSeedIcons[i] = otherMenus[i].Meta.Icon
	}

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(system.SysBaseMenu{}).Where("name = ? ", parentMenu.Name).FirstOrCreate(&parentMenu).Error
		if err != nil {
			zap.L().Error("Failed to register menu", zap.Error(err))
			return errors.Wrap(err, "Failed to register menu")
		}
		// FirstOrCreate does not refresh arbitrary columns on an existing row;
		// always sync key presentation + routing columns so seed changes
		// (e.g. label i18n key, re-parenting, turning a leaf page into a
		// routerHolder parent) propagate to live rows. Seeds are the source of
		// truth for plugin-managed menus.
		if err = tx.Model(&system.SysBaseMenu{}).Where("id = ?", parentMenu.ID).
			Updates(map[string]any{
				"parent_id":        parentMenu.ParentId,
				"path":             parentMenu.Path,
				"hidden":           parentMenu.Hidden,
				"component":        parentMenu.Component,
				"sort":             parentMenu.Sort,
				"title":            parentSeedTitle,
				"icon":             parentSeedIcon,
				"active_name":      parentMenu.ActiveName,
				"keep_alive":       parentMenu.KeepAlive,
				"default_menu":     parentMenu.DefaultMenu,
				"close_tab":        parentMenu.CloseTab,
				"transition_type":  parentMenu.TransitionType,
			}).Error; err != nil {
			zap.L().Error("Failed to sync parent menu title/icon", zap.Error(err))
			return errors.Wrap(err, "Failed to sync parent menu title/icon")
		}
		pid := parentMenu.ID
		for i := range otherMenus {
			otherMenus[i].ParentId = pid
			err = tx.Model(system.SysBaseMenu{}).Where("name = ? ", otherMenus[i].Name).FirstOrCreate(&otherMenus[i]).Error
			if err != nil {
				zap.L().Error("Failed to register menu", zap.Error(err))
				return errors.Wrap(err, "Failed to register menu")
			}
			// FirstOrCreate does not refresh arbitrary columns on an existing row;
			// always sync parent_id + routing columns so reparenting and label
			// migrations converge.
			if err = tx.Model(&system.SysBaseMenu{}).Where("id = ?", otherMenus[i].ID).
				Updates(map[string]any{
					"parent_id": pid,
					"path":      otherMenus[i].Path,
					"hidden":    otherMenus[i].Hidden,
					"component": otherMenus[i].Component,
					"sort":      otherMenus[i].Sort,
					"title":     childSeedTitles[i],
					"icon":      childSeedIcons[i],
					"active_name":     otherMenus[i].ActiveName,
					"keep_alive":      otherMenus[i].KeepAlive,
					"default_menu":    otherMenus[i].DefaultMenu,
					"close_tab":       otherMenus[i].CloseTab,
					"transition_type": otherMenus[i].TransitionType,
				}).Error; err != nil {
				zap.L().Error("Failed to reparent menu", zap.Error(err))
				return errors.Wrap(err, "Failed to reparent menu")
			}
		}
		return nil
	})
	if err != nil {
		zap.L().Error("Failed to register menu", zap.Error(err))
		return
	}
	names := make([]string, 0, len(menus))
	for i := range menus {
		names = append(names, menus[i].Name)
	}
	ensureDefaultSuperAdminMenus(names)
}

func RegisterDictionaries(dictionaries ...system.SysDictionary) {
	name := getPluginName()
	if name != "" {
		rw.Lock()
		DictMap[name] = dictionaries
		rw.Unlock()
	}

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, dict := range dictionaries {
			details := dict.SysDictionaryDetails
			dict.SysDictionaryDetails = nil
			err := tx.Model(system.SysDictionary{}).Where("type = ?", dict.Type).FirstOrCreate(&dict).Error
			if err != nil {
				zap.L().Error("Failed to register dictionary", zap.Error(err), zap.String("type", dict.Type))
				return err
			}
			for _, detail := range details {
				detail.SysDictionaryID = int(dict.ID)
				err = tx.Model(system.SysDictionaryDetail{}).Where("sys_dictionary_id = ? AND value = ?", dict.ID, detail.Value).FirstOrCreate(&detail).Error
				if err != nil {
					zap.L().Error("Failed to register dictionary detail", zap.Error(err), zap.String("value", detail.Value))
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		zap.L().Error("Failed to register dictionary", zap.Error(err))
	}
}

func Pointer[T any](in T) *T {
	return &in
}

func GetPluginData(pluginName string) ([]system.SysApi, []system.SysBaseMenu, []system.SysDictionary) {
	rw.Lock()
	defer rw.Unlock()
	return ApiMap[pluginName], MenuMap[pluginName], DictMap[pluginName]
}

