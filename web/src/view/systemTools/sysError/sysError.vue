<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <el-form-item :label="t('admin.systemtools.sys_error.created_at')" prop="createdAtRange">
          <template #label>
            <span>
              {{ t('admin.systemtools.sys_error.created_at') }}
              <el-tooltip
                :content="t('admin.systemtools.sys_error.range_tip')"
              >
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>

          <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            :range-separator="t('admin.systemtools.sys_error.to')"
            :start-placeholder="t('admin.systemtools.sys_error.start_time')"
            :end-placeholder="t('admin.systemtools.sys_error.end_time')"
          />
        </el-form-item>

        <el-form-item :label="t('admin.systemtools.sys_error.source')" prop="form">
          <el-input v-model="searchInfo.form" :placeholder="t('admin.systemtools.sys_error.search_placeholder')" />
        </el-form-item>

        <el-form-item :label="t('admin.systemtools.sys_error.message')" prop="info">
          <el-input v-model="searchInfo.info" :placeholder="t('admin.systemtools.sys_error.search_placeholder')" />
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- Put additional query fields here when needed -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit"
            >{{ t('admin.systemtools.sys_error.search') }}</el-button
          >
          <el-button icon="refresh" @click="onReset">{{ t('admin.systemtools.sys_error.reset') }}</el-button>
          <el-button
            link
            type="primary"
            icon="arrow-down"
            @click="showAllQuery = true"
            v-if="!showAllQuery"
            >{{ t('admin.systemtools.sys_error.expand') }}</el-button
          >
          <el-button
            link
            type="primary"
            icon="arrow-up"
            @click="showAllQuery = false"
            v-else
            >{{ t('admin.systemtools.sys_error.collapse') }}</el-button
          >
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          icon="delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
          >{{ t('admin.systemtools.sys_error.delete') }}</el-button
        >
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

        <el-table-column
          sortable
          align="left"
          :label="t('admin.systemtools.sys_error.date')"
          prop="CreatedAt"
          width="180"
        >
          <template #default="scope">{{
            formatDate(scope.row.CreatedAt)
          }}</template>
        </el-table-column>

        <el-table-column
          align="left"
          :label="t('admin.systemtools.sys_error.source')"
          prop="form"
          width="120"
        />

        <el-table-column
          align="left"
          :label="t('admin.systemtools.sys_error.level')"
          prop="level"
          width="120"
        >
          <template #default="scope">
            <el-tag
              effect="dark"
              :type="levelTagMap[scope.row.level] || 'info'"
            >
              {{ levelLabelMap[scope.row.level] || defaultLevelLabel }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          :label="t('admin.systemtools.sys_error.status')"
          prop="status"
          width="140"
        >
          <template #default="scope">
            <el-tag
              effect="light"
              :type="statusTagMap[scope.row.status] || 'info'"
            >
              {{ statusLabelMap[scope.row.status] || defaultStatusLabel }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          :label="t('admin.systemtools.sys_error.message')"
          prop="info"
          show-overflow-tooltip
          width="240"
        />

        <el-table-column
          align="left"
          :label="t('admin.systemtools.sys_error.solution')"
          show-overflow-tooltip
          prop="solution"
          width="120"
        />

        <el-table-column
          align="left"
          :label="t('admin.systemtools.sys_error.actions')"
          fixed="right"
          :min-width="appStore.operateMinWith"
        >
          <template #default="scope">
            <el-button
              v-if="scope.row.status !== STATUS_PROCESSING"
              type="primary"
              link
              class="table-button"
              @click="getSolution(scope.row.ID)"
            >
              <el-icon><ai-gva /></el-icon>{{ t('admin.systemtools.sys_error.ai') }}
            </el-button>
            <el-button
              type="primary"
              link
              class="table-button"
              @click="getDetails(scope.row)"
              ><el-icon style="margin-right: 5px"><InfoFilled /></el-icon
              >{{ t('admin.systemtools.sys_error.view') }}</el-button
            >
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
              >{{ t('admin.systemtools.sys_error.delete') }}</el-button
            >
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
      destroy-on-close
      :size="appStore.drawerSize"
      v-model="detailShow"
      :show-close="true"
      :before-close="closeDetailShow"
      :title="t('admin.systemtools.sys_error.view')"
    >
      <el-descriptions :column="2" border direction="vertical">
        <el-descriptions-item :label="t('admin.systemtools.sys_error.source')">
          {{ detailForm.form }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('admin.systemtools.sys_error.level')">
          <el-tag
            effect="dark"
            :type="levelTagMap[detailForm.level] || 'info'"
          >
            {{ levelLabelMap[detailForm.level] || defaultLevelLabel }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="t('admin.systemtools.sys_error.status')">
          <el-tag
            effect="light"
            :type="statusTagMap[detailForm.status] || 'info'"
          >
            {{ statusLabelMap[detailForm.status] || defaultStatusLabel }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="t('admin.systemtools.sys_error.message')" :span="2">
          <pre class="whitespace-pre-wrap break-words">{{ detailForm.info }}</pre>
        </el-descriptions-item>
        <el-descriptions-item :label="t('admin.systemtools.sys_error.solution')" :span="2">
          <pre class="whitespace-pre-wrap break-words">{{ detailForm.solution }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    deleteSysError,
    deleteSysErrorByIds,
    findSysError,
    getSysErrorList,
    getSysErrorSolution
  } from '@/api/system/sysError'

  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, computed } from 'vue'
  import { useAppStore } from '@/pinia'
  import { useI18n } from 'vue-i18n'

  const { t } = useI18n()

  defineOptions({
    name: 'SysError'
  })

  const appStore = useAppStore()

  // Toggle advanced query fields
  const showAllQuery = ref(false)

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

  const getSolution = async (id) => {
    const confirmed = await ElMessageBox.confirm(
      t('admin.systemtools.sys_error.ai_confirm'),
      t('admin.systemtools.sys_error.notice_beta'),
      {
        confirmButtonText: t('admin.common.confirm'),
        cancelButtonText: t('admin.common.cancel'),
        type: 'warning'
      }
    ).catch(() => false)
    if (!confirmed) return
    const res = await getSysErrorSolution({ id })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: res.msg || t('admin.systemtools.sys_error.ai_submitted') })
      getTableData()
    }
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
    const table = await getSysErrorList({
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

  // ============== end table ===============

  // Optional dictionary init (keep if needed)
  const setOptions = async () => {}

  // Optional dictionary init (keep if needed)
  setOptions()

  // Selection
  const multipleSelection = ref([])
  // Selection change
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  // Delete row
  const deleteRow = (row) => {
    ElMessageBox.confirm(t('admin.systemtools.sys_error.delete_item_confirm'), t('admin.common.confirms.delete_title'), {
      confirmButtonText: t('admin.common.delete'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(() => {
      deleteSysErrorFunc(row)
    })
  }

  // Bulk delete
  const onDelete = async () => {
    ElMessageBox.confirm(t('admin.systemtools.sys_error.delete_selected_confirm'), t('admin.common.confirms.delete_title'), {
      confirmButtonText: t('admin.common.delete'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: t('admin.systemtools.sys_error.select_items_to_delete')
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map((item) => {
          IDs.push(item.ID)
        })
      const res = await deleteSysErrorByIds({ IDs })
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
  const deleteSysErrorFunc = async (row) => {
    const res = await deleteSysError({ ID: row.ID })
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
    const res = await findSysError({ ID: row.ID })
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

  const STATUS_PENDING = '\u672a\u5904\u7406'
  const STATUS_PROCESSING = '\u5904\u7406\u4e2d'
  const STATUS_DONE = '\u5904\u7406\u5b8c\u6210'
  const STATUS_FAILED = '\u5904\u7406\u5931\u8d25'

  const statusLabelMap = computed(() => ({
    [STATUS_PENDING]: t('admin.systemtools.sys_error.status_pending'),
    [STATUS_PROCESSING]: t('admin.systemtools.sys_error.status_processing'),
    [STATUS_DONE]: t('admin.systemtools.sys_error.status_done'),
    [STATUS_FAILED]: t('admin.systemtools.sys_error.status_failed')
  }))
  const statusTagMap = {
    [STATUS_PENDING]: 'info',
    [STATUS_PROCESSING]: 'warning',
    [STATUS_DONE]: 'success',
    [STATUS_FAILED]: 'danger'
  }
  const defaultStatusLabel = computed(() => t('admin.systemtools.sys_error.status_pending'))

  const levelLabelMap = computed(() => ({
    fatal: t('admin.systemtools.sys_error.level_fatal'),
    error: t('admin.systemtools.sys_error.level_error')
  }))
  const levelTagMap = {
    fatal: 'danger',
    error: 'warning'
  }
  const defaultLevelLabel = computed(() => t('admin.systemtools.sys_error.level_error'))
</script>
