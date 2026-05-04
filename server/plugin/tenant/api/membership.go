package api

import (
	"errors"

	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model/request"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"

	"github.com/gin-gonic/gin"
)

type membershipApi struct{}

// AssignUser
// @Tags     TenantMembership
// @Summary  assign a user to a tenant
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.AssignUserReq true "userID, tenantID, isPrimary"
// @Router   /tenantMembership/assign [post]
func (a *membershipApi) AssignUser(c *gin.Context) {
	var req request.AssignUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceMembership.Assign(req.UserID, req.TenantID, req.IsPrimary); err != nil {
		switch {
		case errors.Is(err, service.ErrAccountLimitReached):
			response.FailWithCode(c, "admin.plugin.tenant.account_limit_reached")
		case errors.Is(err, service.ErrTenantDisabled):
			response.FailWithCode(c, "admin.plugin.tenant.disabled")
		case errors.Is(err, service.ErrTenantExpired):
			response.FailWithCode(c, "admin.plugin.tenant.expired")
		default:
			response.FailWithMessage(err.Error(), c)
		}
		return
	}
	response.OkWithCode(c, "admin.plugin.tenant.assigned")
}

// UnassignUser
// @Tags     TenantMembership
// @Router   /tenantMembership/unassign [delete]
func (a *membershipApi) UnassignUser(c *gin.Context) {
	var req request.UnassignUserReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceMembership.Unassign(req.UserID, req.TenantID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithCode(c, "admin.plugin.tenant.removed")
}

// MembersOfTenant
// @Tags     TenantMembership
// @Router   /tenantMembership/members [get]
func (a *membershipApi) MembersOfTenant(c *gin.Context) {
	var req request.IdReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := serviceMembership.MembersOfTenant(req.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}
