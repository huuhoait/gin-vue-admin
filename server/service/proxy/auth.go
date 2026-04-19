package proxy

import (
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHeaders extracts the admin user's identity from the JWT claims
// (set by GVA middleware as "claims" in context) and returns headers
// to inject into downstream requests.
//
// X-Maker-ID / X-Checker-ID carry the admin's UUID so the Core service
// can populate audit fields (created_by / updated_by) and satisfy its
// `uuid` validator on maker_id/checker_id body fields.
func AuthHeaders(c *gin.Context) map[string]string {
	id := CurrentUserUUID(c)
	headers := map[string]string{
		"X-Maker-ID":   id,
		"X-Checker-ID": id,
	}
	// Propagate trace id if present to allow end-to-end request tracing
	// across admin proxy -> core/order services.
	if traceID := c.GetHeader("X-Trace-Id"); traceID != "" {
		headers["X-Trace-Id"] = traceID
	}
	return headers
}

// CurrentUserUUID returns the admin user's UUID from JWT claims, or the
// empty-UUID string when claims are missing (e.g. anonymous public routes).
func CurrentUserUUID(c *gin.Context) string {
	if claims, ok := c.Get("claims"); ok {
		if cl, ok := claims.(*systemReq.CustomClaims); ok && cl.UUID != (uuid.UUID{}) {
			return cl.UUID.String()
		}
	}
	return uuid.UUID{}.String()
}
