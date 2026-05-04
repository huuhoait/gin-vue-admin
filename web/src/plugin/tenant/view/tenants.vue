<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.plugin.tenant.keyword_label')">
          <el-input v-model="searchInfo.keyword" :placeholder="t('admin.plugin.tenant.keyword_placeholder')" clearable @keyup.enter="onSubmit" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.enabled_label')">
          <el-select v-model="searchInfo.enabled" :placeholder="t('admin.plugin.tenant.enabled_any')" clearable style="width: 120px">
            <el-option :value="true" :label="t('admin.plugin.tenant.enabled_yes')" />
            <el-option :value="false" :label="t('admin.plugin.tenant.enabled_no')" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">{{ t('admin.plugin.tenant.search') }}</el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.plugin.tenant.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openCreate">{{ t('admin.plugin.tenant.new_tenant') }}</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" style="width: 100%">
        <el-table-column :label="t('admin.plugin.tenant.col_code')" prop="code" width="160" />
        <el-table-column :label="t('admin.plugin.tenant.col_name')" prop="name" min-width="180" />
        <el-table-column :label="t('admin.plugin.tenant.col_description')" prop="description" min-width="200" show-overflow-tooltip />
        <el-table-column :label="t('admin.plugin.tenant.col_enabled')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.enabled ? 'success' : 'danger'" size="small">
              {{ scope.row.enabled ? t('admin.plugin.tenant.yes') : t('admin.plugin.tenant.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.col_expire_at')" width="160">
          <template #default="scope">
            <span v-if="scope.row.expireAt">{{ formatDate(scope.row.expireAt) }}</span>
            <el-tag v-else type="info" size="small">{{ t('admin.plugin.tenant.expire_at_never') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.col_account_limit')" width="120">
          <template #default="scope">
            <span v-if="scope.row.accountLimit > 0">{{ scope.row.accountLimit }}</span>
            <el-tag v-else type="info" size="small">{{ t('admin.plugin.tenant.account_limit_unlimited') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.col_created')" width="160">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.col_actions')" width="240" fixed="right">
          <template #default="scope">
            <el-button link icon="user" @click="openMembers(scope.row)">{{ t('admin.plugin.tenant.members') }}</el-button>
            <el-button link icon="edit" @click="openEdit(scope.row)">{{ t('admin.plugin.tenant.edit') }}</el-button>
            <el-button link type="danger" icon="delete" @click="onDelete(scope.row)">{{ t('admin.plugin.tenant.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :page-sizes="[10, 30, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="(v) => { page = v; loadList() }"
        @size-change="(v) => { pageSize = v; page = 1; loadList() }"
      />
    </div>

    <el-drawer v-model="formDrawer" :title="formMode === 'create' ? t('admin.plugin.tenant.drawer_new') : t('admin.plugin.tenant.drawer_edit')" size="500" destroy-on-close>
      <el-form :model="formData" label-position="top">
        <el-form-item :label="t('admin.plugin.tenant.code_label')" required>
          <el-input v-model="formData.code" :disabled="formMode === 'edit'" :placeholder="t('admin.plugin.tenant.code_placeholder')" />
          <div class="text-xs text-gray-500 mt-1">{{ t('admin.plugin.tenant.code_hint') }}</div>
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.name_label')" required>
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.description_label')">
          <el-input v-model="formData.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.field_contact_name')">
          <el-input v-model="formData.contactName" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.field_contact_phone')">
          <el-input v-model="formData.contactPhone" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.field_domain')">
          <el-input v-model="formData.domain" :placeholder="t('admin.plugin.tenant.domain_placeholder')" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.field_expire_at')">
          <el-date-picker
            v-model="formData.expireAt"
            type="datetime"
            :placeholder="t('admin.plugin.tenant.expire_at_never')"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            clearable
            class="w-full"
          />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.field_account_limit')">
          <el-input-number v-model="formData.accountLimit" :min="0" />
          <div class="text-xs text-gray-500 mt-1">{{ t('admin.plugin.tenant.account_limit_hint') }}</div>
        </el-form-item>
        <el-form-item v-if="formMode === 'edit'" :label="t('admin.plugin.tenant.enabled_field')">
          <el-switch v-model="formData.enabled" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSave">{{ t('admin.plugin.tenant.save') }}</el-button>
          <el-button @click="formDrawer = false">{{ t('admin.plugin.tenant.cancel') }}</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer v-model="membersDrawer" :title="t('admin.plugin.tenant.members_drawer_title', { name: activeTenant?.name || '' })" size="600" destroy-on-close>
      <el-form :inline="true" :model="memberForm">
        <el-form-item :label="t('admin.plugin.tenant.member_user_id')">
          <el-input-number v-model="memberForm.userID" :min="1" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.member_primary')">
          <el-switch v-model="memberForm.isPrimary" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onAssign">{{ t('admin.plugin.tenant.member_assign') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="memberRows" size="small" style="width: 100%">
        <el-table-column :label="t('admin.plugin.tenant.member_user_id')" prop="userID" width="120" />
        <el-table-column :label="t('admin.plugin.tenant.member_primary')" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.isPrimary" type="success" size="small">{{ t('admin.plugin.tenant.primary_tag') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.member_joined')">
          <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.col_actions')" width="120">
          <template #default="scope">
            <el-button link type="danger" @click="onUnassign(scope.row)">{{ t('admin.plugin.tenant.member_remove') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/format'
import {
  listTenants,
  createTenant,
  updateTenant,
  deleteTenant,
  assignUser,
  unassignUser,
  membersOfTenant
} from '@/plugin/tenant/api/tenant'

defineOptions({ name: 'Tenants' })

const { t } = useI18n()
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchInfo = ref({ keyword: '', enabled: null })

const formDrawer = ref(false)
const formMode = ref('create')
const formData = ref(emptyForm())

const membersDrawer = ref(false)
const activeTenant = ref(null)
const memberRows = ref([])
const memberForm = ref({ userID: 1, isPrimary: false })

function emptyForm() {
  return {
    id: null,
    code: '',
    name: '',
    description: '',
    contactName: '',
    contactPhone: '',
    domain: '',
    expireAt: null,
    accountLimit: 0,
    enabled: true
  }
}

const loadList = async () => {
  const res = await listTenants({
    page: page.value,
    pageSize: pageSize.value,
    keyword: searchInfo.value.keyword,
    enabled: searchInfo.value.enabled
  })
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
  } else {
    ElMessage.error(res.msg || t('admin.plugin.tenant.load_failed'))
  }
}

const onSubmit = () => { page.value = 1; loadList() }
const onReset = () => { searchInfo.value = { keyword: '', enabled: null }; page.value = 1; loadList() }

const openCreate = () => {
  formMode.value = 'create'
  formData.value = emptyForm()
  formDrawer.value = true
}

const openEdit = (row) => {
  formMode.value = 'edit'
  formData.value = {
    id: row.ID,
    code: row.code,
    name: row.name,
    description: row.description || '',
    contactName: row.contactName || '',
    contactPhone: row.contactPhone || '',
    domain: row.domain || '',
    expireAt: row.expireAt || null,
    accountLimit: row.accountLimit ?? 0,
    enabled: row.enabled
  }
  formDrawer.value = true
}

const onSave = async () => {
  if (formMode.value === 'create') {
    const res = await createTenant({
      code: formData.value.code,
      name: formData.value.name,
      description: formData.value.description,
      contactName: formData.value.contactName,
      contactPhone: formData.value.contactPhone,
      domain: formData.value.domain,
      expireAt: formData.value.expireAt || null,
      accountLimit: formData.value.accountLimit ?? 0
    })
    if (res.code === 0) {
      ElMessage.success(t('admin.plugin.tenant.created_msg'))
      formDrawer.value = false
      loadList()
    } else {
      ElMessage.error(res.msg || t('admin.plugin.tenant.create_failed'))
    }
  } else {
    // Tri-state semantics for ExpireAt: when the picker is cleared we send
    // clearExpireAt=true so the backend NULLs the column instead of leaving
    // it untouched.
    const payload = {
      id: formData.value.id,
      name: formData.value.name,
      description: formData.value.description,
      contactName: formData.value.contactName,
      contactPhone: formData.value.contactPhone,
      domain: formData.value.domain,
      accountLimit: formData.value.accountLimit ?? 0,
      enabled: formData.value.enabled
    }
    if (formData.value.expireAt) {
      payload.expireAt = formData.value.expireAt
    } else {
      payload.clearExpireAt = true
    }
    const res = await updateTenant(payload)
    if (res.code === 0) {
      ElMessage.success(t('admin.plugin.tenant.updated_msg'))
      formDrawer.value = false
      loadList()
    } else {
      ElMessage.error(res.msg || t('admin.plugin.tenant.update_failed'))
    }
  }
}

const onDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.plugin.tenant.delete_confirm_body', { name: row.name }),
      t('admin.plugin.tenant.delete_confirm_title'),
      { type: 'warning' }
    )
  } catch { return }
  const res = await deleteTenant({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success(t('admin.plugin.tenant.deleted_msg'))
    loadList()
  } else {
    ElMessage.error(res.msg || t('admin.plugin.tenant.delete_failed'))
  }
}

const openMembers = async (row) => {
  activeTenant.value = row
  await refreshMembers()
  membersDrawer.value = true
}

const refreshMembers = async () => {
  if (!activeTenant.value) return
  const res = await membersOfTenant({ id: activeTenant.value.ID })
  if (res.code === 0) memberRows.value = res.data || []
  else ElMessage.error(res.msg || t('admin.plugin.tenant.members_load_failed'))
}

const onAssign = async () => {
  if (!activeTenant.value) return
  const res = await assignUser({
    userID: memberForm.value.userID,
    tenantID: activeTenant.value.ID,
    isPrimary: memberForm.value.isPrimary
  })
  if (res.code === 0) {
    ElMessage.success(t('admin.plugin.tenant.assigned_msg'))
    refreshMembers()
  } else {
    // Surface the dedicated cap-reached message when the backend signals it
    // via the i18n key. The server returns the localized text in `msg` and
    // the response interceptor surfaces it; we still provide a tailored
    // fallback so the operator always knows why the assign failed.
    const isLimit = typeof res.msg === 'string' && (
      res.msg.includes('account limit') || res.msg.includes('giới hạn')
    )
    ElMessage.error(
      res.msg || (isLimit
        ? t('admin.plugin.tenant.assign_failed_limit')
        : t('admin.plugin.tenant.assign_failed'))
    )
  }
}

const onUnassign = async (row) => {
  if (!activeTenant.value) return
  const res = await unassignUser({ userID: row.userID, tenantID: activeTenant.value.ID })
  if (res.code === 0) {
    ElMessage.success(t('admin.plugin.tenant.removed_msg'))
    refreshMembers()
  } else {
    ElMessage.error(res.msg || t('admin.plugin.tenant.remove_failed'))
  }
}

onMounted(loadList)
</script>
