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
        <el-form-item label="Created at" prop="createdAtRange">
          <template #label>
            <span>
              Created at
              <el-tooltip
                content="Range is start time (inclusive) to end time (exclusive)."
              >
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>

          <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="to"
            start-placeholder="Start time"
            end-placeholder="End time"
          />
        </el-form-item>

        <el-form-item label="Source" prop="form">
          <el-input v-model="searchInfo.form" placeholder="Search" />
        </el-form-item>

        <el-form-item label="Message" prop="info">
          <el-input v-model="searchInfo.info" placeholder="Search" />
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- Put additional query fields here when needed -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit"
            >Search</el-button
          >
          <el-button icon="refresh" @click="onReset">Reset</el-button>
          <el-button
            link
            type="primary"
            icon="arrow-down"
            @click="showAllQuery = true"
            v-if="!showAllQuery"
            >Expand</el-button
          >
          <el-button
            link
            type="primary"
            icon="arrow-up"
            @click="showAllQuery = false"
            v-else
            >Collapse</el-button
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
          >Delete</el-button
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
          label="Date"
          prop="CreatedAt"
          width="180"
        >
          <template #default="scope">{{
            formatDate(scope.row.CreatedAt)
          }}</template>
        </el-table-column>

        <el-table-column
          align="left"
          label="Source"
          prop="form"
          width="120"
        />

        <el-table-column
          align="left"
          label="Level"
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
          label="Status"
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
          label="Message"
          prop="info"
          show-overflow-tooltip
          width="240"
        />

        <el-table-column
          align="left"
          label="Solution"
          show-overflow-tooltip
          prop="solution"
          width="120"
        />

        <el-table-column
          align="left"
          label="Actions"
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
              <el-icon><ai-gva /></el-icon>AI
            </el-button>
            <el-button
              type="primary"
              link
              class="table-button"
              @click="getDetails(scope.row)"
              ><el-icon style="margin-right: 5px"><InfoFilled /></el-icon
              >View</el-button
            >
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
              >Delete</el-button
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
      title="View"
    >
      <el-descriptions :column="2" border direction="vertical">
        <el-descriptions-item label="Source">
          {{ detailForm.form }}
        </el-descriptions-item>
        <el-descriptions-item label="Level">
          <el-tag
            effect="dark"
            :type="levelTagMap[detailForm.level] || 'info'"
          >
            {{ levelLabelMap[detailForm.level] || defaultLevelLabel }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="Status">
          <el-tag
            effect="light"
            :type="statusTagMap[detailForm.status] || 'info'"
          >
            {{ statusLabelMap[detailForm.status] || defaultStatusLabel }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="Message" :span="2">
          <pre class="whitespace-pre-wrap break-words">{{ detailForm.info }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="Solution" :span="2">
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
  import { ref } from 'vue'
  import { useAppStore } from '@/pinia'

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
      'Logs will be sent via AI-PATH to GVA AI for analysis and temporarily stored on the official GVA platform as context. Continue? (Authorized users only)',
      'Notice (Beta)',
      {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    ).catch(() => false)
    if (!confirmed) return
    const res = await getSysErrorSolution({ id })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: res.msg || 'Submitted. It will complete in ~1 minute.' })
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
    ElMessageBox.confirm('Delete this item?', 'Confirm', {
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      type: 'warning'
    }).then(() => {
      deleteSysErrorFunc(row)
    })
  }

  // Bulk delete
  const onDelete = async () => {
    ElMessageBox.confirm('Delete selected items?', 'Confirm', {
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      type: 'warning'
    }).then(async () => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: 'Please select items to delete'
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
          message: 'Deleted'
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
        message: 'Deleted'
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

  const statusLabelMap = {
    [STATUS_PENDING]: 'Pending',
    [STATUS_PROCESSING]: 'Processing',
    [STATUS_DONE]: 'Done',
    [STATUS_FAILED]: 'Failed'
  }
  const statusTagMap = {
    [STATUS_PENDING]: 'info',
    [STATUS_PROCESSING]: 'warning',
    [STATUS_DONE]: 'success',
    [STATUS_FAILED]: 'danger'
  }
  const defaultStatusLabel = 'Pending'

  const levelLabelMap = {
    fatal: 'Fatal',
    error: 'Error'
  }
  const levelTagMap = {
    fatal: 'danger',
    error: 'warning'
  }
  const defaultLevelLabel = 'Error'
</script>
