package api

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/plugin/onlineusers/model/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type sessionApi struct{}

// ListSessions list online sessions
// @Tags     OnlineUsers
// @Summary  list online sessions
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    data query request.ListSessionsReq true "page info + username filter"
// @Success  200 {object} response.Response{data=response.PageResult,msg=string} "Retrieved successfully"
// @Router   /onlineUsers/list [get]
func (a *sessionApi) ListSessions(c *gin.Context) {
	var req request.ListSessionsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceSession.List(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("list online sessions failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDataCode(c, response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "admin.common.get_success")
}

// KickSession force-logout an online user
// @Tags     OnlineUsers
// @Summary  kick a session by uuid
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.KickSessionReq true "uuid of the session to terminate"
// @Success  200 {object} response.Response{msg=string} "Session terminated"
// @Router   /onlineUsers/kick [post]
func (a *sessionApi) KickSession(c *gin.Context) {
	var req request.KickSessionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceSession.Kick(c.Request.Context(), req.UUID); err != nil {
		global.GVA_LOG.Error("kick session failed", zap.Error(err), zap.String("uuid", req.UUID))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithCode(c, "admin.plugin.online_users.session_terminated")
}
