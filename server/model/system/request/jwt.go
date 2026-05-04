package request

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	NickName    string
	AuthorityId uint
	// TenantID is the active tenant for this session. 0 means the system
	// tenant — used for super-admin / cross-tenant access where no scoping
	// is applied. Embedded here so the per-request tenant resolution
	// middleware can read it without hitting the database. Pre-existing
	// tokens issued before this field was added will deserialize as 0
	// (system tenant), which matches the previous unscoped behaviour.
	TenantID uint
}
