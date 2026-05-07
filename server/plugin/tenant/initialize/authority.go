package initialize

import (
	"context"
	"fmt"
	"sync"

	"github.com/huuhoait/gin-vue-admin/server/global"
	systemModel "github.com/huuhoait/gin-vue-admin/server/model/system"
	tenantservice "github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"go.uber.org/zap"
)

// tenantAuthorityMenus lists the menus every tenant user gets access to.
// `dashboard` is required to satisfy the FE's defaultRouter check; the
// `tenants` parent + `tenantMembers` submenu let primary tenant admins manage
// their own membership list. Non-primary tenant users hit the runtime
// ensureCanManageMembers gate when they click into management actions.
var tenantAuthorityMenus = []string{"dashboard", "tenants", "tenantMembers"}

// tenantAuthorityCasbinPolicies lists the (path, method) pairs the Tenant
// authority needs to operate the member-management UI. Bootstrap reads
// (/menu/getMenu, /user/getUserInfo, /tenant/mine, /jwt/jsonInBlacklist) are
// already in middleware.isCasbinBypassPath so they're omitted here.
var tenantAuthorityCasbinPolicies = [][2]string{
	{"/tenantMembership/members", "GET"},
	{"/tenantMembership/assign", "POST"},
	{"/tenantMembership/unassign", "DELETE"},
	{"/tenantMembership/createUser", "POST"},
	{"/user/getUserList", "POST"},
	{"/user/setUserInfo", "PUT"},
	{"/user/resetPassword", "POST"},
}

// authorityOnce guards the bootstrap path so plugin reloads (tests, hot
// reload) only re-run the seed work once per process. Idempotent SQL keeps
// it correct even on the first run after a partial failure.
var authorityOnce sync.Once

// Authority seeds the Tenant authority row, links its baseline menus, and
// installs its Casbin policies — exactly once per process. Hoisted out of
// the request path so CreateUserAndAssign no longer pays the 10+ idempotent
// SELECTs per user create.
//
// Must run after initialize.Menu so the tenant-plugin menus this function
// references are already present in sys_base_menus. The dashboard menu is
// seeded by the core system bootstrap, which runs before plugin Register.
func Authority(_ context.Context) {
	authorityOnce.Do(func() {
		if err := provisionTenantAuthority(); err != nil {
			// Don't fail the whole plugin boot — log and continue. CreateUser
			// flows will surface the specific failure (e.g. missing dashboard
			// row) when they hit it.
			zap.L().Error("tenant: failed to provision Tenant authority", zap.Error(err))
		}
	})
}

func provisionTenantAuthority() error {
	if global.GVA_DB == nil {
		return nil
	}

	// 1. Ensure the authority row itself exists.
	authority := systemModel.SysAuthority{
		AuthorityId:   tenantservice.TenantAuthorityID,
		AuthorityName: "Tenant",
		ParentId:      utils.Pointer[uint](0),
		DefaultRouter: "dashboard",
	}
	if err := global.GVA_DB.Where("authority_id = ?", tenantservice.TenantAuthorityID).
		FirstOrCreate(&authority).Error; err != nil {
		return err
	}

	aid := fmt.Sprintf("%d", tenantservice.TenantAuthorityID)

	// 2. Link baseline menus. Skip names that haven't been seeded yet — the
	// next process boot will fill them in once the missing menu lands.
	for _, name := range tenantAuthorityMenus {
		var menuID uint
		if err := global.GVA_DB.Model(&systemModel.SysBaseMenu{}).
			Select("id").Where("name = ?", name).
			Limit(1).Scan(&menuID).Error; err != nil {
			return err
		}
		if menuID == 0 {
			continue
		}
		mid := fmt.Sprintf("%d", menuID)
		var n int64
		if err := global.GVA_DB.Model(&systemModel.SysAuthorityMenu{}).
			Where("sys_authority_authority_id = ? AND sys_base_menu_id = ?", aid, mid).
			Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			continue
		}
		if err := global.GVA_DB.Create(&systemModel.SysAuthorityMenu{
			AuthorityId: aid,
			MenuId:      mid,
		}).Error; err != nil {
			return err
		}
	}

	// 3. Install Casbin policies for the member-management surface.
	for _, p := range tenantAuthorityCasbinPolicies {
		var n int64
		if err := global.GVA_DB.Table("casbin_rule").
			Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", "p", aid, p[0], p[1]).
			Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			continue
		}
		if err := global.GVA_DB.Exec(
			"INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)",
			"p", aid, p[0], p[1],
		).Error; err != nil {
			return err
		}
	}
	return nil
}
