<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="Username">
          <el-input v-model="searchInfo.username" placeholder="Search username" />
        </el-form-item>
        <el-form-item label="Status">
             <el-select v-model="searchInfo.status" placeholder="Select" clearable>
                 <el-option label="Success" :value="true" />
                 <el-option label="Failure" :value="false" />
             </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">Search</el-button>
          <el-button icon="refresh" @click="onReset">Reset</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          icon="delete"
          style="margin-left: 10px;"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >Delete</el-button>
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="Username" prop="username" width="150" />
        <el-table-column align="left" label="Login IP" prop="ip" width="150" />
        <el-table-column align="left" label="Status" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status ? 'success' : 'danger'">
              {{ scope.row.status ? 'Success' : 'Failure' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="Details" show-overflow-tooltip>
             <template #default="scope">
                 {{ scope.row.status ? 'Login success' : scope.row.errorMessage }}
             </template>
        </el-table-column>
        <el-table-column align="left" label="Browser / device" prop="agent" show-overflow-tooltip />
        <el-table-column align="left" label="Login time" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="Actions" width="120">
          <template #default="scope">
            <el-popover v-model:visible="scope.row.visible" placement="top" width="160">
              <p>Delete this item?</p>
              <div style="text-align: right; margin: 0">
                <el-button size="small" type="primary" link @click="scope.row.visible = false">Cancel</el-button>
                <el-button size="small" type="primary" @click="deleteRow(scope.row)">Confirm</el-button>
              </div>
              <template #reference>
                <el-button icon="delete" type="primary" link @click="scope.row.visible = true">Delete</el-button>
              </template>
            </el-popover>
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
  </div>
</template>

<script setup>
import {
  getLoginLogList,
  deleteLoginLog,
  deleteLoginLogByIds
} from '@/api/sysLoginLog'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const multipleSelection = ref([])

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const getTableData = async () => {
  const table = await getLoginLogList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

const deleteRow = async (row) => {
  row.visible = false
  const res = await deleteLoginLog(row)
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

const onDelete = async() => {
    ElMessageBox.confirm('Delete selected items?', 'Confirm', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
    }).then(async() => {
        const ids = multipleSelection.value.map(item => item.ID)
        const res = await deleteLoginLogByIds({ ids })
        if (res.code === 0) {
            ElMessage({
                type: 'success',
                message: 'Deleted'
            })
            if (tableData.value.length === ids.length && page.value > 1) {
                page.value--
            }
            getTableData()
        }
    })
}

const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// Initial load
getTableData()
</script>

<style scoped>
</style>
