package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	sysService "github.com/huuhoait/gin-vue-admin/server/service/system"
	"github.com/huuhoait/gin-vue-admin/server/utils"
)

type GDPRApi struct{}

// ExportMyData
// @Tags      GDPR
// @Summary   Export all personal data for the authenticated user (GDPR Art. 20)
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}
// @Router    /user/gdpr/export [get]
func (g *GDPRApi) ExportMyData(c *gin.Context) {
	userID := utils.GetUserID(c)
	ctx := sysService.WithRequestContext(c.Request.Context(), c)
	export, err := gdprService.ExportUserData(ctx, userID)
	if err != nil {
		global.GVA_LOG.Error("GDPR export failed", zap.Error(err), zap.Uint("user_id", userID))
		response.FailWithMessage("Export failed: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(export, "Data export successful", c)
}

// EraseMyData
// @Tags      GDPR
// @Summary   Erase all personal data for the authenticated user (GDPR Art. 17)
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}
// @Router    /user/gdpr/erase [delete]
func (g *GDPRApi) EraseMyData(c *gin.Context) {
	userID := utils.GetUserID(c)
	ctx := sysService.WithRequestContext(c.Request.Context(), c)
	if err := gdprService.EraseUserData(ctx, userID); err != nil {
		global.GVA_LOG.Error("GDPR erasure failed", zap.Error(err), zap.Uint("user_id", userID))
		response.FailWithMessage("Erasure failed: "+err.Error(), c)
		return
	}
	// Invalidate the current JWT after erasure
	tokenString := c.Request.Header.Get("x-token")
	if tokenString != "" {
		_ = utils.BlacklistAdd(tokenString)
	}
	c.Header("x-token", "")
	response.OkWithMessage("Your data has been anonymized and your account closed", c)
}

// AdminEraseUser
// @Tags      GDPR
// @Summary   Admin: erase a specific user's personal data (GDPR Art. 17)
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      object  true  "User ID to erase"
// @Success   200   {object}  response.Response{msg=string}
// @Router    /user/gdpr/adminErase [delete]
func (g *GDPRApi) AdminEraseUser(c *gin.Context) {
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Prevent admin from erasing themselves via this endpoint
	if req.UserID == utils.GetUserID(c) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Cannot erase your own account via admin endpoint"})
		return
	}
	ctx := sysService.WithRequestContext(c.Request.Context(), c)
	if err := gdprService.EraseUserData(ctx, req.UserID); err != nil {
		global.GVA_LOG.Error("Admin GDPR erasure failed", zap.Error(err), zap.Uint("target_user_id", req.UserID))
		response.FailWithMessage("Erasure failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("User data anonymized successfully", c)
}
