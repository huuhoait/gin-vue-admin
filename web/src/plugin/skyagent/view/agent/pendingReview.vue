<template>
  <div>
    <div class="gva-table-box">
      <el-table :data="tableData" row-key="id" v-loading="loading">
        <el-table-column :label="t('admin.agent.code')" prop="code" min-width="100" />
        <el-table-column :label="t('admin.agent.full_name')" prop="full_name" min-width="150" />
        <el-table-column :label="t('admin.agent.level')" prop="level" min-width="80" align="center" />
        <el-table-column :label="t('admin.agent.created_at')" prop="created_at" min-width="160" />
        <el-table-column :label="t('admin.common.operation')" min-width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDetail(row)">{{ t('admin.agent.detail') }}</el-button>
            <el-popconfirm :title="t('admin.agent.approve_confirm')" @confirm="handleApprove(row)">
              <template #reference>
                <el-button type="success" link>{{ t('admin.agent.approve') }}</el-button>
              </template>
            </el-popconfirm>
            <el-button type="danger" link @click="openReject(row)">{{ t('admin.agent.reject') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="gva-pagination"
        :current-page="pagination.page"
        :page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        @current-change="(p) => { pagination.page = p; fetchList() }"
        @size-change="(s) => { pagination.pageSize = s; pagination.page = 1; fetchList() }"
      />
    </div>

    <agent-detail-drawer v-model="showDetail" :agent="currentAgent" />

    <el-dialog v-model="showReject" :title="t('admin.agent.reject')" width="400px" destroy-on-close>
      <el-input v-model="rejectReason" type="textarea" :rows="3" :placeholder="t('admin.agent.reject_reason')" />
      <template #footer>
        <el-button @click="showReject = false">{{ t('admin.common.cancel') }}</el-button>
        <el-button type="danger" :loading="rejecting" @click="handleReject">{{ t('admin.agent.reject') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getAgentList, updateAgentStatus } from '@/plugin/skyagent/api/agent'
import AgentDetailDrawer from './components/agentDetailDrawer.vue'

const { t } = useI18n()

const loading = ref(false)
const tableData = ref([])
const pagination = reactive({ total: 0, page: 1, pageSize: 20 })

const showDetail = ref(false)
const showReject = ref(false)
const currentAgent = ref(null)
const rejectReason = ref('')
const rejecting = ref(false)

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getAgentList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: 'pending_approval',
    })
    if (res.code === 0 && res.data) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

const openDetail = (row) => { currentAgent.value = row; showDetail.value = true }

const handleApprove = async (row) => {
  const res = await updateAgentStatus(row.id, { status: 'active' })
  if (res.code === 0) {
    ElMessage.success(t('admin.agent.approve_success'))
    fetchList()
  } else {
    ElMessage.error(res.msg || t('admin.common.fail'))
  }
}

const openReject = (row) => {
  currentAgent.value = row
  rejectReason.value = ''
  showReject.value = true
}

const handleReject = async () => {
  rejecting.value = true
  try {
    const res = await updateAgentStatus(currentAgent.value.id, {
      status: 'rejected',
      reason: rejectReason.value,
    })
    if (res.code === 0) {
      ElMessage.success(t('admin.agent.reject_success'))
      showReject.value = false
      fetchList()
    } else {
      ElMessage.error(res.msg || t('admin.common.fail'))
    }
  } finally {
    rejecting.value = false
  }
}

onMounted(fetchList)
</script>
