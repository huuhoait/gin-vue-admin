import axios from 'axios'
import { useUserStore } from '@/pinia/modules/user'
import { useTenantStore } from '@/pinia/modules/tenant'
import { ElLoading, ElMessage } from 'element-plus'
import { emitter } from '@/utils/bus'
import router from '@/router/index'

const DEFAULT_REQUEST_TIMEOUT = 1000 * 60 * 10
const DEFAULT_LOADING_FORCE_CLOSE_DELAY = 30000

const service = axios.create()

let activeAxios = 0
let persistentLoadingCount = 0
let timer = null
let forceCloseTimer = null
let loadingInstance = null
let isLoadingVisible = false

const clearLoadingTimers = () => {
  if (timer) {
    clearTimeout(timer)
    timer = null
  }

  if (forceCloseTimer) {
    clearTimeout(forceCloseTimer)
    forceCloseTimer = null
  }
}

const closeLoadingInstance = () => {
  if (isLoadingVisible && loadingInstance) {
    loadingInstance.close()
  }
  loadingInstance = null
  isLoadingVisible = false
}

const scheduleForceClose = () => {
  if (!isLoadingVisible || activeAxios <= 0 || persistentLoadingCount > 0) {
    return
  }

  forceCloseTimer = setTimeout(() => {
    if (isLoadingVisible && loadingInstance) {
      console.warn(
        `Loading force closed after ${DEFAULT_LOADING_FORCE_CLOSE_DELAY}ms`
      )
      closeLoadingInstance()
      activeAxios = 0
      persistentLoadingCount = 0
    }
  }, DEFAULT_LOADING_FORCE_CLOSE_DELAY)
}

const showLoading = (
  option = {
    target: null
  }
) => {
  const loadDom = document.getElementById('gva-base-load-dom')
  const loadingOption = {
    target: null,
    ...option
  }
  const persistLoading = Boolean(loadingOption.persistLoading)

  delete loadingOption.persistLoading

  activeAxios++
  if (persistLoading) {
    persistentLoadingCount++
  }

  clearLoadingTimers()

  timer = setTimeout(() => {
    if (activeAxios > 0 && !isLoadingVisible) {
      if (!loadingOption.target) {
        loadingOption.target = loadDom
      }
      loadingInstance = ElLoading.service(loadingOption)
      isLoadingVisible = true
    }

    scheduleForceClose()
  }, 400)
}

const closeLoading = (option = {}) => {
  activeAxios--
  if (option?.persistLoading && persistentLoadingCount > 0) {
    persistentLoadingCount--
  }

  if (activeAxios <= 0) {
    activeAxios = 0
    persistentLoadingCount = 0
    clearLoadingTimers()
    closeLoadingInstance()
    return
  }

  if (forceCloseTimer) {
    clearTimeout(forceCloseTimer)
    forceCloseTimer = null
  }

  scheduleForceClose()
}

const resetLoading = () => {
  activeAxios = 0
  persistentLoadingCount = 0
  clearLoadingTimers()
  closeLoadingInstance()
}

// Lightweight UUID v4 generator for idempotency keys. Kept inline to avoid a
// circular import with utils/format.js (which imports this module indirectly).
const uuidv4 = () => {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) return crypto.randomUUID()
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0
    const v = c === 'x' ? r : (r & 0x3) | 0x8
    return v.toString(16)
  })
}

// Routes proxied to SkyAgent Core/Order follow the FE↔BFF contract in
// external-frontend-integration.md: Bearer auth (§2) + X-Idempotency-Key on
// every mutation (§6, code 3002). Native GVA endpoints keep the x-token
// scheme they were built around.
const isSkyAgentRoute = (url = '') => url.startsWith('/admin-api/v1/')
const isMutation = (method = '') =>
  ['post', 'put', 'patch', 'delete'].includes(method.toLowerCase())

