package system

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// For one-time token storage
var (
	exportTokenCache      = make(map[string]interface{})
	exportTokenExpiration = make(map[string]time.Time)
	tokenMutex            sync.RWMutex
)

// Clean up expired tokens every five minutes
func cleanupExpiredTokens() {
	for {
		time.Sleep(5 * time.Minute)
		tokenMutex.Lock()
		now := time.Now()
		for token, expiry := range exportTokenExpiration {
			if now.After(expiry) {
				delete(exportTokenCache, token)
				delete(exportTokenExpiration, token)
			}
		}
		tokenMutex.Unlock()
	}
}

func init() {
	go cleanupExpiredTokens()
}

type SysExportTemplateApi struct {
}

var sysExportTemplateService = service.ServiceGroupApp.SystemServiceGroup.SysExportTemplateService

// PreviewSQL Preview the final generated SQL
// @Tags     SysExportTemplate
// @Summary  Preview the final generated SQL (does not execute query, only returns SQL string)
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    templateID query string true  "Export template ID"
// @Param    params     query string false "Encoded query parameter string, refer to ExportExcel component"
// @Success  200  {object}  response.Response{data=map[string]string} "Retrieved successfully"
// @Router   /sysExportTemplate/previewSQL [get]
func (sysExportTemplateApi *SysExportTemplateApi) PreviewSQL(c *gin.Context) {
    templateID := c.Query("templateID")
    if templateID == "" {
        response.FailWithMessage("Template ID cannot be empty", c)
        return
    }

    // Reuse the export API's parameter organization: use URL Query, where params is the internally encoded query string
    queryParams := c.Request.URL.Query()

    if sqlPreview, err := sysExportTemplateService.PreviewSQL(templateID, queryParams); err != nil {
        global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
        response.FailWithMessage("Failed to retrieve", c)
    } else {
        response.OkWithData(gin.H{"sql": sqlPreview}, c)
    }
}

