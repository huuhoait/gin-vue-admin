<!--
  LanguageSwitcher (Story 8.3)
  Element Plus dropdown to toggle between vi-VN and en-US.
  Persists via i18n.setLocale() -> localStorage['gva_locale'].
  Also flips Element Plus's own locale so built-in component text matches.
-->
<template>
  <el-dropdown trigger="click" @command="handleChange">
    <el-icon class="cursor-pointer" :size="16"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg></el-icon>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="vi-VN" :disabled="locale === 'vi-VN'">
          {{ t('admin.language.vi_vn') }}
        </el-dropdown-item>
        <el-dropdown-item command="en-US" :disabled="locale === 'en-US'">
          {{ t('admin.language.en_us') }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { setLocale } from '@/i18n'
import { useAppStore } from '@/pinia/modules/app'

const { t, locale } = useI18n()
const appStore = useAppStore?.() // optional: if the store isn't present this is a no-op

function handleChange(next) {
  setLocale(next)
  // Notify any Element Plus locale provider mounted at root.
  if (appStore && typeof appStore.setLocale === 'function') {
    appStore.setLocale(next)
  }
}
</script>