service.interceptors.request.use(
  (config) => {
    if (typeof config.timeout === 'undefined') {
      config.timeout = DEFAULT_REQUEST_TIMEOUT
    }

    if (!config.donNotShowLoading) {
      showLoading(config.loadingOption)
    }

    config.baseURL = config.baseURL || import.meta.env.VITE_BASE_API

    const userStore = useUserStore()
    config.headers = {
      'Content-Type': 'application/json',
      'x-token': userStore.token,
      'x-user-id': userStore.userInfo.ID,
      ...config.headers
    }

    // Attach X-Tenant-ID for super-admin tenant switching. The store import
    // is at the top of the file, but useTenantStore() is invoked here
    // per-request so the call resolves lazily after pinia is installed by
    // main.js. The header is only sent when a non-zero tenant is selected;
    // the backend treats absence of the header as "use my primary tenant".
    try {
      const tenantStore = useTenantStore()
      const tid = Number(tenantStore.activeTenantId)
      if (tid > 0 && !config.headers['X-Tenant-ID']) {
        config.headers['X-Tenant-ID'] = String(tid)
      }
    } catch (_) {
      // Pinia not yet initialised (very early startup) — skip header.
    }

    if (isSkyAgentRoute(config.url)) {
      if (userStore.token && !config.headers.Authorization) {
        config.headers.Authorization = `Bearer ${userStore.token}`
      }
      if (isMutation(config.method) && !config.headers['X-Idempotency-Key']) {
        config.headers['X-Idempotency-Key'] = uuidv4()
      }
    }

    return config
  },
  (error) => {
    if (!error.config?.donNotShowLoading) {
      closeLoading(error.config?.loadingOption)
    }

    emitter.emit('show-error', {
      code: 'request',
      message: error.message || 'Request failed to send'
    })

    return error
  }
)

function getErrorMessage(error) {
  return error.response?.data?.msg || error.response?.statusText || 'Request failed'
}

service.interceptors.response.use(
  (response) => {
    const userStore = useUserStore()

    if (!response.config.donNotShowLoading) {
      closeLoading(response.config.loadingOption)
    }

    if (response.headers['new-token']) {
      userStore.setToken(response.headers['new-token'])
    }

    if (typeof response.data.code === 'undefined') {
      return response
    }

    if (response.data.code === 0 || response.headers.success === 'true') {
      if (response.headers.msg) {
        response.data.msg = decodeURI(response.headers.msg)
      }
      return response.data
    }

    ElMessage({
      showClose: true,
      message: response.data.msg || decodeURI(response.headers.msg),
      type: 'error'
    })

    return response.data.msg ? response.data : response
  },
  (error) => {
    if (!error.config?.donNotShowLoading) {
      closeLoading(error.config?.loadingOption)
    }

    if (!error.response) {
      resetLoading()
      emitter.emit('show-error', {
        code: 'network',
        message: getErrorMessage(error)
      })
      return Promise.reject(error)
    }

    if (error.response.status === 401) {
      emitter.emit('show-error', {
        code: '401',
        message: getErrorMessage(error),
        fn: () => {
          const userStore = useUserStore()
          userStore.ClearStorage()
          router.push({ name: 'Login', replace: true })
        }
      })
      return Promise.reject(error)
    }

    if (error.response.status === 403) {
      // Permissions or menu/policy might have just been updated by an admin.
      // A reload forces re-fetch of dynamic routes and re-runs guards.
      emitter.emit('show-error', {
        code: '403',
        message: getErrorMessage(error),
        fn: () => window.location.reload()
      })
      return Promise.reject(error)
    }

    emitter.emit('show-error', {
      code: error.response.status,
      message: getErrorMessage(error)
    })
    return Promise.reject(error)
  }
)

if (typeof window !== 'undefined') {
  window.addEventListener('beforeunload', resetLoading)
  window.addEventListener('unload', resetLoading)
}

export { resetLoading }
export default service