// CreateSysExportTemplate Create export template
// @Tags SysExportTemplate
// @Summary Create export template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysExportTemplate true "Create export template"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Created successfully"}"
// @Router /sysExportTemplate/createSysExportTemplate [post]
func (sysExportTemplateApi *SysExportTemplateApi) CreateSysExportTemplate(c *gin.Context) {
	var sysExportTemplate system.SysExportTemplate
	err := c.ShouldBindJSON(&sysExportTemplate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Name": {utils.NotEmpty()},
	}
	if err := utils.Verify(sysExportTemplate, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysExportTemplateService.CreateSysExportTemplate(&sysExportTemplate); err != nil {
		global.GVA_LOG.Error("Failed to create!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
	} else {
		response.OkWithMessage("Created successfully", c)
	}
}

// DeleteSysExportTemplate Delete export template
// @Tags SysExportTemplate
// @Summary Delete export template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysExportTemplate true "Delete export template"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted successfully"}"
// @Router /sysExportTemplate/deleteSysExportTemplate [delete]
func (sysExportTemplateApi *SysExportTemplateApi) DeleteSysExportTemplate(c *gin.Context) {
	var sysExportTemplate system.SysExportTemplate
	err := c.ShouldBindJSON(&sysExportTemplate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysExportTemplateService.DeleteSysExportTemplate(sysExportTemplate); err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed", c)
	} else {
		response.OkWithMessage("Deleted successfully", c)
	}
}

// DeleteSysExportTemplateByIds Batch delete export templates
// @Tags SysExportTemplate
// @Summary Batch delete export templates
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "Batch delete export templates"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Batch deleted successfully"}"
// @Router /sysExportTemplate/deleteSysExportTemplateByIds [delete]
func (sysExportTemplateApi *SysExportTemplateApi) DeleteSysExportTemplateByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysExportTemplateService.DeleteSysExportTemplateByIds(IDS); err != nil {
		global.GVA_LOG.Error("Failed to batch delete!", zap.Error(err))
		response.FailWithMessage("Batch deletion failed", c)
	} else {
		response.OkWithMessage("Batch deleted successfully", c)
	}
}

// UpdateSysExportTemplate Update export template
// @Tags SysExportTemplate
// @Summary Update export template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysExportTemplate true "Update export template"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated successfully"}"
// @Router /sysExportTemplate/updateSysExportTemplate [put]
func (sysExportTemplateApi *SysExportTemplateApi) UpdateSysExportTemplate(c *gin.Context) {
	var sysExportTemplate system.SysExportTemplate
	err := c.ShouldBindJSON(&sysExportTemplate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Name": {utils.NotEmpty()},
	}
	if err := utils.Verify(sysExportTemplate, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysExportTemplateService.UpdateSysExportTemplate(sysExportTemplate); err != nil {
		global.GVA_LOG.Error("Failed to update!", zap.Error(err))
		response.FailWithMessage("Update failed", c)
	} else {
		response.OkWithMessage("Updated successfully", c)
	}
}

// FindSysExportTemplate Find export template by ID
// @Tags SysExportTemplate
// @Summary Find export template by ID
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query system.SysExportTemplate true "Find export template by ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Query successful"}"
// @Router /sysExportTemplate/findSysExportTemplate [get]
func (sysExportTemplateApi *SysExportTemplateApi) FindSysExportTemplate(c *gin.Context) {
	var sysExportTemplate system.SysExportTemplate
	err := c.ShouldBindQuery(&sysExportTemplate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if resysExportTemplate, err := sysExportTemplateService.GetSysExportTemplate(sysExportTemplate.ID); err != nil {
		global.GVA_LOG.Error("Failed to query!", zap.Error(err))
		response.FailWithMessage("Query failed", c)
	} else {
		response.OkWithData(gin.H{"resysExportTemplate": resysExportTemplate}, c)
	}
}

// GetSysExportTemplateList Get export template list with pagination
// @Tags SysExportTemplate
// @Summary Get export template list with pagination
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query systemReq.SysExportTemplateSearch true "Get export template list with pagination"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Retrieved successfully"}"
// @Router /sysExportTemplate/getSysExportTemplateList [get]
func (sysExportTemplateApi *SysExportTemplateApi) GetSysExportTemplateList(c *gin.Context) {
	var pageInfo systemReq.SysExportTemplateSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := sysExportTemplateService.GetSysExportTemplateInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "Retrieved successfully", c)
	}
}

// ExportExcel Export spreadsheet token
// @Tags SysExportTemplate
// @Summary Export spreadsheet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/exportExcel [get]
func (sysExportTemplateApi *SysExportTemplateApi) ExportExcel(c *gin.Context) {
	templateID := c.Query("templateID")
	if templateID == "" {
		response.FailWithMessage("Template ID cannot be empty", c)
		return
	}

	queryParams := c.Request.URL.Query()

	// Create one-time token
	token := utils.RandomString(32) // Random 32 characters

	// Record current request parameters
	exportParams := map[string]interface{}{
		"templateID":  templateID,
		"queryParams": queryParams,
	}

	// Store parameters to complete authentication
	tokenMutex.Lock()
	exportTokenCache[token] = exportParams
	exportTokenExpiration[token] = time.Now().Add(30 * time.Minute)
	tokenMutex.Unlock()

	// Generate one-time link
	exportUrl := fmt.Sprintf("/sysExportTemplate/exportExcelByToken?token=%s", token)
	response.OkWithData(exportUrl, c)
}

