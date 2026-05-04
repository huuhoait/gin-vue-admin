package api

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type packageApi struct{}

// CreatePackage
// @Tags     TenantPackage
// @Summary  create tenant package
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreatePackageReq true "package fields"
// @Success  200 {object} response.Response{msg=string}
// @Router   /tenantPackage/create [post]
func (a *packageApi) CreatePackage(c *gin.Context) {
	var req request.CreatePackageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, err := serviceTenantPackage.Create(req)
	if err != nil {
		global.GVA_LOG.Error("create tenant package failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDataCode(c, row, "admin.common.create_success")
}

// UpdatePackage
// @Tags     TenantPackage
// @Summary  update tenant package
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdatePackageReq true "package fields"
// @Success  200 {object} response.Response{msg=string}
// @Router   /tenantPackage/update [put]
func (a *packageApi) UpdatePackage(c *gin.Context) {
	var req request.UpdatePackageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, err := serviceTenantPackage.Update(req)
	if err != nil {
		global.GVA_LOG.Error("update tenant package failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDataCode(c, row, "admin.common.update_success")
}

// DeletePackage
// @Tags     TenantPackage
// @Summary  delete tenant package
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    id query uint true "package id"
// @Success  200 {object} response.Response{msg=string}
// @Router   /tenantPackage/delete [delete]
func (a *packageApi) DeletePackage(c *gin.Context) {
	var req request.IdReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceTenantPackage.Delete(req.ID); err != nil {
		global.GVA_LOG.Error("delete tenant package failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithCode(c, "admin.common.delete_success")
}

// FindPackage
// @Tags     TenantPackage
// @Summary  find tenant package by id
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    id query uint true "package id"
// @Success  200 {object} response.Response{data=model.TenantPackage,msg=string}
// @Router   /tenantPackage/find [get]
func (a *packageApi) FindPackage(c *gin.Context) {
	var req request.IdReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, err := serviceTenantPackage.FindByID(req.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDataCode(c, row, "admin.common.get_success")
}

// ListPackages
// @Tags     TenantPackage
// @Summary  list tenant packages
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    data query request.PackageListReq true "list filters"
// @Success  200 {object} response.Response{data=response.PageResult,msg=string}
// @Router   /tenantPackage/list [get]
func (a *packageApi) ListPackages(c *gin.Context) {
	var req request.PackageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceTenantPackage.List(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDataCode(c, response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "admin.common.get_success")
}
