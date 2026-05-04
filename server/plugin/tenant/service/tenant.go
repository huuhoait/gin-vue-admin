package service

import (
	"strings"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model/request"
)

type tenantService struct{}

func (s *tenantService) Create(req request.CreateTenantReq) (model.Tenant, error) {
	row := model.Tenant{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Domain:       req.Domain,
		ExpireAt:     req.ExpireAt,
		AccountLimit: req.AccountLimit,
		Enabled:      true,
	}
	err := global.GVA_DB.Create(&row).Error
	return row, err
}

func (s *tenantService) Update(req request.UpdateTenantReq) (model.Tenant, error) {
	var row model.Tenant
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&row).Error; err != nil {
		return row, err
	}
	updates := map[string]any{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ContactName != "" {
		updates["contact_name"] = req.ContactName
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.Domain != "" {
		updates["domain"] = req.Domain
	}
	// ExpireAt has tri-state semantics:
	//   - ClearExpireAt=true → set NULL (never expires).
	//   - ExpireAt non-nil   → set to that timestamp.
	//   - else               → leave unchanged.
	switch {
	case req.ClearExpireAt:
		updates["expire_at"] = nil
	case req.ExpireAt != nil:
		updates["expire_at"] = *req.ExpireAt
	}
	if req.AccountLimit != nil {
		updates["account_limit"] = *req.AccountLimit
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if len(updates) == 0 {
		return row, nil
	}
	if err := global.GVA_DB.Model(&row).Updates(updates).Error; err != nil {
		return row, err
	}
	return row, global.GVA_DB.Where("id = ?", req.ID).First(&row).Error
}

// Delete refuses to drop a tenant that still has user assignments to avoid
// orphaning rows. Caller must unassign users first.
func (s *tenantService) Delete(id uint) error {
	var count int64
	if err := global.GVA_DB.Model(&model.UserTenant{}).Where("tenant_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrTenantHasMembers
	}
	return global.GVA_DB.Delete(&model.Tenant{}, "id = ?", id).Error
}

// FindByID loads a tenant by primary key. The result is returned regardless
// of Enabled / ExpireAt — callers must use IsActive to decide whether to
// allow runtime use.
func (s *tenantService) FindByID(id uint) (model.Tenant, error) {
	var t model.Tenant
	err := global.GVA_DB.Where("id = ?", id).First(&t).Error
	return t, err
}

// FindByCode loads a tenant by code. Same lifecycle policy as FindByID.
func (s *tenantService) FindByCode(code string) (model.Tenant, error) {
	var t model.Tenant
	err := global.GVA_DB.Where("code = ?", code).First(&t).Error
	return t, err
}

// IsActive reports whether the tenant is currently usable: it must be
// Enabled and (if an expiration is set) not yet expired.
//
// Callers (login flow, tenant middleware, membership.Assign) typically
// translate the false case into ErrTenantDisabled / ErrTenantExpired by
// inspecting the tenant directly; this helper is the single source of truth
// for "can the tenant be used right now".
func (s *tenantService) IsActive(t model.Tenant) bool {
	if !t.Enabled {
		return false
	}
	if t.ExpireAt != nil && !t.ExpireAt.After(time.Now()) {
		return false
	}
	return true
}

func (s *tenantService) List(req request.TenantListReq) ([]model.Tenant, int64, error) {
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	db := global.GVA_DB.Model(&model.Tenant{})
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		db = db.Where("code LIKE ? OR name LIKE ?", "%"+kw+"%", "%"+kw+"%")
	}
	if req.Enabled != nil {
		db = db.Where("enabled = ?", *req.Enabled)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Tenant
	err := db.Limit(req.PageSize).Offset(req.PageSize * (req.Page - 1)).Order("id DESC").Find(&list).Error
	return list, total, err
}

// Sentinel errors so callers can branch without parsing message text.
var (
	ErrTenantHasMembers    = stringErr("cannot delete tenant with active members; unassign first")
	ErrTenantDisabled      = stringErr("tenant is disabled")
	ErrTenantExpired       = stringErr("tenant has expired")
	ErrAccountLimitReached = stringErr("tenant has reached its account limit")
)

type stringErr string

func (s stringErr) Error() string { return string(s) }
