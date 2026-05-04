package service

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model"

	"gorm.io/gorm"
)

type membershipService struct{}

// Assign adds a user to a tenant. When isPrimary=true, demotes any other
// primary entry for this user atomically — a user has at most one primary
// tenant at a time.
func (s *membershipService) Assign(userID, tenantID uint, isPrimary bool) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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

func (s *membershipService) MembersOfTenant(tenantID uint) ([]model.UserTenant, error) {
	var list []model.UserTenant
	err := global.GVA_DB.Where("tenant_id = ?", tenantID).Find(&list).Error
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
