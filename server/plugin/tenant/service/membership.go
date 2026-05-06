package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/huuhoait/gin-vue-admin/server/global"
	systemModel "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model/request"
	systemService "github.com/huuhoait/gin-vue-admin/server/service/system"
	"github.com/huuhoait/gin-vue-admin/server/utils"
)

// TenantAuthorityID is the SysAuthority row used for users created via the
// "create user in tenant" flow. Idempotently seeded on first use — no schema
// migration required for existing DBs. Tenant-scoped users get no system
// permissions by default; admins grant via the Authority management UI.
const TenantAuthorityID uint = 9300

// userTenantTargetID encodes a UserTenant composite key for the data-change
// log's targetID column ("<userID>:<tenantID>"). Keeps the audit trail
// queryable per (user, tenant) pair without needing a JSON join.
func userTenantTargetID(userID, tenantID uint) string {
	return fmt.Sprintf("%d:%d", userID, tenantID)
}

type membershipService struct{}

// IsPrimaryMember reports whether userID is the primary member of tenantID.
func (s *membershipService) IsPrimaryMember(userID, tenantID uint) bool {
	if userID == 0 || tenantID == 0 {
		return false
	}
	var n int64
	if err := global.GVA_DB.Model(&model.UserTenant{}).
		Where("user_id = ? AND tenant_id = ? AND is_primary = ?", userID, tenantID, true).
		Count(&n).Error; err != nil {
		return false
	}
	return n > 0
}

// Assign adds a user to a tenant. When isPrimary=true, demotes any other
// primary entry for this user atomically — a user has at most one primary
// tenant at a time.
//
// Enforces the tenant's AccountLimit when present (>0). The check is done
// inside the transaction so concurrent assigns cannot both squeeze past a
// stale read; if the membership row already exists for this (user,tenant)
// pair the limit is NOT re-checked (re-assigning is idempotent and does not
// grow the count).
func (s *membershipService) Assign(ctx context.Context, userID, tenantID uint, isPrimary bool) error {
	// before/after snapshots populated inside the TX so we can emit one audit
	// row after a successful commit. nil before means it was a net-new add.
	var (
		hadBefore     bool
		beforePrimary bool
	)
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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
		if alreadyMember {
			hadBefore = true
			beforePrimary = existing.IsPrimary
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
	if err != nil {
		return err
	}

	// Audit trail — recorded after commit so a failed mutation never leaves a
	// dangling log row. Action distinguishes net-new vs idempotent re-assign
	// so reviewers can spot accidental duplicates.
	action := "assign"
	if hadBefore {
		action = "reassign"
	}
	var before any
	if hadBefore {
		before = map[string]any{
			"userID":    userID,
			"tenantID":  tenantID,
			"isPrimary": beforePrimary,
		}
	}
	after := map[string]any{
		"userID":    userID,
		"tenantID":  tenantID,
		"isPrimary": isPrimary,
	}
	systemService.RecordDataChange(ctx, "UserTenant", userTenantTargetID(userID, tenantID), action, before, after)
	return nil
}

// CreateUserAndAssign provisions a fresh SysUser with the default Tenant
// authority and adds them to the target tenant in a single transaction. The
// AccountLimit cap is enforced inside the TX to avoid concurrent over-fills,
// and a unique-username collision short-circuits before any side effect.
//
// Audit emits one row keyed on the new user's id ("create_user_in_tenant"
// action) so reviewers see the (creator → user → tenant) link in the data
// change log without needing a follow-up assign row.
func (s *membershipService) CreateUserAndAssign(ctx context.Context, req request.CreateUserAndAssignReq) (systemModel.SysUser, error) {
	// Ensure the default Tenant authority row exists. FirstOrCreate is
	// idempotent and works on existing DBs that pre-date this feature, so we
	// don't need a separate migration step.
	authority := systemModel.SysAuthority{
		AuthorityId:   TenantAuthorityID,
		AuthorityName: "Tenant",
		ParentId:      utils.Pointer[uint](0),
		DefaultRouter: "dashboard",
	}
	if err := global.GVA_DB.Where("authority_id = ?", TenantAuthorityID).
		FirstOrCreate(&authority).Error; err != nil {
		return systemModel.SysUser{}, err
	}

	hashedPwd := utils.BcryptHash(req.Password)
	user := systemModel.SysUser{
		UUID:        uuid.New(),
		Username:    req.Username,
		Password:    hashedPwd,
		NickName:    req.NickName,
		Phone:       req.Phone,
		Email:       req.Email,
		AuthorityId: TenantAuthorityID,
		Authorities: []systemModel.SysAuthority{{AuthorityId: TenantAuthorityID}},
		Enable:      1,
	}

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var tenant model.Tenant
		if err := tx.Where("id = ?", req.TenantID).First(&tenant).Error; err != nil {
			return err
		}

		// Username uniqueness — checked inside the TX so two concurrent
		// creates can't both win the race.
		var dupe systemModel.SysUser
		err := tx.Select("id").Where("username = ?", req.Username).First(&dupe).Error
		if err == nil {
			return errors.New("username already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Cap enforcement applies before user creation so we don't leave an
		// orphaned SysUser row behind on rollback.
		if tenant.AccountLimit > 0 {
			var count int64
			if err := tx.Model(&model.UserTenant{}).
				Where("tenant_id = ?", req.TenantID).Count(&count).Error; err != nil {
				return err
			}
			if count >= int64(tenant.AccountLimit) {
				return ErrAccountLimitReached
			}
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if req.IsPrimary {
			if err := tx.Model(&model.UserTenant{}).
				Where("user_id = ? AND is_primary = ?", user.ID, true).
				Update("is_primary", false).Error; err != nil {
				return err
			}
		}
		membership := model.UserTenant{
			UserID:    user.ID,
			TenantID:  req.TenantID,
			IsPrimary: req.IsPrimary,
		}
		return tx.Create(&membership).Error
	})
	if err != nil {
		return systemModel.SysUser{}, err
	}

	// Audit trail — Password is intentionally scrubbed by RecordDataChange's
	// JSON sensitive-field regex, but we also blank it here as defense-in-depth.
	auditUser := user
	auditUser.Password = ""
	systemService.RecordDataChange(ctx, "SysUser", fmt.Sprintf("%d", user.ID), "create_user_in_tenant",
		nil,
		map[string]any{
			"user":     auditUser,
			"tenantID": req.TenantID,
			"isPrimary": req.IsPrimary,
		},
	)
	return user, nil
}

func (s *membershipService) Unassign(ctx context.Context, userID, tenantID uint) error {
	// Snapshot the row before deletion so the audit trail can show what was
	// removed. Missing row → no-op delete and no log entry; the caller's
	// effect is identical either way.
	var before model.UserTenant
	hadBefore := global.GVA_DB.Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		First(&before).Error == nil

	if err := global.GVA_DB.Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Delete(&model.UserTenant{}).Error; err != nil {
		return err
	}
	if hadBefore {
		systemService.RecordDataChange(ctx, "UserTenant", userTenantTargetID(userID, tenantID), "unassign", before, nil)
	}
	return nil
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