// ExportExcelByToken Export spreadsheet by token
// @Tags ExportExcelByToken
// @Summary Export spreadsheet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/exportExcelByToken [get]
func (sysExportTemplateApi *SysExportTemplateApi) ExportExcelByToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.FailWithMessage("Export token cannot be empty", c)
		return
	}

	// Get token and remove from cache
	tokenMutex.RLock()
	exportParamsRaw, exists := exportTokenCache[token]
	expiry, _ := exportTokenExpiration[token]
	tokenMutex.RUnlock()

	if !exists || time.Now().After(expiry) {
		global.GVA_LOG.Error("Export token is invalid or expired!")
		response.FailWithMessage("Export token is invalid or expired", c)
		return
	}

	// Get parameters from token
	exportParams, ok := exportParamsRaw.(map[string]interface{})
	if !ok {
		global.GVA_LOG.Error("Failed to parse export parameters!")
		response.FailWithMessage("Failed to parse export parameters", c)
		return
	}

	// Get export parameters
	templateID := exportParams["templateID"].(string)
	queryParams := exportParams["queryParams"].(url.Values)

	// Clean up one-time token
	tokenMutex.Lock()
	delete(exportTokenCache, token)
	delete(exportTokenExpiration, token)
	tokenMutex.Unlock()

	// Export
	if file, name, err := sysExportTemplateService.ExportExcel(templateID, queryParams); err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
	} else {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name+utils.RandomString(6)+".xlsx"))
		c.Header("success", "true")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file.Bytes())
	}
}

// ExportTemplate Export spreadsheet template
// @Tags SysExportTemplate
// @Summary Export spreadsheet template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/exportTemplate [get]
func (sysExportTemplateApi *SysExportTemplateApi) ExportTemplate(c *gin.Context) {
	templateID := c.Query("templateID")
	if templateID == "" {
		response.FailWithMessage("Template ID cannot be empty", c)
		return
	}

	// Create one-time token
	token := utils.RandomString(32) // Random 32 characters

	// Record current request parameters
	exportParams := map[string]interface{}{
		"templateID": templateID,
		"isTemplate": true,
	}

	// Store parameters to complete authentication
	tokenMutex.Lock()
	exportTokenCache[token] = exportParams
	exportTokenExpiration[token] = time.Now().Add(30 * time.Minute)
	tokenMutex.Unlock()

	// Generate one-time link
	exportUrl := fmt.Sprintf("/sysExportTemplate/exportTemplateByToken?token=%s", token)
	response.OkWithData(exportUrl, c)
}

// ExportTemplateByToken Export spreadsheet template by token
// @Tags ExportTemplateByToken
// @Summary Export spreadsheet template by token
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/exportTemplateByToken [get]
func (sysExportTemplateApi *SysExportTemplateApi) ExportTemplateByToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.FailWithMessage("Export token cannot be empty", c)
		return
	}

	// Get token and remove from cache
	tokenMutex.RLock()
	exportParamsRaw, exists := exportTokenCache[token]
	expiry, _ := exportTokenExpiration[token]
	tokenMutex.RUnlock()

	if !exists || time.Now().After(expiry) {
		global.GVA_LOG.Error("Export token is invalid or expired!")
		response.FailWithMessage("Export token is invalid or expired", c)
		return
	}

	// Get parameters from token
	exportParams, ok := exportParamsRaw.(map[string]interface{})
	if !ok {
		global.GVA_LOG.Error("Failed to parse export parameters!")
		response.FailWithMessage("Failed to parse export parameters", c)
		return
	}

	// Check if this is a template export
	isTemplate, _ := exportParams["isTemplate"].(bool)
	if !isTemplate {
		global.GVA_LOG.Error("Invalid token type!")
		response.FailWithMessage("Invalid token type", c)
		return
	}

	// Get export parameters
	templateID := exportParams["templateID"].(string)

	// Clean up one-time token
	tokenMutex.Lock()
	delete(exportTokenCache, token)
	delete(exportTokenExpiration, token)
	tokenMutex.Unlock()

	// Export template
	if file, name, err := sysExportTemplateService.ExportTemplate(templateID); err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
	} else {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name+"_template.xlsx"))
		c.Header("success", "true")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file.Bytes())
	}
}

// ImportExcel Import spreadsheet
// @Tags SysImportTemplate
// @Summary Import spreadsheet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/importExcel [post]
func (sysExportTemplateApi *SysExportTemplateApi) ImportExcel(c *gin.Context) {
	templateID := c.Query("templateID")
	if templateID == "" {
		response.FailWithMessage("Template ID cannot be empty", c)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("Failed to get file!", zap.Error(err))
		response.FailWithMessage("Failed to get file", c)
		return
	}
	if err := sysExportTemplateService.ImportExcel(templateID, file); err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("Imported successfully", c)
	}
}
