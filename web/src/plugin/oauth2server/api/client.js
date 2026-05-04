import service from '@/utils/request'

export const listOAuth2Clients = (params) =>
  service({ url: '/oauth2Client/list', method: 'get', params })

export const findOAuth2Client = (params) =>
  service({ url: '/oauth2Client/find', method: 'get', params })

export const createOAuth2Client = (data) =>
  service({ url: '/oauth2Client/create', method: 'post', data })

export const updateOAuth2Client = (data) =>
  service({ url: '/oauth2Client/update', method: 'put', data })

export const deleteOAuth2Client = (params) =>
  service({ url: '/oauth2Client/delete', method: 'delete', params })

export const regenerateOAuth2ClientSecret = (data) =>
  service({ url: '/oauth2Client/regenerateSecret', method: 'post', data })
