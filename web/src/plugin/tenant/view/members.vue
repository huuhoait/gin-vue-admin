<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true">
        <el-form-item :label="t('admin.plugin.tenant.col_name')">
          <el-tag v-if="displayTenant" type="success">{{ displayTenant.name }}</el-tag>
          <el-tag v-else type="info">{{ t('admin.plugin.tenant.enabled_any') }}</el-tag>
        </el-form-item>
        <el-form-item v-if="isAdmin">
          <el-button type="primary" icon="office-building" @click="openTenantPicker">
            {{ t('admin.plugin.tenant.tenant_picker_button') }}
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button icon="refresh" @click="loadMembers">{{ t('admin.plugin.tenant.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-form :inline="true" :model="memberForm">
        <el-form-item :label="t('admin.plugin.tenant.member_user')">
          <div class="flex items-center gap-2 flex-wrap">
            <el-button icon="plus" type="success" plain :disabled="!activeTenantId" @click="openCreateUser">{{ t('admin.plugin.tenant.member_create_user') }}</el-button>
            <span v-if="memberForm.userID" class="text-sm text-gray-700">
              #{{ memberForm.userID }} · {{ memberForm.username }}<span v-if="memberForm.nickName"> ({{ memberForm.nickName }})</span>
            </span>
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
        <el-table-column :label="t('admin.plugin.tenant.col_actions')" width="260" fixed="right">
          <template #default="scope">
            <el-button link type="primary" @click="openEditUser(scope.row)">{{ t('admin.plugin.tenant.member_edit') }}</el-button>
            <el-button link type="warning" @click="openResetPassword(scope.row)">{{ t('admin.plugin.tenant.member_reset_password') }}</el-button>
            <el-button link type="danger" @click="onUnassign(scope.row)">{{ t('admin.plugin.tenant.member_remove') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog
      v-model="tenantPickerVisible"
      :title="t('admin.plugin.tenant.tenant_picker_title')"
      width="720"
      destroy-on-close
      append-to-body
    >
      <el-form :inline="true" :model="tenantPickerSearch">
        <el-form-item :label="t('admin.plugin.tenant.keyword_label')">
          <el-input v-model="tenantPickerSearch.keyword" :placeholder="t('admin.plugin.tenant.keyword_placeholder')" clearable @keyup.enter="onTenantPickerSearch" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onTenantPickerSearch">{{ t('admin.plugin.tenant.search') }}</el-button>
          <el-button icon="refresh" @click="onTenantPickerReset">{{ t('admin.plugin.tenant.reset') }}</el-button>
        </el-form-item>
      </el-form>
      <el-table :data="tenantPickerRows" size="small" style="width: 100%" highlight-current-row @row-click="onTenantRowSelect">
        <el-table-column :label="t('admin.plugin.tenant.col_code')" prop="code" width="160" />
        <el-table-column :label="t('admin.plugin.tenant.col_name')" prop="name" min-width="180" />
        <el-table-column :label="t('admin.plugin.tenant.col_enabled')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.enabled ? 'success' : 'danger'" size="small">
              {{ scope.row.enabled ? t('admin.plugin.tenant.yes') : t('admin.plugin.tenant.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.col_actions')" width="120" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click.stop="onTenantRowSelect(scope.row)">{{ t('admin.plugin.tenant.user_picker_select') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        :current-page="tenantPickerPage"
        :page-size="tenantPickerPageSize"
        :total="tenantPickerTotal"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        class="mt-2"
        @current-change="(v) => { tenantPickerPage = v; loadTenantPicker() }"
        @size-change="(v) => { tenantPickerPageSize = v; tenantPickerPage = 1; loadTenantPicker() }"
      />
    </el-dialog>

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

    <el-dialog
      v-model="createUserVisible"
      :title="t('admin.plugin.tenant.create_user_title')"
      width="520"
      destroy-on-close
      append-to-body
    >
      <el-form ref="createUserFormRef" :model="createUserForm" :rules="createUserRules" label-position="top">
        <el-form-item :label="t('admin.plugin.tenant.create_user_username')" prop="userName">
          <el-input v-model="createUserForm.userName" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.create_user_password')" prop="password">
          <el-input v-model="createUserForm.password" type="password" show-password autocomplete="new-password" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.create_user_nickname')">
          <el-input v-model="createUserForm.nickName" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.create_user_email')">
          <el-input v-model="createUserForm.email" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.create_user_phone')">
          <el-input v-model="createUserForm.phone" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.member_primary')">
          <el-switch v-model="createUserForm.isPrimary" />
        </el-form-item>
        <div class="text-xs text-gray-500">{{ t('admin.plugin.tenant.create_user_role_hint') }}</div>
      </el-form>
      <template #footer>
        <el-button @click="createUserVisible = false">{{ t('admin.plugin.tenant.cancel') }}</el-button>
        <el-button type="primary" :loading="createUserSubmitting" @click="onCreateUserSubmit">{{ t('admin.plugin.tenant.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="editUserVisible"
      :title="t('admin.plugin.tenant.edit_user_title')"
      width="520"
      destroy-on-close
      append-to-body
    >
      <el-form ref="editUserFormRef" :model="editUserForm" label-position="top">
        <el-form-item :label="t('admin.plugin.tenant.edit_user_username')">
          <el-input v-model="editUserForm.userName" disabled />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.edit_user_nickname')">
          <el-input v-model="editUserForm.nickName" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.edit_user_email')">
          <el-input v-model="editUserForm.email" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.edit_user_phone')">
          <el-input v-model="editUserForm.phone" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.edit_user_status')">
          <el-select v-model="editUserForm.enable" style="width: 100%">
            <el-option :value="1" :label="t('admin.plugin.tenant.edit_user_status_active')" />
            <el-option :value="2" :label="t('admin.plugin.tenant.edit_user_status_locked')" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editUserVisible = false">{{ t('admin.plugin.tenant.cancel') }}</el-button>
        <el-button type="primary" :loading="editUserSubmitting" @click="onEditUserSubmit">{{ t('admin.plugin.tenant.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="resetPasswordVisible"
      :title="t('admin.plugin.tenant.reset_password_title', { name: resetPasswordForm.userName })"
      width="480"
      destroy-on-close
      append-to-body
    >
      <el-form ref="resetPasswordFormRef" :model="resetPasswordForm" :rules="resetPasswordRules" label-position="top">
        <el-form-item :label="t('admin.plugin.tenant.reset_password_new')" prop="password">
          <el-input v-model="resetPasswordForm.password" type="password" show-password autocomplete="new-password" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.reset_password_confirm')" prop="confirm">
          <el-input v-model="resetPasswordForm.confirm" type="password" show-password autocomplete="new-password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPasswordVisible = false">{{ t('admin.plugin.tenant.cancel') }}</el-button>
        <el-button type="primary" :loading="resetPasswordSubmitting" @click="onResetPasswordSubmit">{{ t('admin.plugin.tenant.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useTenantStore } from '@/pinia/modules/tenant'
import { useUserStore } from '@/pinia/modules/user'
import service from '@/utils/request'
import { createUserInTenant, listTenants } from '@/plugin/tenant/api/tenant'

const { t } = useI18n()
const tenantStore = useTenantStore()
const userStore = useUserStore()
const activeTenantId = computed(() => Number(tenantStore.activeTenantId) || 0)
const activeTenant = computed(() => tenantStore.activeTenant)

// Super-admin (authorityId 888) can switch into any tenant, including ones
// they're not a member of — used to gate the inline tenant picker.
const isAdmin = computed(() => Number(userStore.userInfo?.authorityId) === 888)

// Local override so the tag can show the picked tenant's name even when the
// admin has no membership in it (which means it isn't in tenantStore.myTenants).
const pickedTenant = ref(null)
const displayTenant = computed(() => pickedTenant.value || activeTenant.value)

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

const tenantPickerVisible = ref(false)
const tenantPickerSearch = ref({ keyword: '' })
const tenantPickerRows = ref([])
const tenantPickerPage = ref(1)
const tenantPickerPageSize = ref(10)
const tenantPickerTotal = ref(0)

const createUserVisible = ref(false)
const createUserSubmitting = ref(false)
const createUserFormRef = ref(null)
const createUserForm = ref({ userName: '', password: '', nickName: '', email: '', phone: '', isPrimary: false })

const editUserVisible = ref(false)
const editUserSubmitting = ref(false)
const editUserFormRef = ref(null)
const editUserForm = ref({ ID: 0, userName: '', nickName: '', email: '', phone: '', enable: 1 })

const resetPasswordVisible = ref(false)
const resetPasswordSubmitting = ref(false)
const resetPasswordFormRef = ref(null)
const resetPasswordForm = ref({ ID: 0, userName: '', password: '', confirm: '' })
const resetPasswordRules = {
  password: [
    { required: true, message: () => t('admin.plugin.tenant.reset_password_required'), trigger: 'blur' },
    { min: 6, max: 128, message: () => t('admin.plugin.tenant.reset_password_length'), trigger: 'blur' }
  ],
  confirm: [
    {
      validator: (_rule, value, cb) => {
        if (value !== resetPasswordForm.value.password) {
          cb(new Error(t('admin.plugin.tenant.reset_password_mismatch')))
        } else {
          cb()
        }
      },
      trigger: 'blur'
    }
  ]
}
const createUserRules = {
  userName: [
    { required: true, message: () => t('admin.plugin.tenant.create_user_username_required'), trigger: 'blur' },
    { min: 3, max: 64, message: () => t('admin.plugin.tenant.create_user_username_length'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: () => t('admin.plugin.tenant.create_user_password_required'), trigger: 'blur' },
    { min: 6, max: 128, message: () => t('admin.plugin.tenant.create_user_password_length'), trigger: 'blur' }
  ]
}

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

const openTenantPicker = async () => {
  tenantPickerSearch.value = { keyword: '' }
  tenantPickerPage.value = 1
  tenantPickerVisible.value = true
  await loadTenantPicker()
}

const loadTenantPicker = async () => {
  const res = await listTenants({
    page: tenantPickerPage.value,
    pageSize: tenantPickerPageSize.value,
    keyword: tenantPickerSearch.value.keyword || undefined
  })
  if (res?.code === 0) {
    tenantPickerRows.value = res.data?.list || []
    tenantPickerTotal.value = res.data?.total || 0
  } else {
    ElMessage.error(res?.msg || t('admin.plugin.tenant.tenant_picker_load_failed'))
  }
}

const onTenantPickerSearch = async () => {
  tenantPickerPage.value = 1
  await loadTenantPicker()
}

const onTenantPickerReset = async () => {
  tenantPickerSearch.value = { keyword: '' }
  tenantPickerPage.value = 1
  await loadTenantPicker()
}

const onTenantRowSelect = async (row) => {
  if (!row?.ID) return
  pickedTenant.value = { ID: row.ID, name: row.name, code: row.code }
  tenantStore.setActiveTenant(row.ID)
  tenantPickerVisible.value = false
  await loadMembers()
}

const openCreateUser = () => {
  if (!activeTenantId.value) return
  createUserForm.value = { userName: '', password: '', nickName: '', email: '', phone: '', isPrimary: false }
  createUserVisible.value = true
}

const onCreateUserSubmit = async () => {
  if (!activeTenantId.value) return
  try {
    await createUserFormRef.value?.validate()
  } catch {
    return
  }
  createUserSubmitting.value = true
  try {
    const res = await createUserInTenant({
      tenantID: activeTenantId.value,
      userName: createUserForm.value.userName,
      password: createUserForm.value.password,
      nickName: createUserForm.value.nickName,
      email: createUserForm.value.email,
      phone: createUserForm.value.phone,
      isPrimary: createUserForm.value.isPrimary
    })
    if (res?.code === 0) {
      ElMessage.success(t('admin.plugin.tenant.create_user_success'))
      createUserVisible.value = false
      await loadMembers()
    } else {
      ElMessage.error(res?.msg || t('admin.plugin.tenant.create_user_failed'))
    }
  } finally {
    createUserSubmitting.value = false
  }
}

const openEditUser = (row) => {
  if (!row?.userID) return
  editUserForm.value = {
    ID: row.userID,
    userName: row.username || '',
    nickName: row.nickName || '',
    email: row.email || '',
    phone: row.phone || '',
    enable: row.enable === 2 ? 2 : 1
  }
  editUserVisible.value = true
}

const onEditUserSubmit = async () => {
  if (!editUserForm.value.ID) return
  editUserSubmitting.value = true
  try {
    const res = await service({
      url: '/user/setUserInfo',
      method: 'put',
      data: {
        ID: editUserForm.value.ID,
        nickName: editUserForm.value.nickName,
        email: editUserForm.value.email,
        phone: editUserForm.value.phone,
        enable: editUserForm.value.enable
      }
    })
    if (res?.code === 0) {
      ElMessage.success(t('admin.plugin.tenant.edit_user_success'))
      editUserVisible.value = false
      await loadMembers()
    } else {
      ElMessage.error(res?.msg || t('admin.plugin.tenant.edit_user_failed'))
    }
  } finally {
    editUserSubmitting.value = false
  }
}

const openResetPassword = (row) => {
  if (!row?.userID) return
  resetPasswordForm.value = {
    ID: row.userID,
    userName: row.username || '',
    password: '',
    confirm: ''
  }
  resetPasswordVisible.value = true
}

const onResetPasswordSubmit = async () => {
  if (!resetPasswordForm.value.ID) return
  try {
    await resetPasswordFormRef.value?.validate()
  } catch {
    return
  }
  resetPasswordSubmitting.value = true
  try {
    const res = await service({
      url: '/user/resetPassword',
      method: 'post',
      data: {
        ID: resetPasswordForm.value.ID,
        password: resetPasswordForm.value.password
      }
    })
    if (res?.code === 0) {
      ElMessage.success(t('admin.plugin.tenant.reset_password_success'))
      resetPasswordVisible.value = false
    } else {
      ElMessage.error(res?.msg || t('admin.plugin.tenant.reset_password_failed'))
    }
  } finally {
    resetPasswordSubmitting.value = false
  }
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

