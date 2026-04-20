<template>
  <el-drawer :model-value="modelValue" :title="t('admin.agent.detail')" size="45%" destroy-on-close @close="emit('update:modelValue', false)">
    <div v-loading="loading">
      <template v-if="agent">
        <el-descriptions :title="t('admin.agent.basic_info')" :column="2" border>
          <el-descriptions-item :label="t('admin.agent.code')">{{ agent.code }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.full_name')">{{ agent.full_name }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.phone')">{{ agent.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.email')">{{ agent.email || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.level')">{{ agent.level }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.status')">
            <el-tag :type="statusTagType[agent.status] || 'info'">{{ t(`admin.agent.status_${agent.status}`) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.kyc_tier')">
            <el-tag :type="kycTagType[agent.kyc_tier]" size="small">{{ t(`admin.agent.kyc_tier_${agent.kyc_tier}`) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.province')">{{ agent.province || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.partner_id')">{{ agent.partner_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.user_id')">{{ agent.user_id || '-' }}</el-descriptions-item>
        </el-descriptions>

        <el-descriptions :title="t('admin.agent.hierarchy')" :column="1" border class="mt-4">
          <el-descriptions-item :label="t('admin.agent.path')">
            <el-tag v-if="agent.path" effect="plain" size="small">{{ agent.path }}</el-tag>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.parent_id')">{{ agent.parent_id || '-' }}</el-descriptions-item>
        </el-descriptions>

        <el-descriptions :title="t('admin.agent.audit')" :column="2" border class="mt-4">
          <el-descriptions-item :label="t('admin.agent.created_at')">{{ agent.created_at }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.updated_at')">{{ agent.updated_at }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.created_by')">
            <span v-if="agent.created_by_name">{{ agent.created_by_name }}</span>
            <el-tooltip v-else-if="agent.created_by" :content="agent.created_by" placement="top">
              <span class="font-mono">{{ agent.created_by }}</span>
            </el-tooltip>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.updated_by')">
            <span v-if="agent.updated_by_name">{{ agent.updated_by_name }}</span>
            <el-tooltip v-else-if="agent.updated_by" :content="agent.updated_by" placement="top">
              <span class="font-mono">{{ agent.updated_by }}</span>
            </el-tooltip>
            <span v-else>-</span>
          </el-descriptions-item>
        </el-descriptions>
      </template>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { getAgentDetail } from '@/api/skyagent/agent'

const { t } = useI18n()
const props = defineProps({ modelValue: Boolean, agentId: String })
const emit = defineEmits(['update:modelValue'])

const statusTagType = { pending_approval: 'warning', active: 'success', suspended: 'danger', terminated: 'info' }
const kycTagType = { 0: 'info', 1: '', 2: 'success' }

const loading = ref(false)
const agent = ref(null)

watch(() => [props.modelValue, props.agentId], async ([visible, id]) => {
  if (visible && id) {
    loading.value = true
    try {
      const res = await getAgentDetail(id)
      if (res.code === 0) {
        agent.value = res.data
      }
    } finally {
      loading.value = false
    }
  } else {
    agent.value = null
  }
}, { immediate: true })
</script>
