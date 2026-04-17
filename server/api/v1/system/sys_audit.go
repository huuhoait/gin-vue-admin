package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
)

type AuditApi struct{}

// VerifyAuditChain
// @Tags      Audit
// @Summary   Verify policy-change audit log hash chain integrity
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{data=system.ChainVerifyResult,msg=string}
// @Router    /audit/verifyChain [get]
func (a *AuditApi) VerifyAuditChain(c *gin.Context) {
	result, err := auditService.VerifyAuditChain(c.Request.Context())
	if err != nil {
		global.GVA_LOG.Error("audit chain verification failed", zap.Error(err))
		response.FailWithMessage("Verification error: "+err.Error(), c)
		return
	}
	if !result.OK {
		response.FailWithDetailed(result, "Audit chain integrity failure — possible tampering detected", c)
		return
	}
	response.OkWithDetailed(result, "Audit chain intact", c)
}
