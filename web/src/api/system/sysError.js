import service from '@/utils/request'
// @Tags SysError
// @Summary Create error log
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysError true "Create error log"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Created"}"
// @Router /sysError/createSysError [post]
export const createSysError = (data) => {
  return service({
    url: '/sysError/createSysError',
    method: 'post',
    data
  })
}

// @Tags SysError
// @Summary Delete error log
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysError true "Delete error log"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysError/deleteSysError [delete]
export const deleteSysError = (params) => {
  return service({
    url: '/sysError/deleteSysError',
    method: 'delete',
    params
  })
}

// @Tags SysError
// @Summary Delete error logs in bulk
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "Bulk delete error logs"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysError/deleteSysError [delete]
export const deleteSysErrorByIds = (params) => {
  return service({
    url: '/sysError/deleteSysErrorByIds',
    method: 'delete',
    params
  })
}

// @Tags SysError
// @Summary Update error log
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysError true "Update error log"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated"}"
// @Router /sysError/updateSysError [put]
export const updateSysError = (data) => {
  return service({
    url: '/sysError/updateSysError',
    method: 'put',
    data
  })
}

// @Tags SysError
// @Summary Find error log by id
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysError true "Find error log by id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysError/findSysError [get]
export const findSysError = (params) => {
  return service({
    url: '/sysError/findSysError',
    method: 'get',
    params
  })
}

// @Tags SysError
// @Summary Paginated error log list
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "Paginated error log list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysError/getSysErrorList [get]
export const getSysErrorList = (params) => {
  return service({
    url: '/sysError/getSysErrorList',
    method: 'get',
    params
  })
}

// @Tags SysError
// @Summary Public error log endpoint (no auth required)
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysErrorSearch true "Paginated error log list"
// @Success 200 {object} response.Response{data=object,msg=string} "OK"
// @Router /sysError/getSysErrorPublic [get]
export const getSysErrorPublic = () => {
  return service({
    url: '/sysError/getSysErrorPublic',
    method: 'get',
  })
}

// @Tags SysError
// @Summary Trigger error handling (async)
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query string true "Error log ID"
// @Success 200 {string} string "{\"success\":true,\"data\":{},\"msg\":\"Submitted\"}"
// @Router /sysError/getSysErrorSolution [get]
export const getSysErrorSolution = (params) => {
  return service({
    url: '/sysError/getSysErrorSolution',
    method: 'get',
    params
  })
}