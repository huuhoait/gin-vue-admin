package proxy

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// AuthHeaders extracts the admin user's identity from the JWT claims
// (set by GVA middleware as "claims" in context) and returns headers
// to inject into downstream requests.
//
// X-Maker-ID / X-Checker-ID carry the admin's numeric user ID so the
// Core service can populate audit fields (created_by / updated_by).
func AuthHeaders(c *gin.Context) map[string]string {
	var userID uint
	if claims, ok := c.Get("claims"); ok {
		// The claims type from GVA middleware has BaseClaims.ID (uint).
		// Use interface to avoid importing the full system model chain.
		type idProvider interface{ GetID() uint }
		if p, ok := claims.(idProvider); ok {
			userID = p.GetID()
		}
	}
	id := fmt.Sprintf("%d", userID)
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
