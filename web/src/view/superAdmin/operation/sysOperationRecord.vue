<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.superadmin.operation.method')">
          <el-input v-model="searchInfo.method" :placeholder="t('admin.common.search')" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.operation.path')">
          <el-input v-model="searchInfo.path" :placeholder="t('admin.common.search')" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.operation.status')">
          <el-input v-model="searchInfo.status" :placeholder="t('admin.common.search')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit"
            >{{ t('admin.common.search') }}</el-button
          >
          <el-button icon="refresh" @click="onReset">{{ t('admin.common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          icon="delete"
          :disabled="!multipleSelection.length"
          @click="onDelete"
          >{{ t('admin.common.delete') }}</el-button
        >
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column align="left" type="selection" width="55" />
        <el-table-column align="left" :label="t('admin.superadmin.operation.operator')" width="140">
          <template #default="scope">
            <div>
              {{ scope.row.user.userName }}({{ scope.row.user.nickName }})
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.superadmin.operation.date')" width="180">
          <template #default="scope">{{
            formatDate(scope.row.CreatedAt)
          }}</template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.superadmin.operation.status')" prop="status" width="120">
          <template #default="scope">
            <div>
              <el-tag type="success">{{ scope.row.status }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.superadmin.operation.ip')" prop="ip" width="120" />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.operation.method')"
          prop="method"
          width="120"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.operation.path')"
          prop="path"
          width="240"
        />
        <el-table-column align="left" :label="t('admin.superadmin.operation.request')" prop="path" width="80">
          <template #default="scope">
            <div>
              <el-popover
                v-if="scope.row.body"
                placement="left-start"
                :width="444"
              >
                <div class="popover-box">
                  <pre>{{ fmtBody(scope.row.body) }}</pre>
                </div>
                <template #reference>
                  <el-icon style="cursor: pointer"><warning /></el-icon>
                </template>
              </el-popover>

              <span v-else>{{ t('admin.superadmin.operation.none') }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.superadmin.operation.response')" prop="path" width="80">
          <template #default="scope">
            <div>
              <el-popover
                v-if="scope.row.resp"
                placement="left-start"
                :width="444"
              >
                <div class="popover-box">
                  <pre>{{ fmtBody(scope.row.resp) }}</pre>
                </div>
                <template #reference>
                  <el-icon style="cursor: pointer"><warning /></el-icon>
                </template>
              </el-popover>
              <span v-else>{{ t('admin.superadmin.operation.none') }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.common.operation')">
          <template #default="scope">
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteSysOperationRecordFunc(scope.row)"
              >{{ t('admin.common.delete') }}</el-button
            >
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
    deleteSysOperationRecord,
    getSysOperationRecordList,
    deleteSysOperationRecordByIds
  } from '@/api/sysOperationRecord'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useI18n } from 'vue-i18n'

  defineOptions({
    name: 'SysOperationRecord'
  })

  const { t } = useI18n()

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})
  const onReset = () => {
    searchInfo.value = {}
  }
  // submit search
  const onSubmit = () => {
    page.value = 1
    if (searchInfo.value.status === '') {
      searchInfo.value.status = null
    }
    getTableData()
  }

  // pagination
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // fetch
  const getTableData = async () => {
    const table = await getSysOperationRecordList({
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

  const multipleSelection = ref([])
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }
  const onDelete = async () => {
    ElMessageBox.confirm(t('admin.common.delete_confirm'), t('admin.common.are_you_sure'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const ids = []
      multipleSelection.value &&
        multipleSelection.value.forEach((item) => {
          ids.push(item.ID)
        })
      const res = await deleteSysOperationRecordByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: t('admin.superadmin.operation.delete_success')
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }
  const deleteSysOperationRecordFunc = async (row) => {
    ElMessageBox.confirm(t('admin.common.delete_confirm'), t('admin.common.are_you_sure'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const res = await deleteSysOperationRecord({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: t('admin.superadmin.operation.delete_success')
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }
  const fmtBody = (value) => {
    try {
      return JSON.parse(value)
    } catch (_) {
      return value
    }
  }
</script>

<style lang="scss">
  .table-expand {
    padding-left: 60px;
    font-size: 0;
    label {
      width: 90px;
      color: #99a9bf;
      .el-form-item {
        margin-right: 0;
        margin-bottom: 0;
        width: 50%;
      }
    }
  }
  .popover-box {
    background: #112435;
    color: #f08047;
    height: 600px;
    width: 420px;
    overflow: auto;
  }
  .popover-box::-webkit-scrollbar {
    display: none; /* Chrome Safari */
  }
</style>
