import service from '@/utils/request'
// @Router /authority/getAuthorityList [post]
export const getAuthorityList = (data) => {
  return service({
    url: '/authority/getAuthorityList',
    method: 'post',
    data
  })
}

// @Summary Delete role
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {authorityId uint} true "Delete role"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"OK"}"
// @Router /authority/deleteAuthority [post]
export const deleteAuthority = (data) => {
  return service({
    url: '/authority/deleteAuthority',
    method: 'post',
    data
  })
}

// @Summary Create role
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateAuthorityPatams true "Create role"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"OK"}"
// @Router /authority/createAuthority [post]
export const createAuthority = (data) => {
  return service({
    url: '/authority/createAuthority',
    method: 'post',
    data
  })
}

// @Tags authority
// @Summary Copy role
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateAuthorityPatams true "Copy role"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"Copied"}"
// @Router /authority/copyAuthority [post]
export const copyAuthority = (data) => {
  return service({
    url: '/authority/copyAuthority',
    method: 'post',
    data
  })
}

// @Summary Set role data authority
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body sysModel.SysAuthority true "Set role data authority"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated"}"
// @Router /authority/setDataAuthority [post]
export const setDataAuthority = (data) => {
  return service({
    url: '/authority/setDataAuthority',
    method: 'post',
    data
  })
}

// @Summary Update role
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysAuthority true "Update role"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"Updated"}"
// @Router /authority/setDataAuthority [post]
export const updateAuthority = (data) => {
  return service({
    url: '/authority/updateAuthority',
    method: 'put',
    data
  })
}

/**
 * Get user IDs that have a given authority
 * @param {number} authorityId authority ID
 * @returns {Promise<number[]>} user ID list
 */
export const getUsersByAuthorityId = (authorityId) => {
  return service({
    url: '/authority/getUsersByAuthority',
    method: 'get',
    params: { authorityId }
  })
}

/**
 * Replace users bound to an authority
 * @param {Object} data
 * @param {number} data.authorityId authority ID
 * @param {number[]} data.userIds user ID list
 * @returns {Promise}
 */
export const setRoleUsers = (data) => {
  return service({
    url: '/authority/setRoleUsers',
    method: 'post',
    data
  })
}
