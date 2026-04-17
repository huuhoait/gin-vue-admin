package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	systemRes "github.com/huuhoait/gin-vue-admin/server/model/system/response"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysVersionApi struct{}

// buildMenuTree builds a menu tree structure
func buildMenuTree(menus []system.SysBaseMenu) []system.SysBaseMenu {
	// Create menu map
	menuMap := make(map[uint]*system.SysBaseMenu)
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	// Build tree structure
	var rootMenus []system.SysBaseMenu
	for _, menu := range menus {
		if menu.ParentId == 0 {
			// Root menu
			menuData := convertMenuToStruct(menu, menuMap)
			rootMenus = append(rootMenus, menuData)
		}
	}

	// Sort root menus by sort field
	sort.Slice(rootMenus, func(i, j int) bool {
		return rootMenus[i].Sort < rootMenus[j].Sort
	})

	return rootMenus
}

// convertMenuToStruct converts a menu to a struct and recursively processes child menus
func convertMenuToStruct(menu system.SysBaseMenu, menuMap map[uint]*system.SysBaseMenu) system.SysBaseMenu {
	result := system.SysBaseMenu{
		Path:      menu.Path,
		Name:      menu.Name,
		Hidden:    menu.Hidden,
		Component: menu.Component,
		Sort:      menu.Sort,
		Meta:      menu.Meta,
	}

	// Clean and copy parameter data
	if len(menu.Parameters) > 0 {
		cleanParameters := make([]system.SysBaseMenuParameter, 0, len(menu.Parameters))
		for _, param := range menu.Parameters {
			cleanParam := system.SysBaseMenuParameter{
				Type:  param.Type,
				Key:   param.Key,
				Value: param.Value,
				// Do not copy ID, CreatedAt, UpdatedAt, SysBaseMenuID
			}
			cleanParameters = append(cleanParameters, cleanParam)
		}
		result.Parameters = cleanParameters
	}

	// Clean and copy menu button data
	if len(menu.MenuBtn) > 0 {
		cleanMenuBtns := make([]system.SysBaseMenuBtn, 0, len(menu.MenuBtn))
		for _, btn := range menu.MenuBtn {
			cleanBtn := system.SysBaseMenuBtn{
				Name: btn.Name,
				Desc: btn.Desc,
				// Do not copy ID, CreatedAt, UpdatedAt, SysBaseMenuID
			}
			cleanMenuBtns = append(cleanMenuBtns, cleanBtn)
		}
		result.MenuBtn = cleanMenuBtns
	}

	// Find and process child menus
	var children []system.SysBaseMenu
	for _, childMenu := range menuMap {
		if childMenu.ParentId == menu.ID {
			childData := convertMenuToStruct(*childMenu, menuMap)
			children = append(children, childData)
		}
	}

	// Sort child menus by sort field
	if len(children) > 0 {
		sort.Slice(children, func(i, j int) bool {
			return children[i].Sort < children[j].Sort
		})
		result.Children = children
	}

	return result
}

// DeleteSysVersion Delete version management record
// @Tags SysVersion
// @Summary Delete version management record
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysVersion true "Delete version management record"
// @Success 200 {object} response.Response{msg=string} "Deleted successfully"
// @Router /sysVersion/deleteSysVersion [delete]
func (sysVersionApi *SysVersionApi) DeleteSysVersion(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := sysVersionService.DeleteSysVersion(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// DeleteSysVersionByIds Batch delete version management records
// @Tags SysVersion
// @Summary Batch delete version management records
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "Batch deleted successfully"
// @Router /sysVersion/deleteSysVersionByIds [delete]
func (sysVersionApi *SysVersionApi) DeleteSysVersionByIds(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := sysVersionService.DeleteSysVersionByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("Failed to batch delete!", zap.Error(err))
		response.FailWithMessage("Batch deletion failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Batch deleted successfully", c)
}

// FindSysVersion Find version management record by ID
// @Tags SysVersion
// @Summary Find version management record by ID
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "Find version management record by ID"
// @Success 200 {object} response.Response{data=system.SysVersion,msg=string} "Query successful"
// @Router /sysVersion/findSysVersion [get]
func (sysVersionApi *SysVersionApi) FindSysVersion(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	resysVersion, err := sysVersionService.GetSysVersion(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to query!", zap.Error(err))
		response.FailWithMessage("Query failed: "+err.Error(), c)
		return
	}
	response.OkWithData(resysVersion, c)
}

// GetSysVersionList Get version management list with pagination
// @Tags SysVersion
// @Summary Get version management list with pagination
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysVersionSearch true "Get version management list with pagination"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "Retrieved successfully"
// @Router /sysVersion/getSysVersionList [get]
func (sysVersionApi *SysVersionApi) GetSysVersionList(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	var pageInfo systemReq.SysVersionSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysVersionService.GetSysVersionInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Retrieved successfully", c)
}

// GetSysVersionPublic Public version management API (no authentication required)
// @Tags SysVersion
// @Summary Public version management API (no authentication required)
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "Retrieved successfully"
// @Router /sysVersion/getSysVersionPublic [get]
func (sysVersionApi *SysVersionApi) GetSysVersionPublic(c *gin.Context) {
	// Create business context
	ctx := c.Request.Context()

	// This API does not require authentication
	// Example returns a fixed message; typically used for client-side services, implement your own business logic
	sysVersionService.GetSysVersionPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "Public version management API info (no authentication required)",
	}, "Retrieved successfully", c)
}

