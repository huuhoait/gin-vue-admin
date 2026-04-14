<template>
  <div class="flex items-center mx-4 gap-4">
    <el-tooltip v-if="isDev && videoList.length" class="" effect="dark" :content="t('admin.layout.tools.video_tutorial')" placement="bottom">
      <el-dropdown @command="toDoc">
        <span class="w-8 h-8 p-2 rounded-full flex items-center justify-center shadow border border-gray-200 dark:border-gray-600 cursor-pointer border-solid">
          <el-icon>
          <Film />
        </el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item
              v-for="item in videoList"
              :key="item.link"
              :command="item.link"
              >{{ item.title }}</el-dropdown-item
            >
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </el-tooltip>

    <el-tooltip class="" effect="dark" :content="t('admin.layout.tools.search')" placement="bottom">
        <span class="w-8 h-8 p-2 rounded-full flex items-center justify-center shadow border border-gray-200 dark:border-gray-600 cursor-pointer border-solid">
        <el-icon
            @click="handleCommand"
        >
        <Search />
      </el-icon>
        </span>

    </el-tooltip>

    <el-tooltip class="" effect="dark" :content="t('admin.layout.tools.settings')" placement="bottom">
        <span class="w-8 h-8 p-2 rounded-full flex items-center justify-center shadow border border-gray-200 dark:border-gray-600 cursor-pointer border-solid">
         <el-icon
             @click="toggleSetting"
         >
        <Setting />
      </el-icon>
        </span>

    </el-tooltip>

    <el-tooltip class="" effect="dark" :content="t('admin.layout.tools.refresh')" placement="bottom">
      <span class="w-8 h-8 p-2 rounded-full flex items-center justify-center shadow border border-gray-200 dark:border-gray-600 cursor-pointer border-solid">
      <el-icon
          :class="showRefreshAnmite ? 'animate-spin' : ''"
          @click="toggleRefresh"
      >
        <Refresh />
      </el-icon>
      </span>

    </el-tooltip>
    <el-tooltip
      class=""
      effect="dark"
      :content="t('admin.layout.tools.toggle_theme')"
      placement="bottom"
    >
      <span class="w-8 h-8 p-2 rounded-full flex items-center justify-center shadow border border-gray-200 dark:border-gray-600 cursor-pointer border-solid">
        <el-icon
            v-if="appStore.isDark"
            @click="appStore.toggleTheme(false)"
        >
        <Sunny />
      </el-icon>
      <el-icon
          v-else
          @click="appStore.toggleTheme(true)"
      >
        <Moon />
      </el-icon>
      </span>

    </el-tooltip>

    <el-tooltip class="" effect="dark" :content="t('admin.layout.tools.change_language')" placement="bottom">
      <span class="w-8 h-8 p-2 rounded-full flex items-center justify-center shadow border border-gray-200 dark:border-gray-600 cursor-pointer border-solid">
        <language-switcher />
      </span>
    </el-tooltip>

    <gva-setting v-model:drawer="showSettingDrawer"></gva-setting>
    <command-menu ref="command" />
  </div>
</template>

<script setup>
  import { useAppStore } from '@/pinia'
  import GvaSetting from '@/view/layout/setting/index.vue'
  import { ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { emitter } from '@/utils/bus.js'
  import CommandMenu from '@/components/commandMenu/index.vue'
  import LanguageSwitcher from '@/components/languageSwitcher/index.vue'
  import { toDoc } from '@/utils/doc'
  import { isDev } from '@/utils/env.js'

  const { t } = useI18n()
  const appStore = useAppStore()
  const showSettingDrawer = ref(false)
  const showRefreshAnmite = ref(false)
  const toggleRefresh = () => {
    showRefreshAnmite.value = true
    emitter.emit('reload')
    setTimeout(() => {
      showRefreshAnmite.value = false
    }, 1000)
  }

  const toggleSetting = () => {
    showSettingDrawer.value = true
  }

  const first = ref('')
  const command = ref()

  const handleCommand = () => {
    command.value.open()
  }
  const initPage = () => {
    // detect the user's OS to show the right shortcut key
    if (window.localStorage.getItem('osType') === 'WIN') {
      first.value = 'Ctrl'
    } else {
      first.value = '⌘'
    }
    // open command palette on Ctrl+K
    const handleKeyDown = (e) => {
      if (e.ctrlKey && e.key === 'k') {
        // suppress browser default
        e.preventDefault()
        handleCommand()
      }
    }
    window.addEventListener('keydown', handleKeyDown)
  }

  initPage()

  // Story 8.4: upstream gin-vue-admin bundled a list of Chinese Bilibili
  // tutorial links here. Those videos are not relevant to SkyAgent operators,
  // so the list is empty — the dropdown auto-hides via v-if above. If a
  // SkyAgent-specific tutorial library is ever added, fill this array with
  // { title: t('…'), link } objects.
  const videoList = []
</script>

<style scoped lang="scss"></style>
