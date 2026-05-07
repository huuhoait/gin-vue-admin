<template>
  <el-drawer :model-value="modelValue" :title="t('admin.order.detail')" size="50%" destroy-on-close @close="emit('update:modelValue', false)">
    <div v-loading="loading">
      <template v-if="order">
        <el-descriptions :title="t('admin.order.info')" :column="2" border>
          <el-descriptions-item :label="t('admin.order.order_id')">{{ order.id }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.order.agent')">{{ order.agent_code }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.order.product')">{{ order.product_name }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.order.amount')">{{ formatVND(order.total_amount) }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.order.status')">
            <el-tag :type="statusTagType[order.status] || 'info'">{{ t(`admin.order.status_${order.status}`) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.order.created_at')">{{ order.created_at }}</el-descriptions-item>
        </el-descriptions>

        <h4 class="mt-4">{{ t('admin.order.request_payload') }}</h4>
        <el-input type="textarea" :rows="5" :model-value="jsonStr(order.request_payload)" readonly />

        <h4 class="mt-4">{{ t('admin.order.supplier_response') }}</h4>
        <el-input type="textarea" :rows="5" :model-value="jsonStr(order.supplier_response)" readonly />

        <h4 class="mt-4">{{ t('admin.order.webhook_history') }}</h4>
        <el-timeline v-if="order.webhook_logs?.length">
          <el-timeline-item
            v-for="(log, i) in order.webhook_logs"
            :key="i"
            :timestamp="log.timestamp"
            placement="top"
          >
            <strong>{{ log.event_type }}</strong> — Status {{ log.status_code }}
            <p v-if="log.body" class="text-xs text-gray-500">{{ truncate(log.body, 200) }}</p>
          </el-timeline-item>
        </el-timeline>
        <el-empty v-else :description="t('admin.order.no_webhooks')" />

        <template v-if="order.status === 'failed' && order.error_detail">
          <h4 class="mt-4">{{ t('admin.order.error_detail') }}</h4>
          <el-alert type="error" :closable="false">{{ order.error_detail }}</el-alert>
        </template>
      </template>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { getOrderDetail } from '@/api/skyagent/order'

const { t } = useI18n()
const props = defineProps({ modelValue: Boolean, orderId: String })
const emit = defineEmits(['update:modelValue'])

const statusTagType = { pending: 'warning', success: 'success', failed: 'danger', refunded: 'info' }
const loading = ref(false)
const order = ref(null)

const formatVND = (amount) => {
  if (!amount && amount !== 0) return '-'
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(amount)
}

const jsonStr = (obj) => {
  if (!obj) return '-'
  try { return typeof obj === 'string' ? obj : JSON.stringify(obj, null, 2) } catch { return String(obj) }
}

const truncate = (str, len) => str && str.length > len ? str.slice(0, len) + '...' : str

watch(() => [props.modelValue, props.orderId], async ([visible, id]) => {
  if (visible && id) {
    loading.value = true
    try {
      const res = await getOrderDetail(id)
      if (res.code === 0) {
        order.value = res.data
      }
    } finally {
      loading.value = false
    }
  } else {
    order.value = null
  }
}, { immediate: true })
</script>
