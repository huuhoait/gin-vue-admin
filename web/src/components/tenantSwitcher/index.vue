<!--
  TenantSwitcher (Priority 6 — multi-tenancy)
  Topbar dropdown that lists tenants the current user has membership in and
  switches the active X-Tenant-ID. Hidden when the user has 0 or 1 tenants
  (no point in showing a single-option selector).

  On change we trigger a full page reload — every tenant-scoped store and
  cached query in the SPA would otherwise need bespoke invalidation, and a
  reload guarantees a clean state under the new header. Cheaper to design
  for than to debug stale data.
-->
<template>
  <el-dropdown
    v-if="tenantStore.myTenants.length > 1"
    trigger="click"
    placement="bottom"
    @command="handleSwitch"
  >
    <span
      class="flex items-center cursor-pointer text-black dark:text-gray-100 mr-2"
    >
      <el-icon :size="16" class="mr-1">
        <office-building />
      </el-icon>
      <span class="hidden sm:inline-block max-w-32 truncate">{{
        activeLabel
      }}</span>
      <el-icon class="ml-1">
        <arrow-down />
      </el-icon>
    </span>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="t in tenantStore.myTenants"
          :key="t.ID"
          :command="t.ID"
          :disabled="t.ID === tenantStore.activeTenantId"
        >
          <span class="flex items-center justify-between w-full gap-2">
            <span class="truncate">{{ t.name }}</span>
            <el-tag
              v-if="t.ID === tenantStore.activeTenantId"
              size="small"
              type="success"
            >
              {{ $t('admin.layout.tenant_switcher.current_tag') }}
            </el-tag>
            <el-tag
              v-else-if="t.isPrimary"
              size="small"
              type="info"
            >
              {{ $t('admin.layout.tenant_switcher.primary_tag') }}
            </el-tag>
          </span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
/**
 * Topbar tenant switcher.
 * @component TenantSwitcher
 * @description Lists the user's accessible tenants and switches the active
 *              tenant id (which feeds the X-Tenant-ID request header).
 */
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowDown, OfficeBuilding } from '@element-plus/icons-vue'
import { useTenantStore } from '@/pinia/modules/tenant'

const { t } = useI18n()
const tenantStore = useTenantStore()

const activeLabel = computed(
  () =>
    tenantStore.activeTenant?.name ||
    t('admin.layout.tenant_switcher.placeholder')
)

/**
 * Apply a new tenant selection. We reload the window so every store/query
 * is rebuilt against the new tenant scope — the safe option until each
 * feature explicitly supports live re-fetch on tenant change.
 * @param {number} id Tenant ID chosen from the dropdown.
 */
const handleSwitch = (id) => {
  if (id === tenantStore.activeTenantId) return
  tenantStore.setActiveTenant(id)
  window.location.reload()
}
</script>

<style scoped lang="scss"></style>
