<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.plugin.oauth2.keyword_label')">
          <el-input
            v-model="searchInfo.keyword"
            :placeholder="t('admin.plugin.oauth2.keyword_placeholder')"
            clearable
            @keyup.enter="onSubmit"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">{{ t('admin.plugin.oauth2.search') }}</el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.plugin.oauth2.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openCreate">{{ t('admin.plugin.oauth2.new_client') }}</el-button>
      </div>

      <el-table :data="tableData" row-key="ID" style="width: 100%">
        <el-table-column :label="t('admin.plugin.oauth2.col_name')" prop="name" min-width="160" />
        <el-table-column :label="t('admin.plugin.oauth2.col_client_id')" prop="clientID" width="280">
          <template #default="scope">
            <code class="text-xs">{{ scope.row.clientID }}</code>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.oauth2.col_grants')" min-width="180">
          <template #default="scope">
            <el-tag v-for="g in scope.row.grantTypes || []" :key="g" size="small" class="mr-1">
              {{ g }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.oauth2.col_enabled')" width="90">
          <template #default="scope">
            <el-tag :type="scope.row.enabled ? 'success' : 'danger'" size="small">
              {{ scope.row.enabled ? t('admin.plugin.oauth2.yes') : t('admin.plugin.oauth2.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.oauth2.col_created')" width="160">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.oauth2.col_actions')" width="280" fixed="right">
          <template #default="scope">
            <el-button link icon="edit" @click="openEdit(scope.row)">{{ t('admin.plugin.oauth2.edit') }}</el-button>
            <el-button link icon="key" @click="onRegenerate(scope.row)">{{ t('admin.plugin.oauth2.rotate_secret') }}</el-button>
            <el-button link type="danger" icon="delete" @click="onDelete(scope.row)">{{ t('admin.plugin.oauth2.delete') }}</el-button>
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

    <el-drawer v-model="drawerVisible" :title="formMode === 'create' ? t('admin.plugin.oauth2.drawer_new') : t('admin.plugin.oauth2.drawer_edit')" size="600" destroy-on-close>
      <el-form :model="formData" label-position="top">
        <el-form-item :label="t('admin.plugin.oauth2.name_label')" required>
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.oauth2.description_label')">
          <el-input v-model="formData.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.oauth2.redirect_uris_label')" required>
          <el-input v-model="formData.redirectUrisRaw" type="textarea" :rows="3" :placeholder="t('admin.plugin.oauth2.redirect_uri_placeholder')" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.oauth2.grant_types_label')">
          <el-checkbox-group v-model="formData.grantTypes">
            <el-checkbox label="authorization_code" />
            <el-checkbox label="refresh_token" />
            <el-checkbox label="client_credentials" />
          </el-checkbox-group>
        </el-form-item>
        <el-form-item :label="t('admin.plugin.oauth2.scopes_label')">
          <el-input v-model="formData.scopesRaw" :placeholder="t('admin.plugin.oauth2.scopes_placeholder')" />
        </el-form-item>
        <el-form-item v-if="formMode === 'edit'" :label="t('admin.plugin.oauth2.enabled_label')">
          <el-switch v-model="formData.enabled" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSave">{{ t('admin.plugin.oauth2.save') }}</el-button>
          <el-button @click="drawerVisible = false">{{ t('admin.plugin.oauth2.cancel') }}</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-dialog v-model="secretDialog" :title="t('admin.plugin.oauth2.secret_dialog_title')" width="540" destroy-on-close>
      <el-alert type="warning" :closable="false" show-icon>
        {{ t('admin.plugin.oauth2.secret_dialog_warning') }}
      </el-alert>
      <div class="mt-3">
        <div class="text-xs text-gray-500 mb-1">{{ t('admin.plugin.oauth2.field_client_id') }}</div>
        <el-input v-model="newCreds.clientId" readonly>
          <template #append><el-button @click="copy(newCreds.clientId)">{{ t('admin.plugin.oauth2.copy') }}</el-button></template>
        </el-input>
      </div>
      <div class="mt-3">
        <div class="text-xs text-gray-500 mb-1">{{ t('admin.plugin.oauth2.field_client_secret') }}</div>
        <el-input v-model="newCreds.clientSecret" readonly type="textarea" :rows="2">
          <template #append><el-button @click="copy(newCreds.clientSecret)">{{ t('admin.plugin.oauth2.copy') }}</el-button></template>
        </el-input>
      </div>
      <template #footer><el-button @click="secretDialog = false">{{ t('admin.plugin.oauth2.got_it') }}</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/format'
import {
  listOAuth2Clients,
  createOAuth2Client,
  updateOAuth2Client,
  deleteOAuth2Client,
  regenerateOAuth2ClientSecret
} from '@/plugin/oauth2server/api/client'

defineOptions({ name: 'OAuth2Clients' })

const { t } = useI18n()
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchInfo = ref({ keyword: '' })

const drawerVisible = ref(false)
const formMode = ref('create')
const formData = ref(emptyForm())
const secretDialog = ref(false)
const newCreds = ref({ clientId: '', clientSecret: '' })

function emptyForm() {
  return {
    id: null,
    name: '',
    description: '',
    redirectUrisRaw: '',
    grantTypes: ['authorization_code', 'refresh_token'],
    scopesRaw: '',
    enabled: true
  }
}

const loadList = async () => {
  const res = await listOAuth2Clients({
    page: page.value,
    pageSize: pageSize.value,
    keyword: searchInfo.value.keyword
  })
  if (res.code === 0) {
    const list = res.data.list || []
    tableData.value = list.map((c) => ({
      ...c,
      grantTypes: parseJSONArray(c.grantTypes)
    }))
    total.value = res.data.total || 0
  } else {
    ElMessage.error(res.msg || t('admin.plugin.oauth2.load_failed'))
  }
}

function parseJSONArray(v) {
  if (Array.isArray(v)) return v
  if (!v) return []
  try { return JSON.parse(v) } catch { return [] }
}

const onSubmit = () => { page.value = 1; loadList() }
const onReset = () => { searchInfo.value.keyword = ''; page.value = 1; loadList() }

const openCreate = () => {
  formMode.value = 'create'
  formData.value = emptyForm()
  drawerVisible.value = true
}

const openEdit = (row) => {
  formMode.value = 'edit'
  formData.value = {
    id: row.ID,
    name: row.name,
    description: row.description || '',
    redirectUrisRaw: parseJSONArray(row.redirectUris).join('\n'),
    grantTypes: row.grantTypes,
    scopesRaw: parseJSONArray(row.scopes).join(' '),
    enabled: row.enabled
  }
  drawerVisible.value = true
}

const onSave = async () => {
  const uris = formData.value.redirectUrisRaw.split('\n').map((s) => s.trim()).filter(Boolean)
  const scopes = formData.value.scopesRaw.split(/\s+/).filter(Boolean)
  const payload = {
    name: formData.value.name,
    description: formData.value.description,
    redirectUris: uris,
    grantTypes: formData.value.grantTypes,
    scopes
  }
  if (formMode.value === 'create') {
    const res = await createOAuth2Client(payload)
    if (res.code === 0) {
      drawerVisible.value = false
      newCreds.value = { clientId: res.data.clientId, clientSecret: res.data.clientSecret }
      secretDialog.value = true
      loadList()
    } else {
      ElMessage.error(res.msg || t('admin.plugin.oauth2.create_failed'))
    }
  } else {
    const res = await updateOAuth2Client({
      id: formData.value.id,
      ...payload,
      enabled: formData.value.enabled
    })
    if (res.code === 0) {
      ElMessage.success(t('admin.plugin.oauth2.updated_msg'))
      drawerVisible.value = false
      loadList()
    } else {
      ElMessage.error(res.msg || t('admin.plugin.oauth2.update_failed'))
    }
  }
}

const onDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.plugin.oauth2.delete_confirm_body', { name: row.name }),
      t('admin.plugin.oauth2.delete_confirm_title'),
      { type: 'warning' }
    )
  } catch { return }
  const res = await deleteOAuth2Client({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success(t('admin.plugin.oauth2.deleted_msg'))
    loadList()
  } else {
    ElMessage.error(res.msg || t('admin.plugin.oauth2.delete_failed'))
  }
}

const onRegenerate = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.plugin.oauth2.rotate_confirm_body'),
      t('admin.plugin.oauth2.rotate_confirm_title'),
      { type: 'warning' }
    )
  } catch { return }
  const res = await regenerateOAuth2ClientSecret({ id: row.ID })
  if (res.code === 0) {
    newCreds.value = { clientId: row.clientID, clientSecret: res.data.clientSecret }
    secretDialog.value = true
  } else {
    ElMessage.error(res.msg || t('admin.plugin.oauth2.rotate_failed'))
  }
}

const copy = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(t('admin.plugin.oauth2.copy_success'))
  } catch {
    ElMessage.warning(t('admin.plugin.oauth2.copy_failed'))
  }
}

onMounted(loadList)
</script>
