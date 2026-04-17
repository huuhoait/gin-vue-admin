package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DictionaryApi struct{}

// CreateSysDictionary
// @Tags      SysDictionary
// @Summary   Create SysDictionary
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysDictionary           true  "SysDictionary model"
// @Success   200   {object}  response.Response{msg=string}  "Create SysDictionary"
// @Router    /sysDictionary/createSysDictionary [post]
func (s *DictionaryApi) CreateSysDictionary(c *gin.Context) {
	var dictionary system.SysDictionary
	err := c.ShouldBindJSON(&dictionary)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryService.CreateSysDictionary(dictionary)
	if err != nil {
		global.GVA_LOG.Error("Failed to create!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// DeleteSysDictionary
// @Tags      SysDictionary
// @Summary   Delete SysDictionary
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysDictionary           true  "SysDictionary model"
// @Success   200   {object}  response.Response{msg=string}  "Delete SysDictionary"
// @Router    /sysDictionary/deleteSysDictionary [delete]
func (s *DictionaryApi) DeleteSysDictionary(c *gin.Context) {
	var dictionary system.SysDictionary
	err := c.ShouldBindJSON(&dictionary)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryService.DeleteSysDictionary(dictionary)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// UpdateSysDictionary
// @Tags      SysDictionary
// @Summary   Update SysDictionary
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysDictionary           true  "SysDictionary model"
// @Success   200   {object}  response.Response{msg=string}  "Update SysDictionary"
// @Router    /sysDictionary/updateSysDictionary [put]
func (s *DictionaryApi) UpdateSysDictionary(c *gin.Context) {
	var dictionary system.SysDictionary
	err := c.ShouldBindJSON(&dictionary)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryService.UpdateSysDictionary(&dictionary)
	if err != nil {
		global.GVA_LOG.Error("Failed to update!", zap.Error(err))
		response.FailWithMessage("Update failed", c)
		return
	}
	response.OkWithMessage("Updated successfully", c)
}

// FindSysDictionary
// @Tags      SysDictionary
// @Summary   Find SysDictionary by ID
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     system.SysDictionary                                       true  "ID or dictionary type name"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Find SysDictionary by ID"
// @Router    /sysDictionary/findSysDictionary [get]
func (s *DictionaryApi) FindSysDictionary(c *gin.Context) {
	var dictionary system.SysDictionary
	err := c.ShouldBindQuery(&dictionary)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sysDictionary, err := dictionaryService.GetSysDictionary(dictionary.Type, dictionary.ID, dictionary.Status)
	if err != nil {
		global.GVA_LOG.Error("Dictionary not created or not enabled!", zap.Error(err))
		response.FailWithMessage("Dictionary not created or not enabled", c)
		return
	}
	response.OkWithDetailed(gin.H{"resysDictionary": sysDictionary}, "Query successful", c)
}

// GetSysDictionaryList
// @Tags      SysDictionary
// @Summary   Get SysDictionary list with pagination
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.SysDictionarySearch                                    true  "Dictionary name or type"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Get paginated SysDictionary list, returns list, total, page, page size"
// @Router    /sysDictionary/getSysDictionaryList [get]
func (s *DictionaryApi) GetSysDictionaryList(c *gin.Context) {
	var dictionary request.SysDictionarySearch
	err := c.ShouldBindQuery(&dictionary)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := dictionaryService.GetSysDictionaryInfoList(c, dictionary)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Retrieval failed", c)
		return
	}
	response.OkWithDetailed(list, "Retrieved successfully", c)
}

// ExportSysDictionary
// @Tags      SysDictionary
// @Summary   Export dictionary JSON (including dictionary details)
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     system.SysDictionary                                       true  "Dictionary ID"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Export dictionary JSON"
// @Router    /sysDictionary/exportSysDictionary [get]
func (s *DictionaryApi) ExportSysDictionary(c *gin.Context) {
	var dictionary system.SysDictionary
	err := c.ShouldBindQuery(&dictionary)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if dictionary.ID == 0 {
		response.FailWithMessage("Dictionary ID cannot be empty", c)
		return
	}
	exportData, err := dictionaryService.ExportSysDictionary(dictionary.ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to export!", zap.Error(err))
		response.FailWithMessage("Export failed", c)
		return
	}
	response.OkWithDetailed(exportData, "Exported successfully", c)
}

// ImportSysDictionary
// @Tags      SysDictionary
// @Summary   Import dictionary JSON (including dictionary details)
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.ImportSysDictionaryRequest     true  "Dictionary JSON data"
// @Success   200   {object}  response.Response{msg=string}          "Import dictionary"
// @Router    /sysDictionary/importSysDictionary [post]
func (s *DictionaryApi) ImportSysDictionary(c *gin.Context) {
	var req request.ImportSysDictionaryRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryService.ImportSysDictionary(req.Json)
	if err != nil {
		global.GVA_LOG.Error("Failed to import!", zap.Error(err))
		response.FailWithMessage("Import failed:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Imported successfully", c)
}
