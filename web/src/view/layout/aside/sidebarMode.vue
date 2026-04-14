<template>
  <div class="flex h-full">
    <!-- primary menu (always visible) -->
    <div
      class="relative !h-full bg-white text-slate-700 dark:text-slate-300 dark:bg-slate-900 shadow dark:shadow-gray-700"
      :style="{
        width: config.layout_side_collapsed_width + 'px'
      }"
    >
      <el-scrollbar>
        <el-menu
          :collapse="true"
          :collapse-transition="false"
          :default-active="topActive"
          class="!border-r-0 w-full"
          unique-opened
          @select="selectTopMenuItem"
        >
          <template v-for="item in routerStore.asyncRouters[0]?.children || []">
            <el-menu-item
              v-if="!item.hidden && (!item.children || item.children.length === 0)"
              :key="item.name"
              :index="item.name"
              class="dark:text-slate-300 overflow-hidden"
              :style="{
                height: config.layout_side_item_height + 'px'
              }"
            >
              <el-icon v-if="item.meta.icon">
                <component :is="item.meta.icon" />
              </el-icon>
              <template v-else>
                {{ $t(item.meta.title)[0] }}
              </template>
              <template #title>
                {{ $t(item.meta.title) }}
              </template>
            </el-menu-item>
            <template v-else-if="!item.hidden" >
             <el-menu-item
              :key="item.name"
              :index="item.name"
              :class="{'is-active': topActive === item.name}"
              class="dark:text-slate-300 overflow-hidden"
              :style="{
                height: config.layout_side_item_height + 'px'
              }"
            >
              <el-icon v-if="item.meta.icon">
                <component :is="item.meta.icon" />
              </el-icon>
              <template v-else>
                {{ $t(item.meta.title)[0] }}
                </template>
              <template #title>
                {{ $t(item.meta.title) }}
              </template>
            </el-menu-item>
            </template>
          </template>
        </el-menu>
      </el-scrollbar>
    </div>

    <!-- secondary menu (columnar) -->
    <div
      class="relative h-full bg-white text-slate-700 dark:text-slate-300 dark:bg-slate-900 shadow dark:shadow-gray-700 px-2"
      :style="{
        width: layoutSideWidth + 'px'
      }"
    >
      <el-scrollbar>
        <el-menu
          :collapse="isCollapse"
          :collapse-transition="false"
          :default-active="active"
          class="!border-r-0 w-full"
          unique-opened
          @select="selectMenuItem"
        >
          <template v-for="item in secondLevelMenus">
            <aside-component
              v-if="!item.hidden"
              :key="item.name"
              :router-info="item"
            />
          </template>
        </el-menu>
      </el-scrollbar>
      <div
        class="absolute bottom-8 right-2 w-8 h-8 bg-gray-50 dark:bg-slate-800 flex items-center justify-center rounded cursor-pointer"
        :class="isCollapse ? 'right-0 left-0 mx-auto' : 'right-2'"
        @click="toggleCollapse"
      >
        <el-icon v-if="!isCollapse">
          <DArrowLeft />
        </el-icon>
        <el-icon v-else>
          <DArrowRight />
        </el-icon>
      </div>
    </div>
  </div>
</template>