// ExportVersion Create release version data
// @Tags SysVersion
// @Summary Create release version data
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body systemReq.ExportVersionRequest true "Create release version data"
// @Success 200 {object} response.Response{msg=string} "Created successfully"
// @Router /sysVersion/exportVersion [post]
func (sysVersionApi *SysVersionApi) ExportVersion(c *gin.Context) {
	ctx := c.Request.Context()

	var req systemReq.ExportVersionRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Get selected menu data
	var menuData []system.SysBaseMenu
	if len(req.MenuIds) > 0 {
		menuData, err = sysVersionService.GetMenusByIds(ctx, req.MenuIds)
		if err != nil {
			global.GVA_LOG.Error("Failed to get menu data!", zap.Error(err))
			response.FailWithMessage("Failed to get menu data: "+err.Error(), c)
			return
		}
	}

	// Get selected API data
	var apiData []system.SysApi
	if len(req.ApiIds) > 0 {
		apiData, err = sysVersionService.GetApisByIds(ctx, req.ApiIds)
		if err != nil {
			global.GVA_LOG.Error("Failed to get API data!", zap.Error(err))
			response.FailWithMessage("Failed to get API data: "+err.Error(), c)
			return
		}
	}

	// Get selected dictionary data
	var dictData []system.SysDictionary
	if len(req.DictIds) > 0 {
		dictData, err = sysVersionService.GetDictionariesByIds(ctx, req.DictIds)
		if err != nil {
			global.GVA_LOG.Error("Failed to get dictionary data!", zap.Error(err))
			response.FailWithMessage("Failed to get dictionary data: "+err.Error(), c)
			return
		}
	}

	// Process menu data, build recursive children structure
	processedMenus := buildMenuTree(menuData)

	// Process API data, clear ID and timestamp fields
	processedApis := make([]system.SysApi, 0, len(apiData))
	for _, api := range apiData {
		cleanApi := system.SysApi{
			Path:        api.Path,
			Description: api.Description,
			ApiGroup:    api.ApiGroup,
			Method:      api.Method,
		}
		processedApis = append(processedApis, cleanApi)
	}

	// Process dictionary data, clear ID and timestamp fields, include dictionary details
	processedDicts := make([]system.SysDictionary, 0, len(dictData))
	for _, dict := range dictData {
		cleanDict := system.SysDictionary{
			Name:   dict.Name,
			Type:   dict.Type,
			Status: dict.Status,
			Desc:   dict.Desc,
		}
		
		// Process dictionary detail data, clear ID and timestamp fields
		cleanDetails := make([]system.SysDictionaryDetail, 0, len(dict.SysDictionaryDetails))
		for _, detail := range dict.SysDictionaryDetails {
			cleanDetail := system.SysDictionaryDetail{
				Label:  detail.Label,
				Value:  detail.Value,
				Extend: detail.Extend,
				Status: detail.Status,
				Sort:   detail.Sort,
				// Do not copy ID, CreatedAt, UpdatedAt, SysDictionaryID
			}
			cleanDetails = append(cleanDetails, cleanDetail)
		}
		cleanDict.SysDictionaryDetails = cleanDetails
		
		processedDicts = append(processedDicts, cleanDict)
	}

	// Build export data
	exportData := systemRes.ExportVersionResponse{
		Version: systemReq.VersionInfo{
			Name:        req.VersionName,
			Code:        req.VersionCode,
			Description: req.Description,
			ExportTime:  time.Now().Format("2006-01-02 15:04:05"),
		},
		Menus:        processedMenus,
		Apis:         processedApis,
		Dictionaries: processedDicts,
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		global.GVA_LOG.Error("Failed to serialize JSON!", zap.Error(err))
		response.FailWithMessage("JSON serialization failed: "+err.Error(), c)
		return
	}

	// Save version record
	version := system.SysVersion{
		VersionName: utils.Pointer(req.VersionName),
		VersionCode: utils.Pointer(req.VersionCode),
		Description: utils.Pointer(req.Description),
		VersionData: utils.Pointer(string(jsonData)),
	}

	err = sysVersionService.CreateSysVersion(ctx, &version)
	if err != nil {
		global.GVA_LOG.Error("Failed to save version record!", zap.Error(err))
		response.FailWithMessage("Failed to save version record: "+err.Error(), c)
		return
	}

	response.OkWithMessage("Release version created successfully", c)
}

