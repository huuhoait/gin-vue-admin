package api

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/plugin/announcement/model"
	"github.com/huuhoait/gin-vue-admin/server/plugin/announcement/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Info = new(info)

type info struct{}

// CreateInfo create announcement
// @Tags Info
// @Summary create announcement
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Info true "create announcement"
// @Success 200 {object} response.Response{msg=string} "Created successfully"
// @Router /info/createInfo [post]
func (a *info) CreateInfo(c *gin.Context) {
	var info model.Info
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = serviceInfo.CreateInfo(&info)
	if err != nil {
		global.GVA_LOG.Error("Creation failed!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// DeleteInfo delete announcement
// @Tags Info
// @Summary delete announcement
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Info true "delete announcement"
// @Success 200 {object} response.Response{msg=string} "Deleted successfully"
// @Router /info/deleteInfo [delete]
func (a *info) DeleteInfo(c *gin.Context) {
	ID := c.Query("ID")
	err := serviceInfo.DeleteInfo(ID)
	if err != nil {
		global.GVA_LOG.Error("Deletion failed!", zap.Error(err))
		response.FailWithMessage("Deletion failed", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// DeleteInfoByIds batch delete announcements
// @Tags Info
// @Summary batch delete announcements
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "Batch delete succeeded"
// @Router /info/deleteInfoByIds [delete]
func (a *info) DeleteInfoByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := serviceInfo.DeleteInfoByIds(IDs); err != nil {
		global.GVA_LOG.Error("Batch delete failed!", zap.Error(err))
		response.FailWithMessage("Batch delete failed", c)
		return
	}
	response.OkWithMessage("Batch delete succeeded", c)
}

// UpdateInfo update announcement
// @Tags Info
// @Summary update announcement
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Info true "update announcement"
// @Success 200 {object} response.Response{msg=string} "Updated successfully"
// @Router /info/updateInfo [put]
func (a *info) UpdateInfo(c *gin.Context) {
	var info model.Info
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = serviceInfo.UpdateInfo(info)
	if err != nil {
		global.GVA_LOG.Error("Update failed!", zap.Error(err))
		response.FailWithMessage("Update failed", c)
		return
	}
	response.OkWithMessage("Updated successfully", c)
}

// FindInfo find announcement by ID
// @Tags Info
// @Summary find announcement by ID
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Info true "find announcement by ID"
// @Success 200 {object} response.Response{data=model.Info,msg=string} "Query succeeded"
// @Router /info/findInfo [get]
func (a *info) FindInfo(c *gin.Context) {
	ID := c.Query("ID")
	reinfo, err := serviceInfo.GetInfo(ID)
	if err != nil {
		global.GVA_LOG.Error("Query failed!", zap.Error(err))
		response.FailWithMessage("Query failed", c)
		return
	}
	response.OkWithData(reinfo, c)
}

// GetInfoList get paginated announcement list
// @Tags Info
// @Summary get paginated announcement list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.InfoSearch true "get paginated announcement list"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "Retrieved successfully"
// @Router /info/getInfoList [get]
func (a *info) GetInfoList(c *gin.Context) {
	var pageInfo request.InfoSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.GetInfoInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Retrieval failed!", zap.Error(err))
		response.FailWithMessage("Retrieval failed", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Retrieved successfully", c)
}

// GetInfoDataSource getInfodata source
// @Tags Info
// @Summary getInfodata source
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "Query succeeded"
// @Router /info/getInfoDataSource [get]
func (a *info) GetInfoDataSource(c *gin.Context) {
	// This API retrieves data source definitions
	dataSource, err := serviceInfo.GetInfoDataSource()
	if err != nil {
		global.GVA_LOG.Error("Query failed!", zap.Error(err))
		response.FailWithMessage("Query failed", c)
		return
	}
	response.OkWithData(dataSource, c)
}

// GetInfoPublic Public announcement API (no auth required)
// @Tags Info
// @Summary Public announcement API (no auth required)
// @accept application/json
// @Produce application/json
// @Param data query request.InfoSearch true "get paginated announcement list"
// @Success 200 {object} response.Response{data=object,msg=string} "Retrieved successfully"
// @Router /info/getInfoPublic [get]
func (a *info) GetInfoPublic(c *gin.Context) {
	// This API does not require authentication ExampleForReturnDoneOnePieceFixedofMessageAPI, GeneralThisAPIUsed forCEndService, NeedSelfImplementBusinessLogic
	response.OkWithDetailed(gin.H{"info": "Public announcement API (no auth required)Information"}, "Retrieved successfully", c)
}
