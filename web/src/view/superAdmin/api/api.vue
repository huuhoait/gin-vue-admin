<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.superadmin.api.fields.path')">
          <el-input v-model="searchInfo.path" :placeholder="t('admin.superadmin.api.fields.path')" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.api.fields.description')">
          <el-input v-model="searchInfo.description" :placeholder="t('admin.superadmin.api.fields.description')" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.api.fields.group')">
          <el-select
            v-model="searchInfo.apiGroup"
            clearable
            :placeholder="t('admin.common.select')"
          >
            <el-option
              v-for="item in apiGroupOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.api.fields.method')">
          <el-select v-model="searchInfo.method" clearable :placeholder="t('admin.common.select')">
            <el-option
              v-for="item in methodOptions"
              :key="item.value"
              :label="`${item.label}(${item.value})`"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            {{ t('admin.common.search') }}
          </el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog('addApi')">
          {{ t('admin.common.add') }}
        </el-button>
        <el-button icon="delete" :disabled="!apis.length" @click="onDelete">
          {{ t('admin.common.delete') }}
        </el-button>
        <el-button icon="Refresh" @click="onFresh">{{ t('admin.superadmin.api.actions.refresh_cache') }}</el-button>
        <el-button icon="Compass" @click="onSync">{{ t('admin.superadmin.api.actions.sync_api') }}</el-button>
        <ExportTemplate template-id="api" />
        <ExportExcel template-id="api" :limit="9999" />
        <ImportExcel template-id="api" @on-success="getTableData" />
      </div>
      <el-table
        :data="tableData"
        @sort-change="sortChange"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column
          align="left"
          label="id"
          min-width="60"
          prop="ID"
          sortable="custom"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_path')"
          min-width="150"
          prop="path"
          sortable="custom"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_group')"
          min-width="150"
          prop="apiGroup"
          sortable="custom"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_desc')"
          min-width="150"
          prop="description"
          sortable="custom"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.method')"
          min-width="150"
          prop="method"
          sortable="custom"
        >
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>

        <el-table-column align="left" fixed="right" :label="t('admin.common.operation')" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button
              icon="edit"
              type="primary"
              link
              @click="editApiFunc(scope.row)"
            >
              {{ t('admin.common.edit') }}
            </el-button>
            <el-button
              icon="user"
              type="primary"
              link
              @click="openAssignRoleDrawer(scope.row)"
            >
              {{ t('admin.superadmin.api.actions.assign_roles') }}
            </el-button>
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteApiFunc(scope.row)"
            >
              {{ t('admin.common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      v-model="syncApiFlag"
      :size="appStore.drawerSize"
      :before-close="closeSyncDialog"
      :show-close="false"
    >
      <warning-bar
        :title="t('admin.superadmin.api.sync.warning')"
      />
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.superadmin.api.sync.header') }}</span>
          <div>
            <el-button :loading="apiCompletionLoading" @click="closeSyncDialog">
              {{ t('admin.common.cancel') }}
            </el-button>
            <el-button
              type="primary"
              :loading="syncing || apiCompletionLoading"
              @click="enterSyncDialog"
            >
              {{ t('admin.common.confirm') }}
            </el-button>
          </div>
        </div>
      </template>

      <h4>
        {{ t('admin.superadmin.api.sync.new_routes') }}
        <span class="text-xs text-gray-500 mx-2 font-normal"
          >{{ t('admin.superadmin.api.sync.new_routes_hint') }}</span
        >
        <el-button type="primary" size="small" @click="apiCompletion">
          <el-icon size="18">
            <ai-gva />
          </el-icon>
          {{ t('admin.superadmin.api.sync.auto_fill') }}
        </el-button>
      </h4>
      <el-table
        v-loading="syncing || apiCompletionLoading"
        :element-loading-text="t('admin.superadmin.api.sync.working')"
        :data="syncApiData.newApis"
      >
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_path')"
          min-width="150"
          prop="path"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_group')"
          min-width="150"
          prop="apiGroup"
        >
          <template #default="{ row }">
            <el-select
              v-model="row.apiGroup"
              :placeholder="t('admin.superadmin.api.sync.select_or_create')"
              allow-create
              filterable
              default-first-option
            >
              <el-option
                v-for="item in apiGroupOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_desc')"
          min-width="150"
          prop="description"
        >
          <template #default="{ row }">
            <el-input v-model="row.description" autocomplete="off" />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.method')"
          min-width="150"
          prop="method"
        >
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.common.operation')" min-width="150" fixed="right">
          <template #default="{ row }">
            <el-button icon="plus" type="primary" link @click="addApiFunc(row)">
              {{ t('admin.superadmin.api.sync.add_one') }}
            </el-button>
            <el-button
              icon="sunrise"
              type="primary"
              link
              @click="ignoreApiFunc(row, true)"
            >
              {{ t('admin.superadmin.api.sync.ignore') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <h4>
        {{ t('admin.superadmin.api.sync.removed_routes') }}
        <span class="text-xs text-gray-500 ml-2 font-normal"
          >{{ t('admin.superadmin.api.sync.removed_routes_hint') }}</span
        >
      </h4>
      <el-table :data="syncApiData.deleteApis">
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_path')"
          min-width="150"
          prop="path"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_group')"
          min-width="150"
          prop="apiGroup"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_desc')"
          min-width="150"
          prop="description"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.method')"
          min-width="150"
          prop="method"
        >
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>
      </el-table>

      <h4>
        {{ t('admin.superadmin.api.sync.ignored_routes') }}
        <span class="text-xs text-gray-500 ml-2 font-normal"
          >{{ t('admin.superadmin.api.sync.ignored_routes_hint') }}</span
        >
      </h4>
      <el-table :data="syncApiData.ignoreApis">
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_path')"
          min-width="150"
          prop="path"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_group')"
          min-width="150"
          prop="apiGroup"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.api_desc')"
          min-width="150"
          prop="description"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.api.columns.method')"
          min-width="150"
          prop="method"
        >
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.common.operation')" min-width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              icon="sunny"
              type="primary"
              link
              @click="ignoreApiFunc(row, false)"
            >
              {{ t('admin.superadmin.api.sync.unignore') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <el-drawer
      v-model="dialogFormVisible"
      :size="appStore.drawerSize"
      :before-close="closeDialog"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle }}</span>
          <div>
            <el-button @click="closeDialog">{{ t('admin.common.cancel') }}</el-button>
            <el-button type="primary" @click="enterDialog">{{ t('admin.common.confirm') }}</el-button>
          </div>
        </div>
      </template>

      <warning-bar :title="t('admin.superadmin.api.dialog.warning_add_api')" />
      <el-form ref="apiForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item :label="t('admin.superadmin.api.fields.path')" prop="path">
          <el-input v-model="form.path" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.api.fields.method')" prop="method">
          <el-select
            v-model="form.method"
            :placeholder="t('admin.common.select')"
            style="width: 100%"
          >
            <el-option
              v-for="item in methodOptions"
              :key="item.value"
              :label="`${item.label}(${item.value})`"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.api.fields.group')" prop="apiGroup">
          <el-select
            v-model="form.apiGroup"
            :placeholder="t('admin.superadmin.api.sync.select_or_create')"
            allow-create
            filterable
            default-first-option
          >
            <el-option
              v-for="item in apiGroupOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.api.fields.description')" prop="description">
          <el-input v-model="form.description" autocomplete="off" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Assign to roles drawer -->
    <el-drawer
      v-model="assignRoleDrawerVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.superadmin.api.actions.assign_roles') }} - {{ assignApiRow.description }}</span>
          <div>
            <el-button @click="assignRoleDrawerVisible = false">{{ t('admin.common.cancel') }}</el-button>
            <el-button type="primary" :loading="assignRoleSubmitting" @click="confirmAssignRole">{{ t('admin.common.confirm') }}</el-button>
          </div>
        </div>
      </template>
      <warning-bar :title="t('admin.superadmin.api.dialog.assign_warning')" />
      <el-tree
        ref="roleTreeRef"
        v-loading="assignRoleLoading"
        :data="authorityTreeData"
        :props="{ label: 'authorityName', children: 'children' }"
        node-key="authorityId"
        show-checkbox
        check-strictly
        default-expand-all
      />
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    getApiById,
    getApiList,
    createApi,
    updateApi,
    deleteApi,
    deleteApisByIds,
    freshCasbin,
    syncApi,
    getApiGroups,
    ignoreApi,
    enterSyncApi,
    getApiRoles,
    setApiRoles
  } from '@/api/api'
  import { getAuthorityList } from '@/api/authority'
  import { toSQLLine } from '@/utils/stringFun'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { ref, nextTick } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import ExportExcel from '@/components/exportExcel/exportExcel.vue'
  import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'
  import ImportExcel from '@/components/exportExcel/importExcel.vue'
  import { llmAuto } from '@/api/autoCode'
  import { useAppStore } from "@/pinia";

  defineOptions({
    name: 'Api'
  })

  const appStore = useAppStore()
  const { t } = useI18n()

  const methodFilter = (value) => {
    const target = methodOptions.value.filter((item) => item.value === value)[0]
    return target && `${target.label}`
  }

  const apis = ref([])
  const form = ref({
    path: '',
    apiGroup: '',
    method: '',
    description: ''
  })
  const methodOptions = ref([
    {
      value: 'POST',
      label: 'Create',
      type: 'success'
    },
    {
      value: 'GET',
      label: 'Read',
      type: ''
    },
    {
      value: 'PUT',
      label: 'Update',
      type: 'warning'
    },
    {
      value: 'DELETE',
      label: 'Delete',
      type: 'danger'
    }
  ])

  const type = ref('')
  const rules = ref({
    path: [{ required: true, message: () => t('admin.superadmin.api.rules.path_required'), trigger: 'blur' }],
    apiGroup: [{ required: true, message: () => t('admin.superadmin.api.rules.group_required'), trigger: 'blur' }],
    method: [{ required: true, message: () => t('admin.superadmin.api.rules.method_required'), trigger: 'blur' }],
    description: [{ required: true, message: () => t('admin.superadmin.api.rules.description_required'), trigger: 'blur' }]
  })

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})
  const apiGroupOptions = ref([])
  const apiGroupMap = ref({})

  const getGroup = async () => {
    const res = await getApiGroups()
    if (res.code === 0) {
      const groups = res.data.groups
      apiGroupOptions.value = groups.map((item) => ({
        label: item,
        value: item
      }))
      apiGroupMap.value = res.data.apiGroupMap
    }
  }

  const ignoreApiFunc = async (row, flag) => {
    const res = await ignoreApi({ path: row.path, method: row.method, flag })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg
      })
      if (flag) {
        syncApiData.value.newApis = syncApiData.value.newApis.filter(
          (item) => !(item.path === row.path && item.method === row.method)
        )
        syncApiData.value.ignoreApis.push(row)
        return
      }
      syncApiData.value.ignoreApis = syncApiData.value.ignoreApis.filter(
        (item) => !(item.path === row.path && item.method === row.method)
      )
      syncApiData.value.newApis.push(row)
    }
  }

  const addApiFunc = async (row) => {
    if (!row.apiGroup) {
      ElMessage({
        type: 'error',
        message: t('admin.superadmin.api.messages.select_group_first')
      })
      return
    }
    if (!row.description) {
      ElMessage({
        type: 'error',
        message: t('admin.superadmin.api.messages.fill_description_first')
      })
      return
    }
    const res = await createApi(row)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: t('admin.superadmin.api.messages.added_assign_hint'),
        showClose: true
      })
      syncApiData.value.newApis = syncApiData.value.newApis.filter(
        (item) => !(item.path === row.path && item.method === row.method)
      )
    }
    getTableData()
    getGroup()
  }

  const closeSyncDialog = () => {
    syncApiFlag.value = false
  }

  const syncing = ref(false)

  const enterSyncDialog = async () => {
    if (
      syncApiData.value.newApis.some(
        (item) => !item.apiGroup || !item.description
      )
    ) {
      ElMessage({
        type: 'error',
        message: t('admin.superadmin.api.messages.missing_group_or_desc')
      })
      return
    }

    syncing.value = true
    const res = await enterSyncApi(syncApiData.value)
    syncing.value = false
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg
      })
      syncApiFlag.value = false
      getTableData()
    }
  }

  const onReset = () => {
    searchInfo.value = {}
    getTableData()
  }
  // Search

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  // Pagination
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // Sorting
  const sortChange = ({ prop, order }) => {
    if (prop) {
      if (prop === 'ID') {
        prop = 'id'
      }
      searchInfo.value.orderKey = toSQLLine(prop)
      searchInfo.value.desc = order === 'descending'
    }
    getTableData()
  }

  // Query
  const getTableData = async () => {
    const table = await getApiList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()
  getGroup()
  // Batch operations
  const handleSelectionChange = (val) => {
    apis.value = val
  }

  const onDelete = async () => {
    ElMessageBox.confirm(t('admin.superadmin.api.messages.delete_selected_confirm'), t('admin.common.confirms.delete_title'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const ids = apis.value.map((item) => item.ID)
      const res = await deleteApisByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: res.msg
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }
  const onFresh = async () => {
    ElMessageBox.confirm(t('admin.superadmin.api.messages.refresh_cache_confirm'), t('admin.common.confirms.delete_title'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const res = await freshCasbin()
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: res.msg
        })
      }
    })
  }

  const syncApiData = ref({
    newApis: [],
    deleteApis: [],
    ignoreApis: []
  })

  const syncApiFlag = ref(false)

  const onSync = async () => {
    const res = await syncApi()
    if (res.code === 0) {
      res.data.newApis.forEach((item) => {
        item.apiGroup = apiGroupMap.value[item.path.split('/')[1]]
      })

      syncApiData.value = res.data
      syncApiFlag.value = true
    }
  }

  // Dialog
  const apiForm = ref(null)
  const initForm = () => {
    apiForm.value.resetFields()
    form.value = {
      path: '',
      apiGroup: '',
      method: '',
      description: ''
    }
  }

  const dialogTitle = ref(t('admin.superadmin.api.dialog.add_title'))
  const dialogFormVisible = ref(false)
  const openDialog = (key) => {
    switch (key) {
      case 'addApi':
        dialogTitle.value = t('admin.superadmin.api.dialog.add_title')
        break
      case 'edit':
        dialogTitle.value = t('admin.superadmin.api.dialog.edit_title')
        break
      default:
        break
    }
    type.value = key
    dialogFormVisible.value = true
  }
  const closeDialog = () => {
    initForm()
    dialogFormVisible.value = false
  }

  const editApiFunc = async (row) => {
    const res = await getApiById({ id: row.ID })
    form.value = res.data.api
    openDialog('edit')
  }

  const enterDialog = async () => {
    apiForm.value.validate(async (valid) => {
      if (valid) {
        switch (type.value) {
          case 'addApi':
            {
              const res = await createApi(form.value)
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: t('admin.common.messages.created'),
                  showClose: true
                })
              }
              getTableData()
              getGroup()
              closeDialog()
            }

            break
          case 'edit':
            {
              const res = await updateApi(form.value)
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: t('admin.common.messages.updated'),
                  showClose: true
                })
              }
              getTableData()
              closeDialog()
            }
            break
          default:
            {
              ElMessage({
                type: 'error',
                message: t('admin.superadmin.api.messages.unknown_operation'),
                showClose: true
              })
            }
            break
        }
      }
    })
  }

  const deleteApiFunc = async (row) => {
    ElMessageBox.confirm(t('admin.superadmin.api.messages.delete_confirm_detail'), t('admin.common.confirms.delete_title'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const res = await deleteApi(row)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: t('admin.common.messages.deleted')
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
        getGroup()
      }
    })
  }
  const apiCompletionLoading = ref(false)
  const apiCompletion = async () => {
    apiCompletionLoading.value = true
    const routerPaths = syncApiData.value.newApis
      .filter((item) => !item.apiGroup || !item.description)
      .map((item) => item.path)
    const res = await llmAuto({ data: String(routerPaths), mode: 'apiCompletion' })
    apiCompletionLoading.value = false
    if (res.code === 0) {
      try {
        const data = JSON.parse(res.data)
        syncApiData.value.newApis.forEach((item) => {
          const target = data.find((d) => d.path === item.path)
          if (target) {
            if (!item.apiGroup) {
              item.apiGroup = target.apiGroup
            }
            if (!item.description) {
              item.description = target.description
            }
          }
        })
      } catch (_) {
        ElMessage({
          type: 'error',
          message: t('admin.superadmin.api.messages.auto_fill_failed')
        })
      }
    }
  }

  // Assign roles
  const assignRoleDrawerVisible = ref(false)
  const assignApiRow = ref({})
  const authorityTreeData = ref([])
  const assignRoleLoading = ref(false)
  const assignRoleSubmitting = ref(false)
  const roleTreeRef = ref(null)

  const openAssignRoleDrawer = async (row) => {
    assignApiRow.value = row
    assignRoleDrawerVisible.value = true
    assignRoleLoading.value = true
    const [authRes, rolesRes] = await Promise.all([
      getAuthorityList(),
      getApiRoles(row.path, row.method)
    ])
    if (authRes.code === 0) {
      authorityTreeData.value = authRes.data
    }
    if (rolesRes.code === 0 && rolesRes.data) {
      nextTick(() => {
        roleTreeRef.value?.setCheckedKeys(rolesRes.data)
      })
    }
    assignRoleLoading.value = false
  }

  const confirmAssignRole = async () => {
    assignRoleSubmitting.value = true
    try {
      const checkedKeys = roleTreeRef.value?.getCheckedKeys(false) || []
      const res = await setApiRoles({
        path: assignApiRow.value.path,
        method: assignApiRow.value.method,
        authorityIds: checkedKeys
      })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: t('admin.superadmin.api.messages.assigned') })
        assignRoleDrawerVisible.value = false
      }
    } catch {
      ElMessage({ type: 'error', message: t('admin.superadmin.api.messages.assign_failed') })
    }
    assignRoleSubmitting.value = false
  }
</script>

<style scoped lang="scss">
  .warning {
    color: #dc143c;
  }
</style>
