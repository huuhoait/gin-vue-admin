package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginLogApi struct{}

func (s *LoginLogApi) DeleteLoginLog(c *gin.Context) {
	var loginLog system.SysLoginLog
	err := c.ShouldBindJSON(&loginLog)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = loginLogService.DeleteLoginLog(loginLog)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Failed to delete", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

func (s *LoginLogApi) DeleteLoginLogByIds(c *gin.Context) {
	var SDS request.IdsReq
	err := c.ShouldBindJSON(&SDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = loginLogService.DeleteLoginLogByIds(SDS)
	if err != nil {
		global.GVA_LOG.Error("Failed to batch delete!", zap.Error(err))
		response.FailWithMessage("Failed to batch delete", c)
		return
	}
	response.OkWithMessage("Batch deleted successfully", c)
}

func (s *LoginLogApi) FindLoginLog(c *gin.Context) {
	var loginLog system.SysLoginLog
	err := c.ShouldBindQuery(&loginLog)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	reLoginLog, err := loginLogService.GetLoginLog(loginLog.ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to query!", zap.Error(err))
		response.FailWithMessage("Failed to query", c)
		return
	}
	response.OkWithDetailed(reLoginLog, "Query successful", c)
}

func (s *LoginLogApi) GetLoginLogList(c *gin.Context) {
	var pageInfo systemReq.SysLoginLogSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := loginLogService.GetLoginLogInfoList(pageInfo)
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
