import { login, getUserInfo } from '@/api/user'
import { jsonInBlacklist } from '@/api/jwt'
import router from '@/router/index'
import { ElLoading, ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouterStore } from './router'
import { useCookies } from '@vueuse/integrations/useCookies'
import { useStorage } from '@vueuse/core'

import { useAppStore } from '@/pinia'

export const useUserStore = defineStore('user', () => {
  const appStore = useAppStore()
  const loadingInstance = ref(null)

  const userInfo = ref({
    uuid: '',
    nickName: '',
    headerImg: '',
    authority: {}
  })
  const token = useStorage('token', '')
  const xToken = useCookies()
  const currentToken = computed(() => token.value || xToken.get('x-token') || '')

  const setUserInfo = (val) => {
    userInfo.value = val
    if (val.originSetting) {
      Object.keys(appStore.config).forEach((key) => {
        if (val.originSetting[key] !== undefined) {
          appStore.config[key] = val.originSetting[key]
        }
      })
    }
  }

  const setToken = (val) => {
    token.value = val
    xToken.value = val
  }

  const NeedInit = async () => {
    await ClearStorage()
    await router.push({ name: 'Init', replace: true })
  }

  const ResetUserInfo = (value = {}) => {
    userInfo.value = {
      ...userInfo.value,
      ...value
    }
  }
  /* Get user info */
  const GetUserInfo = async () => {
    const res = await getUserInfo()
    if (res.code === 0) {
      setUserInfo(res.data.userInfo)
    }
    return res
  }
  /* Sign in */
  const LoginIn = async (loginInfo) => {
    try {
      loadingInstance.value = ElLoading.service({
        fullscreen: true,
        text: 'Signing in, please wait...'
      })

      const res = await login(loginInfo)

      if (res.code !== 0) {
        return false
      }
      // Signed in: set user info and authority
      setUserInfo(res.data.user)
      setToken(res.data.token)

      // Initialize route info
      const routerStore = useRouterStore()
      await routerStore.SetAsyncRouter()
      const asyncRouters = routerStore.asyncRouters

      // Register to router
      asyncRouters.forEach((asyncRouter) => {
        router.addRoute(asyncRouter)
      })

      if(router.currentRoute.value.query.redirect) {
        await router.replace(router.currentRoute.value.query.redirect)
        return true
      }

      if (!router.hasRoute(userInfo.value.authority.defaultRouter)) {
        ElMessage.error('No default home route is available. Please contact an administrator.')
      } else {
        await router.replace({ name: userInfo.value.authority.defaultRouter })
      }

      const isWindows = /windows/i.test(navigator.userAgent)
      window.localStorage.setItem('osType', isWindows ? 'WIN' : 'MAC')

      // Done: close loading and return
      return true
    } catch (error) {
      console.error('LoginIn error:', error)
      return false
    } finally {
      loadingInstance.value?.close()
    }
  }
  /* Sign out */
  const LoginOut = async () => {
    const res = await jsonInBlacklist()

    // Sign out failed
    if (res.code !== 0) {
      return
    }

    await ClearStorage()

    // Redirect to login and reload
    router.push({ name: 'Login', replace: true })
    window.location.reload()
  }
  /* Clear storage */
  const ClearStorage = async () => {
    token.value = ''
    // Use remove() to delete cookie
    xToken.remove()
    sessionStorage.clear()
    // Clear related localStorage items
    localStorage.removeItem('originSetting')
    localStorage.removeItem('token')
  }

  return {
    userInfo,
    token: currentToken,
    NeedInit,
    ResetUserInfo,
    GetUserInfo,
    LoginIn,
    LoginOut,
    setToken,
    loadingInstance,
    ClearStorage
  }
})
