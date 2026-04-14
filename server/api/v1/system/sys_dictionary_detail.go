package system

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DictionaryDetailApi struct{}

// CreateSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   Create SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysDictionaryDetail     true  "SysDictionaryDetail model"
// @Success   200   {object}  response.Response{msg=string}  "Create SysDictionaryDetail"
// @Router    /sysDictionaryDetail/createSysDictionaryDetail [post]
func (s *DictionaryDetailApi) CreateSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindJSON(&detail)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryDetailService.CreateSysDictionaryDetail(detail)
	if err != nil {
		global.GVA_LOG.Error("Failed to create!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// DeleteSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   Delete SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysDictionaryDetail     true  "SysDictionaryDetail model"
// @Success   200   {object}  response.Response{msg=string}  "Delete SysDictionaryDetail"
// @Router    /sysDictionaryDetail/deleteSysDictionaryDetail [delete]
func (s *DictionaryDetailApi) DeleteSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindJSON(&detail)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryDetailService.DeleteSysDictionaryDetail(detail)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// UpdateSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   Update SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysDictionaryDetail     true  "Update SysDictionaryDetail"
// @Success   200   {object}  response.Response{msg=string}  "Update SysDictionaryDetail"
// @Router    /sysDictionaryDetail/updateSysDictionaryDetail [put]
func (s *DictionaryDetailApi) UpdateSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindJSON(&detail)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryDetailService.UpdateSysDictionaryDetail(&detail)
	if err != nil {
		global.GVA_LOG.Error("Failed to update!", zap.Error(err))
		response.FailWithMessage("Update failed", c)
		return
	}
	response.OkWithMessage("Updated successfully", c)
}

// FindSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   Find SysDictionaryDetail by ID
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     system.SysDictionaryDetail                                 true  "Find SysDictionaryDetail by ID"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Find SysDictionaryDetail by ID"
// @Router    /sysDictionaryDetail/findSysDictionaryDetail [get]
func (s *DictionaryDetailApi) FindSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindQuery(&detail)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(detail, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	reSysDictionaryDetail, err := dictionaryDetailService.GetSysDictionaryDetail(detail.ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to query!", zap.Error(err))
		response.FailWithMessage("Query failed", c)
		return
	}
	response.OkWithDetailed(gin.H{"reSysDictionaryDetail": reSysDictionaryDetail}, "Query successful", c)
}

// GetSysDictionaryDetailList
// @Tags      SysDictionaryDetail
// @Summary   Get SysDictionaryDetail list with pagination
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.SysDictionaryDetailSearch                       true  "Page number, page size, search criteria"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Get paginated SysDictionaryDetail list, returns list, total, page, page size"
// @Router    /sysDictionaryDetail/getSysDictionaryDetailList [get]
func (s *DictionaryDetailApi) GetSysDictionaryDetailList(c *gin.Context) {
	var pageInfo request.SysDictionaryDetailSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := dictionaryDetailService.GetSysDictionaryDetailInfoList(pageInfo)
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

// GetDictionaryTreeList
// @Tags      SysDictionaryDetail
// @Summary   Get dictionary detail tree structure
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     sysDictionaryID  query     int                                                true  "Dictionary ID"
// @Success   200              {object}  response.Response{data=[]system.SysDictionaryDetail,msg=string}  "Get dictionary detail tree structure"
// @Router    /sysDictionaryDetail/getDictionaryTreeList [get]
func (s *DictionaryDetailApi) GetDictionaryTreeList(c *gin.Context) {
	sysDictionaryID := c.Query("sysDictionaryID")
	if sysDictionaryID == "" {
		response.FailWithMessage("Dictionary ID cannot be empty", c)
		return
	}

	var id uint
	if idUint64, err := strconv.ParseUint(sysDictionaryID, 10, 32); err != nil {
		response.FailWithMessage("Invalid dictionary ID format", c)
		return
	} else {
		id = uint(idUint64)
	}
	
	list, err := dictionaryDetailService.GetDictionaryTreeList(id)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "Retrieved successfully", c)
}

// GetDictionaryTreeListByType
// @Tags      SysDictionaryDetail
// @Summary   Get dictionary detail tree structure by dictionary type
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     type  query     string                                                true  "Dictionary type"
// @Success   200   {object}  response.Response{data=[]system.SysDictionaryDetail,msg=string}  "Get dictionary detail tree structure"
// @Router    /sysDictionaryDetail/getDictionaryTreeListByType [get]
func (s *DictionaryDetailApi) GetDictionaryTreeListByType(c *gin.Context) {
	dictType := c.Query("type")
	if dictType == "" {
		response.FailWithMessage("Dictionary type cannot be empty", c)
		return
	}
	
	list, err := dictionaryDetailService.GetDictionaryTreeListByType(dictType)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "Retrieved successfully", c)
}

// GetDictionaryDetailsByParent
// @Tags      SysDictionaryDetail
// @Summary   Get dictionary details by parent ID
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetDictionaryDetailsByParentRequest                true  "Query parameters"
// @Success   200   {object}  response.Response{data=[]system.SysDictionaryDetail,msg=string}  "Get dictionary detail list"
// @Router    /sysDictionaryDetail/getDictionaryDetailsByParent [get]
func (s *DictionaryDetailApi) GetDictionaryDetailsByParent(c *gin.Context) {
	var req request.GetDictionaryDetailsByParentRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	list, err := dictionaryDetailService.GetDictionaryDetailsByParent(req)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "Retrieved successfully", c)
}

// GetDictionaryPath
// @Tags      SysDictionaryDetail
// @Summary   Get the full path of a dictionary detail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id  query     uint                                                true  "Dictionary detail ID"
// @Success   200 {object}  response.Response{data=[]system.SysDictionaryDetail,msg=string}  "Get dictionary detail path"
// @Router    /sysDictionaryDetail/getDictionaryPath [get]
func (s *DictionaryDetailApi) GetDictionaryPath(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("Dictionary detail ID cannot be empty", c)
		return
	}
	
	var id uint
	if idUint64, err := strconv.ParseUint(idStr, 10, 32); err != nil {
		response.FailWithMessage("Invalid dictionary detail ID format", c)
		return
	} else {
		id = uint(idUint64)
	}
	
	path, err := dictionaryDetailService.GetDictionaryPath(id)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"path": path}, "Retrieved successfully", c)
}
