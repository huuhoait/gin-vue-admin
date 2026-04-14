import service from '@/utils/request'
// @Tags authority
// @Summary Update role API permissions
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateAuthorityPatams true "Update role API permissions"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"OK"}"
// @Router /casbin/UpdateCasbin [post]
export const UpdateCasbin = (data) => {
  return service({
    url: '/casbin/updateCasbin',
    method: 'post',
    data
  })
}

// @Tags casbin
// @Summary Get permission list
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateAuthorityPatams true "Get permission list"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"OK"}"
// @Router /casbin/getPolicyPathByAuthorityId [post]
export const getPolicyPathByAuthorityId = (data) => {
  return service({
    url: '/casbin/getPolicyPathByAuthorityId',
    method: 'post',
    data
  })
}
