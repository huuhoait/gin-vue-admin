package service

import (
	"strings"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model/request"
)

type tenantService struct{}

func (s *tenantService) Create(req request.CreateTenantReq) (model.Tenant, error) {
	row := model.Tenant{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     true,
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

func (s *tenantService) FindByID(id uint) (model.Tenant, error) {
	var t model.Tenant
	err := global.GVA_DB.Where("id = ?", id).First(&t).Error
	return t, err
}

func (s *tenantService) FindByCode(code string) (model.Tenant, error) {
	var t model.Tenant
	err := global.GVA_DB.Where("code = ?", code).First(&t).Error
	return t, err
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
	ErrTenantHasMembers = stringErr("cannot delete tenant with active members; unassign first")
)

type stringErr string

func (s stringErr) Error() string { return string(s) }
