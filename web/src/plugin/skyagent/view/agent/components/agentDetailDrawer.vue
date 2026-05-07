<template>
  <el-drawer :model-value="modelValue" :title="t('admin.agent.detail')" size="60%" destroy-on-close @close="emit('update:modelValue', false)">
    <div v-loading="loading">
      <template v-if="agent">
        <!-- Header row with avatar + quick actions -->
        <div class="agent-header">
          <el-avatar :size="64" :src="agent.avatar_url || undefined">
            {{ (agent.full_name || '?').slice(0, 1).toUpperCase() }}
          </el-avatar>
          <div class="agent-header__meta">
            <div class="agent-header__name">
              {{ agent.full_name || '-' }}
              <el-tag size="small" effect="plain" class="ml-2">L{{ agent.level }}</el-tag>
              <el-tag :type="statusTagType[agent.status] || 'info'" size="small" class="ml-1">
                {{ t(`admin.agent.status_${agent.status}`) }}
              </el-tag>
            </div>
            <div class="agent-header__code">{{ agent.code }}</div>
          </div>
          <div class="flex-1" />
          <el-tooltip v-if="!showingPlain" :content="t('admin.agent.reveal_pii_hint')" placement="top">
            <el-button type="warning" :icon="Unlock" :loading="loadingPlain" @click="loadPlain">
              {{ t('admin.agent.reveal_pii') }}
            </el-button>
          </el-tooltip>
          <el-tag v-else type="warning" effect="dark" :icon="Warning">{{ t('admin.agent.plaintext_mode') }}</el-tag>
        </div>

        <!-- Basic info -->
        <el-descriptions :title="t('admin.agent.basic_info')" :column="2" border class="mt-4">
          <el-descriptions-item :label="t('admin.agent.phone')">{{ display.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.email')">{{ display.email || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.kyc_tier')">
            <el-tag :type="kycTagType[agent.kyc_tier]" size="small">{{ t(`admin.agent.kyc_tier_${agent.kyc_tier}`) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.partner_id')">{{ agent.partner_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.user_id')">
            <span class="font-mono">{{ agent.user_id || '-' }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.referral_code')">{{ agent.referral_code || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.referral_name')">{{ agent.referral_name || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- Business -->
        <el-descriptions v-if="hasBusiness" :title="t('admin.onboarding.business_profile')" :column="2" border class="mt-4">
          <el-descriptions-item :label="t('admin.onboarding.business_name')">{{ agent.business_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.business_type')">{{ agent.business_type || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.tax_code')">{{ agent.tax_code || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.representative_name')">{{ agent.representative_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.representative_cccd')" :span="2">{{ display.representative_cccd || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- Bank -->
        <el-descriptions v-if="agent.bank_name || agent.bank_account" :title="t('admin.agent.bank_info')" :column="2" border class="mt-4">
          <el-descriptions-item :label="t('admin.onboarding.bank_name')">{{ agent.bank_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.bank_account')">{{ display.bank_account || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- Address -->
        <el-descriptions :title="t('admin.agent.address')" :column="2" border class="mt-4">
          <el-descriptions-item :label="t('admin.agent.province')">{{ agent.province || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.district')">{{ agent.district || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.permanent_address')" :span="2">{{ agent.permanent_address || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.contact_address')" :span="2">{{ agent.contact_address || '-' }}</el-descriptions-item>
        </el-descriptions>

        <!-- Hierarchy -->
        <el-descriptions :title="t('admin.agent.hierarchy')" :column="2" border class="mt-4">
          <el-descriptions-item :label="t('admin.agent.path')" :span="2">
            <el-tag v-if="agent.path" effect="plain" size="small" class="font-mono">{{ agent.path }}</el-tag>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.parent_id')">
            <span class="font-mono">{{ agent.parent_id || '-' }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.children_count')">
            {{ (agent.children_count ?? 0) }}
            <span v-if="agent.max_children != null" class="text-gray-400"> / {{ agent.max_children }}</span>
          </el-descriptions-item>
        </el-descriptions>

        <!-- Attachments -->
        <h4 class="mt-4 mb-2">{{ t('admin.onboarding.attachments') }}</h4>
        <el-table v-if="agent.attachments?.length" :data="agent.attachments" size="small">
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
        <el-empty v-else :description="t('admin.onboarding.no_attachments')" :image-size="80" />

        <!-- Audit -->
        <el-descriptions :title="t('admin.agent.audit')" :column="2" border class="mt-4">
          <el-descriptions-item :label="t('admin.agent.created_at')">{{ agent.created_at || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.updated_at')">{{ agent.updated_at || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.created_by')">
            <UserRef :uuid="agent.created_by" :name="agent.created_by_name" />
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.updated_by')">
            <UserRef :uuid="agent.updated_by" :name="agent.updated_by_name" />
          </el-descriptions-item>
        </el-descriptions>
      </template>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, computed, watch, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElTooltip } from 'element-plus'
import { Unlock, Warning } from '@element-plus/icons-vue'
import { getAgentDetail, getAgentAdminDetail } from '@/api/skyagent/agent'

const { t } = useI18n()
const props = defineProps({ modelValue: Boolean, agentId: String })
const emit = defineEmits(['update:modelValue'])

const statusTagType = { pending_approval: 'warning', active: 'success', suspended: 'danger', terminated: 'info' }
const kycTagType = { 0: 'info', 1: '', 2: 'success' }

const loading = ref(false)
const loadingPlain = ref(false)
const agent = ref(null)
const showingPlain = ref(false)

const hasBusiness = computed(() => {
  const a = agent.value
  return a && (a.business_name || a.business_type || a.tax_code || a.representative_name || a.representative_cccd)
})

// Merge plaintext PII (from /full endpoint) over the masked values when the
// operator has explicitly opted in. Keeps masked defaults otherwise.
const display = computed(() => {
  const a = agent.value || {}
  return {
    phone: a.phone_plain || a.phone,
    email: a.email_plain || a.email,
    representative_cccd: a.representative_cccd_plain || a.representative_cccd,
    bank_account: a.bank_account_plain || a.bank_account,
  }
})

const loadPlain = async () => {
  if (!agent.value?.id) return
  loadingPlain.value = true
  try {
    const res = await getAgentAdminDetail(agent.value.id)
    if (res.code === 0 && res.data) {
      // Keep masked as fallback; overlay plaintext fields.
      agent.value = { ...agent.value, ...res.data }
      showingPlain.value = true
    } else {
      ElMessage.error(res.msg || t('admin.common.fail'))
    }
  } finally {
    loadingPlain.value = false
  }
}

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

watch(() => [props.modelValue, props.agentId], async ([visible, id]) => {
  if (visible && id) {
    showingPlain.value = false
    loading.value = true
    try {
      const res = await getAgentDetail(id)
      if (res.code === 0) agent.value = res.data
    } finally {
      loading.value = false
    }
  } else {
    agent.value = null
    showingPlain.value = false
  }
}, { immediate: true })
</script>

<style scoped>
.agent-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  background: #f7f8fa;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
}
.agent-header__meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.agent-header__name {
  font-size: 18px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}
.agent-header__code {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 13px;
  color: #606266;
}
.flex-1 { flex: 1; }
</style>
