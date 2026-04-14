package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AIWorkflowSessionApi struct{}

func (a *AIWorkflowSessionApi) Save(c *gin.Context) {
	var info systemReq.SysAIWorkflowSessionUpsert
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session, err := aiWorkflowSessionService.Save(c.Request.Context(), utils.GetUserID(c), info)
	if err != nil {
		global.GVA_LOG.Error("Failed to save AI workflow session", zap.Error(err))
		response.FailWithMessage("Failed to save session", c)
		return
	}

	response.OkWithDetailed(gin.H{"session": session}, "Saved successfully", c)
}

func (a *AIWorkflowSessionApi) GetList(c *gin.Context) {
	var info systemReq.SysAIWorkflowSessionSearch
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := aiWorkflowSessionService.GetList(c.Request.Context(), utils.GetUserID(c), info)
	if err != nil {
		global.GVA_LOG.Error("Failed to get AI workflow session list", zap.Error(err))
		response.FailWithMessage("Failed to get session list", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "Retrieved successfully", c)
}

func (a *AIWorkflowSessionApi) GetDetail(c *gin.Context) {
	var info commonReq.GetById
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session, err := aiWorkflowSessionService.GetDetail(c.Request.Context(), utils.GetUserID(c), info.Uint())
	if err != nil {
		global.GVA_LOG.Error("Failed to get AI workflow session detail", zap.Error(err))
		response.FailWithMessage("Failed to get session details", c)
		return
	}

	response.OkWithDetailed(gin.H{"session": session}, "Retrieved successfully", c)
}

func (a *AIWorkflowSessionApi) Delete(c *gin.Context) {
	var info commonReq.GetById
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := aiWorkflowSessionService.Delete(c.Request.Context(), utils.GetUserID(c), info.Uint()); err != nil {
		global.GVA_LOG.Error("Failed to delete AI workflow session", zap.Error(err))
		response.FailWithMessage("Failed to delete session", c)
		return
	}

	response.OkWithMessage("Deleted successfully", c)
}

func (a *AIWorkflowSessionApi) DumpMarkdown(c *gin.Context) {
	var info systemReq.SysAIWorkflowMarkdownDump
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := aiWorkflowSessionService.DumpMarkdown(c.Request.Context(), utils.GetUserID(c), info)
	if err != nil {
		global.GVA_LOG.Error("Failed to dump AI workflow Markdown to disk", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"result": result}, "Saved to disk successfully", c)
}
