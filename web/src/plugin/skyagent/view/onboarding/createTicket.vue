<template>
  <div>
    <div class="gva-search-box">
      <h3>{{ t('admin.onboarding.create_ticket') }}</h3>
      <p class="text-sm text-gray-500">{{ t('admin.onboarding.create_desc') }}</p>
    </div>

    <div class="gva-table-box">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="160px" style="max-width: 600px">
        <el-form-item :label="t('admin.onboarding.agent_id')" prop="agent_id">
          <el-input v-model="form.agent_id" placeholder="UUID" />
        </el-form-item>

        <el-divider>{{ t('admin.onboarding.attachments') }}</el-divider>

        <el-form-item :label="t('admin.onboarding.att_business_license')">
          <el-input v-model="form.business_license_url" placeholder="https://..." />
          <el-tag v-if="form.business_license_url" type="success" size="small" class="ml-2">{{ t('admin.onboarding.uploaded') }}</el-tag>
        </el-form-item>

        <el-form-item :label="t('admin.onboarding.att_agency_contract')">
          <el-input v-model="form.agency_contract_url" placeholder="https://..." />
          <el-tag v-if="form.agency_contract_url" type="success" size="small" class="ml-2">{{ t('admin.onboarding.uploaded') }}</el-tag>
        </el-form-item>

        <el-form-item :label="t('admin.onboarding.att_ekyc_photos')">
          <el-input v-model="form.ekyc_photos_url" placeholder="https://..." />
          <el-tag v-if="form.ekyc_photos_url" type="success" size="small" class="ml-2">{{ t('admin.onboarding.uploaded') }}</el-tag>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="creating" @click="handleCreate">{{ t('admin.onboarding.create_and_save') }}</el-button>
          <el-button v-if="ticketId" type="success" :loading="submitting" @click="handleSubmit">{{ t('admin.onboarding.submit') }}</el-button>
        </el-form-item>
      </el-form>

      <el-result v-if="ticketId" icon="success" :title="t('admin.onboarding.create_success')">
        <template #sub-title>
          <span>Ticket ID: <strong>{{ ticketId }}</strong></span>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { createTicket, uploadAttachment, submitTicket } from '@/api/skyagent/ticket'

const { t } = useI18n()

const formRef = ref(null)
const creating = ref(false)
const submitting = ref(false)
const ticketId = ref('')
const form = reactive({
  agent_id: '',
  business_license_url: '',
  agency_contract_url: '',
  ekyc_photos_url: '',
})
const rules = {
  agent_id: [{ required: true, message: 'Agent ID required', trigger: 'blur' }],
}

const handleCreate = async () => {
  await formRef.value.validate()
  creating.value = true
  try {
    const res = await createTicket({ agent_id: form.agent_id })
    if (res.code === 0) {
      ticketId.value = res.data.ticket_id
      ElMessage.success(t('admin.onboarding.create_success'))

      // Upload attachments if provided
      const attachments = [
        { type: 'business_license', url: form.business_license_url },
        { type: 'agency_contract', url: form.agency_contract_url },
        { type: 'ekyc_photos', url: form.ekyc_photos_url },
      ].filter(a => a.url)

      for (const att of attachments) {
        await uploadAttachment(ticketId.value, att)
      }
    } else {
      ElMessage.error(res.msg || t('admin.common.fail'))
    }
  } finally {
    creating.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    const res = await submitTicket(ticketId.value)
    if (res.code === 0) {
      ElMessage.success(t('admin.onboarding.submit_success'))
    } else {
      ElMessage.error(res.msg || t('admin.common.fail'))
    }
  } finally {
    submitting.value = false
  }
}
</script>
