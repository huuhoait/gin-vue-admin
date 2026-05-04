<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.plugin.tenant.package.keyword_label')">
          <el-input
            v-model="searchInfo.keyword"
            :placeholder="t('admin.plugin.tenant.package.keyword_placeholder')"
            clearable
            @keyup.enter="onSubmit"
          />
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
        <el-button type="primary" icon="plus" @click="openCreate">{{ t('admin.plugin.tenant.package.new_package') }}</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" style="width: 100%">
        <el-table-column :label="t('admin.plugin.tenant.package.col_code')" prop="code" width="160" />
        <el-table-column :label="t('admin.plugin.tenant.package.col_name')" prop="name" min-width="180" />
        <el-table-column :label="t('admin.plugin.tenant.package.col_description')" prop="description" min-width="220" show-overflow-tooltip />
        <el-table-column :label="t('admin.plugin.tenant.package.col_menus')" width="100">
          <template #default="scope">{{ countIDs(scope.row.menuIDs) }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.package.col_apis')" width="100">
          <template #default="scope">{{ countIDs(scope.row.apiIDs) }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.package.col_enabled')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.enabled ? 'success' : 'danger'" size="small">
              {{ scope.row.enabled ? t('admin.plugin.tenant.yes') : t('admin.plugin.tenant.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.col_created')" width="160">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.tenant.package.col_actions')" width="200" fixed="right">
          <template #default="scope">
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

    <el-drawer
      v-model="formDrawer"
      :title="formMode === 'create' ? t('admin.plugin.tenant.package.drawer_new') : t('admin.plugin.tenant.package.drawer_edit')"
      size="640"
      destroy-on-close
      @open="onDrawerOpen"
    >
      <el-form :model="formData" label-position="top">
        <el-form-item :label="t('admin.plugin.tenant.package.code_label')" required>
          <el-input
            v-model="formData.code"
            :disabled="formMode === 'edit'"
            :placeholder="t('admin.plugin.tenant.package.code_placeholder')"
          />
          <div class="text-xs text-gray-500 mt-1">{{ t('admin.plugin.tenant.package.code_hint') }}</div>
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.package.name_label')" required>
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.package.description_label')">
          <el-input v-model="formData.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.package.menus_label')">
          <div class="w-full border border-gray-200 rounded p-2 max-h-72 overflow-auto">
            <el-tree
              ref="menuTreeRef"
              :data="menuTreeData"
              :default-checked-keys="formData.menuIDs"
              node-key="ID"
              show-checkbox
              default-expand-all
              :props="{ children: 'children', label: (d) => d?.meta?.title || d?.name || '' }"
            />
          </div>
        </el-form-item>
        <el-form-item :label="t('admin.plugin.tenant.package.apis_label')">
          <el-select
            v-model="formData.apiIDs"
            multiple
            filterable
            collapse-tags
            collapse-tags-tooltip
            class="w-full"
          >
            <el-option-group
              v-for="g in apiGroupedOptions"
              :key="g.label"
              :label="g.label"
            >
              <el-option
                v-for="item in g.options"
                :key="item.ID"
                :value="item.ID"
                :label="`${item.method} ${item.path}`"
              >
                <span>{{ item.method }} {{ item.path }}</span>
                <span class="text-xs text-gray-400 ml-2">{{ item.description }}</span>
              </el-option>
            </el-option-group>
          </el-select>
        </el-form-item>
        <el-form-item v-if="formMode === 'edit'" :label="t('admin.plugin.tenant.package.enabled_label')">
          <el-switch v-model="formData.enabled" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSave">{{ t('admin.plugin.tenant.package.save') }}</el-button>
          <el-button @click="formDrawer = false">{{ t('admin.plugin.tenant.package.cancel') }}</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/format'
import { getBaseMenuTree } from '@/api/menu'
import { getAllApis } from '@/api/api'
import {
  listTenantPackages,
  createTenantPackage,
  updateTenantPackage,
  deleteTenantPackage
} from '@/plugin/tenant/api/package'

defineOptions({ name: 'TenantPackages' })

const { t } = useI18n()

const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchInfo = ref({ keyword: '', enabled: null })

const formDrawer = ref(false)
const formMode = ref('create')
const formData = ref(emptyForm())
const menuTreeRef = ref(null)
const menuTreeData = ref([])
const apiList = ref([])

function emptyForm() {
  return { id: null, code: '', name: '', description: '', menuIDs: [], apiIDs: [], enabled: true }
}

// countIDs renders the size of a JSON-encoded id array. Backend stores
// menuIDs / apiIDs as datatypes.JSON, which serialises through to []number.
const countIDs = (val) => {
  if (Array.isArray(val)) return val.length
  if (val == null || val === '') return 0
  try {
    const arr = typeof val === 'string' ? JSON.parse(val) : val
    return Array.isArray(arr) ? arr.length : 0
  } catch {
    return 0
  }
}

// Group APIs by apiGroup for visual organisation in the multi-select.
const apiGroupedOptions = computed(() => {
  const groups = {}
  for (const api of apiList.value) {
    const key = api.apiGroup || 'default'
    if (!groups[key]) groups[key] = { label: key, options: [] }
    groups[key].options.push(api)
  }
  return Object.values(groups)
})

const loadList = async () => {
  const res = await listTenantPackages({
    page: page.value,
    pageSize: pageSize.value,
    keyword: searchInfo.value.keyword,
    enabled: searchInfo.value.enabled
  })
  if (res.code === 0) {
    const list = res.data.list || []
    // Normalise menuIDs/apiIDs from backend (datatypes.JSON serialises as
    // either an array or a JSON-encoded string depending on dialect).
    tableData.value = list.map((row) => ({
      ...row,
      menuIDs: parseIDs(row.menuIDs),
      apiIDs: parseIDs(row.apiIDs)
    }))
    total.value = res.data.total || 0
  } else {
    ElMessage.error(res.msg || t('admin.plugin.tenant.package.load_failed'))
  }
}

const parseIDs = (val) => {
  if (Array.isArray(val)) return val
  if (val == null || val === '') return []
  try {
    const arr = typeof val === 'string' ? JSON.parse(val) : val
    return Array.isArray(arr) ? arr : []
  } catch {
    return []
  }
}

const onSubmit = () => { page.value = 1; loadList() }
const onReset = () => { searchInfo.value = { keyword: '', enabled: null }; page.value = 1; loadList() }

const ensureSelectorData = async () => {
  if (menuTreeData.value.length === 0) {
    try {
      const res = await getBaseMenuTree()
      menuTreeData.value = res?.data?.menus || []
    } catch {
      menuTreeData.value = []
    }
  }
  if (apiList.value.length === 0) {
    try {
      const res = await getAllApis()
      apiList.value = res?.data?.apis || []
    } catch {
      apiList.value = []
    }
  }
}

const openCreate = async () => {
  formMode.value = 'create'
  formData.value = emptyForm()
  await ensureSelectorData()
  formDrawer.value = true
}

const openEdit = async (row) => {
  formMode.value = 'edit'
  formData.value = {
    id: row.ID,
    code: row.code,
    name: row.name,
    description: row.description || '',
    menuIDs: parseIDs(row.menuIDs),
    apiIDs: parseIDs(row.apiIDs),
    enabled: row.enabled
  }
  await ensureSelectorData()
  formDrawer.value = true
}

// onDrawerOpen syncs the el-tree's checked state to the current form value.
// el-tree honours default-checked-keys only on first mount, so manual sync
// is needed when reopening the drawer for a different row.
const onDrawerOpen = () => {
  // Defer to next tick so the tree component is mounted.
  setTimeout(() => {
    if (menuTreeRef.value) {
      menuTreeRef.value.setCheckedKeys(formData.value.menuIDs || [])
    }
  }, 0)
}

const collectMenuIDs = () => {
  if (!menuTreeRef.value) return formData.value.menuIDs || []
  // Include half-checked (parent of partially-selected children) so the
  // parent menu is granted along with its leaves.
  const checked = menuTreeRef.value.getCheckedKeys() || []
  const half = menuTreeRef.value.getHalfCheckedKeys() || []
  return [...checked, ...half]
}

const onSave = async () => {
  const menuIDs = collectMenuIDs()
  const apiIDs = formData.value.apiIDs || []
  if (formMode.value === 'create') {
    const res = await createTenantPackage({
      code: formData.value.code,
      name: formData.value.name,
      description: formData.value.description,
      menuIDs,
      apiIDs
    })
    if (res.code === 0) {
      ElMessage.success(t('admin.plugin.tenant.package.created_msg'))
      formDrawer.value = false
      loadList()
    } else {
      ElMessage.error(res.msg || t('admin.plugin.tenant.package.create_failed'))
    }
  } else {
    const res = await updateTenantPackage({
      id: formData.value.id,
      name: formData.value.name,
      description: formData.value.description,
      menuIDs,
      apiIDs,
      enabled: formData.value.enabled
    })
    if (res.code === 0) {
      ElMessage.success(t('admin.plugin.tenant.package.updated_msg'))
      formDrawer.value = false
      loadList()
    } else {
      ElMessage.error(res.msg || t('admin.plugin.tenant.package.update_failed'))
    }
  }
}

const onDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.plugin.tenant.package.delete_confirm_body', { name: row.name }),
      t('admin.plugin.tenant.package.delete_confirm_title'),
      { type: 'warning' }
    )
  } catch { return }
  const res = await deleteTenantPackage({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success(t('admin.plugin.tenant.package.deleted_msg'))
    loadList()
  } else {
    ElMessage.error(res.msg || t('admin.plugin.tenant.package.delete_failed'))
  }
}

onMounted(loadList)
</script>
