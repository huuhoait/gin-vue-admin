<template>
  <div class="gva-theme-font">
    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">{{ t('admin.layout.settings.sections.layout_mode') }}</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <LayoutModeCard
          v-model="config.side_mode"
          @update:modelValue="appStore.toggleSideMode"
        />
      </div>
    </div>

    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">{{ t('admin.layout.settings.sections.ui_config') }}</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <div class="gva-theme-card-bg">
          <SettingItem :label="t('admin.layout.settings.items.show_tabs.label')">
            <template #suffix>
              <span class="text-xs text-gray-400 dark:text-gray-500 ml-2">{{ t('admin.layout.settings.items.show_tabs.hint') }}</span>
            </template>
            <el-switch
              v-model="config.showTabs"
              @change="appStore.toggleTabs"
            />
          </SettingItem>

          <SettingItem :label="t('admin.layout.settings.items.transition.label')">
            <template #suffix>
              <span class="text-xs text-gray-400 dark:text-gray-500 ml-2">{{ t('admin.layout.settings.items.transition.hint') }}</span>
            </template>
            <el-select
              v-model="config.transition_type"
              @change="appStore.toggleTransition"
              class="w-32"
              size="small"
            >
              <el-option value="fade" :label="t('admin.layout.settings.transitions.fade')" />
              <el-option value="slide" :label="t('admin.layout.settings.transitions.slide')" />
              <el-option value="zoom" :label="t('admin.layout.settings.transitions.zoom')" />
              <el-option value="none" :label="t('admin.layout.settings.transitions.none')" />
            </el-select>
          </SettingItem>
        </div>
      </div>
    </div>

    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">{{ t('admin.layout.settings.sections.size_config') }}</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <div class="gva-theme-card-bg">
          <div class="space-y-4">
            <div class="gva-theme-card-white">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="text-sm font-medium gva-theme-text-main">{{ t('admin.layout.settings.items.side_width.label') }}</h4>
                  <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ t('admin.layout.settings.items.side_width.hint') }}</p>
                </div>
                <div class="flex items-center gap-2">
                  <el-input-number
                    v-model="config.layout_side_width"
                    :min="150"
                    :max="400"
                    :step="10"
                    size="small"
                    class="w-24"
                  />
                  <span class="text-xs font-medium text-gray-500 dark:text-gray-400">px</span>
                </div>
              </div>
            </div>

            <div class="gva-theme-card-white">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="text-sm font-medium gva-theme-text-main">{{ t('admin.layout.settings.items.side_collapsed.label') }}</h4>
                  <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ t('admin.layout.settings.items.side_collapsed.hint') }}</p>
                </div>
                <div class="flex items-center gap-2">
                  <el-input-number
                    v-model="config.layout_side_collapsed_width"
                    :min="60"
                    :max="100"
                    size="small"
                    class="w-24"
                  />
                  <span class="text-xs font-medium text-gray-500 dark:text-gray-400">px</span>
                </div>
              </div>
            </div>

            <div class="gva-theme-card-white">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="text-sm font-medium gva-theme-text-main">{{ t('admin.layout.settings.items.side_item_height.label') }}</h4>
                  <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ t('admin.layout.settings.items.side_item_height.hint') }}</p>
                </div>
                <div class="flex items-center gap-2">
                  <el-input-number
                    v-model="config.layout_side_item_height"
                    :min="30"
                    :max="50"
                    size="small"
                    class="w-24"
                  />
                  <span class="text-xs font-medium text-gray-500 dark:text-gray-400">px</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/pinia'
import LayoutModeCard from '../../components/layoutModeCard.vue'
import SettingItem from '../../components/settingItem.vue'

defineOptions({
  name: 'LayoutSettings'
})

const { t } = useI18n()
const appStore = useAppStore()
const { config } = storeToRefs(appStore)
</script>
