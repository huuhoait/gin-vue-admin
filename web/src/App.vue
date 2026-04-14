<template>
  <div
    id="app"
    class="bg-gray-50 text-slate-700 !dark:text-slate-500 dark:bg-slate-800"
  >
    <el-config-provider :locale="elementLocale" :size="appStore.config.global_size">
      <router-view />
      <Application />
    </el-config-provider>
  </div>
</template>

<script setup>
  import { computed } from 'vue'
  import { useI18n } from 'vue-i18n'
  import en from 'element-plus/dist/locale/en.mjs'
  import vi from 'element-plus/dist/locale/vi.mjs'
  import Application from '@/components/application/index.vue'
  import { useAppStore } from '@/pinia'

  const appStore = useAppStore()
  const { locale } = useI18n()

  const elementLocale = computed(() => {
    // Keep Element Plus built-in component strings in sync with vue-i18n.
    // Element Plus locale codes: "en" and "vi" (not "en-US"/"vi-VN").
    return locale.value === 'vi-VN' ? vi : en
  })
  defineOptions({
    name: 'App'
  })
</script>
<style lang="scss">
  // Base styles
  #app {
    height: 100vh;
    overflow: hidden;
    font-weight: 400 !important;
  }

  .el-button {
    font-weight: 400 !important;
  }

  .gva-body-h {
    min-height: calc(100% - 3rem);
  }

  .gva-container {
    height: calc(100% - 2.5rem);
  }

  .gva-container2 {
    height: calc(100% - 4.5rem);
  }
</style>
