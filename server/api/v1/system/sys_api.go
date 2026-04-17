package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	systemRes "github.com/huuhoait/gin-vue-admin/server/model/system/response"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemApiApi struct{}

// CreateApi
// @Tags      SysApi
// @Summary   Create base API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "API path, API description, API group, method"
// @Success   200   {object}  response.Response{msg=string}  "Create base API"
// @Router    /api/createApi [post]
func (s *SystemApiApi) CreateApi(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.CreateApi(api)
	if err != nil {
		global.GVA_LOG.Error("Failed to create!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// SyncApi
// @Tags      SysApi
// @Summary   Sync APIs
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{msg=string}  "Sync APIs"
// @Router    /api/syncApi [get]
func (s *SystemApiApi) SyncApi(c *gin.Context) {
	newApis, deleteApis, ignoreApis, err := apiService.SyncApi()
	if err != nil {
		global.GVA_LOG.Error("Failed to sync!", zap.Error(err))
		response.FailWithMessage("Sync failed", c)
		return
	}
	response.OkWithData(gin.H{
		"newApis":    newApis,
		"deleteApis": deleteApis,
		"ignoreApis": ignoreApis,
	}, c)
}

// GetApiGroups
// @Tags      SysApi
// @Summary   Get API groups
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{msg=string}  "Get API groups"
// @Router    /api/getApiGroups [get]
func (s *SystemApiApi) GetApiGroups(c *gin.Context) {
	groups, apiGroupMap, err := apiService.GetApiGroups()
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithData(gin.H{
		"groups":      groups,
		"apiGroupMap": apiGroupMap,
	}, c)
}

// IgnoreApi
// @Tags      IgnoreApi
// @Summary   Ignore API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{msg=string}  "Ignore API"
// @Router    /api/ignoreApi [post]
func (s *SystemApiApi) IgnoreApi(c *gin.Context) {
	var ignoreApi system.SysIgnoreApi
	err := c.ShouldBindJSON(&ignoreApi)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.IgnoreApi(ignoreApi)
	if err != nil {
		global.GVA_LOG.Error("Failed to ignore!", zap.Error(err))
		response.FailWithMessage("Failed to ignore", c)
		return
	}
	response.Ok(c)
}

// EnterSyncApi
// @Tags      SysApi
// @Summary   Confirm sync APIs
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{msg=string}  "Confirm sync APIs"
// @Router    /api/enterSyncApi [post]
func (s *SystemApiApi) EnterSyncApi(c *gin.Context) {
	var syncApi systemRes.SysSyncApis
	err := c.ShouldBindJSON(&syncApi)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.EnterSyncApi(syncApi)
	if err != nil {
		global.GVA_LOG.Error("Failed to ignore!", zap.Error(err))
		response.FailWithMessage("Failed to ignore", c)
		return
	}
	response.Ok(c)
}

// DeleteApi
// @Tags      SysApi
// @Summary   Delete API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "Delete API"
// @Router    /api/deleteApi [post]
func (s *SystemApiApi) DeleteApi(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApi(api)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// GetApiList
// @Tags      SysApi
// @Summary   Get API list with pagination
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SearchApiParams                               true  "Get API list with pagination"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Get paginated API list, returns list, total, page, page size"
// @Router    /api/getApiList [post]
func (s *SystemApiApi) GetApiList(c *gin.Context) {
	var pageInfo systemReq.SearchApiParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := apiService.GetAPIInfoList(pageInfo.SysApi, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
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

// GetApiById
// @Tags      SysApi
// @Summary   Get API by ID
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                                   true  "Get API by ID"
// @Success   200   {object}  response.Response{data=systemRes.SysAPIResponse}  "Get API by ID, returns API details"
// @Router    /api/getApiById [post]
func (s *SystemApiApi) GetApiById(c *gin.Context) {
	var idInfo request.GetById
	err := c.ShouldBindJSON(&idInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	api, err := apiService.GetApiById(idInfo.ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(systemRes.SysAPIResponse{Api: api}, "Retrieved successfully", c)
}

// UpdateApi
// @Tags      SysApi
// @Summary   Update base API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "API path, API description, API group, method"
// @Success   200   {object}  response.Response{msg=string}  "Update base API"
// @Router    /api/updateApi [post]
func (s *SystemApiApi) UpdateApi(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.UpdateApi(api)
	if err != nil {
		global.GVA_LOG.Error("Failed to update!", zap.Error(err))
		response.FailWithMessage("Modification failed", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}

// GetAllApis
// @Tags      SysApi
// @Summary   Get all APIs without pagination
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysAPIListResponse,msg=string}  "Get all APIs without pagination, returns API list"
// @Router    /api/getAllApis [post]
func (s *SystemApiApi) GetAllApis(c *gin.Context) {
	authorityID := utils.GetUserAuthorityId(c)
	apis, err := apiService.GetAllApis(authorityID)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(systemRes.SysAPIListResponse{Apis: apis}, "Retrieved successfully", c)
}

// DeleteApisByIds
// @Tags      SysApi
// @Summary   Delete selected APIs
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "Delete selected APIs"
// @Router    /api/deleteApisByIds [delete]
func (s *SystemApiApi) DeleteApisByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApisByIds(ids)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// FreshCasbin
// @Tags      SysApi
// @Summary   Refresh casbin cache
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{msg=string}  "Refreshed successfully"
// @Router    /api/freshCasbin [get]
func (s *SystemApiApi) FreshCasbin(c *gin.Context) {
	err := casbinService.FreshCasbin()
	if err != nil {
		global.GVA_LOG.Error("Failed to refresh!", zap.Error(err))
		response.FailWithMessage("Refresh failed", c)
		return
	}
	response.OkWithMessage("Refreshed successfully", c)
}

// GetApiRoles
// @Tags      SysApi
// @Summary   Get role IDs that have access to specified API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     path    query     string                                                    true  "API path"
// @Param     method  query     string                                                    true  "Request method"
// @Success   200     {object}  response.Response{data=map[string]interface{},msg=string}  "Retrieved successfully"
// @Router    /api/getApiRoles [get]
func (s *SystemApiApi) GetApiRoles(c *gin.Context) {
	path := c.Query("path")
	method := c.Query("method")
	if path == "" || method == "" {
		response.FailWithMessage("API path and request method cannot be empty", c)
		return
	}
	authorityIds, err := casbinService.GetAuthoritiesByApi(path, method)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve"+err.Error(), c)
		return
	}
	if authorityIds == nil {
		authorityIds = []uint{}
	}
	response.OkWithDetailed(authorityIds, "Retrieved successfully", c)
}

// SetApiRoles
// @Tags      SysApi
// @Summary   Fully replace the role list associated with an API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetApiAuthorities    true  "API path, request method, and authority ID list"
// @Success   200   {object}  response.Response{msg=string}  "Set successfully"
// @Router    /api/setApiRoles [post]
func (s *SystemApiApi) SetApiRoles(c *gin.Context) {
	var req systemReq.SetApiAuthorities
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Path == "" || req.Method == "" {
		response.FailWithMessage("API path and request method cannot be empty", c)
		return
	}
	if err := casbinService.SetApiAuthorities(req.Path, req.Method, req.AuthorityIds); err != nil {
		global.GVA_LOG.Error("Failed to set!", zap.Error(err))
		response.FailWithMessage("Failed to set"+err.Error(), c)
		return
	}
	// Refresh casbin cache to apply policies immediately
	_ = casbinService.FreshCasbin()
	response.OkWithMessage("Set successfully", c)
}
