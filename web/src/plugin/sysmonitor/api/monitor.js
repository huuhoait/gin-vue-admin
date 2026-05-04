import service from '@/utils/request'

export const getServerStats = () => service({ url: '/sysmonitor/server', method: 'get' })
export const getRuntimeStats = () => service({ url: '/sysmonitor/runtime', method: 'get' })
export const getCacheStats = () => service({ url: '/sysmonitor/cache', method: 'get' })
