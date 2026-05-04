package service

import (
	"encoding/json"
	"strings"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model/request"

	"gorm.io/datatypes"
)

type packageService struct{}

// Create persists a new TenantPackage. MenuIDs/ApiIDs are normalised to
// non-nil JSON arrays so consumers always see "[]" instead of NULL.
func (s *packageService) Create(req request.CreatePackageReq) (model.TenantPackage, error) {
	row := model.TenantPackage{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		MenuIDs:     marshalUintSlice(req.MenuIDs),
		ApiIDs:      marshalUintSlice(req.ApiIDs),
		Enabled:     true,
	}
	err := global.GVA_DB.Create(&row).Error
	return row, err
}

// Update applies a partial change set. Code is immutable; nil slice fields
// are left untouched, empty slices clear the relation list.
func (s *packageService) Update(req request.UpdatePackageReq) (model.TenantPackage, error) {
	var row model.TenantPackage
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
	if req.MenuIDs != nil {
		updates["menu_ids"] = marshalUintSlice(*req.MenuIDs)
	}
	if req.ApiIDs != nil {
		updates["api_ids"] = marshalUintSlice(*req.ApiIDs)
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

// Delete removes a package by id. Tenants pointing at this package via
// PackageCode keep their stale reference — the loose FK contract makes that
// the caller's problem to clean up.
func (s *packageService) Delete(id uint) error {
	return global.GVA_DB.Delete(&model.TenantPackage{}, "id = ?", id).Error
}

func (s *packageService) FindByID(id uint) (model.TenantPackage, error) {
	var p model.TenantPackage
	err := global.GVA_DB.Where("id = ?", id).First(&p).Error
	return p, err
}

func (s *packageService) FindByCode(code string) (model.TenantPackage, error) {
	var p model.TenantPackage
	err := global.GVA_DB.Where("code = ?", code).First(&p).Error
	return p, err
}

func (s *packageService) List(req request.PackageListReq) ([]model.TenantPackage, int64, error) {
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	db := global.GVA_DB.Model(&model.TenantPackage{})
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
	var list []model.TenantPackage
	err := db.Limit(req.PageSize).Offset(req.PageSize * (req.Page - 1)).Order("id DESC").Find(&list).Error
	return list, total, err
}

// marshalUintSlice converts a []uint to a non-nil datatypes.JSON. A nil
// input produces "[]" so DB rows never contain NULL for these columns.
func marshalUintSlice(in []uint) datatypes.JSON {
	if in == nil {
		in = []uint{}
	}
	b, _ := json.Marshal(in)
	return datatypes.JSON(b)
}
