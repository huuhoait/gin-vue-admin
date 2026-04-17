package system

import (
	"context"
	"errors"
	"strconv"

	"gorm.io/gorm"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	_ "github.com/go-sql-driver/mysql"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateCasbin
//@description: Update casbin permissions
//@param: authorityId string, casbinInfos []request.CasbinInfo
//@return: error

type CasbinService struct{}

var CasbinServiceApp = new(CasbinService)

func (casbinService *CasbinService) UpdateCasbin(adminAuthorityID, AuthorityID uint, casbinInfos []request.CasbinInfo) error {
	authorityIDStr := strconv.Itoa(int(AuthorityID))
	before := casbinService.GetPolicyPathByAuthorityId(AuthorityID)

	err := AuthorityServiceApp.CheckAuthorityIDAuth(adminAuthorityID, AuthorityID)
	if err != nil {
		return err
	}

	if global.GVA_CONFIG.System.UseStrictAuth {
		apis, e := ApiServiceApp.GetAllApis(adminAuthorityID)
		if e != nil {
			return e
		}

		for i := range casbinInfos {
			hasApi := false
			for j := range apis {
				if apis[j].Path == casbinInfos[i].Path && apis[j].Method == casbinInfos[i].Method {
					hasApi = true
					break
				}
			}
			if !hasApi {
				return errors.New("some APIs are not in the authorized list")
			}
		}
	}

	authorityId := authorityIDStr
	casbinService.ClearCasbin(0, authorityId)
	rules := [][]string{}
	// deduplicate permissions
	deduplicateMap := make(map[string]bool)
	for _, v := range casbinInfos {
		key := authorityId + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
	}
	if len(rules) == 0 {
		RecordPolicyChange(context.Background(), "update", authorityId, before, casbinInfos, "cleared all policies")
		return nil
	} // no need to call AddPolicies when setting empty permissions
	e := utils.GetCasbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("duplicate API exists, failed to add, please contact administrator")
	}
	RecordPolicyChange(context.Background(), "update", authorityId, before, casbinInfos, "")
	return nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateCasbinApi
//@description: Update casbin when API changes
//@param: oldPath string, newPath string, oldMethod string, newMethod string
//@return: error

func (casbinService *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := global.GVA_DB.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	if err != nil {
		return err
	}
	RecordPolicyChange(context.Background(), "update_api", "",
		map[string]string{"path": oldPath, "method": oldMethod},
		map[string]string{"path": newPath, "method": newMethod},
		"casbin api rename",
	)
	e := utils.GetCasbin()
	return e.LoadPolicy()
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetPolicyPathByAuthorityId
//@description: Get policy path list by authority id
//@param: authorityId string
//@return: pathMaps []request.CasbinInfo

func (casbinService *CasbinService) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []request.CasbinInfo) {
	e := utils.GetCasbin()
	authorityId := strconv.Itoa(int(AuthorityID))
	list, _ := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, request.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ClearCasbin
//@description: Clear matching permissions
//@param: v int, p ...string
//@return: bool

func (casbinService *CasbinService) ClearCasbin(v int, p ...string) bool {
	e := utils.GetCasbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: RemoveFilteredPolicy
//@description: Remove filtered policy using database method; requires calling FreshCasbin to take effect immediately
//@param: db *gorm.DB, authorityId string
//@return: error

func (casbinService *CasbinService) RemoveFilteredPolicy(db *gorm.DB, authorityId string) error {
	return db.Delete(&gormadapter.CasbinRule{}, "v0 = ?", authorityId).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SyncPolicy
//@description: Sync current database policy; requires calling FreshCasbin to take effect immediately
//@param: db *gorm.DB, authorityId string, rules [][]string
//@return: error

func (casbinService *CasbinService) SyncPolicy(db *gorm.DB, authorityId string, rules [][]string) error {
	err := casbinService.RemoveFilteredPolicy(db, authorityId)
	if err != nil {
		return err
	}
	return casbinService.AddPolicies(db, rules)
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: AddPolicies
//@description: Add matching permissions
//@param: v int, p ...string
//@return: bool

func (casbinService *CasbinService) AddPolicies(db *gorm.DB, rules [][]string) error {
	var casbinRules []gormadapter.CasbinRule
	for i := range rules {
		casbinRules = append(casbinRules, gormadapter.CasbinRule{
			Ptype: "p",
			V0:    rules[i][0],
			V1:    rules[i][1],
			V2:    rules[i][2],
		})
	}
	return db.Create(&casbinRules).Error
}

func (casbinService *CasbinService) FreshCasbin() (err error) {
	e := utils.GetCasbin()
	err = e.LoadPolicy()
	return err
}

// GetAuthoritiesByApi gets all role IDs that have the specified API permission
func (casbinService *CasbinService) GetAuthoritiesByApi(path, method string) (authorityIds []uint, err error) {
	var rules []gormadapter.CasbinRule
	err = global.GVA_DB.Where("ptype = 'p' AND v1 = ? AND v2 = ?", path, method).Find(&rules).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rules {
		id, e := strconv.Atoi(r.V0)
		if e == nil {
			authorityIds = append(authorityIds, uint(id))
		}
	}
	return authorityIds, nil
}

// SetApiAuthorities fully replaces the role list associated with an API
func (casbinService *CasbinService) SetApiAuthorities(path, method string, authorityIds []uint) error {
	before, _ := casbinService.GetAuthoritiesByApi(path, method)
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 1. delete all existing role associations for this API
		if err := tx.Where("ptype = 'p' AND v1 = ? AND v2 = ?", path, method).Delete(&gormadapter.CasbinRule{}).Error; err != nil {
			return err
		}
		// 2. batch insert new association records
		if len(authorityIds) > 0 {
			newRules := make([]gormadapter.CasbinRule, 0, len(authorityIds))
			for _, authorityId := range authorityIds {
				newRules = append(newRules, gormadapter.CasbinRule{
					Ptype: "p",
					V0:    strconv.Itoa(int(authorityId)),
					V1:    path,
					V2:    method,
				})
			}
			if err := tx.Create(&newRules).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err == nil {
		RecordPolicyChange(context.Background(), "set_api_authorities", "",
			map[string]any{"path": path, "method": method, "authorities": before},
			map[string]any{"path": path, "method": method, "authorities": authorityIds},
			"",
		)
	}
	return err
}