// DownloadVersionJson Download version JSON data
// @Tags SysVersion
// @Summary Download version JSON data
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "Version ID"
// @Success 200 {object} response.Response{data=object,msg=string} "Downloaded successfully"
// @Router /sysVersion/downloadVersionJson [get]
func (sysVersionApi *SysVersionApi) DownloadVersionJson(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("Version ID cannot be empty", c)
		return
	}

	// Get version record
	version, err := sysVersionService.GetSysVersion(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to get version record!", zap.Error(err))
		response.FailWithMessage("Failed to get version record: "+err.Error(), c)
		return
	}

	// Build JSON data
	var jsonData []byte
	if version.VersionData != nil && *version.VersionData != "" {
		jsonData = []byte(*version.VersionData)
	} else {
		// If no stored JSON data, build a basic structure
		basicData := systemRes.ExportVersionResponse{
			Version: systemReq.VersionInfo{
				Name:        *version.VersionName,
				Code:        *version.VersionCode,
				Description: *version.Description,
				ExportTime:  version.CreatedAt.Format("2006-01-02 15:04:05"),
			},
			Menus: []system.SysBaseMenu{},
			Apis:  []system.SysApi{},
		}
		jsonData, _ = json.MarshalIndent(basicData, "", "  ")
	}

	// Set download response headers
	filename := fmt.Sprintf("version_%s_%s.json", *version.VersionCode, time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Length", strconv.Itoa(len(jsonData)))

	c.Data(http.StatusOK, "application/json", jsonData)
}

// ImportVersion Import version data
// @Tags SysVersion
// @Summary Import version data
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body systemReq.ImportVersionRequest true "Version JSON data"
// @Success 200 {object} response.Response{msg=string} "Imported successfully"
// @Router /sysVersion/importVersion [post]
func (sysVersionApi *SysVersionApi) ImportVersion(c *gin.Context) {
	ctx := c.Request.Context()

	// Get JSON data
	var importData systemReq.ImportVersionRequest
	err := c.ShouldBindJSON(&importData)
	if err != nil {
		response.FailWithMessage("Failed to parse JSON data: "+err.Error(), c)
		return
	}

	// Validate data format
	if importData.VersionInfo.Name == "" || importData.VersionInfo.Code == "" {
		response.FailWithMessage("Invalid version info format", c)
		return
	}

	// Import menu data
	if len(importData.ExportMenu) > 0 {
		if err := sysVersionService.ImportMenus(ctx, importData.ExportMenu); err != nil {
			global.GVA_LOG.Error("Failed to import menus!", zap.Error(err))
			response.FailWithMessage("Failed to import menus: "+err.Error(), c)
			return
		}
	}

	// Import API data
	if len(importData.ExportApi) > 0 {
		if err := sysVersionService.ImportApis(importData.ExportApi); err != nil {
			global.GVA_LOG.Error("Failed to import APIs!", zap.Error(err))
			response.FailWithMessage("Failed to import APIs: "+err.Error(), c)
			return
		}
	}

	// Import dictionary data
	if len(importData.ExportDictionary) > 0 {
		if err := sysVersionService.ImportDictionaries(importData.ExportDictionary); err != nil {
			global.GVA_LOG.Error("Failed to import dictionaries!", zap.Error(err))
			response.FailWithMessage("Failed to import dictionaries: "+err.Error(), c)
			return
		}
	}

	// Create import record
	jsonData, _ := json.Marshal(importData)
	version := system.SysVersion{
		VersionName: utils.Pointer(importData.VersionInfo.Name),
		VersionCode: utils.Pointer(fmt.Sprintf("%s_imported_%s", importData.VersionInfo.Code, time.Now().Format("20060102150405"))),
		Description: utils.Pointer(fmt.Sprintf("Imported version: %s", importData.VersionInfo.Description)),
		VersionData: utils.Pointer(string(jsonData)),
	}

	err = sysVersionService.CreateSysVersion(ctx, &version)
	if err != nil {
		global.GVA_LOG.Error("Failed to save import record!", zap.Error(err))
		// Do not return error here since data has already been imported successfully
	}

	response.OkWithMessage("Imported successfully", c)
}