<script setup>
  import AsideComponent from '@/view/layout/aside/asideComponent/index.vue'
  import { ref, provide, watchEffect, computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useRouterStore } from '@/pinia/modules/router'
  import { useAppStore } from '@/pinia'
  import { storeToRefs } from 'pinia'

  const appStore = useAppStore()
  const { device, config } = storeToRefs(appStore)

  defineOptions({
    name: 'SidebarMode'
  })

  const route = useRoute()
  const router = useRouter()
  const routerStore = useRouterStore()
  const isCollapse = ref(false)
  const active = ref('')
  const topActive = ref('')
  const secondLevelMenus = ref([])

  const layoutSideWidth = computed(() => {
    if (!isCollapse.value) {
      return config.value.layout_side_width
    } else {
      return config.value.layout_side_collapsed_width
    }
  })


  provide('isCollapse', isCollapse)

  // update secondary menu list
  const updateSecondLevelMenus = (menuName) => {
    const menu = routerStore.asyncRouters[0]?.children.find(item => item.name === menuName)
    if (menu && menu.children && menu.children.length > 0) {
      secondLevelMenus.value = menu.children
    }
  }

  // click handler for top-level menu
  const selectTopMenuItem = (index) => {
    topActive.value = index

    // fetch the selected menu item
    const menu = routerStore.asyncRouters[0]?.children.find(item => item.name === index)

    // only refresh the secondary pane when the selected menu has children
    if (menu && menu.children && menu.children.length > 0) {
      updateSecondLevelMenus(index)

      // navigate to the first visible child
      const firstVisibleChild = menu.children.find(child => !child.hidden)
      if (firstVisibleChild) {
        navigateToMenuItem(firstVisibleChild.name)
      }
    } else {
      // leaf menu: navigate directly without touching the secondary pane
      navigateToMenuItem(index)
    }
  }

  // click handler for secondary (or deeper) menu
  const selectMenuItem = (index) => {
    navigateToMenuItem(index)
  }

  // navigate to the given menu name
  const navigateToMenuItem = (index) => {
    const query = {}
    const params = {}
    routerStore.routeMap[index]?.parameters &&
      routerStore.routeMap[index]?.parameters.forEach((item) => {
        if (item.type === 'query') {
          query[item.key] = item.value
        } else {
          params[item.key] = item.value
        }
      })
    if (index === route.name) return
    if (index.indexOf('http://') > -1 || index.indexOf('https://') > -1) {
        window.open(index, '_blank')
        return
    } else {
      router.push({ name: index, query, params })
    }
  }

  const toggleCollapse = () => {
    isCollapse.value = !isCollapse.value
  }



  watchEffect(() => {
    if (route.name === 'gvaLayoutIframe') {
      active.value = decodeURIComponent(route.query.url)
      return
    }
    active.value = route.meta.activeName || route.name

    // find the top-level menu that owns the current route
    const findParentMenu = () => {
      // first check whether the current route itself is a top-level menu
      const isTopMenu = routerStore.asyncRouters[0]?.children.some(
        item => !item.hidden && item.name === route.name
      )

      if (isTopMenu) {
        return route.name
      }

      for (const topMenu of routerStore.asyncRouters[0]?.children || []) {
        if (topMenu.hidden) continue

        // check whether the current route is a direct child
        if (topMenu.children && topMenu.children.some(child => child.name === route.name)) {
          return topMenu.name
        }

        // recurse for deeper levels
        const checkChildren = (items) => {
          for (const item of items || []) {
            if (item.name === route.name) {
              return true
            }
            if (item.children && checkChildren(item.children)) {
              return true
            }
          }
          return false
        }

        if (topMenu.children && checkChildren(topMenu.children)) {
          return topMenu.name
        }
      }
      return null
    }

    const parentMenu = findParentMenu()
    if (parentMenu) {
      topActive.value = parentMenu

      // only refresh the secondary pane when the parent has children
      const menu = routerStore.asyncRouters[0]?.children.find(item => item.name === parentMenu)
      if (menu && menu.children && menu.children.length > 0) {
        updateSecondLevelMenus(parentMenu)
      } else {
        // parent has no children: keep its highlight but show the first
        // non-empty sibling's children in the secondary pane so the UI
        // doesn't go empty.
        const firstMenuWithChildren = routerStore.asyncRouters[0].children.find(
          item => !item.hidden && item.children && item.children.length > 0
        )

        if (firstMenuWithChildren) {
          // refresh secondary pane only; keep top highlight unchanged.
          updateSecondLevelMenus(firstMenuWithChildren.name)
        }
      }
    } else if (routerStore.asyncRouters[0]?.children?.length > 0) {
      // no parent found: highlight the current route, but still populate the
      // secondary pane with the first non-empty top menu's children.
      const firstMenuWithChildren = routerStore.asyncRouters[0].children.find(
        item => !item.hidden && item.children && item.children.length > 0
      )

      if (firstMenuWithChildren) {
        // update secondary pane only; highlight stays on current route name.
        topActive.value = route.name
        secondLevelMenus.value = firstMenuWithChildren.children
      }
    }
  })

  watchEffect(() => {
    if (device.value === 'mobile') {
      isCollapse.value = true
    } else {
      isCollapse.value = false
    }
  })
</script>
