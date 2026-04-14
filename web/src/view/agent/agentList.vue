<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.agent.filter_status')">
          <el-select v-model="searchInfo.status" :placeholder="t('admin.common.select')" clearable>
            <el-option v-for="s in statusOptions" :key="s.value" :label="t(s.label)" :value="s.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.agent.filter_level')">
          <el-select v-model="searchInfo.level" :placeholder="t('admin.common.select')" clearable>
            <el-option v-for="l in 5" :key="l - 1" :label="`Level ${l - 1}`" :value="l - 1" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.agent.filter_kyc')">
          <el-select v-model="searchInfo.kyc_tier" :placeholder="t('admin.common.select')" clearable>
            <el-option v-for="k in 3" :key="k - 1" :label="t(`admin.agent.kyc_tier_${k - 1}`)" :value="k - 1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="searchInfo.keyword" :placeholder="t('admin.agent.search_keyword')" clearable />
        </el-form-item>
        <el-form-item :label="t('admin.agent.referral_code')">
          <el-input v-model="searchInfo.referral_code" :placeholder="t('admin.agent.referral_code')" clearable style="width: 140px" />
        </el-form-item>
        <el-form-item :label="t('admin.agent.filter_date')">
          <el-date-picker v-model="dateRange" type="daterange" value-format="YYYY-MM-DD" @change="onDateChange" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">{{ t('admin.common.search') }}</el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="showCreate = true">{{ t('admin.agent.create') }}</el-button>
        <el-switch
          v-model="treeView"
          :active-text="t('admin.agent.tree_view')"
          :inactive-text="t('admin.agent.flat_view')"
          style="margin-left: 16px"
          @change="fetchList"
        />
      </div>

      <el-table
        :data="displayData"
        row-key="id"
        v-loading="loading"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        :default-expand-all="treeView"
      >
        <el-table-column :label="t('admin.agent.full_name')" min-width="200">
          <template #default="{ row }">
            <span>
              <el-tag size="small" class="mr-2" effect="plain">L{{ row.level }}</el-tag>
              {{ row.full_name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.agent.code')" prop="code" min-width="110" />
        <el-table-column :label="t('admin.agent.status')" min-width="130">
          <template #default="{ row }">
            <el-tag :type="statusTagType[row.status] || 'info'" size="small">{{ t(`admin.agent.status_${row.status}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.agent.kyc_tier')" min-width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="kycTagType[row.kyc_tier]" size="small">{{ t(`admin.agent.kyc_tier_${row.kyc_tier}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.agent.province')" prop="province" min-width="120" />
        <el-table-column :label="t('admin.agent.created_at')" prop="created_at" min-width="160" />
        <el-table-column :label="t('admin.common.operation')" min-width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDetail(row)">{{ t('admin.agent.detail') }}</el-button>
            <el-button type="primary" link @click="openEdit(row)">{{ t('admin.common.edit') }}</el-button>
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
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>

    <create-agent-dialog v-model="showCreate" @success="fetchList" />
    <edit-agent-dialog v-model="showEdit" :agent="currentAgent" @success="fetchList" />
    <agent-detail-drawer v-model="showDetail" :agent-id="currentAgent?.id" />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getAgentList } from '@/api/skyagent/agent'
import CreateAgentDialog from './components/createAgentDialog.vue'
import EditAgentDialog from './components/editAgentDialog.vue'
import AgentDetailDrawer from './components/agentDetailDrawer.vue'

const { t } = useI18n()

const statusOptions = [
  { value: 'pending_approval', label: 'admin.agent.status_pending_approval' },
  { value: 'active', label: 'admin.agent.status_active' },
  { value: 'suspended', label: 'admin.agent.status_suspended' },
  { value: 'terminated', label: 'admin.agent.status_terminated' },
]
const statusTagType = {
  pending_approval: 'warning',
  active: 'success',
  suspended: 'danger',
  terminated: 'info',
}
const kycTagType = { 0: 'info', 1: '', 2: 'success' }

const loading = ref(false)
const flatData = ref([])
const treeView = ref(true)
const pagination = reactive({ total: 0, page: 1, pageSize: 50 })
const searchInfo = reactive({ status: '', level: '', kyc_tier: '', keyword: '', referral_code: '', created_from: '', created_to: '' })
const dateRange = ref(null)

const onDateChange = (val) => {
  if (val && val.length === 2) {
    searchInfo.created_from = val[0]
    searchInfo.created_to = val[1]
  } else {
    searchInfo.created_from = ''
    searchInfo.created_to = ''
  }
}

const showCreate = ref(false)
const showEdit = ref(false)
const showDetail = ref(false)
const currentAgent = ref(null)

// Build tree from flat list using parent_id
const buildTree = (items) => {
  const map = {}
  const roots = []

  // Index all items by id
  items.forEach((item) => {
    map[item.id] = { ...item, children: [] }
  })

  // Link children to parents
  items.forEach((item) => {
    const node = map[item.id]
    if (item.parent_id && map[item.parent_id]) {
      map[item.parent_id].children.push(node)
    } else {
      roots.push(node)
    }
  })

  // Remove empty children arrays (el-table uses hasChildren)
  const clean = (nodes) => {
    nodes.forEach((n) => {
      if (n.children.length === 0) {
        delete n.children
      } else {
        clean(n.children)
      }
    })
  }
  clean(roots)
  return roots
}

const displayData = computed(() => {
  if (treeView.value) {
    return buildTree(flatData.value)
  }
  return flatData.value
})

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getAgentList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...Object.fromEntries(Object.entries(searchInfo).filter(([, v]) => v !== '' && v !== null)),
    })
    if (res.code === 0 && res.data) {
      flatData.value = res.data.list || []
      pagination.total = res.data.total || 0
      pagination.page = res.data.page || 1
      pagination.pageSize = res.data.pageSize || 50
    }
  } finally {
    loading.value = false
  }
}

const onSearch = () => { pagination.page = 1; fetchList() }
const onReset = () => {
  Object.assign(searchInfo, { status: '', level: '', kyc_tier: '', keyword: '', referral_code: '', created_from: '', created_to: '' })
  dateRange.value = null
  pagination.page = 1
  fetchList()
}
const handlePageChange = (page) => { pagination.page = page; fetchList() }
const handleSizeChange = (size) => { pagination.pageSize = size; pagination.page = 1; fetchList() }

const openDetail = (row) => { currentAgent.value = row; showDetail.value = true }
const openEdit = (row) => { currentAgent.value = row; showEdit.value = true }

onMounted(fetchList)
</script>
