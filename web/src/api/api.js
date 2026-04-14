import service from '@/utils/request'

// @Tags api
// @Summary Paginated role list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "Paginated list"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"OK"}"
// @Router /api/getApiList [post]
// {
//  page     int
//	pageSize int
// }
export const getApiList = (data) => {
  return service({
    url: '/api/getApiList',
    method: 'post',
    data
  })
}

// @Tags Api
// @Summary Create base API
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateApiParams true "Create API"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"OK"}"
// @Router /api/createApi [post]
export const createApi = (data) => {
  return service({
    url: '/api/createApi',
    method: 'post',
    data
  })
}

// @Tags menu
// @Summary Get by ID
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.GetById true "Get by ID"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"OK"}"
// @Router /menu/getApiById [post]
export const getApiById = (data) => {
  return service({
    url: '/api/getApiById',
    method: 'post',
    data
  })
}

// @Tags Api
// @Summary Update API
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateApiParams true "Update API"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"Updated"}"
// @Router /api/updateApi [post]
export const updateApi = (data) => {
  return service({
    url: '/api/updateApi',
    method: 'post',
    data
  })
}

// @Tags Api
// @Summary Update API
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateApiParams true "Update API"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"Updated"}"
// @Router /api/setAuthApi [post]
export const setAuthApi = (data) => {
  return service({
    url: '/api/setAuthApi',
    method: 'post',
    data
  })
}

// @Tags Api
// @Summary List all APIs (no pagination)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} json "{"success":true,"data":{},"msg":"OK"}"
// @Router /api/getAllApis [post]
export const getAllApis = (data) => {
  return service({
    url: '/api/getAllApis',
    method: 'post',
    data
  })
}

// @Tags Api
// @Summary Delete API
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Api true "Delete API"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /api/deleteApi [post]
export const deleteApi = (data) => {
  return service({
    url: '/api/deleteApi',
    method: 'post',
    data
  })
}

// @Tags SysApi
// @Summary Delete selected APIs
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /api/deleteApisByIds [delete]
export const deleteApisByIds = (data) => {
  return service({
    url: '/api/deleteApisByIds',
    method: 'delete',
    data
  })
}

// FreshCasbin
// @Tags      SysApi
// @Summary   Refresh Casbin cache
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{msg=string}  "Refreshed"
// @Router    /api/freshCasbin [get]
export const freshCasbin = () => {
  return service({
    url: '/api/freshCasbin',
    method: 'get'
  })
}

export const syncApi = () => {
  return service({
    url: '/api/syncApi',
    method: 'get'
  })
}

export const getApiGroups = () => {
  return service({
    url: '/api/getApiGroups',
    method: 'get'
  })
}

export const ignoreApi = (data) => {
  return service({
    url: '/api/ignoreApi',
    method: 'post',
    data
  })
}

export const enterSyncApi = (data) => {
  return service({
    url: '/api/enterSyncApi',
    method: 'post',
    data
  })
}

/**
 * Get authority IDs that have access to a given API
 * @param {string} path API path
 * @param {string} method HTTP method
 * @returns {Promise<number[]>} authority ID list
 */
export const getApiRoles = (path, method) => {
  return service({
    url: '/api/getApiRoles',
    method: 'get',
    params: { path, method }
  })
}

/**
 * Replace the authority list bound to an API
 * @param {Object} data
 * @param {string} data.path API path
 * @param {string} data.method HTTP method
 * @param {number[]} data.authorityIds authority ID list
 * @returns {Promise}
 */
export const setApiRoles = (data) => {
  return service({
    url: '/api/setApiRoles',
    method: 'post',
    data
  })
}
