package api

import (
	"errors"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model/request"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type tenantApi struct{}

// CreateTenant
// @Tags     Tenant
// @Summary  create tenant
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateTenantReq true "tenant fields (code, name, description, contact info, domain, expireAt, accountLimit)"
// @Router   /tenant/create [post]
func (a *tenantApi) CreateTenant(c *gin.Context) {
	var req request.CreateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, err := serviceTenant.Create(req)
	if err != nil {
		global.GVA_LOG.Error("create tenant failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(row, c)
}

// UpdateTenant
// @Tags     Tenant
// @Summary  update tenant
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdateTenantReq true "tenant fields (id required; other fields are partial-update)"
// @Router   /tenant/update [put]
func (a *tenantApi) UpdateTenant(c *gin.Context) {
	var req request.UpdateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, err := serviceTenant.Update(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(row, c)
}

// DeleteTenant
// @Tags     Tenant
// @Router   /tenant/delete [delete]
func (a *tenantApi) DeleteTenant(c *gin.Context) {
	var req request.IdReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceTenant.Delete(req.ID); err != nil {
		if errors.Is(err, service.ErrTenantHasMembers) {
			response.FailWithCode(c, "admin.plugin.tenant.has_members")
			return
		}
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithCode(c, "admin.common.delete_success")
}

// FindTenant
// @Tags     Tenant
// @Router   /tenant/find [get]
func (a *tenantApi) FindTenant(c *gin.Context) {
	var req request.IdReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, err := serviceTenant.FindByID(req.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(row, c)
}

// ListTenants
// @Tags     Tenant
// @Router   /tenant/list [get]
func (a *tenantApi) ListTenants(c *gin.Context) {
	global.GVA_LOG.Info("Current role", zap.Uint("authorityId", utils.GetUserAuthorityId(c)))
	var req request.TenantListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceTenant.List(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDataCode(c, response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "admin.common.get_success")
}
