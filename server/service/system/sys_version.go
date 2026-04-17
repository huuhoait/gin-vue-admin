package system

import (
	"context"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"gorm.io/gorm"
)

type SysVersionService struct{}

// CreateSysVersion creates a version management record
// Author [yourname](https://github.com/yourname)
func (sysVersionService *SysVersionService) CreateSysVersion(ctx context.Context, sysVersion *system.SysVersion) (err error) {
	err = global.GVA_DB.Create(sysVersion).Error
	return err
}

// DeleteSysVersion deletes a version management record
// Author [yourname](https://github.com/yourname)
func (sysVersionService *SysVersionService) DeleteSysVersion(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&system.SysVersion{}, "id = ?", ID).Error
	return err
}

// DeleteSysVersionByIds batch deletes version management records
// Author [yourname](https://github.com/yourname)
func (sysVersionService *SysVersionService) DeleteSysVersionByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Where("id in ?", IDs).Delete(&system.SysVersion{}).Error
	return err
}

// GetSysVersion retrieves a version management record by ID
// Author [yourname](https://github.com/yourname)
func (sysVersionService *SysVersionService) GetSysVersion(ctx context.Context, ID string) (sysVersion system.SysVersion, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sysVersion).Error
	return
}

// GetSysVersionInfoList retrieves paginated version management records
// Author [yourname](https://github.com/yourname)
func (sysVersionService *SysVersionService) GetSysVersionInfoList(ctx context.Context, info systemReq.SysVersionSearch) (list []system.SysVersion, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// create db
	db := global.GVA_DB.Model(&system.SysVersion{})
	var sysVersions []system.SysVersion
	// if search conditions exist, search queries will be automatically built below
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.VersionName != nil && *info.VersionName != "" {
		db = db.Where("version_name LIKE ?", "%"+*info.VersionName+"%")
	}
	if info.VersionCode != nil && *info.VersionCode != "" {
		db = db.Where("version_code = ?", *info.VersionCode)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&sysVersions).Error
	return sysVersions, total, err
}
func (sysVersionService *SysVersionService) GetSysVersionPublic(ctx context.Context) {
	// this method retrieves data defined by the data source
	// please implement it yourself
}

// GetMenusByIds retrieves menu data by ID list
func (sysVersionService *SysVersionService) GetMenusByIds(ctx context.Context, ids []uint) (menus []system.SysBaseMenu, err error) {
	err = global.GVA_DB.Where("id in ?", ids).Preload("Parameters").Preload("MenuBtn").Find(&menus).Error
	return
}

// GetApisByIds retrieves API data by ID list
func (sysVersionService *SysVersionService) GetApisByIds(ctx context.Context, ids []uint) (apis []system.SysApi, err error) {
	err = global.GVA_DB.Where("id in ?", ids).Find(&apis).Error
	return
}

// GetDictionariesByIds retrieves dictionary data by ID list
func (sysVersionService *SysVersionService) GetDictionariesByIds(ctx context.Context, ids []uint) (dictionaries []system.SysDictionary, err error) {
	err = global.GVA_DB.Where("id in ?", ids).Preload("SysDictionaryDetails").Find(&dictionaries).Error
	return
}

// ImportMenus imports menu data
func (sysVersionService *SysVersionService) ImportMenus(ctx context.Context, menus []system.SysBaseMenu) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// recursively create menus
		return sysVersionService.createMenusRecursively(tx, menus, 0)
	})
}

// createMenusRecursively recursively creates menus
func (sysVersionService *SysVersionService) createMenusRecursively(tx *gorm.DB, menus []system.SysBaseMenu, parentId uint) error {
	for _, menu := range menus {
		// check if menu already exists
		var existingMenu system.SysBaseMenu
		if err := tx.Where("name = ? AND path = ?", menu.Name, menu.Path).First(&existingMenu).Error; err == nil {
			// menu already exists, continue processing child menus with the existing menu ID
			if len(menu.Children) > 0 {
				if err := sysVersionService.createMenusRecursively(tx, menu.Children, existingMenu.ID); err != nil {
					return err
				}
			}
			continue
		}

		// save parameter and button data for later processing
		parameters := menu.Parameters
		menuBtns := menu.MenuBtn
		children := menu.Children

		// create new menu (without associated data)
		newMenu := system.SysBaseMenu{
			ParentId:  parentId,
			Path:      menu.Path,
			Name:      menu.Name,
			Hidden:    menu.Hidden,
			Component: menu.Component,
			Sort:      menu.Sort,
			Meta:      menu.Meta,
		}

		if err := tx.Create(&newMenu).Error; err != nil {
			return err
		}

		// create parameters
		if len(parameters) > 0 {
			for _, param := range parameters {
				newParam := system.SysBaseMenuParameter{
					SysBaseMenuID: newMenu.ID,
					Type:          param.Type,
					Key:           param.Key,
					Value:         param.Value,
				}
				if err := tx.Create(&newParam).Error; err != nil {
					return err
				}
			}
		}

		// create menu buttons
		if len(menuBtns) > 0 {
			for _, btn := range menuBtns {
				newBtn := system.SysBaseMenuBtn{
					SysBaseMenuID: newMenu.ID,
					Name:          btn.Name,
					Desc:          btn.Desc,
				}
				if err := tx.Create(&newBtn).Error; err != nil {
					return err
				}
			}
		}

		// recursively process child menus
		if len(children) > 0 {
			if err := sysVersionService.createMenusRecursively(tx, children, newMenu.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// ImportApis imports API data
func (sysVersionService *SysVersionService) ImportApis(apis []system.SysApi) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, api := range apis {
			// check if API already exists
			var existingApi system.SysApi
			if err := tx.Where("path = ? AND method = ?", api.Path, api.Method).First(&existingApi).Error; err == nil {
				// API already exists, skip
				continue
			}

			// create new API
			newApi := system.SysApi{
				Path:        api.Path,
				Description: api.Description,
				ApiGroup:    api.ApiGroup,
				Method:      api.Method,
			}

			if err := tx.Create(&newApi).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ImportDictionaries imports dictionary data
func (sysVersionService *SysVersionService) ImportDictionaries(dictionaries []system.SysDictionary) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, dict := range dictionaries {
			// check if dictionary already exists
			var existingDict system.SysDictionary
			if err := tx.Where("type = ?", dict.Type).First(&existingDict).Error; err == nil {
				// dictionary already exists, skip
				continue
			}

			// create new dictionary
			newDict := system.SysDictionary{
				Name:                 dict.Name,
				Type:                 dict.Type,
				Status:               dict.Status,
				Desc:                 dict.Desc,
				SysDictionaryDetails: dict.SysDictionaryDetails,
			}

			if err := tx.Create(&newDict).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
