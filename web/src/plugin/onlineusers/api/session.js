import service from '@/utils/request'

export const listOnlineSessions = (params) => {
  return service({
    url: '/onlineUsers/list',
    method: 'get',
    params
  })
}

export const kickOnlineSession = (data) => {
  return service({
    url: '/onlineUsers/kick',
    method: 'post',
    data
  })
}
