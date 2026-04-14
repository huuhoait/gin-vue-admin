<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.catalog.filter_group')">
          <el-select v-model="searchInfo.group" :placeholder="t('admin.common.select')" clearable>
            <el-option v-for="g in groupOptions" :key="g" :label="t(`admin.catalog.group_${g}`)" :value="g" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.catalog.filter_status')">
          <el-select v-model="searchInfo.status" :placeholder="t('admin.common.select')" clearable>
            <el-option label="Active" value="active" />
            <el-option label="Deactive" value="deactive" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">{{ t('admin.common.search') }}</el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" row-key="id" v-loading="loading">
        <el-table-column label="ID" prop="id" min-width="80" />
        <el-table-column :label="t('admin.catalog.product_name')" prop="name" min-width="200" />
        <el-table-column :label="t('admin.catalog.product_image')" min-width="80">
          <template #default="{ row }">
            <el-image v-if="row.image" :src="row.image" style="width: 40px; height: 40px" fit="cover" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.catalog.product_group')" min-width="120">
          <template #default="{ row }">
            <el-tag>{{ t(`admin.catalog.group_${row.product_group}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.catalog.product_status')" min-width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="gva-pagination"
        :current-page="pagination.page"
        :page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        @current-change="(p) => { pagination.page = p; fetchList() }"
        @size-change="(s) => { pagination.pageSize = s; pagination.page = 1; fetchList() }"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getProductList } from '@/api/skyagent/catalog'

const { t } = useI18n()
const groupOptions = ['telco', 'bank', 'travel', 'insurance']

const loading = ref(false)
const tableData = ref([])
const pagination = reactive({ total: 0, page: 1, pageSize: 20 })
const searchInfo = reactive({ group: '', status: '' })

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getProductList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...Object.fromEntries(Object.entries(searchInfo).filter(([, v]) => v !== '')),
    })
    if (res.code === 0 && res.data) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

const onSearch = () => { pagination.page = 1; fetchList() }
const onReset = () => { Object.assign(searchInfo, { group: '', status: '' }); pagination.page = 1; fetchList() }

onMounted(fetchList)
</script>
