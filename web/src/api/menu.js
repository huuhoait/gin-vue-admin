import service from '@/utils/request'
// @Summary User login: get async routes
// @Produce  application/json
// @Param No params required
// @Router /menu/getMenu [post]
export const asyncMenu = () => {
  return service({
    url: '/menu/getMenu',
    method: 'post'
  })
}

// @Summary Get menu list
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /menu/getMenuList [post]
export const getMenuList = (data) => {
  return service({
    url: '/menu/getMenuList',
    method: 'post',
    data
  })
}

// @Summary Add base menu
// @Produce  application/json
// @Param menu Object
// @Router /menu/getMenuList [post]
export const addBaseMenu = (data) => {
  return service({
    url: '/menu/addBaseMenu',
    method: 'post',
    data
  })
}

// @Summary Get base menu tree
// @Produce  application/json
// @Param No params required
// @Router /menu/getBaseMenuTree [post]
export const getBaseMenuTree = () => {
  return service({
    url: '/menu/getBaseMenuTree',
    method: 'post'
  })
}

// @Summary Add menu-authority relation
// @Produce  application/json
// @Param menus Object authorityId string
// @Router /menu/getMenuList [post]
export const addMenuAuthority = (data) => {
  return service({
    url: '/menu/addMenuAuthority',
    method: 'post',
    data
  })
}

// @Summary Get menu-authority relation
// @Produce  application/json
// @Param authorityId string
// @Router /menu/getMenuAuthority [post]
export const getMenuAuthority = (data) => {
  return service({
    url: '/menu/getMenuAuthority',
    method: 'post',
    data
  })
}

// @Summary Delete menu
// @Produce  application/json
// @Param ID float64
// @Router /menu/deleteBaseMenu [post]
export const deleteBaseMenu = (data) => {
  return service({
    url: '/menu/deleteBaseMenu',
    method: 'post',
    data
  })
}

// @Summary Update menu
// @Produce  application/json
// @Param menu Object
// @Router /menu/updateBaseMenu [post]
export const updateBaseMenu = (data) => {
  return service({
    url: '/menu/updateBaseMenu',
    method: 'post',
    data
  })
}

// @Tags menu
// @Summary Get menu by id
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.GetById true "Get menu by id"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"OK"}"
// @Router /menu/getBaseMenuById [post]
export const getBaseMenuById = (data) => {
  return service({
    url: '/menu/getBaseMenuById',
    method: 'post',
    data
  })
}

/**
 * Get authority IDs that have access to a menu
 * @param {number} menuId menu ID
 * @returns {Promise<number[]>} authority ID list
 */
export const getMenuRoles = (menuId) => {
  return service({
    url: '/menu/getMenuRoles',
    method: 'get',
    params: { menuId }
  })
}

/**
 * Replace authorities bound to a menu
 * @param {Object} data
 * @param {number} data.menuId menu ID
 * @param {number[]} data.authorityIds authority ID list
 * @returns {Promise}
 */
export const setMenuRoles = (data) => {
  return service({
    url: '/menu/setMenuRoles',
    method: 'post',
    data
  })
}
