import { defineStore } from 'pinia'
import { ref, watchEffect, reactive } from 'vue'
import { setBodyPrimaryColor } from '@/utils/format'
import { useDark, usePreferredDark } from '@vueuse/core'
import { getLocale, setLocale as setI18nLocale } from '@/i18n'

export const useAppStore = defineStore('app', () => {
  const device = ref('')
  const drawerSize = ref('')
  const operateMinWith = ref('240')
  const locale = ref(getLocale())
  const config = reactive({
    weakness: false,
    grey: false,
    primaryColor: '#3b82f6',
    showTabs: true,
    darkMode: 'light',
    layout_side_width: 256,
    layout_side_collapsed_width: 80,
    layout_side_item_height: 48,
    show_watermark: false,
    side_mode: 'normal',
    // Page transition animation
    transition_type: 'slide',
    global_size: 'default'
  })

  const isDark = useDark({
    selector: 'html',
    attribute: 'class',
    valueDark: 'dark',
    valueLight: 'light'
  })

  const preferredDark = usePreferredDark()

  const toggleTheme = (darkMode) => {
    isDark.value = darkMode
  }

  const toggleWeakness = (e) => {
    config.weakness = e
  }

  const toggleGrey = (e) => {
    config.grey = e
  }

  const togglePrimaryColor = (e) => {
    config.primaryColor = e
  }

  const toggleTabs = (e) => {
    config.showTabs = e
  }

  const toggleDevice = (e) => {
    if (e === 'mobile') {
      drawerSize.value = '100%'
      operateMinWith.value = '80'
    } else {
      drawerSize.value = '800'
      operateMinWith.value = '240'
    }
    device.value = e
  }

  const toggleDarkMode = (e) => {
    config.darkMode = e
  }

  // Watch system theme changes
  watchEffect(() => {
    if (config.darkMode === 'auto') {
      isDark.value = preferredDark.value
      return
    }
    isDark.value = config.darkMode === 'dark'
  })

  const toggleConfigSideWidth = (e) => {
    config.layout_side_width = e
  }

  const toggleConfigSideCollapsedWidth = (e) => {
    config.layout_side_collapsed_width = e
  }

  const toggleConfigSideItemHeight = (e) => {
    config.layout_side_item_height = e
  }

  const toggleConfigWatermark = (e) => {
    config.show_watermark = e
  }

  const toggleSideMode = (e) => {
    config.side_mode = e
  }

  const toggleTransition = (e) => {
    config.transition_type = e
  }

  const toggleGlobalSize = (e) => {
    config.global_size = e
  }

  const setLocale = (nextLocale) => {
    locale.value = nextLocale
    setI18nLocale(nextLocale)
  }

  const baseCoinfg = {
    weakness: false,
    grey: false,
    primaryColor: '#3b82f6',
    showTabs: true,
    darkMode: 'light',
    layout_side_width: 256,
    layout_side_collapsed_width: 80,
    layout_side_item_height: 48,
    show_watermark: false,
    side_mode: 'normal',
    // Page transition animation
    transition_type: 'slide',
    global_size: 'default'
  }

  const resetConfig = () => {
    for (let baseCoinfgKey in baseCoinfg) {
      config[baseCoinfgKey] = baseCoinfg[baseCoinfgKey]
    }
  }

  // Watch weakness/grey mode flags
  watchEffect(() => {
    document.documentElement.classList.toggle('html-weakenss', config.weakness)
    document.documentElement.classList.toggle('html-grey', config.grey)
  })

  // Watch primary color
  watchEffect(() => {
    setBodyPrimaryColor(config.primaryColor, isDark.value ? 'dark' : 'light')
  })

  return {
    isDark,
    device,
    drawerSize,
    operateMinWith,
    locale,
    config,
    toggleTheme,
    toggleDevice,
    toggleWeakness,
    toggleGrey,
    togglePrimaryColor,
    toggleTabs,
    toggleDarkMode,
    toggleConfigSideWidth,
    toggleConfigSideCollapsedWidth,
    toggleConfigSideItemHeight,
    toggleConfigWatermark,
    toggleSideMode,
    toggleTransition,
    resetConfig,
    toggleGlobalSize,
    setLocale
  }
})
