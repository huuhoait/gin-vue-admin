import service from '@/utils/request'

export const listTenants = (params) =>
  service({ url: '/tenant/list', method: 'get', params })

export const findTenant = (params) =>
  service({ url: '/tenant/find', method: 'get', params })

export const createTenant = (data) =>
  service({ url: '/tenant/create', method: 'post', data })

export const updateTenant = (data) =>
  service({ url: '/tenant/update', method: 'put', data })

export const deleteTenant = (params) =>
  service({ url: '/tenant/delete', method: 'delete', params })

export const assignUser = (data) =>
  service({ url: '/tenantMembership/assign', method: 'post', data })

export const unassignUser = (params) =>
  service({ url: '/tenantMembership/unassign', method: 'delete', params })

export const membersOfTenant = (params) =>
  service({ url: '/tenantMembership/members', method: 'get', params })
