import service from '@/utils/request'
// @Tags systrm
// @Summary Get system config
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /system/getSystemConfig [post]
export const getSystemConfig = () => {
  return service({
    url: '/system/getSystemConfig',
    method: 'post'
  })
}

// @Tags system
// @Summary Set system config
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body sysModel.System true
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /system/setSystemConfig [post]
export const setSystemConfig = (data) => {
  return service({
    url: '/system/setSystemConfig',
    method: 'post',
    data
  })
}

// @Tags system
// @Summary Get server state
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /system/getServerInfo [post]
export const getSystemState = () => {
  return service({
    url: '/system/getServerInfo',
    method: 'post',
    donNotShowLoading: true
  })
}

/**
 * Reload service
 * @param data
 * @returns {*}
 */
export const reloadSystem = (data) => {
  return service({
    url: '/system/reloadSystem',
    method: 'post',
    data
  })
}
