import service from '@/utils/request'

// @Tags SysVersion
// @Summary Delete version
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysVersion true "Delete version"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysVersion/deleteSysVersion [delete]
export const deleteSysVersion = (params) => {
  return service({
    url: '/sysVersion/deleteSysVersion',
    method: 'delete',
    params
  })
}

// @Tags SysVersion
// @Summary Delete versions in bulk
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "Bulk delete versions"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysVersion/deleteSysVersion [delete]
export const deleteSysVersionByIds = (params) => {
  return service({
    url: '/sysVersion/deleteSysVersionByIds',
    method: 'delete',
    params
  })
}

// @Tags SysVersion
// @Summary Find version by id
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysVersion true "Find version by id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysVersion/findSysVersion [get]
export const findSysVersion = (params) => {
  return service({
    url: '/sysVersion/findSysVersion',
    method: 'get',
    params
  })
}

// @Tags SysVersion
// @Summary Paginated version list
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "Paginated version list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysVersion/getSysVersionList [get]
export const getSysVersionList = (params) => {
  return service({
    url: '/sysVersion/getSysVersionList',
    method: 'get',
    params
  })
}

// @Tags SysVersion
// @Summary Export versions
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body object true "Export versions"
// @Success 200 {string} string "{\"success\":true,\"data\":{},\"msg\":\"Exported\"}"
// @Router /sysVersion/exportVersion [post]
export const exportVersion = (data) => {
  return service({
    url: '/sysVersion/exportVersion',
    method: 'post',
    data
  })
}

// @Tags SysVersion
// @Summary Download version JSON
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "Version ID"
// @Success 200 {string} string "{\"success\":true,\"data\":{},\"msg\":\"Downloaded\"}"
// @Router /sysVersion/downloadVersionJson [get]
export const downloadVersionJson = (params) => {
  return service({
    url: '/sysVersion/downloadVersionJson',
    method: 'get',
    params,
    responseType: 'blob'
  })
}

// @Tags SysVersion
// @Summary Import versions
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body object true "Version JSON"
// @Success 200 {string} string "{\"success\":true,\"data\":{},\"msg\":\"Imported\"}"
// @Router /sysVersion/importVersion [post]
export const importVersion = (data) => {
  return service({
    url: '/sysVersion/importVersion',
    method: 'post',
    data
  })
}
