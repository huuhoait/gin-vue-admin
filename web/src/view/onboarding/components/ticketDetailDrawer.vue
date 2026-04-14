<template>
  <el-drawer :model-value="modelValue" :title="t('admin.onboarding.detail')" size="50%" destroy-on-close @close="emit('update:modelValue', false)">
    <div v-loading="loading">
      <template v-if="ticket">
        <el-descriptions :title="t('admin.onboarding.ticket_info')" :column="2" border>
          <el-descriptions-item :label="t('admin.onboarding.ticket_id')">
            <el-tag effect="plain">{{ ticket.ticket_id }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.agent_id')">{{ ticket.agent_id }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.status')">
            <el-tag :type="statusTagType[ticket.status] || 'info'">{{ t(`admin.onboarding.status_${ticket.status}`) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.common.created_at')">{{ ticket.created_at }}</el-descriptions-item>
        </el-descriptions>

        <el-descriptions v-if="ticket.workflow_id" :title="t('admin.onboarding.workflow_info')" :column="2" border class="mt-4">
          <el-descriptions-item :label="t('admin.onboarding.workflow_id')">{{ ticket.workflow_id }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.run_id')">{{ ticket.run_id }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.current_step')">
            <el-tag>{{ ticket.current_step }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.common.updated_at')">{{ ticket.updated_at }}</el-descriptions-item>
        </el-descriptions>

        <h4 class="mt-4">{{ t('admin.onboarding.attachments') }}</h4>
        <el-table v-if="ticket.attachments?.length" :data="ticket.attachments" size="small">
          <el-table-column :label="t('admin.onboarding.att_type')" prop="type" min-width="120" />
          <el-table-column :label="t('admin.onboarding.att_url')" min-width="300">
            <template #default="{ row }">
              <el-link type="primary" :href="row.url" target="_blank">{{ row.url }}</el-link>
            </template>
          </el-table-column>
          <el-table-column :label="t('admin.common.created_at')" prop="uploaded_at" min-width="160" />
        </el-table>
        <el-empty v-else :description="t('admin.onboarding.no_attachments')" />

        <template v-if="ticket.reject_reason">
          <h4 class="mt-4">{{ t('admin.agent.reject_reason') }}</h4>
          <el-alert type="error" :closable="false">{{ ticket.reject_reason }}</el-alert>
        </template>
      </template>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { getTicketDetail } from '@/api/skyagent/ticket'

const { t } = useI18n()
const props = defineProps({ modelValue: Boolean, ticketId: String })
const emit = defineEmits(['update:modelValue'])

const statusTagType = { draft: 'info', pending_upload: 'warning', pending_review: '', approved: 'success', rejected: 'danger' }
const loading = ref(false)
const ticket = ref(null)

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
