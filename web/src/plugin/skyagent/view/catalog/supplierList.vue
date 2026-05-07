<template>
  <div>
    <div class="gva-table-box">
      <el-table :data="tableData" row-key="id" v-loading="loading">
        <el-table-column :label="t('admin.catalog.supplier_name')" prop="name" min-width="200" />
        <el-table-column :label="t('admin.catalog.adapter_type')" prop="adapter_type" min-width="150" />
        <el-table-column :label="t('admin.catalog.health_status')" min-width="120">
          <template #default="{ row }">
            <el-tag :type="healthTagType[row.health_status] || 'info'">
              {{ t(`admin.catalog.health_${row.health_status || 'na'}`) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="gva-pagination"
        :current-page="pagination.page"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="(p) => { pagination.page = p; fetchList() }"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getSupplierList } from '@/plugin/skyagent/api/catalog'

const { t } = useI18n()
const healthTagType = { online: 'success', offline: 'danger', degraded: 'warning' }

const loading = ref(false)
const tableData = ref([])
const pagination = reactive({ total: 0, page: 1, pageSize: 20 })

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getSupplierList({ page: pagination.page, pageSize: pagination.pageSize })
    if (res.code === 0 && res.data) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
</script>
