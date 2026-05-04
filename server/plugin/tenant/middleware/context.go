package middleware

import (
	"strconv"

	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
)

// TenantContext resolves the active tenant for the current request and
// stores it on the gin/request context so handlers can call
// service.FromContext(c.Request.Context()).
//
// Resolution order:
//  1. X-Tenant-ID header — admin-style explicit selection. The middleware
//     verifies the authenticated user has access to that tenant.
//  2. The user's primary tenant from gva_user_tenants.
//  3. Unscoped (no value set) — the request runs without tenant filtering.
//
// Mount this AFTER middleware.JWTAuth so claims are available.
func TenantContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := utils.GetUserID(c)
		if userID == 0 {
			c.Next()
			return
		}

		// Explicit header selection.
		if raw := c.GetHeader("X-Tenant-ID"); raw != "" {
			id, err := strconv.ParseUint(raw, 10, 64)
			if err != nil || id == 0 {
				response.FailWithCode(c, "admin.plugin.tenant.invalid_header")
				c.Abort()
				return
			}
			tid := uint(id)
			if !service.Service.Membership.HasAccess(userID, tid) {
				response.FailWithCode(c, "admin.plugin.tenant.access_denied")
				c.Abort()
				return
			}
			c.Request = c.Request.WithContext(service.WithTenant(c.Request.Context(), tid))
			c.Set("tenantID", tid)
			c.Next()
			return
		}

		// Implicit primary-tenant fallback.
		if tid, ok := service.Service.Membership.PrimaryTenantForUser(userID); ok {
			c.Request = c.Request.WithContext(service.WithTenant(c.Request.Context(), tid))
			c.Set("tenantID", tid)
		}
		c.Next()
	}
}
