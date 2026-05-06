<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true">
        <el-form-item :label="t('admin.plugin.tenant.col_name')">
          <el-tag v-if="activeTenant" type="success">{{ activeTenant.name }}</el-tag>
          <el-tag v-else type="info">{{ t('admin.plugin.tenant.enabled_any') }}</el-tag>
        </el-form-item>
        <el-form-item>
          <el-button icon="refresh" @click="loadMembers">{{ t('admin.plugin.tenant.reset') }}</el-button>
        </el-form-item>
      </el-form>
      <div class="text-xs text-gray-500" v-if="!activeTenant">
        {{ t('admin.plugin.tenant.member_no_user_chosen') }}
      </div>
    </div>

    <div class="gva-table-box">
      <el-form :inline="true" :model="memberForm">
        <el-form-item :label="t('admin.plugin.tenant.member_user')">
          <div class="flex items-center gap-2">
            <el-button icon="user" @click="openUserPicker">{{ t('admin.plugin.tenant.member_pick_user') }}</el-button>
            <span v-if="memberForm.userID" class="text-sm text-gray-700">
              #{{ memberForm.userID }} · {{ memberForm.username }}<span v-if="memberForm.nickName"> ({{ memberForm.nickName }})</span>
            </span>
            <span v-else class="text-sm text-gray-400">{{ t('admin.plugin.tenant.member_no_user_chosen') }}</span>
          </div>
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.member_primary')">
          <el-switch v-model="memberForm.isPrimary" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :disabled="!activeTenantId || !memberForm.userID" @click="onAssign">
            {{ t('admin.plugin.tenant.member_assign') }}
          </el-button>
        </el-form-item>
      </el-form>

      <el-table :data="memberRows" size="small" style="width: 100%">
        <el-table-column :label="t('admin.plugin.tenant.member_user_id')" prop="userID" width="80" />
        <el-table-column :label="t('admin.plugin.tenant.member_username')" prop="username" min-width="140" />
        <el-table-column :label="t('admin.plugin.tenant.member_nickname')" prop="nickName" min-width="140" />
        <el-table-column :label="t('admin.plugin.tenant.member_primary')" width="90">
          <template #default="scope">
            <el-tag v-if="scope.row.isPrimary" type="success" size="small">{{ t('admin.plugin.tenant.primary_tag') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.member_joined')" width="160">
          <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.col_actions')" width="100" fixed="right">
          <template #default="scope">
            <el-button link type="danger" @click="onUnassign(scope.row)">{{ t('admin.plugin.tenant.member_remove') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog
      v-model="userPickerVisible"
      :title="t('admin.plugin.tenant.user_picker_title')"
      width="720"
      destroy-on-close
      append-to-body
    >
      <el-form :inline="true" :model="userPickerSearch">
        <el-form-item :label="t('admin.plugin.tenant.user_picker_username')">
          <el-input v-model="userPickerSearch.username" clearable @keyup.enter="onUserPickerSearch" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.user_picker_nickname')">
          <el-input v-model="userPickerSearch.nickName" clearable @keyup.enter="onUserPickerSearch" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onUserPickerSearch">{{ t('admin.plugin.tenant.search') }}</el-button>
          <el-button icon="refresh" @click="onUserPickerReset">{{ t('admin.plugin.tenant.reset') }}</el-button>
        </el-form-item>
      </el-form>
      <el-table :data="userPickerRows" size="small" style="width: 100%" highlight-current-row @row-click="onUserRowSelect">
        <el-table-column :label="t('admin.plugin.tenant.member_user_id')" prop="ID" width="80" />
        <el-table-column :label="t('admin.plugin.tenant.member_username')" prop="userName" min-width="140" />
        <el-table-column :label="t('admin.plugin.tenant.member_nickname')" prop="nickName" min-width="140" />
        <el-table-column :label="t('admin.plugin.tenant.col_actions')" width="120" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click.stop="onUserRowSelect(scope.row)">{{ t('admin.plugin.tenant.user_picker_select') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        :current-page="userPickerPage"
        :page-size="userPickerPageSize"
        :total="userPickerTotal"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        class="mt-2"
        @current-change="(v) => { userPickerPage = v; loadUserPicker() }"
        @size-change="(v) => { userPickerPageSize = v; userPickerPage = 1; loadUserPicker() }"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useTenantStore } from '@/pinia/modules/tenant'
import service from '@/utils/request'

const { t } = useI18n()
const tenantStore = useTenantStore()
const activeTenantId = computed(() => Number(tenantStore.activeTenantId) || 0)
const activeTenant = computed(() => tenantStore.activeTenant)

const memberRows = ref([])
const memberForm = ref({
  userID: 0,
  username: '',
  nickName: '',
  isPrimary: false
})

const userPickerVisible = ref(false)
const userPickerSearch = ref({ username: '', nickName: '' })
const userPickerRows = ref([])
const userPickerPage = ref(1)
const userPickerPageSize = ref(10)
const userPickerTotal = ref(0)

function formatDate(val) {
  if (!val) return ''
  const d = new Date(val)
  return isNaN(d.getTime()) ? String(val) : d.toLocaleString()
}

const loadMembers = async () => {
  if (!activeTenantId.value) {
    memberRows.value = []
    return
  }
  const res = await service({
    url: '/tenantMembership/members',
    method: 'get',
    params: { id: activeTenantId.value }
  })
  if (res?.code === 0) {
    memberRows.value = Array.isArray(res.data) ? res.data : []
  } else {
    ElMessage.error(res?.msg || t('admin.plugin.tenant.members_load_failed'))
  }
}

const onAssign = async () => {
  if (!activeTenantId.value || !memberForm.value.userID) return
  const res = await service({
    url: '/tenantMembership/assign',
    method: 'post',
    data: {
      userID: memberForm.value.userID,
      tenantID: activeTenantId.value,
      isPrimary: !!memberForm.value.isPrimary
    }
  })
  if (res?.code === 0) {
    await loadMembers()
  } else {
    ElMessage.error(res?.msg || t('admin.plugin.tenant.remove_failed'))
  }
}

const onUnassign = async (row) => {
  if (!activeTenantId.value || !row?.userID) return
  const res = await service({
    url: '/tenantMembership/unassign',
    method: 'delete',
    params: {
      userID: row.userID,
      tenantID: activeTenantId.value
    }
  })
  if (res?.code === 0) {
    await loadMembers()
  } else {
    ElMessage.error(res?.msg || t('admin.plugin.tenant.remove_failed'))
  }
}

const openUserPicker = async () => {
  userPickerVisible.value = true
  await loadUserPicker()
}

const loadUserPicker = async () => {
  const res = await service({
    url: '/user/getUserList',
    method: 'post',
    data: {
      page: userPickerPage.value,
      pageSize: userPickerPageSize.value,
      username: userPickerSearch.value.username || undefined,
      nickName: userPickerSearch.value.nickName || undefined
    },
    donNotShowLoading: true
  })
  if (res?.code === 0) {
    userPickerRows.value = res.data?.list || []
    userPickerTotal.value = res.data?.total || 0
  } else {
    ElMessage.error(res?.msg || t('admin.plugin.tenant.user_picker_load_failed'))
  }
}

const onUserPickerSearch = async () => {
  userPickerPage.value = 1
  await loadUserPicker()
}

const onUserPickerReset = async () => {
  userPickerSearch.value = { username: '', nickName: '' }
  userPickerPage.value = 1
  await loadUserPicker()
}

const onUserRowSelect = (row) => {
  memberForm.value.userID = row?.ID || 0
  memberForm.value.username = row?.userName || ''
  memberForm.value.nickName = row?.nickName || ''
  userPickerVisible.value = false
}

watch(activeTenantId, async () => {
  await loadMembers()
})

onMounted(async () => {
  // Ensure memberships are loaded so activeTenantId is meaningful.
  if (!tenantStore.loaded) {
    await tenantStore.loadMyTenants()
  }
  await loadMembers()
})
</script>

