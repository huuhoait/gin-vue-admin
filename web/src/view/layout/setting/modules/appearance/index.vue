<template>
  <div class="gva-theme-font">
    <!-- Theme Mode Section -->
    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">{{ t('admin.layout.settings.sections.theme_mode') }}</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <ThemeModeSelector v-model="config.darkMode" @update:modelValue="appStore.toggleDarkMode" />
      </div>
    </div>

    <!-- Theme Color Section -->
    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">{{ t('admin.layout.settings.sections.theme_color') }}</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <ThemeColorPicker v-model="config.primaryColor" @update:modelValue="appStore.togglePrimaryColor" />
      </div>
    </div>

    <!-- Global Size Section -->
    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">{{ t('admin.layout.settings.sections.global_size') }}</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <div class="gva-theme-card-bg">
          <SettingItem :label="t('admin.layout.settings.sections.global_size')">
            <template #suffix>
              <span class="text-xs text-gray-400 dark:text-gray-500 ml-2">{{ t('admin.layout.settings.items.global_size_hint') }}</span>
            </template>
            <div class="w-39">
              <el-select v-model="config.global_size" :placeholder="t('admin.common.select')" @change="appStore.toggleGlobalSize">
                <el-option :label="t('admin.layout.settings.sizes.default')" value="default" />
                <el-option :label="t('admin.layout.settings.sizes.large')" value="large" />
                <el-option :label="t('admin.layout.settings.sizes.small')" value="small" />
              </el-select>
            </div>
          </SettingItem>
        </div>
      </div>
    </div>

    <!-- Visual Accessibility Section -->
    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">{{ t('admin.layout.settings.sections.accessibility') }}</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <div class="gva-theme-card-bg">
          <SettingItem :label="t('admin.layout.settings.items.grey.label')">
            <template #suffix>
              <span class="text-xs text-gray-400 dark:text-gray-500 ml-2">{{ t('admin.layout.settings.items.grey.hint') }}</span>
            </template>
            <el-switch v-model="config.grey" @change="appStore.toggleGrey" />
          </SettingItem>

          <SettingItem :label="t('admin.layout.settings.items.weakness.label')">
            <template #suffix>
              <span class="text-xs text-gray-400 dark:text-gray-500 ml-2">{{ t('admin.layout.settings.items.weakness.hint') }}</span>
            </template>
            <el-switch v-model="config.weakness" @change="appStore.toggleWeakness" />
          </SettingItem>

          <SettingItem :label="t('admin.layout.settings.items.watermark.label')">
            <template #suffix>
              <span class="text-xs text-gray-400 dark:text-gray-500 ml-2">{{ t('admin.layout.settings.items.watermark.hint') }}</span>
            </template>
            <el-switch v-model="config.show_watermark" @change="appStore.toggleConfigWatermark" />
          </SettingItem>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/pinia'
import ThemeModeSelector from '../../components/themeModeSelector.vue'
import ThemeColorPicker from '../../components/themeColorPicker.vue'
import SettingItem from '../../components/settingItem.vue'

defineOptions({
  name: 'AppearanceSettings'
})

const { t } = useI18n()
const appStore = useAppStore()
const { config } = storeToRefs(appStore)
</script>
