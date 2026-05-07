<template>
  <div>
    <div class="gva-table-box">
      <el-table :data="tableData" row-key="id" v-loading="loading">
        <el-table-column :label="t('admin.onboarding.ticket_id')" prop="ticket_id" min-width="160" />
        <el-table-column :label="t('admin.onboarding.agent_name')" prop="agent_name" min-width="150" />
        <el-table-column :label="t('admin.onboarding.current_step')" min-width="140">
          <template #default="{ row }">
            <el-tag>{{ row.current_step || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.onboarding.workflow_id')" prop="workflow_id" min-width="180" />
        <el-table-column :label="t('admin.common.updated_at')" prop="updated_at" min-width="160" />
        <el-table-column :label="t('admin.common.operation')" min-width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDetail(row)">{{ t('admin.onboarding.detail') }}</el-button>
            <el-popconfirm :title="t('admin.onboarding.approve_confirm')" @confirm="handleReview(row, 'approve')">
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
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="(p) => { pagination.page = p; fetchList() }"
      />
    </div>

    <el-dialog v-model="showReject" :title="t('admin.agent.reject')" width="400px" destroy-on-close>
      <el-input v-model="rejectReason" type="textarea" :rows="3" :placeholder="t('admin.agent.reject_reason')" />
      <template #footer>
        <el-button @click="showReject = false">{{ t('admin.common.cancel') }}</el-button>
        <el-button type="danger" :loading="reviewing" @click="handleReview(currentTicket, 'reject')">{{ t('admin.agent.reject') }}</el-button>
      </template>
    </el-dialog>

    <ticket-detail-drawer v-model="showDetail" :ticket-id="currentTicket?.ticket_id" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getTicketList, reviewTicket } from '@/plugin/skyagent/api/ticket'
import TicketDetailDrawer from './components/ticketDetailDrawer.vue'

const { t } = useI18n()

const loading = ref(false)
const tableData = ref([])
const pagination = reactive({ total: 0, page: 1, pageSize: 20 })
const showDetail = ref(false)
const showReject = ref(false)
const reviewing = ref(false)
const currentTicket = ref(null)
const rejectReason = ref('')

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getTicketList({ page: pagination.page, pageSize: pagination.pageSize, status: 'pending_review' })
    if (res.code === 0 && res.data) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } finally { loading.value = false }
}

const openDetail = (row) => { currentTicket.value = row; showDetail.value = true }
const openReject = (row) => { currentTicket.value = row; rejectReason.value = ''; showReject.value = true }

const handleReview = async (row, action) => {
  reviewing.value = true
  try {
    const res = await reviewTicket(row.ticket_id, { action, reason: rejectReason.value })
    if (res.code === 0) {
      ElMessage.success(action === 'approve' ? t('admin.agent.approve_success') : t('admin.agent.reject_success'))
      showReject.value = false
      fetchList()
    } else { ElMessage.error(res.msg) }
  } finally { reviewing.value = false }
}

onMounted(fetchList)
</script>
