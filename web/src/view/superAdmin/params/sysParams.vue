<template>
  <div>
    <warning-bar title="Parameter caching is already implemented in frontend utils/params. See comments in that file for usage." />
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="Created at" prop="createdAt">
          <template #label>
            <span>
              Created at
              <el-tooltip
                content="Search range: start time (inclusive) to end time (exclusive)"
              >
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.startCreatedAt"
            type="datetime"
            placeholder="Start time"
            :disabled-date="
              (time) =>
                searchInfo.endCreatedAt
                  ? time.getTime() > searchInfo.endCreatedAt.getTime()
                  : false
            "
          ></el-date-picker>
          —
          <el-date-picker
            v-model="searchInfo.endCreatedAt"
            type="datetime"
            placeholder="End time"
            :disabled-date="
              (time) =>
                searchInfo.startCreatedAt
                  ? time.getTime() < searchInfo.startCreatedAt.getTime()
                  : false
            "
          ></el-date-picker>
        </el-form-item>

        <el-form-item label="Name" prop="name">
          <el-input v-model="searchInfo.name" placeholder="Search" />
        </el-form-item>
        <el-form-item label="Key" prop="key">
          <el-input v-model="searchInfo.key" placeholder="Search" />
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- Add extra query conditions here if needed -->
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
        <el-button type="primary" icon="plus" @click="openDialog"
          >Add</el-button
        >
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

        <el-table-column align="left" label="Created at" prop="createdAt" width="180">
          <template #default="scope">{{
            formatDate(scope.row.CreatedAt)
          }}</template>
        </el-table-column>

        <el-table-column
          align="left"
          label="Name"
          prop="name"
          width="120"
        />
        <el-table-column align="left" label="Key" prop="key" width="120" />
        <el-table-column align="left" label="Value" prop="value" width="120" />
        <el-table-column
          align="left"
          label="Description"
          prop="desc"
          width="120"
        />
        <el-table-column
          align="left"
          label="Actions"
          fixed="right"
          min-width="240"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              class="table-button"
              @click="getDetails(scope.row)"
              ><el-icon style="margin-right: 5px"><InfoFilled /></el-icon
              >Details</el-button
            >
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateSysParamsFunc(scope.row)"
              >Edit</el-button
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
      size="800"
      v-model="dialogFormVisible"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Add' : 'Edit' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">Confirm</el-button>
            <el-button @click="closeDialog">Cancel</el-button>
          </div>
        </div>
      </template>

      <el-form
        :model="formData"
        label-position="top"
        ref="elFormRef"
        :rules="rule"
        label-width="80px"
      >
        <el-form-item label="Name:" prop="name">
          <el-input
            v-model="formData.name"
            :clearable="true"
            placeholder="Enter name"
          />
        </el-form-item>
        <el-form-item label="Key:" prop="key">
          <el-input
            v-model="formData.key"
            :clearable="true"
            placeholder="Enter key"
          />
        </el-form-item>
        <el-form-item label="Value:" prop="value">
          <el-input
            type="textarea"
            :rows="5"
            v-model="formData.value"
            :clearable="true"
            placeholder="Enter value"
          />
        </el-form-item>
        <el-form-item label="Description:" prop="desc">
          <el-input
            v-model="formData.desc"
            :clearable="true"
            placeholder="Enter description"
          />
        </el-form-item>
      </el-form>

      <div
        class="usage-instructions bg-gray-100 border border-gray-300 rounded-lg p-4 mt-5"
      >
        <h3 class="mb-3 text-lg text-gray-800">Usage</h3>
        <p class="mb-2 text-sm text-gray-600">
          On the frontend, import
          <code class="bg-blue-100 px-1 py-0.5 rounded"
            >import { getParams } from '@/utils/params'</code
          >
          then call
          <code class="bg-blue-100 px-1 py-0.5 rounded"
            >await getParams("{{ formData.key }}")</code
          >
          to fetch the parameter.
        </p>
        <p class="text-sm text-gray-600">
          On the backend, register
          <code class="bg-blue-100 px-1 py-0.5 rounded"
            >import
            "github.com/flipped-aurora/gin-vue-admin/server/service/system"</code
          >
        </p>
        <p class="mb-2 text-sm text-gray-600">
          then call
          <code class="bg-blue-100 px-1 py-0.5 rounded"
            >new(system.SysParamsService).GetSysParam("{{
              formData.key
            }}")</code
          >
          to fetch the value.
        </p>
      </div>
    </el-drawer>

    <el-drawer
      destroy-on-close
      size="800"
      v-model="detailShow"
      :show-close="true"
      :before-close="closeDetailShow"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="Name">
          {{ detailForm.name }}
        </el-descriptions-item>
        <el-descriptions-item label="Key">
          {{ detailForm.key }}
        </el-descriptions-item>
        <el-descriptions-item label="Value">
          {{ detailForm.value }}
        </el-descriptions-item>
        <el-descriptions-item label="Description">
          {{ detailForm.desc }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createSysParams,
    deleteSysParams,
    deleteSysParamsByIds,
    updateSysParams,
    findSysParams,
    getSysParamsList
  } from '@/api/sysParams'

  // Optional: import all formatters; keep as needed
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, reactive } from 'vue'
  import WarningBar from "@/components/warningBar/warningBar.vue";

  defineOptions({
    name: 'SysParams'
  })

  // Toggle additional query conditions
  const showAllQuery = ref(false)

  // Auto-generated dictionary (may be empty) and fields
  const formData = ref({
    name: '',
    key: '',
    value: '',
    desc: ''
  })

  // Validation rules
  const rule = reactive({
    name: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: 'Cannot be only whitespace',
        trigger: ['input', 'blur']
      }
    ],
    key: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: 'Cannot be only whitespace',
        trigger: ['input', 'blur']
      }
    ],
    value: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: 'Cannot be only whitespace',
        trigger: ['input', 'blur']
      }
    ]
  })

  const searchRule = reactive({
    createdAt: [
      {
        validator: (rule, value, callback) => {
          if (
            searchInfo.value.startCreatedAt &&
            !searchInfo.value.endCreatedAt
          ) {
            callback(new Error('End time is required'))
          } else if (
            !searchInfo.value.startCreatedAt &&
            searchInfo.value.endCreatedAt
          ) {
            callback(new Error('Start time is required'))
          } else if (
            searchInfo.value.startCreatedAt &&
            searchInfo.value.endCreatedAt &&
            (searchInfo.value.startCreatedAt.getTime() ===
              searchInfo.value.endCreatedAt.getTime() ||
              searchInfo.value.startCreatedAt.getTime() >
                searchInfo.value.endCreatedAt.getTime())
          ) {
            callback(new Error('Start time must be earlier than end time'))
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
    const table = await getSysParamsList({
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
    ElMessageBox.confirm('Delete this item?', 'Confirm', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning'
    }).then(() => {
      deleteSysParamsFunc(row)
    })
  }

  // Delete selected
  const onDelete = async () => {
    ElMessageBox.confirm('Delete selected items?', 'Confirm', {
      confirmButtonText: 'Confirm',
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
      const res = await deleteSysParamsByIds({ IDs })
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

  // Action marker (add vs edit)
  const type = ref('')

  // Update row
  const updateSysParamsFunc = async (row) => {
    const res = await findSysParams({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = res.data
      dialogFormVisible.value = true
    }
  }

  // Delete row
  const deleteSysParamsFunc = async (row) => {
    const res = await deleteSysParams({ ID: row.ID })
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
      name: '',
      key: '',
      value: '',
      desc: ''
    }
  }
  // Confirm dialog
  const enterDialog = async () => {
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      let res
      switch (type.value) {
        case 'create':
          res = await createSysParams(formData.value)
          break
        case 'update':
          res = await updateSysParams(formData.value)
          break
        default:
          res = await createSysParams(formData.value)
          break
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: 'Saved'
        })
        closeDialog()
        getTableData()
      }
    })
  }

  const detailForm = ref({})

  // Details dialog flag
  const detailShow = ref(false)

  // Open details dialog
  const openDetailShow = () => {
    detailShow.value = true
  }

  // Open details
  const getDetails = async (row) => {
    // Open dialog
    const res = await findSysParams({ ID: row.ID })
    if (res.code === 0) {
      detailForm.value = res.data
      openDetailShow()
    }
  }

  // Close details dialog
  const closeDetailShow = () => {
    detailShow.value = false
    detailForm.value = {}
  }
</script>

<style></style>
