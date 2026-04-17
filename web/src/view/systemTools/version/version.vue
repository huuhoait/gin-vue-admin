<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline"
        @keyup.enter="onSubmit">
        <el-form-item :label="t('admin.systemtools.version.created_at')" prop="createdAtRange">
          <template #label>
            <span>
              {{ t('admin.systemtools.version.created_at') }}
              <el-tooltip :content="t('admin.systemtools.version.created_at_tooltip')">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>

          <el-date-picker v-model="searchInfo.createdAtRange" class="w-[380px]" type="datetimerange" :range-separator="t('admin.systemtools.version.range_to')"
            :start-placeholder="t('admin.systemtools.version.start_time')" :end-placeholder="t('admin.systemtools.version.end_time')" />
        </el-form-item>

        <el-form-item :label="t('admin.systemtools.version.version_name')" prop="versionName">
          <el-input v-model="searchInfo.versionName" :placeholder="t('admin.common.search')" />
        </el-form-item>

        <el-form-item :label="t('admin.systemtools.version.version_code')" prop="versionCode">
          <el-input v-model="searchInfo.versionCode" :placeholder="t('admin.common.search')" />
        </el-form-item>



        <template v-if="showAllQuery">
          <!-- Put additional query fields here when needed -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">{{ t('admin.common.search') }}</el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.common.reset') }}</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true"
            v-if="!showAllQuery">{{ t('admin.systemtools.version.expand') }}</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>{{ t('admin.systemtools.version.collapse') }}</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="success" icon="download" @click="openExportDialog">{{ t('admin.systemtools.version.create_release') }}</el-button>
        <el-button type="warning" icon="upload" @click="openImportDialog">{{ t('admin.systemtools.version.import_version') }}</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
          @click="onDelete">{{ t('admin.systemtools.version.delete') }}</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />

        <el-table-column sortable align="left" :label="t('admin.systemtools.version.date')" prop="CreatedAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" :label="t('admin.systemtools.version.version_name')" prop="versionName" width="120" />

        <el-table-column align="left" :label="t('admin.systemtools.version.version_code')" prop="versionCode" width="120" />

        <el-table-column align="left" :label="t('admin.systemtools.version.actions')" fixed="right" min-width="320">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>{{ t('admin.systemtools.version.view') }}</el-button>
            <el-button type="success" link icon="download" class="table-button"
              @click="downloadJson(scope.row)">{{ t('admin.systemtools.version.download_package') }}</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">{{ t('admin.systemtools.version.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true"
      :before-close="closeDetailShow" :title="t('admin.systemtools.version.view')">
      <el-descriptions :column="1" border>
        <el-descriptions-item :label="t('admin.systemtools.version.version_name')">
          {{ detailForm.versionName }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('admin.systemtools.version.version_code')">
          {{ detailForm.versionCode }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('admin.systemtools.version.description')">
          {{ detailForm.description }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>

    <!-- Export version drawer -->
    <el-drawer v-model="exportDialogVisible" :title="t('admin.systemtools.version.create_release')" direction="rtl" size="80%" :before-close="closeExportDialog"
      :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.systemtools.version.create_release') }}</span>
          <div>
            <el-button @click="closeExportDialog">{{ t('admin.systemtools.version.cancel') }}</el-button>
            <el-button type="primary" @click="handleExport" :loading="exportLoading">{{ t('admin.systemtools.version.create') }}</el-button>
          </div>
        </div>
      </template>
      <el-form :model="exportForm" label-width="100px">
        <el-form-item :label="t('admin.systemtools.version.version_name')" required>
          <el-input v-model="exportForm.versionName" :placeholder="t('admin.systemtools.version.enter_version_name')" />
        </el-form-item>
        <el-form-item :label="t('admin.systemtools.version.version_code')" required>
          <el-input v-model="exportForm.versionCode" :placeholder="t('admin.systemtools.version.enter_version_code')" />
        </el-form-item>
        <el-form-item :label="t('admin.systemtools.version.description')">
          <el-input v-model="exportForm.description" type="textarea" :placeholder="t('admin.systemtools.version.enter_description')" />
        </el-form-item>
        <el-form-item :label="t('admin.systemtools.version.release_content')">
          <div class="flex gap-3 w-full">
            <!-- Menu selection -->
            <div class="card-col card-vertical">
              <div class="card-header">
                <span class="card-title">{{ t('admin.systemtools.version.menus') }}</span>
              </div>
              <div class="card-filter">
                <el-input v-model="menuFilterText" :placeholder="t('admin.systemtools.version.filter_by_keyword')" clearable size="small" />
              </div>
              <div class="card-body">
                <el-tree ref="menuTreeRef" :data="menuTreeData" :default-checked-keys="selectedMenuIds"
                  :props="menuTreeProps" default-expand-all highlight-current node-key="ID" show-checkbox
                  :filter-node-method="filterMenuNode" @check="onMenuCheck" class="menu-tree">
                  <template #default="{ node }">
                    <span class="flex-1 flex items-center justify-between text-sm pr-2">
                      <span>{{ node.label }}</span>
                    </span>
                  </template>
                </el-tree>
              </div>
            </div>

            <!-- API selection -->
            <div class="card-col card-vertical">
              <div class="card-header">
                <span class="card-title">{{ t('admin.systemtools.version.apis') }}</span>
              </div>
              <div class="card-filter">
                <el-input v-model="apiFilterTextName" :placeholder="t('admin.systemtools.version.filter_by_name')" clearable size="small"
                  style="margin-bottom: 8px" />
                <el-input v-model="apiFilterTextPath" :placeholder="t('admin.systemtools.version.filter_by_path')" clearable size="small" />
              </div>
              <div class="card-body">
                <el-tree ref="apiTreeRef" :data="apiTreeData" :default-checked-keys="selectedApiIds"
                  :props="apiTreeProps" default-expand-all highlight-current node-key="onlyId" show-checkbox
                  :filter-node-method="filterApiNode" @check="onApiCheck" class="api-tree">
                  <template #default="{ data }">
                    <div class="flex items-center justify-between w-full pr-1">
                      <span>{{ data.description }}</span>
                      <el-tooltip :content="data.path">
                        <span
                          class="max-w-[240px] break-all overflow-ellipsis overflow-hidden text-gray-500 dark:text-gray-400">
                          {{ data.path }}
                        </span>
                      </el-tooltip>
                    </div>
                  </template>
                </el-tree>
              </div>
            </div>

            <!-- Dictionary selection -->
            <div class="card-col card-vertical">
              <div class="card-header">
                <span class="card-title">{{ t('admin.systemtools.version.dictionaries') }}</span>
              </div>
              <div class="card-filter">
                <el-input v-model="dictFilterText" :placeholder="t('admin.systemtools.version.filter_by_keyword')" clearable size="small" />
              </div>
              <div class="card-body">
                <el-tree ref="dictTreeRef" :data="dictTreeData" :default-checked-keys="selectedDictIds"
                  :props="dictTreeProps" default-expand-all highlight-current node-key="ID" show-checkbox
                  :filter-node-method="filterDictNode" @check="onDictCheck" class="dict-tree">
                  <template #default="{ data }">
                    <div class="flex items-center justify-between w-full pr-1">
                      <span>{{ data.name || data.label }}</span>
                      <el-tooltip :content="data.desc || (data.value ? `Value: ${data.value}` : '')">
                        <span class="text-gray-500 dark:text-gray-400 text-xs ml-2">
                          {{ data.type || (data.value ? `Value: ${data.value}` : '') }}
                        </span>
                      </el-tooltip>
                    </div>
                  </template>
                </el-tree>
              </div>
            </div>
          </div>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Import version drawer -->
    <el-drawer v-model="importDialogVisible" :title="t('admin.systemtools.version.import_version')" direction="rtl" size="80%" :before-close="closeImportDialog"
      :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.systemtools.version.import_version') }}</span>
          <div>
            <el-button @click="closeImportDialog">{{ t('admin.systemtools.version.cancel') }}</el-button>
            <el-button type="primary" @click="handleImport" :loading="importLoading"
              :disabled="!importJsonContent.trim()">{{ t('admin.systemtools.version.import') }}</el-button>
          </div>
        </div>
      </template>
      <el-form label-width="100px">
        <el-form-item :label="t('admin.systemtools.version.upload_file')">
          <el-upload ref="uploadRef" :auto-upload="false" :show-file-list="true" :limit="1" accept=".json"
            :on-change="handleFileChange" :on-remove="handleFileRemove" drag>
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              {{ t('admin.systemtools.version.drag_json_prefix') }}<em>{{ t('admin.systemtools.version.click_to_upload') }}</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                {{ t('admin.common.messages.json_only') }}
              </div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item :label="t('admin.systemtools.version.version_json')">
          <el-input v-model="importJsonContent" type="textarea" :rows="20" :placeholder="t('admin.systemtools.version.paste_version_json')"
            @input="handleJsonContentChange" />
        </el-form-item>
        <el-form-item :label="t('admin.systemtools.version.preview')" v-if="importPreviewData">
          <div class="preview-wrap">
            <div class="flex gap-3 w-full">
              <div class="card-col">
                <div class="card-vertical">
                  <div class="card-header">
                    <h3 class="card-title">{{ t('admin.systemtools.version.menus') }} ({{ getTotalMenuCount() }})</h3>
                  </div>
                  <div class="card-body">
                    <el-tree :data="previewMenuTreeData" :props="menuTreeProps" node-key="name"
                      :expand-on-click-node="false" :check-on-click-node="false" :show-checkbox="false"
                      default-expand-all>
                      <template #default="{ data }">
                        <div class="flex-1 flex items-center justify-between text-sm pr-2">
                          <span>{{ data.meta?.title || data.title }}</span>
                          <span class="text-gray-500 dark:text-gray-400 text-xs ml-2">{{ data.path }}</span>
                        </div>
                      </template>
                    </el-tree>
                  </div>
                </div>
              </div>
              <div class="card-col">
                <div class="card-vertical">
                  <div class="card-header">
                    <h3 class="card-title">{{ t('admin.systemtools.version.apis') }} ({{ importPreviewData.apis?.length || 0 }})</h3>
                  </div>
                  <div class="card-body">
                    <el-tree :data="previewApiTreeData" :props="apiTreeProps" node-key="ID"
                      :expand-on-click-node="false" :check-on-click-node="false" :show-checkbox="false"
                      default-expand-all>
                      <template #default="{ data }">
                        <div class="flex-1 flex items-center justify-between text-sm pr-2">
                          <span>{{ data.description }}</span>
                          <span class="text-gray-500 dark:text-gray-400 text-xs ml-2">{{ data.path }} [{{ data.method
                            }}]</span>
                        </div>
                      </template>
                    </el-tree>
                  </div>
                </div>
              </div>
              <div class="card-col">
                <div class="card-vertical">
                  <div class="card-header">
                    <h3 class="card-title">{{ t('admin.systemtools.version.dictionaries') }} ({{ importPreviewData.dictionaries?.length || 0 }})</h3>
                  </div>
                  <div class="card-body">
                    <el-tree :data="previewDictTreeData" :props="dictTreeProps" node-key="ID"
                      :expand-on-click-node="false" :check-on-click-node="false" :show-checkbox="false"
                      default-expand-all>
                      <template #default="{ data }">
                        <div class="flex-1 flex items-center justify-between text-sm pr-2">
                          <span>{{ data.name || data.label }}</span>
                          <span class="text-gray-500 dark:text-gray-400 text-xs ml-2">
                            {{ data.type || (data.value ? `Value: ${data.value}` : '') }}
                          </span>
                        </div>
                      </template>
                    </el-tree>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-form-item>
      </el-form>
    </el-drawer>

  </div>
</template>

<script setup>
import {
  deleteSysVersion,
  deleteSysVersionByIds,
  findSysVersion,
  getSysVersionList,
  exportVersion,
  importVersion,
  downloadVersionJson
} from '@/api/version'

// Menu/API/dictionary related APIs
import { getMenuList } from '@/api/menu'
import { getApiList } from '@/api/api'
import { getSysDictionaryList } from '@/api/sysDictionary'

// Import formatting utility (keep if needed)
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import { ref, watch } from 'vue'
import { useAppStore } from "@/pinia"
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineOptions({
  name: 'SysVersion'
})

const appStore = useAppStore()

// Toggle advanced query fields
const showAllQuery = ref(false)

// Export-related state
const exportDialogVisible = ref(false)
const exportLoading = ref(false)
const exportForm = ref({
  versionName: '',
  versionCode: '',
  description: '',
  menuIds: [],
  apiIds: [],
  dictIds: []
})

// Tree-related state
const menuTreeData = ref([])
const apiTreeData = ref([])
const dictTreeData = ref([])
const selectedMenuIds = ref([])
const selectedApiIds = ref([])
const selectedDictIds = ref([])
const menuFilterText = ref('')
const apiFilterTextName = ref('')
const apiFilterTextPath = ref('')
const dictFilterText = ref('')

// Tree refs
const menuTreeRef = ref(null)
const apiTreeRef = ref(null)
const dictTreeRef = ref(null)

// Tree props
const menuTreeProps = ref({
  children: 'children',
  label: function (data) {
    return data.meta?.title || data.title
  }
})

const apiTreeProps = ref({
  children: 'children',
  label: 'description'
})

const dictTreeProps = ref({
  children: 'sysDictionaryDetails',
  label: function (data) {
    // Dictionary root: show dictionary name
    if (data.name) {
      return data.name
    }
    // Dictionary item: show label
    if (data.label) {
      return data.label
    }
    return t('admin.systemtools.version.unknown')
  }
})

// Import-related state
const importDialogVisible = ref(false)
const importLoading = ref(false)
const importJsonContent = ref('')
const importPreviewData = ref(null)
const uploadRef = ref(null)
const previewMenuTreeData = ref([])
const previewApiTreeData = ref([])
const previewDictTreeData = ref([])



const elSearchFormRef = ref()

// =========== table ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
// Reset
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// Search
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

// Pagination
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// Page change
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// Query
const getTableData = async () => {
  const table = await getSysVersionList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== end table ===============

// Selection
const multipleSelection = ref([])
// Selection change
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// Delete row
const deleteRow = (row) => {
  ElMessageBox.confirm(t('admin.systemtools.version.delete_item_confirm'), t('admin.common.confirm'), {
    confirmButtonText: t('admin.systemtools.version.delete'),
    cancelButtonText: t('admin.common.cancel'),
    type: 'warning'
  }).then(() => {
    deleteSysVersionFunc(row)
  })
}

// Bulk delete
const onDelete = async () => {
  ElMessageBox.confirm(t('admin.systemtools.version.delete_selected_confirm'), t('admin.common.confirm'), {
    confirmButtonText: t('admin.systemtools.version.delete'),
    cancelButtonText: t('admin.common.cancel'),
    type: 'warning'
  }).then(async () => {
    const IDs = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: t('admin.systemtools.version.select_items_to_delete')
      })
      return
    }
    multipleSelection.value &&
      multipleSelection.value.map(item => {
        IDs.push(item.ID)
      })
    const res = await deleteSysVersionByIds({ IDs })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: t('admin.common.messages.deleted')
      })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// Delete row
const deleteSysVersionFunc = async (row) => {
  const res = await deleteSysVersion({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: t('admin.common.messages.deleted')
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

const detailForm = ref({})

// Detail drawer visibility
const detailShow = ref(false)


// Open drawer
const openDetailShow = () => {
  detailShow.value = true
}


// Load detail
const getDetails = async (row) => {
  // Load detail then open drawer
  const res = await findSysVersion({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// Close drawer
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}



// Load menu/API lists
const getMenuAndApiList = async () => {
  try {
    // Load menu list
    const menuRes = await getMenuList()
    if (menuRes.code === 0) {
      menuTreeData.value = menuRes.data || []
    }

    // Load API list
    const apiRes = await getApiList({ page: 1, pageSize: 9999 })
    if (apiRes.code === 0) {
      console.log('Raw API data:', apiRes.data)
      const apis = apiRes.data.list || []
      apiTreeData.value = buildApiTree(apis)
    }
  } catch (error) {
    console.error('Failed to load data:', error)
    ElMessage.error(t('admin.systemtools.version.load_menu_api_failed'))
  }
}

// Load dictionary list
const getDictList = async () => {
  try {
    const dictRes = await getSysDictionaryList({ page: 1, pageSize: 9999 })
    if (dictRes.code === 0) {
      dictTreeData.value = dictRes.data || []
    }
  } catch (error) {
    console.error('Failed to load dictionary data:', error)
    ElMessage.error(t('admin.systemtools.version.load_dict_failed'))
  }
}

// Build API tree
const buildApiTree = (apis) => {
  const apiObj = {}
  apis.forEach((item) => {
    item.onlyId = 'p:' + item.path + 'm:' + item.method
    if (Object.prototype.hasOwnProperty.call(apiObj, item.apiGroup)) {
      apiObj[item.apiGroup].push(item)
    } else {
      Object.assign(apiObj, { [item.apiGroup]: [item] })
    }
  })
  const apiTree = []
  for (const key in apiObj) {
    const treeNode = {
      ID: key,
      description: `${key} group`,
      children: apiObj[key]
    }
    apiTree.push(treeNode)
  }
  return apiTree
}

// Tree event handlers
const filterMenuNode = (value, data) => {
  if (!value) return true
  const title = data.meta?.title || data.title || ''
  return title.indexOf(value) !== -1
}

const filterApiNode = (value, data) => {
  if (!apiFilterTextName.value && !apiFilterTextPath.value) return true
  let matchesName, matchesPath
  if (!apiFilterTextName.value) {
    matchesName = true
  } else {
    matchesName = data.description && data.description.includes(apiFilterTextName.value)
  }
  if (!apiFilterTextPath.value) {
    matchesPath = true
  } else {
    matchesPath = data.path && data.path.includes(apiFilterTextPath.value)
  }
  return matchesName && matchesPath
}

const filterDictNode = (value, data) => {
  if (!value) return true
  const name = data.name || ''
  const type = data.type || ''
  const desc = data.desc || ''
  const label = data.label || ''
  const dataValue = data.value || ''
  return name.indexOf(value) !== -1 ||
    type.indexOf(value) !== -1 ||
    desc.indexOf(value) !== -1 ||
    label.indexOf(value) !== -1 ||
    dataValue.indexOf(value) !== -1
}

const onMenuCheck = (data, checked) => {
  if (checked.checkedKeys) {
    selectedMenuIds.value = checked.checkedKeys
  }
}

const onApiCheck = (data, checked) => {
  if (checked.checkedKeys) {
    selectedApiIds.value = checked.checkedKeys
  }
}

const onDictCheck = (data, checked) => {
  if (checked.checkedKeys) {
    selectedDictIds.value = checked.checkedKeys
  }
}

// Watch filter text changes
watch(menuFilterText, (val) => {
  if (menuTreeRef.value) {
    menuTreeRef.value.filter(val)
  }
})

watch([apiFilterTextName, apiFilterTextPath], () => {
  if (apiTreeRef.value) {
    apiTreeRef.value.filter('')
  }
})

watch(dictFilterText, (val) => {
  if (dictTreeRef.value) {
    dictTreeRef.value.filter(val)
  }
})

// Export handlers
const openExportDialog = async () => {
  exportDialogVisible.value = true
  await getMenuAndApiList()
  await getDictList()
}

const closeExportDialog = () => {
  exportDialogVisible.value = false
  exportForm.value = {
    versionName: '',
    versionCode: '',
    description: '',
    menuIds: [],
    apiIds: [],
    dictIds: []
  }
  selectedMenuIds.value = []
  selectedApiIds.value = []
  selectedDictIds.value = []
  menuFilterText.value = ''
  apiFilterTextName.value = ''
  apiFilterTextPath.value = ''
  dictFilterText.value = ''
}

const handleExport = async () => {
  if (!exportForm.value.versionName || !exportForm.value.versionCode) {
    ElMessage.warning(t('admin.systemtools.version.fill_version_name_code'))
    return
  }

  exportLoading.value = true
  try {
    // Collect selected menu/API/dictionary IDs
    const checkedMenus = menuTreeRef.value ? menuTreeRef.value.getCheckedNodes(false, true) : []
    const checkedApis = apiTreeRef.value ? apiTreeRef.value.getCheckedNodes(true) : []
    const checkedDicts = dictTreeRef.value ? dictTreeRef.value.getCheckedNodes(true) : []

    const menuIds = checkedMenus.map(menu => menu.ID)
    const apiIds = checkedApis.map(api => api.ID)
    const dictIds = checkedDicts.map(dict => dict.ID)

    exportForm.value.menuIds = menuIds
    exportForm.value.apiIds = apiIds
    exportForm.value.dictIds = dictIds

    const res = await exportVersion(exportForm.value)
    if (res.code !== 0) {
      ElMessage.error(res.msg || t('admin.systemtools.version.create_release_failed'))
      return
    }

    ElMessage.success(t('admin.systemtools.version.release_created'))
    closeExportDialog()
    getTableData() // Refresh table
  } catch (error) {
    console.error('Create release failed:', error)
    ElMessage.error(t('admin.systemtools.version.create_release_failed'))
  } finally {
    exportLoading.value = false
  }
}

// Import handlers
const openImportDialog = () => {
  importDialogVisible.value = true
}

const closeImportDialog = () => {
  importDialogVisible.value = false
  importJsonContent.value = ''
  importPreviewData.value = null
  previewMenuTreeData.value = []
  previewApiTreeData.value = []
  // Clear uploaded file
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
  }
}

// File upload handler
const handleFileChange = (file) => {
  if (!file.raw) return

  // Validate file type
  if (!file.name.toLowerCase().endsWith('.json')) {
    ElMessage.error(t('admin.common.messages.json_only'))
    uploadRef.value.clearFiles()
    return
  }

  // Read file content
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const content = e.target.result
      // Validate JSON
      JSON.parse(content)
      importJsonContent.value = content
      handleJsonContentChange()
      ElMessage.success(t('admin.systemtools.version.file_uploaded'))
    } catch (error) {
      ElMessage.error(t('admin.common.messages.invalid_json_file'))
      uploadRef.value.clearFiles()
    }
  }
  reader.readAsText(file.raw)
}

const handleFileRemove = () => {
  importJsonContent.value = ''
  importPreviewData.value = null
  previewMenuTreeData.value = []
  previewApiTreeData.value = []
}

// Count menus (recursive)
const getTotalMenuCount = () => {
  if (!importPreviewData.value?.menus) return 0

  const countMenus = (menus) => {
    let count = 0
    menus.forEach(menu => {
      count += 1 // current node
      if (menu.children && menu.children.length > 0) {
        count += countMenus(menu.children) // children
      }
    })
    return count
  }

  return countMenus(importPreviewData.value.menus)
}



const handleJsonContentChange = () => {
  if (!importJsonContent.value.trim()) {
    importPreviewData.value = null
    previewMenuTreeData.value = []
    previewApiTreeData.value = []
    previewDictTreeData.value = []
    return
  }

  try {
    const data = JSON.parse(importJsonContent.value)

    // Build preview data
    importPreviewData.value = {
      menus: data.menus || [],
      apis: data.apis || [],
      dictionaries: data.dictionaries || []
    }

    // Menus are already in a tree structure (with children)
    if (data.menus && data.menus.length > 0) {
      previewMenuTreeData.value = data.menus
    } else {
      previewMenuTreeData.value = []
    }

    // Build API tree grouped by apiGroup
    if (data.apis && data.apis.length > 0) {
      const apiGroups = {}
      data.apis.forEach(api => {
        const group = api.apiGroup || 'Ungrouped'
        if (!apiGroups[group]) {
          apiGroups[group] = {
            ID: `group_${group}`,
            description: group,
            path: '',
            method: '',
            children: []
          }
        }
        apiGroups[group].children.push(api)
      })
      previewApiTreeData.value = Object.values(apiGroups)
    } else {
      previewApiTreeData.value = []
    }

    // Handle dictionaries
    if (data.dictionaries && data.dictionaries.length > 0) {
      previewDictTreeData.value = data.dictionaries
    } else {
      previewDictTreeData.value = []
    }
  } catch (error) {
    console.error('JSON parse failed:', error)
    importPreviewData.value = null
    previewMenuTreeData.value = []
    previewApiTreeData.value = []
    previewDictTreeData.value = []
  }
}

const handleImport = async () => {
  if (!importJsonContent.value.trim()) {
    ElMessage.warning(t('admin.systemtools.version.provide_version_json'))
    return
  }

  try {
    JSON.parse(importJsonContent.value)
  } catch (error) {
    ElMessage.error(t('admin.common.messages.invalid_json'))
    return
  }

  importLoading.value = true
  try {
    const data = JSON.parse(importJsonContent.value)
    const res = await importVersion(data)
    if (res.code === 0) {
      ElMessage.success(t('admin.common.messages.imported'))
      closeImportDialog()
      getTableData() // Refresh table
    } else {
      ElMessage.error(res.msg || t('admin.systemtools.version.import_failed'))
    }
  } catch (error) {
    console.error('Import failed:', error)
    ElMessage.error(t('admin.systemtools.version.import_failed'))
  } finally {
    importLoading.value = false
  }
}

// Download version JSON
const downloadJson = async (row) => {
  try {
    const res = await downloadVersionJson({ ID: row.ID })
    // Handle axios response and resolve the blob
    // When responseType=blob, axios interceptors may return a full response object
    let blob
    if (res instanceof Blob) {
      blob = res
    } else if (res.data instanceof Blob) {
      blob = res.data
    } else {
      // If it's not a blob, it might be an error response; keep as-is.
      blob = res
    }

    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${row.versionName}_${row.versionCode}.json`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success(t('admin.common.messages.download_success'))
  } catch (error) {
    console.error('Download failed:', error)
    ElMessage.error(t('admin.systemtools.version.download_failed'))
  }
}

</script>

<style scoped>
/* Reusable card styles (dark-mode friendly) */
.card-col {
  @apply border border-gray-300 dark:border-gray-600 rounded overflow-hidden flex-1 bg-white dark:bg-gray-900;
}

.card-vertical {
  @apply flex flex-col h-full;
}

.card-header {
  @apply flex justify-between items-center px-4 py-3 bg-gray-50 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-600;
}

.card-title {
  @apply m-0 text-gray-800 dark:text-gray-200 text-base font-medium;
}

.card-filter {
  @apply px-4 py-3 border-b border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-800;
}

.card-body {
  @apply flex-1 p-2 min-h-[300px] max-h-[400px] overflow-y-auto;
}

.preview-wrap {
  @apply flex flex-col flex-1 gap-4 border border-gray-300 dark:border-gray-600 rounded p-4 bg-gray-50 dark:bg-gray-900;
}

/* Element Plus tree style tweaks */
:deep(.el-tree) {
  background-color: transparent;
}

:deep(.el-tree-node__content) {
  height: 32px;
  line-height: 32px;
}

:deep(.el-tree-node__label) {
  font-size: 14px;
}

:deep(.el-scrollbar__view) {
  padding: 0;
}
</style>
