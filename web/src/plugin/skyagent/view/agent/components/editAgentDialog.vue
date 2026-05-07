<template>
  <el-dialog
    :model-value="modelValue"
    :title="t('admin.agent.edit')"
    width="700px"
    destroy-on-close
    @close="emit('update:modelValue', false)"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="140px">
      <el-divider content-position="left">{{ t('admin.agent.basic_info') }}</el-divider>
      <el-form-item :label="t('admin.agent.full_name')" prop="full_name">
        <el-input v-model="form.full_name" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.phone')" prop="phone">
        <el-input v-model="form.phone" :placeholder="maskedHintFor('phone')" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.email')" prop="email">
        <el-input v-model="form.email" :placeholder="maskedHintFor('email')" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.business_type')" prop="business_type">
        <el-select v-model="form.business_type" clearable :placeholder="t('admin.common.select')">
          <el-option label="COMPANY" value="COMPANY" />
          <el-option label="HOUSEHOLD" value="HOUSEHOLD" />
          <el-option label="INDIVIDUAL" value="INDIVIDUAL" />
        </el-select>
      </el-form-item>

      <el-divider content-position="left">{{ t('admin.agent.address') }}</el-divider>
      <el-form-item :label="t('admin.agent.province')" prop="province">
        <el-input v-model="form.province" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.district')" prop="district">
        <el-input v-model="form.district" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.address_line')" prop="address_line">
        <el-input v-model="form.address_line" type="textarea" :rows="2" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.province_code')" prop="province_code">
        <el-input v-model="form.province_code" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.district_code')" prop="district_code">
        <el-input v-model="form.district_code" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.ward_code')" prop="ward_code">
        <el-input v-model="form.ward_code" />
      </el-form-item>

      <el-divider content-position="left">{{ t('admin.agent.referral_section') }}</el-divider>
      <el-form-item :label="t('admin.agent.partner_id')" prop="partner_id">
        <el-input v-model="form.partner_id" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.user_id')" prop="user_id">
        <el-input v-model="form.user_id" :placeholder="t('admin.agent.user_id_placeholder')" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.referral_code')" prop="referral_code">
        <el-input v-model="form.referral_code" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.referral_name')" prop="referral_name">
        <el-input v-model="form.referral_name" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:modelValue', false)">{{ t('admin.common.cancel') }}</el-button>
      <el-button type="primary" :loading="submitting" :disabled="!hasChanges" @click="onSubmit">{{ t('admin.common.save') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { updateAgent } from '@/plugin/skyagent/api/agent'

const { t } = useI18n()
const props = defineProps({ modelValue: Boolean, agent: Object })
const emit = defineEmits(['update:modelValue', 'success'])

const FIELDS = [
  'full_name', 'phone', 'email',
  'province', 'district', 'business_type',
  'address_line', 'province_code', 'district_code', 'ward_code',
  'partner_id', 'user_id', 'referral_code', 'referral_name',
]

// Fields that arrive masked (contain '*') from the detail endpoint. Per
// contract §14.1 / §14.6 rule 1 the FE must never re-submit a masked value —
// treat an unchanged masked value as "no change" and drop it from the patch.
const MASKED_FIELDS = new Set(['phone', 'email'])

const formRef = ref(null)
const submitting = ref(false)

// Two parallel snapshots: `form` is the edit buffer; `original` captures the
// server values so the submit path can compute a minimal diff.
const blank = () => FIELDS.reduce((acc, k) => { acc[k] = ''; return acc }, {})
const form = reactive(blank())
const original = reactive(blank())

const rules = {
  // vn_phone on Core side; accept empty (no change) or a 10-digit VN mobile.
  phone: [{ pattern: /^$|^(03|05|07|08|09)\d{8}$/, message: 'Phone must be 10 digits starting with 03/05/07/08/09', trigger: 'blur' }],
  email: [{ type: 'email', message: 'Invalid email', trigger: 'blur' }],
}

const isMasked = (v) => typeof v === 'string' && v.includes('*')
const maskedHintFor = (key) => {
  if (MASKED_FIELDS.has(key) && isMasked(original[key])) return original[key]
  return ''
}

watch(() => props.agent, (val) => {
  if (!val) return
  FIELDS.forEach((k) => {
    const v = val[k] ?? ''
    original[k] = v
    // Don't pre-populate the input with a masked value — the user would
    // otherwise "edit" a string that's already lossy.
    form[k] = isMasked(v) ? '' : v
  })
}, { immediate: true })

// Minimal patch: only fields the user changed from their initial value. For
// masked fields we require a fresh non-empty value; re-sending '' would clear
// the field, re-sending the mask would corrupt it.
const changedPatch = computed(() => {
  const patch = {}
  for (const k of FIELDS) {
    const current = form[k]
    const seed = isMasked(original[k]) ? '' : original[k]
    if ((current ?? '') !== (seed ?? '')) {
      if (MASKED_FIELDS.has(k) && current === '') continue
      patch[k] = current
    }
  }
  return patch
})

const hasChanges = computed(() => Object.keys(changedPatch.value).length > 0)

const onSubmit = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return
  }
  const patch = changedPatch.value
  if (!Object.keys(patch).length) {
    ElMessage.info(t('admin.common.no_changes'))
    return
  }
  submitting.value = true
  try {
    const res = await updateAgent(props.agent.id, patch)
    if (res.code === 0) {
      ElMessage.success(t('admin.agent.update_success'))
      emit('update:modelValue', false)
      emit('success')
    } else {
      if (Array.isArray(res.data?.details)) {
        res.data.details.forEach((d) => ElMessage.error(`${d.field}: ${d.message}`))
      } else {
        ElMessage.error(res.msg || t('admin.common.fail'))
      }
    }
  } finally {
    submitting.value = false
  }
}
</script>
