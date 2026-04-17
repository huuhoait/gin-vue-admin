<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item :label="t('admin.plugin.announcement.created_at')" prop="createdAt">
          <template #label>
            <span>
              {{ t('admin.plugin.announcement.created_at') }}
              <el-tooltip
                :content="t('admin.plugin.announcement.created_at_tooltip')"
              >
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.startCreatedAt"
            type="datetime"
            :placeholder="t('admin.plugin.announcement.start_time')"
            :disabled-date="
              (time) =>
                searchInfo.endCreatedAt
                  ? time.getTime() > searchInfo.endCreatedAt.getTime()
                  : false
            "
          />
          —
          <el-date-picker
            v-model="searchInfo.endCreatedAt"
            type="datetime"
            :placeholder="t('admin.plugin.announcement.end_time')"
            :disabled-date="
              (time) =>
                searchInfo.startCreatedAt
                  ? time.getTime() < searchInfo.startCreatedAt.getTime()
                  : false
            "
          />
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- Add extra query conditions here if needed -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            {{ t('admin.common.search') }}
          </el-button>
          <el-button icon="refresh" @click="onReset"> {{ t('admin.common.reset') }} </el-button>
          <el-button
            v-if="!showAllQuery"
            link
            type="primary"
            icon="arrow-down"
            @click="showAllQuery = true"
          >
            {{ t('admin.plugin.announcement.expand') }}
          </el-button>
          <el-button
            v-else
            link
            type="primary"
            icon="arrow-up"
            @click="showAllQuery = false"
          >
            {{ t('admin.plugin.announcement.collapse') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">
          {{ t('admin.common.add') }}
        </el-button>
        <el-button
          icon="delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >
          {{ t('admin.common.delete') }}
        </el-button>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column align="left" :label="t('admin.plugin.announcement.created_at')" prop="CreatedAt" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>

        <el-table-column align="left" :label="t('admin.plugin.announcement.title')" prop="title" width="120" />
        <el-table-column align="left" :label="t('admin.plugin.announcement.author')" prop="userID" width="120">
          <template #default="scope">
            <span>{{
              filterDataSource(dataSource.userID, scope.row.userID)
            }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.announcement.attachments')" prop="attachments" width="200">
          <template #default="scope">
            <div class="file-list">
              <el-tag
                v-for="file in scope.row.attachments"
                :key="file.uid"
                @click="downloadFile(file.url)"
              >
                {{ file.name }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          :label="t('admin.common.actions')"
          fixed="right"
          min-width="240"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateInfoFunc(scope.row)"
            >
              {{ t('admin.common.edit') }}
            </el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
            >
              {{ t('admin.common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-drawer
      v-model="dialogFormVisible"
      destroy-on-close
      size="800"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? t('admin.plugin.announcement.type_create') : t('admin.plugin.announcement.type_edit') }}</span>
          <div>
            <el-button type="primary" @click="enterDialog"> {{ t('admin.common.confirm') }} </el-button>
            <el-button @click="closeDialog"> {{ t('admin.common.cancel') }} </el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="elFormRef"
        :model="formData"
        label-position="top"
        :rules="rule"
        label-width="80px"
      >
        <el-form-item :label="t('admin.plugin.announcement.title') + ':'" prop="title">
          <el-input
            v-model="formData.title"
            :clearable="true"
            :placeholder="t('admin.plugin.announcement.title_enter')"
          />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.announcement.content') + ':'" prop="content">
          <RichEdit v-model="formData.content" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.announcement.author') + ':'" prop="userID">
          <el-select
            v-model="formData.userID"
            :placeholder="t('admin.plugin.announcement.author_select')"
            style="width: 100%"
            :clearable="true"
          >
            <el-option
              v-for="(item, key) in dataSource.userID"
              :key="key"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.plugin.announcement.attachments') + ':'" prop="attachments">
          <SelectFile v-model="formData.attachments" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    getInfoDataSource,
    createInfo,
    deleteInfo,
    deleteInfoByIds,
    updateInfo,
    findInfo,
    getInfoList
  } from '@/plugin/announcement/api/info'
  import { getUrl } from '@/utils/image'
  // Rich text component
  import RichEdit from '@/components/richtext/rich-edit.vue'
  // File picker component
  import SelectFile from '@/components/selectFile/selectFile.vue'

  // Optional: import all formatters; keep as needed
  import { formatDate, filterDataSource } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import { ref, reactive } from 'vue'

  const { t } = useI18n()

  defineOptions({
    name: 'Info'
  })

  // Toggle additional query conditions
  const showAllQuery = ref(false)

  // Auto-generated dictionary (may be empty) and fields
  const formData = ref({
    title: '',
    content: '',
    userID: undefined,
    attachments: []
  })
  const dataSource = ref([])
  const getDataSourceFunc = async () => {
    const res = await getInfoDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

  // Validation rules
  const rule = reactive({})

  const searchRule = reactive({
    createdAt: [
      {
        validator: (rule, value, callback) => {
          if (
            searchInfo.value.startCreatedAt &&
            !searchInfo.value.endCreatedAt
          ) {
            callback(new Error(t('admin.plugin.announcement.end_time_required')))
          } else if (
            !searchInfo.value.startCreatedAt &&
            searchInfo.value.endCreatedAt
          ) {
            callback(new Error(t('admin.plugin.announcement.start_time_required')))
          } else if (
            searchInfo.value.startCreatedAt &&
            searchInfo.value.endCreatedAt &&
            (searchInfo.value.startCreatedAt.getTime() ===
              searchInfo.value.endCreatedAt.getTime() ||
              searchInfo.value.startCreatedAt.getTime() >
                searchInfo.value.endCreatedAt.getTime())
          ) {
            callback(new Error(t('admin.plugin.announcement.start_before_end')))
          } else {
            callback()
          }
        },
        trigger: 'change'
      }
    ]
  })

  const elFormRef = ref()
  const elSearchFormRef = ref()

  // =========== Table controls ===========
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

  // Change page size
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // Query
  const getTableData = async () => {
    const table = await getInfoList({
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

  // ============== End table controls ===============

  // Optional: load dictionaries if needed
  const setOptions = async () => {}

  // Optional: load dictionaries if needed
  setOptions()

  // Multi-select data
  const multipleSelection = ref([])
  // Multi-select
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  // Delete row
  const deleteRow = (row) => {
    ElMessageBox.confirm(t('admin.plugin.announcement.delete_confirm'), t('admin.common.confirms.delete_title'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(() => {
      deleteInfoFunc(row)
    })
  }

  // Delete selected
  const onDelete = async () => {
    ElMessageBox.confirm(t('admin.plugin.announcement.delete_selected_confirm'), t('admin.common.confirms.delete_title'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: t('admin.plugin.announcement.please_select_items')
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map((item) => {
          IDs.push(item.ID)
        })
      const res = await deleteInfoByIds({ IDs })
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

  // Action marker (add vs edit)
  const type = ref('')

  // Update row
  const updateInfoFunc = async (row) => {
    const res = await findInfo({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = res.data
      dialogFormVisible.value = true
    }
  }

  // Delete row
  const deleteInfoFunc = async (row) => {
    const res = await deleteInfo({ ID: row.ID })
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

  // Dialog control flag
  const dialogFormVisible = ref(false)

  // Open dialog
  const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
  }

  // Close dialog
  const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
      title: '',
      content: '',
      userID: undefined,
      attachments: []
    }
  }
  // Confirm dialog
  const enterDialog = async () => {
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      let res
      switch (type.value) {
        case 'create':
          res = await createInfo(formData.value)
          break
        case 'update':
          res = await updateInfo(formData.value)
          break
        default:
          res = await createInfo(formData.value)
          break
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: t('admin.common.messages.saved')
        })
        closeDialog()
        getTableData()
      }
    })
  }

  const downloadFile = (url) => {
    window.open(getUrl(url), '_blank')
  }
</script>

<style>
  .file-list {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }

  .fileBtn {
    margin-bottom: 10px;
  }

  .fileBtn:last-child {
    margin-bottom: 0;
  }
</style>
