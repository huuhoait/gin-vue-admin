<template>
  <el-drawer :model-value="modelValue" :title="t('admin.onboarding.detail')" size="60%" destroy-on-close @close="emit('update:modelValue', false)">
    <div v-loading="loading">
      <template v-if="ticket">
        <!-- Ticket -->
        <el-descriptions :title="t('admin.onboarding.ticket_info')" :column="2" border>
          <el-descriptions-item :label="t('admin.onboarding.ticket_id')">
            <el-tag effect="plain">{{ ticket.ticket_id }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.agent_id')">
            <span class="font-mono">{{ ticket.agent_id }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.status')">
            <el-tag :type="statusTagType[ticket.status] || 'info'">{{ t(`admin.onboarding.status_${ticket.status}`) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.current_step')">
            <el-tag v-if="ticket.current_step">{{ ticket.current_step }}</el-tag>
            <span v-else>-</span>
          </el-descriptions-item>
        </el-descriptions>

        <!-- Applicant identity (masked from upstream) -->
        <el-descriptions :title="t('admin.onboarding.applicant_info')" :column="2" border >
          <el-descriptions-item :label="t('admin.agent.full_name')">{{ ticket.full_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.phone')">{{ ticket.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.email')">{{ ticket.email || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.cccd')">{{ ticket.cccd || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- Business profile from agent_payload -->
        <el-descriptions v-if="profile" :title="t('admin.onboarding.business_profile')" :column="2" border >
          <el-descriptions-item :label="t('admin.onboarding.business_name')">{{ profile.business_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.business_type')">{{ profile.business_type || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.tax_code')">{{ profile.tax_code || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.representative_name')">{{ profile.representative_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.representative_cccd')">{{ profile.representative_cccd || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.bank_name')">{{ profile.bank_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.bank_account')">{{ profile.bank_account || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.permanent_address')">{{ profile.permanent_address || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.contact_address')" :span="2">{{ profile.contact_address || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- Identity from agent_payload (non-masked plaintext source) -->
        <el-descriptions v-if="payloadIdentity" :title="t('admin.onboarding.payload_identity')" :column="2" border >
          <el-descriptions-item :label="t('admin.agent.province')">{{ payloadIdentity.province || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.district')">{{ payloadIdentity.district || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.parent_id')">{{ payloadIdentity.parent_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.referral_code')">{{ payloadIdentity.referral_code || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.referral_name')">{{ payloadIdentity.referral_name || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- Workflow -->
        <el-descriptions v-if="ticket.workflow_id" :title="t('admin.onboarding.workflow_info')" :column="2" border >
          <el-descriptions-item :label="t('admin.onboarding.workflow_id')">
            <span class="font-mono">{{ ticket.workflow_id }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.run_id')">
            <span class="font-mono">{{ ticket.run_id }}</span>
          </el-descriptions-item>
        </el-descriptions>

        <!-- Audit / timeline -->
        <el-descriptions :title="t('admin.onboarding.audit_info')" :column="2" border >
          <el-descriptions-item :label="t('admin.agent.created_by')">
            <UserRef :uuid="ticket.maker_id" :name="ticket.maker_id_name" />
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.updated_by')">
            <UserRef :uuid="ticket.checker_id" :name="ticket.checker_id_name" />
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.common.created_at')">{{ ticket.created_at || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.common.updated_at')">{{ ticket.updated_at || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.submitted_at')">{{ ticket.submitted_at || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.reviewed_at')">{{ ticket.reviewed_at || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- Attachments -->
        <h4 >{{ t('admin.onboarding.attachments') }}</h4>
        <el-table v-if="attachmentRows.length" :data="attachmentRows" size="small">
          <el-table-column :label="t('admin.onboarding.att_type')" prop="type" min-width="140">
            <template #default="{ row }">
              <el-tag size="small" effect="plain">{{ row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('admin.onboarding.att_url')" min-width="240">
            <template #default="{ row }">
              <el-link type="primary" :href="row.url" target="_blank">{{ row.url }}</el-link>
            </template>
          </el-table-column>
          <el-table-column :label="t('admin.common.created_at')" prop="uploaded_at" min-width="160" />
        </el-table>
        <el-empty v-else :description="t('admin.onboarding.no_attachments')" />

        <!-- Reject reason -->
        <template v-if="ticket.reject_reason">
          <h4 >{{ t('admin.agent.reject_reason') }}</h4>
          <el-alert type="error" :closable="false">{{ ticket.reject_reason }}</el-alert>
        </template>

        <!-- Raw payload (collapsible) -->
        <el-collapse v-if="parsedPayload" >
          <el-collapse-item :title="t('admin.onboarding.raw_payload')" name="payload">
            <pre class="raw-json">{{ pretty(parsedPayload) }}</pre>
          </el-collapse-item>
        </el-collapse>
      </template>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, computed, watch, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElTooltip } from 'element-plus'
import { getTicketDetail } from '@/api/skyagent/ticket'

const { t } = useI18n()
const props = defineProps({ modelValue: Boolean, ticketId: String })
const emit = defineEmits(['update:modelValue'])

const statusTagType = { draft: 'info', pending_upload: 'warning', pending_review: '', approved: 'success', rejected: 'danger' }
const loading = ref(false)
const ticket = ref(null)

// agent_payload is a JSON-encoded string in the contract; parse defensively.
const parsedPayload = computed(() => {
  const p = ticket.value?.agent_payload
  if (!p) return null
  if (typeof p === 'object') return p
  try { return JSON.parse(p) } catch { return null }
})
const profile = computed(() => parsedPayload.value?.profile || null)
const payloadIdentity = computed(() => parsedPayload.value?.identity || null)
// Prefer ticket-level attachments; fall back to payload.attachments if the
// upstream omits the top-level array (early drafts).
const attachmentRows = computed(() => {
  if (Array.isArray(ticket.value?.attachments) && ticket.value.attachments.length) {
    return ticket.value.attachments
  }
  return parsedPayload.value?.attachments || []
})

const pretty = (obj) => JSON.stringify(obj, null, 2)

// Tiny inline component for uuid-or-name rendering with tooltip fallback.
const UserRef = {
  props: ['uuid', 'name'],
  setup(p) {
    return () => {
      if (p.name) return h('span', p.name)
      if (!p.uuid) return h('span', '-')
      return h(ElTooltip, { content: p.uuid, placement: 'top' }, {
        default: () => h('span', { class: 'font-mono' }, p.uuid),
      })
    }
  },
}

watch(() => [props.modelValue, props.ticketId], async ([visible, id]) => {
  if (visible && id) {
    loading.value = true
    try {
      const res = await getTicketDetail(id)
      if (res.code === 0) { ticket.value = res.data }
    } finally { loading.value = false }
  } else { ticket.value = null }
}, { immediate: true })
</script>

<style scoped>
.raw-json {
  max-height: 320px;
  overflow: auto;
  background: #f7f8fa;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 12px;
  font-size: 12px;
  line-height: 1.5;
  margin: 0;
}
</style>
