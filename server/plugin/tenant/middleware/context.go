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
//  1. JWT claim — utils.GetTenantID(c). The login flow stamps the user's
//     primary tenant into the claims, so the common case requires no DB
//     access. tenantID=0 means the system tenant (super-admin / unscoped).
//     Note: pre-existing tokens issued before TenantID was added to the
//     claim struct deserialize as 0; that grants system-tenant scope by
//     default, which matches the pre-change behaviour (no scoping was
//     happening before this feature shipped).
//  2. X-Tenant-ID header override. When the header is present and either
//     the JWT carries no tenant or the header differs from the claim, the
//     middleware verifies membership against gva_user_tenants and uses
//     the header's tenant. This supports super-admin tenant switching
//     without re-issuing the JWT.
//  3. JWT tenantID=0 with no header — the request runs unscoped (system
//     tenant). FromContext will return ok=false in that case so callers
//     using WithTenantScope see all rows (super-admin behaviour).
//
// Mount this AFTER middleware.JWTAuth so claims are available.
func TenantContext() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := utils.GetUserID(c)
		if userID == 0 || utils.GetUserAuthorityId(c) == 888 {
			c.Next()
			return
		}

		// 1. Read tenant from JWT claims (no DB access).
		tid := utils.GetTenantID(c)

		// 2. Header override — only verify membership when the header is
		// supplied AND it differs from whatever the JWT already says.
		if raw := c.GetHeader("X-Tenant-ID"); raw != "" {
			id, err := strconv.ParseUint(raw, 10, 64)
			if err != nil || id == 0 {
				response.FailWithCode(c, "admin.plugin.tenant.invalid_header")
				c.Abort()
				return
			}
			headerTID := uint(id)
			if headerTID != tid {
				// Authoritative membership check is required when the
				// caller asks to operate as a tenant other than their
				// claim says. This is the only path that hits the DB.
				if !service.Service.Membership.HasAccess(userID, headerTID) {
					response.FailWithCode(c, "admin.plugin.tenant.access_denied")
					c.Abort()
					return
				}
				tid = headerTID
			}
		}

		// 3. Apply scoping. tid==0 means system tenant — leave the context
		// unscoped so super-admin queries see everything (FromContext
		// will return ok=false).
		if tid != 0 {
			c.Request = c.Request.WithContext(service.WithTenant(c.Request.Context(), tid))
			c.Set("tenantID", tid)
		}
		c.Next()
	}
}
