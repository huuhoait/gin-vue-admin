package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
)

// ErrNoUserUUID is returned when the JWT claims are missing or do not
// contain a valid user UUID. Callers that forward downstream should abort
// the request rather than proxying with an empty identity.
var ErrNoUserUUID = errors.New("proxy: current user uuid unavailable from JWT claims")

// AuthHeaders extracts the admin user's identity from the JWT claims
// (set by GVA middleware as "claims" in context) and returns headers
// to inject into downstream requests.
//
// X-Maker-ID / X-Checker-ID carry the admin's UUID so the Core service
// can populate audit fields (created_by / updated_by) and satisfy its
// `uuid` validator on maker_id/checker_id body fields.
//
// Returns ErrNoUserUUID when claims are missing or invalid — the caller
// must abort instead of forwarding unauthenticated requests.
func AuthHeaders(c *gin.Context) (map[string]string, error) {
	id, err := CurrentUserUUID(c)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"X-Maker-ID":   id,
		"X-Checker-ID": id,
	}
	// Propagate trace id if present to allow end-to-end request tracing
	// across admin proxy -> core/order services.
	if traceID := c.GetHeader("X-Trace-Id"); traceID != "" {
		headers["X-Trace-Id"] = traceID
	}
	return headers, nil
}

// CurrentUserUUID returns the admin user's UUID from JWT claims.
// Returns ErrNoUserUUID when claims are absent, of an unexpected type,
// or carry the zero UUID — all of which indicate a misconfigured auth
// chain (e.g. a proxy route mounted outside the JWT-protected group).
func CurrentUserUUID(c *gin.Context) (string, error) {
	claims, ok := c.Get("claims")
	if !ok {
		return "", ErrNoUserUUID
	}
	cl, ok := claims.(*systemReq.CustomClaims)
	if !ok {
		return "", ErrNoUserUUID
	}
	if cl.UUID == (uuid.UUID{}) {
		return "", ErrNoUserUUID
	}
	return cl.UUID.String(), nil
}
