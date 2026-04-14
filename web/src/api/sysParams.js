import service from '@/utils/request'
// @Tags SysParams
// @Summary Create param
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysParams true "Create param"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Created"}"
// @Router /sysParams/createSysParams [post]
export const createSysParams = (data) => {
  return service({
    url: '/sysParams/createSysParams',
    method: 'post',
    data
  })
}

// @Tags SysParams
// @Summary Delete param
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysParams true "Delete param"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysParams/deleteSysParams [delete]
export const deleteSysParams = (params) => {
  return service({
    url: '/sysParams/deleteSysParams',
    method: 'delete',
    params
  })
}

// @Tags SysParams
// @Summary Delete params in bulk
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "Bulk delete params"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysParams/deleteSysParams [delete]
export const deleteSysParamsByIds = (params) => {
  return service({
    url: '/sysParams/deleteSysParamsByIds',
    method: 'delete',
    params
  })
}

// @Tags SysParams
// @Summary Update param
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysParams true "Update param"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated"}"
// @Router /sysParams/updateSysParams [put]
export const updateSysParams = (data) => {
  return service({
    url: '/sysParams/updateSysParams',
    method: 'put',
    data
  })
}

// @Tags SysParams
// @Summary Find param by id
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.SysParams true "Find param by id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysParams/findSysParams [get]
export const findSysParams = (params) => {
  return service({
    url: '/sysParams/findSysParams',
    method: 'get',
    params
  })
}

// @Tags SysParams
// @Summary Paginated param list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "Paginated param list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysParams/getSysParamsList [get]
export const getSysParamsList = (params) => {
  return service({
    url: '/sysParams/getSysParamsList',
    method: 'get',
    params
  })
}

// @Tags SysParams
// @Summary Public param endpoint (no auth required)
// @accept application/json
// @Produce application/json
// @Param data query systemReq.SysParamsSearch true "Paginated param list"
// @Success 200 {object} response.Response{data=object,msg=string} "OK"
// @Router /sysParams/getSysParam [get]
export const getSysParam = (params) => {
  return service({
    url: '/sysParams/getSysParam',
    method: 'get',
    params
  })
}
