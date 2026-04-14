import { useUserStore } from '@/pinia/modules/user'
import { useRouterStore } from '@/pinia/modules/router'
import getPageTitle from '@/utils/page'
import router from '@/router'
import Nprogress from 'nprogress'
import 'nprogress/nprogress.css'

// Configure NProgress
Nprogress.configure({
  showSpinner: false,
  ease: 'ease',
  speed: 500
})

// Route whitelist
const WHITE_LIST = ['Login', 'Init']

function isExternalUrl(val) {
  return typeof val === 'string' && /^(https?:)?\/\//.test(val)
}

// Utility: normalize paths
function normalizeAbsolutePath(p) {
  const s = '/' + String(p || '')
  return s.replace(/\/+/g, '/')
}

function normalizeRelativePath(p) {
  return String(p || '').replace(/^\/+/, '')
}

// Safe registration: only register top-level routes when name does not exist
function addTopLevelIfAbsent(r) {
  if (!router.hasRoute(r.name)) {
    router.addRoute(r)
  }
}

// Flatten N-level menu into:
// - normal: top-level layout + second-level page components
// - if meta.defaultMenu === true: node becomes top-level (not wrapped by layout), and its children become second-level pages
function addRouteByChildren(route, segments = [], parentName = null) {
  // Skip external-link root nodes
  if (isExternalUrl(route?.path) || isExternalUrl(route?.name) || isExternalUrl(route?.component)) {
    return
  }

  // Top-level layout is only a container; it does not participate in path join
  if (route?.name === 'layout') {
    route.children?.forEach((child) => addRouteByChildren(child, [], null))
    return
  }

  // If defaultMenu is set, route should be top-level (not wrapped by layout)
  if (route?.meta?.defaultMenu === true && parentName === null) {
    const fullPath = [...segments, route.path].filter(Boolean).join('/')
    const children = route.children ? [...route.children] : []
    const newRoute = { ...route, path: fullPath }
    delete newRoute.children
    delete newRoute.parent
    // Top-level routes use absolute paths
    newRoute.path = normalizeAbsolutePath(newRoute.path)

    // If route name already exists, skip (children should have been processed)
    if (router.hasRoute(newRoute.name)) return
    addTopLevelIfAbsent(newRoute)

    // If defaultMenu node has children, recurse and mount them under this top-level route
    if (children.length) {
      // Reset segments so children become relative paths under the top-level route
      children.forEach((child) => addRouteByChildren(child, [], newRoute.name))
    }
    return
  }

  // If there are children, continue collecting path segments (ignore external-link segments)
  if (route?.children && route.children.length) {
    if(!parentName){
      const firstChild = route.children[0]
      if (firstChild) {
         const fullParentPath = [...segments, route.path].filter(Boolean).join('/')
         const redirectPath = normalizeRelativePath(
           [fullParentPath, firstChild.path].filter(Boolean).join('/')
         )
         const parentRoute = {
           path: normalizeRelativePath(fullParentPath),
           name: route.name, // keep parent name so defaultRouter can point to it
           meta: route.meta,
           redirect: "/layout/" + redirectPath,
         }
         router.addRoute('layout', parentRoute)
       }
    }
    const nextSegments = isExternalUrl(route.path) ? segments : [...segments, route.path]
    route.children.forEach((child) => addRouteByChildren(child, nextSegments, parentName))
    return
  }

  // Leaf node: register as second-level child route under its parent (defaultMenu top-level or layout)
  const fullPath = [...segments, route.path].filter(Boolean).join('/')
  const newRoute = { ...route, path: fullPath }
  delete newRoute.children
  delete newRoute.parent
  // Child routes use relative paths to avoid /layout/layout/... issues
  newRoute.path = normalizeRelativePath(newRoute.path)

  if (parentName) {
    // Mount under defaultMenu top-level route
    router.addRoute(parentName, newRoute)
  } else {
    // Normal: mount under layout
    router.addRoute('layout', newRoute)
  }
}

// Handle route loading
const setupRouter = async (userStore) => {
  try {
    const routerStore = useRouterStore()
    await Promise.all([routerStore.SetAsyncRouter(), userStore.GetUserInfo()])

    // Ensure parent layout is registered first
    const baseRouters = routerStore.asyncRouters || []
    const layoutRoute = baseRouters[0]
    if (layoutRoute?.name === 'layout' && !router.hasRoute('layout')) {
      const bareLayout = { ...layoutRoute, children: [] }
      router.addRoute(bareLayout)
    }

    // Flatten: register layout.children and other async top-level routes as layout's second-level children
    const toRegister = []
    if (layoutRoute?.children?.length) {
      toRegister.push(...layoutRoute.children)
    }
    if (baseRouters.length > 1) {
      baseRouters.slice(1).forEach((r) => {
        if (r?.name !== 'layout') toRegister.push(r)
      })
    }
  toRegister.forEach((r) => addRouteByChildren(r, [], null))
    return true
  } catch (error) {
    console.error('Setup router failed:', error)
    return false
  }
}

// Remove loading animation
const removeLoading = () => {
  const element = document.getElementById('gva-loading-box')
  element?.remove()
}


// Route guards
router.beforeEach(async (to, from) => {
  const userStore = useUserStore()
  const routerStore = useRouterStore()
  const token = userStore.token

  Nprogress.start()

  // Handle metadata and keep-alive cache
  to.meta.matched = [...to.matched]
  await routerStore.handleKeepAlive(to)
  // Set page title
  document.title = getPageTitle(to.meta.title, to)
  if (to.meta.client) {
    return true
  }

  // Whitelist handling
  if (WHITE_LIST.includes(to.name)) {
    if (token) {
      if(!routerStore.asyncRouterFlag){
        await setupRouter(userStore)
      }
      if(userStore.userInfo.authority.defaultRouter){
        return { name: userStore.userInfo.authority.defaultRouter }
      }
    }
    return  true
  }

  // Auth-required route handling
  if (token) {
    // Handle redirect-to-home case
    if (sessionStorage.getItem('needToHome') === 'true') {
      sessionStorage.removeItem('needToHome')
      return { path: '/' }
    }

    // Handle async routes
    if (!routerStore.asyncRouterFlag && !WHITE_LIST.includes(from.name)) {
      await setupRouter(userStore)
      return to
    }

    return to.matched.length ? true : { path: '/layout/404' }
  }

  // Not logged in: redirect to login
  return {
    name: 'Login',
    query: {
      redirect: to.fullPath
    }
  }
})

// Route load completed
router.afterEach(() => {
  document.querySelector('.main-cont.main-right')?.scrollTo(0, 0)
  Nprogress.done()
})

// Route error handling
router.onError((error) => {
  console.error('Router error:', error)
  Nprogress.remove()
})

// Remove initial loading animation
removeLoading()
