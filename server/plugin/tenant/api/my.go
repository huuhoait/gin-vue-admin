package api

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type myApi struct{}

// MyTenants
// @Tags     Tenant
// @Summary  list tenants accessible to the current user
// @Description Returns the tenants the authenticated user has membership in,
// @Description with the primary tenant first. Used by the FE tenant switcher
// @Description to populate its dropdown without exposing the full tenant list
// @Description (which is admin-only via /tenant/list).
// @Security ApiKeyAuth
// @Produce  application/json
// @Success  200 {object} response.Response{data=[]service.TenantWithMembership} "tenant list"
// @Router   /tenant/mine [get]
func (a *myApi) MyTenants(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithCode(c, "admin.common.unauthorized")
		return
	}
	list, err := serviceMembership.MyTenantsForUser(userID)
	if err != nil {
		global.GVA_LOG.Error("list my tenants failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}
