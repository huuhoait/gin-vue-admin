package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	sysReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiTokenApi struct{}

// CreateApiToken Issue Token
func (s *ApiTokenApi) CreateApiToken(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"userId"`
		AuthorityID uint   `json:"authorityId"`
		Days        int    `json:"days"` // -1 for permanent, otherwise number of days
		Remark      string `json:"remark"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	token := system.SysApiToken{
		UserID:      req.UserID,
		AuthorityID: req.AuthorityID,
		Remark:      req.Remark,
	}

	jwtStr, err := apiTokenService.CreateApiToken(token, req.Days)
	if err != nil {
		global.GVA_LOG.Error("Failed to issue token!", zap.Error(err))
		response.FailWithMessage("Failed to issue token: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"token": jwtStr}, "Token issued successfully", c)
}

// GetApiTokenList Get list
func (s *ApiTokenApi) GetApiTokenList(c *gin.Context) {
	var pageInfo sysReq.SysApiTokenSearch
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := apiTokenService.GetApiTokenList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Retrieved successfully", c)
}

// DeleteApiToken Revoke Token
func (s *ApiTokenApi) DeleteApiToken(c *gin.Context) {
	var req system.SysApiToken
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiTokenService.DeleteApiToken(req.ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to revoke!", zap.Error(err))
		response.FailWithMessage("Failed to revoke token", c)
		return
	}
	response.OkWithMessage("Token revoked successfully", c)
}
