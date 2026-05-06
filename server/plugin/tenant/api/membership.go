package api

import (
	"errors"

	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model/request"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"
	systemService "github.com/huuhoait/gin-vue-admin/server/service/system"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
)

type membershipApi struct{}

// primary member of the current tenant can manage its members; super-admin
// flows (system tenant / unscoped) are allowed to manage any tenant.
func ensureCanManageMembers(c *gin.Context, targetTenantID uint) bool {
	// When the request is tenant-scoped, TenantContext stamps tenantID.
	if raw, ok := c.Get("tenantID"); ok {
		if current, ok := raw.(uint); ok && current != 0 {
			// must not cross-tenant
			if targetTenantID != current {
				response.FailWithCode(c, "admin.plugin.tenant.access_denied")
				return false
			}
			actor := utils.GetUserID(c)
			if !serviceMembership.IsPrimaryMember(actor, current) {
				response.FailWithCode(c, "admin.plugin.tenant.access_denied")
				return false
			}
		}
	}
	return true
}

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
	if !ensureCanManageMembers(c, req.TenantID) {
		return
	}
	ctx := systemService.WithRequestContext(c.Request.Context(), c)
	if err := serviceMembership.Assign(ctx, req.UserID, req.TenantID, req.IsPrimary); err != nil {
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
	if !ensureCanManageMembers(c, req.TenantID) {
		return
	}
	ctx := systemService.WithRequestContext(c.Request.Context(), c)
	if err := serviceMembership.Unassign(ctx, req.UserID, req.TenantID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithCode(c, "admin.plugin.tenant.removed")
}

// CreateUserAndAssign
// @Tags     TenantMembership
// @Summary  create a basic user with the default Tenant role and assign to a tenant
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateUserAndAssignReq true "tenantID, userName, password, nickName, phone, email, isPrimary"
// @Router   /tenantMembership/createUser [post]
func (a *membershipApi) CreateUserAndAssign(c *gin.Context) {
	var req request.CreateUserAndAssignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ctx := systemService.WithRequestContext(c.Request.Context(), c)
	if _, err := serviceMembership.CreateUserAndAssign(ctx, req); err != nil {
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

// MembersOfTenant
// @Tags     TenantMembership
// @Router   /tenantMembership/members [get]
func (a *membershipApi) MembersOfTenant(c *gin.Context) {
	var req request.IdReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !ensureCanManageMembers(c, req.ID) {
		return
	}
	list, err := serviceMembership.MembersOfTenant(req.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}
