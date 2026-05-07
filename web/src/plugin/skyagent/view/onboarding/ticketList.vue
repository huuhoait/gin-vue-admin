<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.onboarding.filter_status')">
          <el-select v-model="searchInfo.status" :placeholder="t('admin.common.select')" clearable>
            <el-option v-for="s in statusOptions" :key="s.value" :label="t(s.label)" :value="s.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="searchInfo.ticket_id" :placeholder="t('admin.onboarding.search_ticket')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">{{ t('admin.common.search') }}</el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="showCreate = true">{{ t('admin.onboarding.create_ticket') }}</el-button>
      </div>

      <el-table :data="tableData" row-key="id" v-loading="loading">
        <el-table-column :label="t('admin.onboarding.ticket_id')" prop="ticket_id" min-width="160" />
        <el-table-column :label="t('admin.onboarding.agent_name')" prop="agent_name" min-width="150" />
        <el-table-column :label="t('admin.agent.full_name')" prop="full_name" min-width="160">
          <template #default="{ row }">{{ row.full_name || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.agent.phone')" prop="phone" min-width="130">
          <template #default="{ row }">{{ row.phone || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.onboarding.status')" min-width="130">
          <template #default="{ row }">
            <el-tag :type="statusTagType[row.status] || 'info'" size="small">{{ t(`admin.onboarding.status_${row.status}`) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.onboarding.current_step')" prop="current_step" min-width="140" />
        <el-table-column :label="t('admin.common.updated_at')" prop="updated_at" min-width="160" />
        <el-table-column :label="t('admin.common.operation')" min-width="360" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDetail(row)">{{ t('admin.onboarding.detail') }}</el-button>
            <el-button v-if="row.status === 'draft' || row.status === 'pending_upload'" type="warning" link @click="openUpload(row)">{{ t('admin.onboarding.upload') }}</el-button>
            <el-button v-if="row.status === 'pending_upload'" type="success" link @click="handleSubmit(row)">{{ t('admin.onboarding.submit') }}</el-button>
            <el-tooltip :disabled="row.status === 'pending_review'" :content="t('admin.onboarding.review_only_pending_review')" placement="top">
              <span>
                <el-popconfirm
                  :title="t('admin.onboarding.approve_confirm')"
                  :disabled="row.status !== 'pending_review'"
                  @confirm="handleReview(row, 'approve')"
                >
                  <template #reference>
                    <el-button type="success" link :disabled="row.status !== 'pending_review'">{{ t('admin.agent.approve') }}</el-button>
                  </template>
                </el-popconfirm>
              </span>
            </el-tooltip>
            <el-tooltip :disabled="row.status === 'pending_review'" :content="t('admin.onboarding.review_only_pending_review')" placement="top">
              <span>
                <el-button type="danger" link :disabled="row.status !== 'pending_review'" @click="openReject(row)">{{ t('admin.agent.reject') }}</el-button>
              </span>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="gva-pagination"
        :current-page="pagination.page"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        @current-change="(p) => { pagination.page = p; fetchList() }"
        @size-change="(s) => { pagination.pageSize = s; pagination.page = 1; fetchList() }"
      />
    </div>

    <!-- Create Dialog -->
    <el-dialog v-model="showCreate" :title="t('admin.onboarding.create_ticket')" width="400px" destroy-on-close>
      <el-form :model="createForm" label-width="100px">
        <el-form-item :label="t('admin.onboarding.agent_id')">
          <el-input v-model="createForm.agent_id" placeholder="UUID" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">{{ t('admin.common.cancel') }}</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreate">{{ t('admin.common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- Upload Dialog -->
    <el-dialog v-model="showUpload" :title="t('admin.onboarding.upload')" width="500px" destroy-on-close>
      <el-form :model="uploadForm" label-width="140px">
        <el-form-item :label="t('admin.onboarding.att_type')">
          <el-select v-model="uploadForm.type">
            <el-option label="Business License" value="business_license" />
            <el-option label="Agency Contract" value="agency_contract" />
            <el-option label="eKYC Photos" value="ekyc_photos" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.att_url')">
          <el-input v-model="uploadForm.url" placeholder="https://..." />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showUpload = false">{{ t('admin.common.cancel') }}</el-button>
        <el-button type="primary" :loading="uploading" @click="handleUpload">{{ t('admin.onboarding.upload') }}</el-button>
      </template>
    </el-dialog>

    <!-- Reject Dialog -->
    <el-dialog v-model="showReject" :title="t('admin.agent.reject')" width="400px" destroy-on-close>
      <el-input v-model="rejectReason" type="textarea" :rows="3" :placeholder="t('admin.agent.reject_reason')" />
      <template #footer>
        <el-button @click="showReject = false">{{ t('admin.common.cancel') }}</el-button>
        <el-button type="danger" :loading="reviewing" @click="handleReview(currentTicket, 'reject')">{{ t('admin.agent.reject') }}</el-button>
      </template>
    </el-dialog>

    <!-- Detail Drawer -->
    <ticket-detail-drawer v-model="showDetail" :ticket-id="currentTicket?.ticket_id" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { createTicket, getTicketList, uploadAttachment, submitTicket, reviewTicket } from '@/plugin/skyagent/api/ticket'
import TicketDetailDrawer from './components/ticketDetailDrawer.vue'

const { t } = useI18n()

const statusOptions = [
  { value: 'draft', label: 'admin.onboarding.status_draft' },
  { value: 'pending_upload', label: 'admin.onboarding.status_pending_upload' },
  { value: 'pending_review', label: 'admin.onboarding.status_pending_review' },
  { value: 'approved', label: 'admin.onboarding.status_approved' },
  { value: 'rejected', label: 'admin.onboarding.status_rejected' },
]
const statusTagType = { draft: 'info', pending_upload: 'warning', pending_review: '', approved: 'success', rejected: 'danger' }

const loading = ref(false)
const tableData = ref([])
const pagination = reactive({ total: 0, page: 1, pageSize: 20 })
const searchInfo = reactive({ status: '', ticket_id: '' })

const showCreate = ref(false)
const showUpload = ref(false)
const showDetail = ref(false)
const showReject = ref(false)
const creating = ref(false)
const uploading = ref(false)
const reviewing = ref(false)
const currentTicket = ref(null)
const rejectReason = ref('')
const createForm = reactive({ agent_id: '' })
const uploadForm = reactive({ type: 'business_license', url: '' })

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getTicketList({ page: pagination.page, pageSize: pagination.pageSize, ...searchInfo })
    if (res.code === 0 && res.data) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } finally { loading.value = false }
}

const onSearch = () => { pagination.page = 1; fetchList() }
const onReset = () => { Object.assign(searchInfo, { status: '', ticket_id: '' }); pagination.page = 1; fetchList() }

const handleCreate = async () => {
  creating.value = true
  try {
    const res = await createTicket(createForm)
    if (res.code === 0) {
      ElMessage.success(t('admin.onboarding.create_success'))
      showCreate.value = false
      fetchList()
    } else { ElMessage.error(res.msg) }
  } finally { creating.value = false }
}

const openUpload = (row) => { currentTicket.value = row; showUpload.value = true }
const openDetail = (row) => { currentTicket.value = row; showDetail.value = true }
const openReject = (row) => { currentTicket.value = row; rejectReason.value = ''; showReject.value = true }

const handleUpload = async () => {
  uploading.value = true
  try {
    const res = await uploadAttachment(currentTicket.value.ticket_id, uploadForm)
    if (res.code === 0) {
      ElMessage.success(t('admin.onboarding.upload_success'))
      showUpload.value = false
      fetchList()
    } else { ElMessage.error(res.msg) }
  } finally { uploading.value = false }
}

const handleSubmit = async (row) => {
  const res = await submitTicket(row.ticket_id)
  if (res.code === 0) {
    ElMessage.success(t('admin.onboarding.submit_success'))
    fetchList()
  } else { ElMessage.error(res.msg) }
}

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
