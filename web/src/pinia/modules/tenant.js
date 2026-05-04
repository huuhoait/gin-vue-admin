import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useStorage } from '@vueuse/core'
import service from '@/utils/request'

/**
 * Tenant store — owns the active tenant id used by the X-Tenant-ID request
 * interceptor in @/utils/request, plus the cached list of tenants the
 * current user has membership in (used by the topbar tenant switcher).
 *
 * State persists in localStorage under `active-tenant-id` so the selection
 * survives full page reloads. A value of 0 means "unset / no tenant header"
 * — the backend tenant middleware will then fall back to the user's primary.
 */
export const useTenantStore = defineStore('tenant', () => {
  // 0 = unset / no header attached. The backend treats absence of
  // X-Tenant-ID as "use my primary tenant", so this is a safe default.
  const activeTenantId = useStorage('active-tenant-id', 0)
  const myTenants = ref([])
  const loaded = ref(false)

  /** Active tenant object, or null when not yet loaded / not in the list. */
  const activeTenant = computed(
    () => myTenants.value.find((t) => t.ID === activeTenantId.value) || null
  )

  /**
   * Fetch the current user's tenant memberships from the BFF and reconcile
   * the persisted activeTenantId. If the persisted id is no longer accessible
   * (or unset), default to the primary membership — falling back to the
   * first tenant in the response if none is flagged primary.
   */
  const loadMyTenants = async () => {
    const res = await service({
      url: '/tenant/mine',
      method: 'get',
      donNotShowLoading: true
    })
    if (res?.code === 0) {
      myTenants.value = Array.isArray(res.data) ? res.data : []
      const stillAccessible = myTenants.value.some(
        (t) => t.ID === activeTenantId.value
      )
      if (!stillAccessible) {
        if (myTenants.value.length > 0) {
          const primary =
            myTenants.value.find((t) => t.isPrimary) || myTenants.value[0]
          activeTenantId.value = primary.ID
        } else {
          activeTenantId.value = 0
        }
      }
      loaded.value = true
    }
    return res
  }

  /** Persist a new tenant selection. Caller is responsible for any reload. */
  const setActiveTenant = (id) => {
    activeTenantId.value = Number(id) || 0
  }

  /** Reset the selection (e.g. on logout). */
  const clearActiveTenant = () => {
    activeTenantId.value = 0
    myTenants.value = []
    loaded.value = false
  }

  return {
    activeTenantId,
    myTenants,
    activeTenant,
    loaded,
    loadMyTenants,
    setActiveTenant,
    clearActiveTenant
  }
})
