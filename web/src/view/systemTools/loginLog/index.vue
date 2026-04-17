<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.systemtools.login_log.username')">
          <el-input v-model="searchInfo.username" :placeholder="t('admin.systemtools.login_log.username_placeholder')" />
        </el-form-item>
        <el-form-item :label="t('admin.systemtools.login_log.status')">
             <el-select v-model="searchInfo.status" :placeholder="t('admin.systemtools.login_log.select_placeholder')" clearable>
                 <el-option :label="t('admin.systemtools.login_log.success')" :value="true" />
                 <el-option :label="t('admin.systemtools.login_log.failure')" :value="false" />
             </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">{{ t('admin.common.search') }}</el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.common.reset') }}</el-button>
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
        >{{ t('admin.systemtools.login_log.delete') }}</el-button>
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
        <el-table-column align="left" :label="t('admin.systemtools.login_log.id')" prop="ID" width="80" />
        <el-table-column align="left" :label="t('admin.systemtools.login_log.username')" prop="username" width="150" />
        <el-table-column align="left" :label="t('admin.systemtools.login_log.login_ip')" prop="ip" width="150" />
        <el-table-column align="left" :label="t('admin.systemtools.login_log.status')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status ? 'success' : 'danger'">
              {{ scope.row.status ? t('admin.systemtools.login_log.success') : t('admin.systemtools.login_log.failure') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.systemtools.login_log.details')" show-overflow-tooltip>
             <template #default="scope">
                 {{ scope.row.status ? t('admin.systemtools.login_log.login_success') : scope.row.errorMessage }}
             </template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.systemtools.login_log.browser_device')" prop="agent" show-overflow-tooltip />
        <el-table-column align="left" :label="t('admin.systemtools.login_log.login_time')" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.systemtools.login_log.actions')" width="120">
          <template #default="scope">
            <el-popover v-model:visible="scope.row.visible" placement="top" width="160">
              <p>{{ t('admin.systemtools.login_log.delete_confirm_item') }}</p>
              <div style="text-align: right; margin: 0">
                <el-button size="small" type="primary" link @click="scope.row.visible = false">{{ t('admin.common.cancel') }}</el-button>
                <el-button size="small" type="primary" @click="deleteRow(scope.row)">{{ t('admin.common.confirm') }}</el-button>
              </div>
              <template #reference>
                <el-button icon="delete" type="primary" link @click="scope.row.visible = true">{{ t('admin.systemtools.login_log.delete') }}</el-button>
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
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

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
      message: t('admin.common.messages.deleted')
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

const onDelete = async() => {
    ElMessageBox.confirm(t('admin.systemtools.login_log.delete_selected'), t('admin.common.confirms.delete_title'), {
        confirmButtonText: t('admin.common.confirm'),
        cancelButtonText: t('admin.common.cancel'),
        type: 'warning'
    }).then(async() => {
        const ids = multipleSelection.value.map(item => item.ID)
        const res = await deleteLoginLogByIds({ ids })
        if (res.code === 0) {
            ElMessage({
                type: 'success',
                message: t('admin.common.messages.deleted')
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
