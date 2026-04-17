package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysParamsApi struct{}

// CreateSysParams Create parameter
// @Tags SysParams
// @Summary Create parameter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysParams true "Create parameter"
// @Success 200 {object} response.Response{msg=string} "Created successfully"
// @Router /sysParams/createSysParams [post]
func (sysParamsApi *SysParamsApi) CreateSysParams(c *gin.Context) {
	var sysParams system.SysParams
	err := c.ShouldBindJSON(&sysParams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysParamsService.CreateSysParams(&sysParams)
	if err != nil {
		global.GVA_LOG.Error("Failed to create!", zap.Error(err))
		response.FailWithMessage("Creation failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// DeleteSysParams Delete parameter
// @Tags SysParams
// @Summary Delete parameter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysParams true "Delete parameter"
// @Success 200 {object} response.Response{msg=string} "Deleted successfully"
// @Router /sysParams/deleteSysParams [delete]
func (sysParamsApi *SysParamsApi) DeleteSysParams(c *gin.Context) {
	ID := c.Query("ID")
	err := sysParamsService.DeleteSysParams(ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// DeleteSysParamsByIds Batch delete parameters
// @Tags SysParams
// @Summary Batch delete parameters
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "Batch deleted successfully"
// @Router /sysParams/deleteSysParamsByIds [delete]
func (sysParamsApi *SysParamsApi) DeleteSysParamsByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	err := sysParamsService.DeleteSysParamsByIds(IDs)
	if err != nil {
		global.GVA_LOG.Error("Failed to batch delete!", zap.Error(err))
		response.FailWithMessage("Batch deletion failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Batch deleted successfully", c)
}

// UpdateSysParams Update parameter
// @Tags SysParams
// @Summary Update parameter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysParams true "Update parameter"
// @Success 200 {object} response.Response{msg=string} "Updated successfully"
// @Router /sysParams/updateSysParams [put]
func (sysParamsApi *SysParamsApi) UpdateSysParams(c *gin.Context) {
	var sysParams system.SysParams
	err := c.ShouldBindJSON(&sysParams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysParamsService.UpdateSysParams(sysParams)
	if err != nil {
		global.GVA_LOG.Error("Failed to update!", zap.Error(err))
		response.FailWithMessage("Update failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Updated successfully", c)
}

// FindSysParams Find parameter by ID
// @Tags SysParams
// @Summary Find parameter by ID
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query system.SysParams true "Find parameter by ID"
// @Success 200 {object} response.Response{data=system.SysParams,msg=string} "Query successful"
// @Router /sysParams/findSysParams [get]
func (sysParamsApi *SysParamsApi) FindSysParams(c *gin.Context) {
	ID := c.Query("ID")
	resysParams, err := sysParamsService.GetSysParams(ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to query!", zap.Error(err))
		response.FailWithMessage("Query failed: "+err.Error(), c)
		return
	}
	response.OkWithData(resysParams, c)
}

// GetSysParamsList Get parameter list with pagination
// @Tags SysParams
// @Summary Get parameter list with pagination
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query systemReq.SysParamsSearch true "Get parameter list with pagination"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "Retrieved successfully"
// @Router /sysParams/getSysParamsList [get]
func (sysParamsApi *SysParamsApi) GetSysParamsList(c *gin.Context) {
	var pageInfo systemReq.SysParamsSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysParamsService.GetSysParamsInfoList(pageInfo)
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

// GetSysParam Get parameter value by key
// @Tags SysParams
// @Summary Get parameter value by key
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param key query string true "key"
// @Success 200 {object} response.Response{data=system.SysParams,msg=string} "Retrieved successfully"
// @Router /sysParams/getSysParam [get]
func (sysParamsApi *SysParamsApi) GetSysParam(c *gin.Context) {
	k := c.Query("key")
	params, err := sysParamsService.GetSysParam(k)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Retrieval failed: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(params, "Retrieved successfully", c)
}
