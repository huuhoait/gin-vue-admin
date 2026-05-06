package service

import (
	"errors"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model"

	"gorm.io/gorm"
)

type membershipService struct{}

// Assign adds a user to a tenant. When isPrimary=true, demotes any other
// primary entry for this user atomically — a user has at most one primary
// tenant at a time.
//
// Enforces the tenant's AccountLimit when present (>0). The check is done
// inside the transaction so concurrent assigns cannot both squeeze past a
// stale read; if the membership row already exists for this (user,tenant)
// pair the limit is NOT re-checked (re-assigning is idempotent and does not
// grow the count).
func (s *membershipService) Assign(userID, tenantID uint, isPrimary bool) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// Existence check first — we need the tenant's AccountLimit anyway,
		// so loading it here doubles as a fast-fail when the tenant id is
		// bogus and lets us surface a meaningful error instead of GORM's
		// foreign-key-style noise downstream.
		var tenant model.Tenant
		if err := tx.Where("id = ?", tenantID).First(&tenant).Error; err != nil {
			return err
		}

		// Detect "already a member" — re-assignment is idempotent w.r.t. the
		// account-limit cap. We still need to handle is_primary toggling
		// further down.
		var existing model.UserTenant
		err := tx.Where("user_id = ? AND tenant_id = ?", userID, tenantID).First(&existing).Error
		alreadyMember := err == nil
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Cap enforcement applies only to net-new assignments and only when
		// the tenant has explicitly set a non-zero AccountLimit.
		if !alreadyMember && tenant.AccountLimit > 0 {
			var count int64
			if err := tx.Model(&model.UserTenant{}).
				Where("tenant_id = ?", tenantID).Count(&count).Error; err != nil {
				return err
			}
			if count >= int64(tenant.AccountLimit) {
				return ErrAccountLimitReached
			}
		}

		if isPrimary {
			if err := tx.Model(&model.UserTenant{}).
				Where("user_id = ? AND is_primary = ?", userID, true).
				Update("is_primary", false).Error; err != nil {
				return err
			}
		}
		row := model.UserTenant{UserID: userID, TenantID: tenantID, IsPrimary: isPrimary}
		// FirstOrCreate keeps Assign idempotent when the row already exists;
		// then if the caller wants to flip is_primary we update explicitly.
		if err := tx.Where("user_id = ? AND tenant_id = ?", userID, tenantID).
			FirstOrCreate(&row).Error; err != nil {
			return err
		}
		if row.IsPrimary != isPrimary {
			return tx.Model(&row).Update("is_primary", isPrimary).Error
		}
		return nil
	})
}

func (s *membershipService) Unassign(userID, tenantID uint) error {
	return global.GVA_DB.Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Delete(&model.UserTenant{}).Error
}

// MembershipsForUser lists every tenant the user can act on. Used by the
// tenant middleware when no X-Tenant-ID is supplied.
func (s *membershipService) MembershipsForUser(userID uint) ([]model.UserTenant, error) {
	var list []model.UserTenant
	err := global.GVA_DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

// PrimaryTenantForUser returns the user's primary tenant id; (0,false) if
// the user has no membership rows at all.
func (s *membershipService) PrimaryTenantForUser(userID uint) (uint, bool) {
	var row model.UserTenant
	err := global.GVA_DB.Where("user_id = ? AND is_primary = ?", userID, true).First(&row).Error
	if err == nil {
		return row.TenantID, true
	}
	// Fallback: any membership.
	err = global.GVA_DB.Where("user_id = ?", userID).First(&row).Error
	if err != nil {
		return 0, false
	}
	return row.TenantID, true
}

// HasAccess returns true when the user is assigned to the given tenant.
func (s *membershipService) HasAccess(userID, tenantID uint) bool {
	var n int64
	if err := global.GVA_DB.Model(&model.UserTenant{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).Count(&n).Error; err != nil {
		return false
	}
	return n > 0
}

// MemberWithUser is the row shape returned by MembersOfTenant — the membership
// fields plus the joined sys_users columns the FE needs to render a row
// without a second round-trip per user.
type MemberWithUser struct {
	UserID    uint      `json:"userID"`
	Username  string    `json:"username"`
	NickName  string    `json:"nickName"`
	IsPrimary bool      `json:"isPrimary"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *membershipService) MembersOfTenant(tenantID uint) ([]MemberWithUser, error) {
	var list []MemberWithUser
	// JOIN sys_users so the admin UI can show username/nickname instead of a
	// raw numeric id. LEFT JOIN guards against orphaned membership rows whose
	// user record was deleted out from under us — those still surface with
	// blank name fields rather than vanishing silently.
	err := global.GVA_DB.Table("gva_user_tenants AS ut").
		Select("ut.user_id AS user_id, ut.is_primary AS is_primary, ut.created_at AS created_at, u.username AS username, u.nick_name AS nick_name").
		Joins("LEFT JOIN sys_users u ON u.id = ut.user_id").
		Where("ut.tenant_id = ?", tenantID).
		Order("ut.is_primary DESC, ut.created_at ASC").
		Scan(&list).Error
	return list, err
}

// TenantWithMembership pairs a Tenant row with the user's per-membership flag
// so the frontend tenant switcher can render the dropdown without a second
// round-trip to discover which entry is the user's primary tenant.
type TenantWithMembership struct {
	model.Tenant
	IsPrimary bool `json:"isPrimary"`
}

// MyTenantsForUser returns every enabled tenant the user has membership in,
// ordered with the primary tenant first and then alphabetically by name. The
// IsPrimary flag is preserved per row so the FE can highlight the default
// selection in the switcher.
func (s *membershipService) MyTenantsForUser(userID uint) ([]TenantWithMembership, error) {
	var list []TenantWithMembership
	if userID == 0 {
		return list, nil
	}
	// SELECT t.*, ut.is_primary FROM gva_tenants t
	//   JOIN gva_user_tenants ut ON t.id = ut.tenant_id
	//   WHERE ut.user_id = ? AND t.enabled = TRUE
	//   ORDER BY ut.is_primary DESC, t.name ASC
	err := global.GVA_DB.Table("gva_tenants AS t").
		Select("t.*, ut.is_primary AS is_primary").
		Joins("JOIN gva_user_tenants ut ON ut.tenant_id = t.id").
		Where("ut.user_id = ? AND t.enabled = ?", userID, true).
		Order("ut.is_primary DESC, t.name ASC").
		Scan(&list).Error
	return list, err
}
