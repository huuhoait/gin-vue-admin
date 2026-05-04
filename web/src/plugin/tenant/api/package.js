import service from '@/utils/request'

/**
 * List tenant packages with pagination and optional filters.
 * @param {Object} params
 * @param {number} params.page Page number (1-based)
 * @param {number} params.pageSize Page size
 * @param {string} [params.keyword] Filter by code or name
 * @param {boolean} [params.enabled] Filter by enabled flag
 * @returns {Promise} envelope { code, data: { list, total, page, pageSize }, msg }
 */
export const listTenantPackages = (params) =>
  service({ url: '/tenantPackage/list', method: 'get', params })

/**
 * Find a single tenant package by ID.
 * @param {Object} params
 * @param {number} params.id Package ID
 * @returns {Promise}
 */
export const findTenantPackage = (params) =>
  service({ url: '/tenantPackage/find', method: 'get', params })

/**
 * Create a tenant package.
 * @param {Object} data
 * @param {string} data.code Stable identifier (immutable after create)
 * @param {string} data.name Display name
 * @param {string} [data.description]
 * @param {number[]} [data.menuIDs] SysBaseMenu IDs
 * @param {number[]} [data.apiIDs] SysApi IDs
 * @returns {Promise}
 */
export const createTenantPackage = (data) =>
  service({ url: '/tenantPackage/create', method: 'post', data })

/**
 * Update a tenant package. Code is immutable; pass null for menuIDs/apiIDs
 * to leave them unchanged or [] to clear them.
 * @param {Object} data
 * @param {number} data.id Package ID
 * @param {string} [data.name]
 * @param {string} [data.description]
 * @param {number[]|null} [data.menuIDs]
 * @param {number[]|null} [data.apiIDs]
 * @param {boolean} [data.enabled]
 * @returns {Promise}
 */
export const updateTenantPackage = (data) =>
  service({ url: '/tenantPackage/update', method: 'put', data })

/**
 * Delete a tenant package by ID.
 * @param {Object} params
 * @param {number} params.id
 * @returns {Promise}
 */
export const deleteTenantPackage = (params) =>
  service({ url: '/tenantPackage/delete', method: 'delete', params })
