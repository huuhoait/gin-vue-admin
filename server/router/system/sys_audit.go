package system

import (
	"github.com/gin-gonic/gin"
	"github.com/huuhoait/gin-vue-admin/server/middleware"
)

type AuditRouter struct{}

// InitAuditRouter registers audit endpoints under /audit.
// All routes require a valid JWT (caller must add JWTAuth middleware in
// the parent group) and are restricted to super-admin by Casbin policy.
func (r *AuditRouter) InitAuditRouter(Router *gin.RouterGroup) {
	auditRouter := Router.Group("audit").Use(middleware.OperationRecord())
	{
		auditRouter.GET("verifyChain", auditApiVar.VerifyAuditChain)

	}
}
