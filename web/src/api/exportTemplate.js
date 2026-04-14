import service from '@/utils/request'

// @Tags SysExportTemplate
// @Summary Create export template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysExportTemplate true "Create export template"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Created"}"
// @Router /sysExportTemplate/createSysExportTemplate [post]
export const createSysExportTemplate = (data) => {
  return service({
    url: '/sysExportTemplate/createSysExportTemplate',
    method: 'post',
    data
  })
}

// @Tags SysExportTemplate
// @Summary Delete export template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysExportTemplate true "Delete export template"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysExportTemplate/deleteSysExportTemplate [delete]
export const deleteSysExportTemplate = (data) => {
  return service({
    url: '/sysExportTemplate/deleteSysExportTemplate',
    method: 'delete',
    data
  })
}

// @Tags SysExportTemplate
// @Summary Delete export templates in bulk
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "Bulk delete export templates"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysExportTemplate/deleteSysExportTemplate [delete]
export const deleteSysExportTemplateByIds = (data) => {
  return service({
    url: '/sysExportTemplate/deleteSysExportTemplateByIds',
    method: 'delete',
    data
  })
}

// @Tags SysExportTemplate
// @Summary Update export template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysExportTemplate true "Update export template"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated"}"
// @Router /sysExportTemplate/updateSysExportTemplate [put]
export const updateSysExportTemplate = (data) => {
  return service({
    url: '/sysExportTemplate/updateSysExportTemplate',
    method: 'put',
    data
  })
}

// @Tags SysExportTemplate
// @Summary Find export template by id
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.SysExportTemplate true "Find export template by id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysExportTemplate/findSysExportTemplate [get]
export const findSysExportTemplate = (params) => {
  return service({
    url: '/sysExportTemplate/findSysExportTemplate',
    method: 'get',
    params
  })
}

// @Tags SysExportTemplate
// @Summary Paginated export template list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "Paginated export template list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysExportTemplate/getSysExportTemplateList [get]
export const getSysExportTemplateList = (params) => {
  return service({
    url: '/sysExportTemplate/getSysExportTemplateList',
    method: 'get',
    params
  })
}


// ExportExcel export token
// @Tags SysExportTemplate
// @Summary Export excel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/exportExcel [get]
export const exportExcel = (params) => {
  return service({
    url: '/sysExportTemplate/exportExcel',
    method: 'get',
    params
  })
}

// ExportTemplate export template
// @Tags SysExportTemplate
// @Summary Export template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/exportTemplate [get]
export const exportTemplate = (params) => {
  return service({
    url: '/sysExportTemplate/exportTemplate',
    method: 'get',
    params
  })
}

// PreviewSQL Preview generated SQL
// @Tags SysExportTemplate
// @Summary Preview generated SQL
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/previewSQL [get]
// @Param templateID query string true  "Export template ID"
// @Param params     query string false "Encoded query param string (see ExportExcel component)"
export const previewSQL = (params) => {
  return service({
    url: '/sysExportTemplate/previewSQL',
    method: 'get',
    params
  })
}
