package api

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/model/request"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type clientApi struct{}

// CreateClient register a new OAuth2 client
// @Tags     OAuth2Client
// @Summary  create OAuth2 client (returns plaintext secret once)
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateClientReq true "client metadata"
// @Success  200 {object} response.Response{data=object,msg=string}
// @Router   /oauth2Client/create [post]
func (a *clientApi) CreateClient(c *gin.Context) {
	var req request.CreateClientReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, secret, err := serviceClient.CreateClient(req, utils.GetUserID(c))
	if err != nil {
		global.GVA_LOG.Error("create oauth2 client failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDataCode(c, gin.H{
		"client":       row,
		"clientId":     row.ClientID,
		"clientSecret": secret, // shown once; persist on the integrator side
	}, "admin.plugin.oauth2.client_created")
}

// UpdateClient
// @Tags     OAuth2Client
// @Summary  update OAuth2 client
// @Security ApiKeyAuth
// @Router   /oauth2Client/update [put]
func (a *clientApi) UpdateClient(c *gin.Context) {
	var req request.UpdateClientReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, err := serviceClient.UpdateClient(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(row, c)
}

// DeleteClient
// @Tags     OAuth2Client
// @Summary  delete OAuth2 client
// @Security ApiKeyAuth
// @Router   /oauth2Client/delete [delete]
func (a *clientApi) DeleteClient(c *gin.Context) {
	var req request.IdReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceClient.DeleteClient(req.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithCode(c, "admin.common.delete_success")
}

// FindClient
// @Tags     OAuth2Client
// @Router   /oauth2Client/find [get]
func (a *clientApi) FindClient(c *gin.Context) {
	var req request.IdReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	row, err := serviceClient.FindByID(req.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(row, c)
}

// ListClients
// @Tags     OAuth2Client
// @Router   /oauth2Client/list [get]
func (a *clientApi) ListClients(c *gin.Context) {
	var req request.ClientListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceClient.ListClients(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDataCode(c, response.PageResult{
		List: list, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "admin.common.get_success")
}

// RegenerateSecret rotates a client's secret
// @Tags     OAuth2Client
// @Router   /oauth2Client/regenerateSecret [post]
func (a *clientApi) RegenerateSecret(c *gin.Context) {
	var req request.IdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	secret, err := serviceClient.RegenerateSecret(req.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDataCode(c, gin.H{"clientSecret": secret},
		"admin.plugin.oauth2.secret_rotated")
}
