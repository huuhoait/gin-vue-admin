package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysErrorApi struct{}

// CreateSysError Create error log
// @Tags SysError
// @Summary Create error log
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysError true "Create error log"
// @Success 200 {object} response.Response{msg=string} "Created successfully"
// @Router /sysError/createSysError [post]
func (sysErrorApi *SysErrorApi) CreateSysError(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	var sysError system.SysError
	err := c.ShouldBindJSON(&sysError)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysErrorService.CreateSysError(ctx, &sysError)
	if err != nil {
		global.GVA_LOG.Error("Failed to create!", zap.Error(err))
		response.FailWithMessage("Creation failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// DeleteSysError Delete error log
// @Tags SysError
// @Summary Delete error log
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysError true "Delete error log"
// @Success 200 {object} response.Response{msg=string} "Deleted successfully"
// @Router /sysError/deleteSysError [delete]
func (sysErrorApi *SysErrorApi) DeleteSysError(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := sysErrorService.DeleteSysError(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// DeleteSysErrorByIds Batch delete error logs
// @Tags SysError
// @Summary Batch delete error logs
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "Batch deleted successfully"
// @Router /sysError/deleteSysErrorByIds [delete]
func (sysErrorApi *SysErrorApi) DeleteSysErrorByIds(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := sysErrorService.DeleteSysErrorByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("Failed to batch delete!", zap.Error(err))
		response.FailWithMessage("Batch deletion failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Batch deleted successfully", c)
}

// UpdateSysError Update error log
// @Tags SysError
// @Summary Update error log
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysError true "Update error log"
// @Success 200 {object} response.Response{msg=string} "Updated successfully"
// @Router /sysError/updateSysError [put]
func (sysErrorApi *SysErrorApi) UpdateSysError(c *gin.Context) {
	// Get standard context from ctx for business operations
	ctx := c.Request.Context()

	var sysError system.SysError
	err := c.ShouldBindJSON(&sysError)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysErrorService.UpdateSysError(ctx, sysError)
	if err != nil {
		global.GVA_LOG.Error("Failed to update!", zap.Error(err))
		response.FailWithMessage("Update failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Updated successfully", c)
}

// FindSysError Find error log by ID
// @Tags SysError
// @Summary Find error log by ID
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "Find error log by ID"
// @Success 200 {object} response.Response{data=system.SysError,msg=string} "Query successful"
// @Router /sysError/findSysError [get]
func (sysErrorApi *SysErrorApi) FindSysError(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	resysError, err := sysErrorService.GetSysError(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to query!", zap.Error(err))
		response.FailWithMessage("Query failed: "+err.Error(), c)
		return
	}
	response.OkWithData(resysError, c)
}

// GetSysErrorList Get error log list with pagination
// @Tags SysError
// @Summary Get error log list with pagination
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysErrorSearch true "Get error log list with pagination"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "Retrieved successfully"
// @Router /sysError/getSysErrorList [get]
func (sysErrorApi *SysErrorApi) GetSysErrorList(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	var pageInfo systemReq.SysErrorSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysErrorService.GetSysErrorInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Retrieval failed: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Retrieved successfully", c)
}

// GetSysErrorSolution Trigger async processing of error log
// @Tags SysError
// @Summary Trigger processing by ID: mark as processing, auto-complete after 1 minute
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query string true "Error log ID"
// @Success 200 {object} response.Response{msg=string} "Processing submitted"
// @Router /sysError/getSysErrorSolution [get]
func (sysErrorApi *SysErrorApi) GetSysErrorSolution(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	// Support both "id" and "ID" parameter names
	ID := c.Query("id")
	if ID == "" {
		response.FailWithMessage("Missing parameter: id", c)
		return
	}

	err := sysErrorService.GetSysErrorSolution(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to trigger processing!", zap.Error(err))
		response.FailWithMessage("Failed to trigger processing: "+err.Error(), c)
		return
	}

	response.OkWithMessage("Submitted for AI processing", c)
}
