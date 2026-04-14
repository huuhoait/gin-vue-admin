<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.order.filter_date')">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            value-format="YYYY-MM-DD"
            @change="onDateChange"
          />
        </el-form-item>
        <el-form-item :label="t('admin.order.filter_status')">
          <el-select v-model="searchInfo.status" :placeholder="t('admin.common.select')" clearable>
            <el-option v-for="s in statusOptions" :key="s.value" :label="t(s.label)" :value="s.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="searchInfo.keyword" :placeholder="t('admin.order.search_placeholder')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">{{ t('admin.common.search') }}</el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-tooltip :content="t('admin.order.export_coming_soon')">
          <el-button disabled icon="download">{{ t('admin.common.export') }}</el-button>
        </el-tooltip>
      </div>

      <el-table :data="tableData" row-key="id" v-loading="loading">
        <el-table-column :label="t('admin.order.order_id')" min-width="140">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDetail(row)">{{ row.id }}</el-button>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.order.agent')" prop="agent_code" min-width="120" />
        <el-table-column :label="t('admin.order.product')" prop="product_name" min-width="160" />
        <el-table-column :label="t('admin.order.amount')" min-width="130" align="right">
          <template #default="{ row }">{{ formatVND(row.total_amount) }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.order.status')" min-width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType[row.status] || 'info'">{{ t(`admin.order.status_${row.status}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.order.created_at')" prop="created_at" min-width="160" />
        <el-table-column :label="t('admin.common.operation')" min-width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDetail(row)">{{ t('admin.order.detail') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="gva-pagination"
        :current-page="pagination.page"
        :page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="(p) => { pagination.page = p; fetchList() }"
        @size-change="(s) => { pagination.pageSize = s; pagination.page = 1; fetchList() }"
      />
    </div>

    <order-detail-drawer v-model="showDetail" :order-id="currentOrderId" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getOrderList } from '@/api/skyagent/order'
import OrderDetailDrawer from './components/orderDetailDrawer.vue'

const { t } = useI18n()

const statusOptions = [
  { value: 'pending', label: 'admin.order.status_pending' },
  { value: 'success', label: 'admin.order.status_success' },
  { value: 'failed', label: 'admin.order.status_failed' },
  { value: 'refunded', label: 'admin.order.status_refunded' },
]
const statusTagType = { pending: 'warning', success: 'success', failed: 'danger', refunded: 'info' }

const formatVND = (amount) => {
  if (!amount && amount !== 0) return '-'
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(amount)
}

const loading = ref(false)
const tableData = ref([])
const pagination = reactive({ total: 0, page: 1, pageSize: 20 })
const searchInfo = reactive({ status: '', keyword: '', from: '', to: '' })
const dateRange = ref(null)

const showDetail = ref(false)
const currentOrderId = ref('')

const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...Object.fromEntries(Object.entries(searchInfo).filter(([, v]) => v !== '')),
    }
    const res = await getOrderList(params)
    if (res.code === 0 && res.data) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

const onDateChange = (val) => {
  if (val && val.length === 2) {
    searchInfo.from = val[0]
    searchInfo.to = val[1]
  } else {
    searchInfo.from = ''
    searchInfo.to = ''
  }
}

const onSearch = () => { pagination.page = 1; fetchList() }
const onReset = () => {
  Object.assign(searchInfo, { status: '', keyword: '', from: '', to: '' })
  dateRange.value = null
  pagination.page = 1
  fetchList()
}

const openDetail = (row) => { currentOrderId.value = row.id; showDetail.value = true }

onMounted(fetchList)
</script>
