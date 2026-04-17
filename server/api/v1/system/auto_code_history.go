package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	common "github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	request "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AutoCodeHistoryApi struct{}

// First
// @Tags      AutoCode
// @Summary   Get meta information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                                            true  "Request parameters"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Get meta information"
// @Router    /autoCode/getMeta [post]
func (a *AutoCodeHistoryApi) First(c *gin.Context) {
	var info common.GetById
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := autoCodeHistoryService.First(c.Request.Context(), info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"meta": data}, "Retrieved successfully", c)
}

// Delete
// @Tags      AutoCode
// @Summary   Delete rollback record
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "Request parameters"
// @Success   200   {object}  response.Response{msg=string}  "Delete rollback record"
// @Router    /autoCode/delSysHistory [post]
func (a *AutoCodeHistoryApi) Delete(c *gin.Context) {
	var info common.GetById
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = autoCodeHistoryService.Delete(c.Request.Context(), info)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Failed to delete", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// RollBack
// @Tags      AutoCode
// @Summary   Rollback auto-generated code
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SysAutoHistoryRollBack             true  "Request parameters"
// @Success   200   {object}  response.Response{msg=string}  "Rollback auto-generated code"
// @Router    /autoCode/rollback [post]
func (a *AutoCodeHistoryApi) RollBack(c *gin.Context) {
	var info request.SysAutoHistoryRollBack
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = autoCodeHistoryService.RollBack(c.Request.Context(), info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Rolled back successfully", c)
}

// GetList
// @Tags      AutoCode
// @Summary   Query rollback records
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      common.PageInfo                                true  "Request parameters"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Query rollback records, returns list, total, page number, page size"
// @Router    /autoCode/getSysHistory [post]
func (a *AutoCodeHistoryApi) GetList(c *gin.Context) {
	var info common.PageInfo
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := autoCodeHistoryService.GetList(c.Request.Context(), info)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "Retrieved successfully", c)
}
