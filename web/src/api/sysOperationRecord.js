import service from '@/utils/request'
// @Tags SysOperationRecord
// @Summary Delete SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysOperationRecord true "Delete SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysOperationRecord/deleteSysOperationRecord [delete]
export const deleteSysOperationRecord = (data) => {
  return service({
    url: '/sysOperationRecord/deleteSysOperationRecord',
    method: 'delete',
    data
  })
}

// @Tags SysOperationRecord
// @Summary Delete SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "Delete SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Deleted"}"
// @Router /sysOperationRecord/deleteSysOperationRecord [delete]
export const deleteSysOperationRecordByIds = (data) => {
  return service({
    url: '/sysOperationRecord/deleteSysOperationRecordByIds',
    method: 'delete',
    data
  })
}

// @Tags SysOperationRecord
// @Summary Paginated SysOperationRecord list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "Paginated SysOperationRecord list"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /sysOperationRecord/getSysOperationRecordList [get]
export const getSysOperationRecordList = (params) => {
  return service({
    url: '/sysOperationRecord/getSysOperationRecordList',
    method: 'get',
    params
  })
}
