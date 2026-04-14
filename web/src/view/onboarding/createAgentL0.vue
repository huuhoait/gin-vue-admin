<template>
  <div>
    <div class="gva-search-box">
      <h3>{{ t('admin.onboarding.create_agent_l0') }}</h3>
      <p class="text-sm text-gray-500">{{ t('admin.onboarding.create_agent_desc') }}</p>
    </div>

    <div class="gva-table-box" style="max-width: 800px">
      <el-steps :active="activeStep" finish-status="success" align-center class="mb-6">
        <el-step :title="t('admin.onboarding.step_basic')" />
        <el-step :title="t('admin.onboarding.step_business')" />
        <el-step :title="t('admin.onboarding.step_bank')" />
        <el-step :title="t('admin.onboarding.step_docs')" />
        <el-step :title="t('admin.onboarding.step_confirm')" />
      </el-steps>

      <!-- Step 1: Basic Info -->
      <el-form v-show="activeStep === 0" ref="step1Ref" :model="form" :rules="step1Rules" label-width="140px">
        <el-form-item :label="t('admin.agent.full_name')" prop="full_name">
          <el-input v-model="form.full_name" />
        </el-form-item>
        <el-form-item :label="t('admin.agent.phone')" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item :label="t('admin.agent.email')" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item :label="t('admin.agent.province')" prop="province">
          <el-input v-model="form.province" />
        </el-form-item>
        <el-form-item :label="t('admin.agent.district')">
          <el-input v-model="form.district" />
        </el-form-item>
        <el-form-item :label="t('admin.agent.parent_id')" prop="parent_id">
          <el-input v-model="form.parent_id" placeholder="UUID" />
        </el-form-item>
        <el-form-item :label="t('admin.agent.referral_code')">
          <el-input v-model="form.referral_code" />
        </el-form-item>
        <el-form-item :label="t('admin.agent.referral_name')">
          <el-input v-model="form.referral_name" />
        </el-form-item>
      </el-form>

      <!-- Step 2: Business Info -->
      <el-form v-show="activeStep === 1" ref="step2Ref" :model="form" label-width="160px">
        <el-form-item :label="t('admin.onboarding.business_name')">
          <el-input v-model="form.business_name" />
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.business_type')">
          <el-input v-model="form.business_type" />
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.tax_code')">
          <el-input v-model="form.tax_code" />
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.representative_name')">
          <el-input v-model="form.representative_name" />
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.representative_cccd')">
          <el-input v-model="form.representative_cccd" />
        </el-form-item>
      </el-form>

      <!-- Step 3: Bank & Address -->
      <el-form v-show="activeStep === 2" ref="step3Ref" :model="form" label-width="160px">
        <el-form-item :label="t('admin.onboarding.bank_name')">
          <el-input v-model="form.bank_name" />
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.bank_account')">
          <el-input v-model="form.bank_account" />
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.permanent_address')">
          <el-input v-model="form.permanent_address" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.contact_address')">
          <el-input v-model="form.contact_address" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>

      <!-- Step 4: Documents -->
      <el-form v-show="activeStep === 3" ref="step4Ref" :model="form" label-width="160px">
        <el-form-item :label="t('admin.onboarding.att_business_license')">
          <el-input v-model="form.business_license_url" placeholder="https://..." />
          <el-tag v-if="form.business_license_url" type="success" size="small" class="ml-2">{{ t('admin.onboarding.uploaded') }}</el-tag>
          <el-tag v-else type="danger" size="small" class="ml-2">{{ t('admin.onboarding.required') }}</el-tag>
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.att_agency_contract')">
          <el-input v-model="form.agency_contract_url" placeholder="https://..." />
          <el-tag v-if="form.agency_contract_url" type="success" size="small" class="ml-2">{{ t('admin.onboarding.uploaded') }}</el-tag>
          <el-tag v-else type="danger" size="small" class="ml-2">{{ t('admin.onboarding.required') }}</el-tag>
        </el-form-item>
        <el-form-item :label="t('admin.onboarding.att_ekyc_photos')">
          <el-input v-model="form.ekyc_photos_url" placeholder="https://..." />
          <el-tag v-if="form.ekyc_photos_url" type="success" size="small" class="ml-2">{{ t('admin.onboarding.uploaded') }}</el-tag>
        </el-form-item>
      </el-form>

      <!-- Step 5: Confirm -->
      <div v-show="activeStep === 4">
        <el-descriptions :title="t('admin.onboarding.step_confirm')" :column="2" border>
          <el-descriptions-item :label="t('admin.agent.full_name')">{{ form.full_name }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.phone')">{{ form.phone }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.email')">{{ form.email }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.agent.province')">{{ form.province }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.business_name')">{{ form.business_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.tax_code')">{{ form.tax_code || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.bank_name')">{{ form.bank_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.bank_account')">{{ form.bank_account || '-' }}</el-descriptions-item>
        </el-descriptions>

        <el-descriptions class="mt-4" :title="t('admin.onboarding.attachments')" :column="1" border>
          <el-descriptions-item :label="t('admin.onboarding.att_business_license')">
            <el-tag :type="form.business_license_url ? 'success' : 'danger'" size="small">{{ form.business_license_url ? t('admin.onboarding.uploaded') : t('admin.onboarding.missing') }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.att_agency_contract')">
            <el-tag :type="form.agency_contract_url ? 'success' : 'danger'" size="small">{{ form.agency_contract_url ? t('admin.onboarding.uploaded') : t('admin.onboarding.missing') }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('admin.onboarding.att_ekyc_photos')">
            <el-tag :type="form.ekyc_photos_url ? 'success' : 'info'" size="small">{{ form.ekyc_photos_url ? t('admin.onboarding.uploaded') : t('admin.onboarding.optional') }}</el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <el-alert v-if="!allRequiredDocs" type="warning" :closable="false" class="mt-4">
          {{ t('admin.onboarding.missing_docs_warning') }}
        </el-alert>
      </div>

      <!-- Navigation -->
      <div class="flex justify-between mt-6">
        <el-button v-if="activeStep > 0" @click="activeStep--">{{ t('admin.onboarding.prev') }}</el-button>
        <div v-else />
        <div>
          <el-button v-if="activeStep < 4" type="primary" @click="nextStep">{{ t('admin.onboarding.next') }}</el-button>
          <template v-if="activeStep === 4">
            <el-button @click="handleSubmit('draft')" :loading="submitting">{{ t('admin.onboarding.save_draft') }}</el-button>
            <el-button type="primary" @click="handleSubmit('submit')" :loading="submitting">{{ t('admin.onboarding.submit_review') }}</el-button>
          </template>
        </div>
      </div>

      <!-- Success result -->
      <el-result v-if="result" icon="success" :title="t('admin.onboarding.create_success')" class="mt-6">
        <template #sub-title>
          <p>Ticket ID: <el-tag effect="plain">{{ result.ticket_id }}</el-tag></p>
          <p>Agent ID: <el-tag effect="plain">{{ result.agent_id }}</el-tag></p>
          <p>{{ t('admin.onboarding.status') }}: <el-tag :type="result.status === 'pending_review' ? 'success' : 'info'">{{ t(`admin.onboarding.status_${result.status}`) }}</el-tag></p>
        </template>
        <template #extra>
          <el-button type="primary" @click="$router.push('/onboarding/tickets')">{{ t('admin.onboarding.go_tickets') }}</el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { onboardingAgent } from '@/api/skyagent/onboarding'

const { t } = useI18n()

const activeStep = ref(0)
const submitting = ref(false)
const result = ref(null)
const step1Ref = ref(null)

const form = reactive({
  full_name: '', phone: '', email: '', province: '', district: '',
  parent_id: '', referral_code: '', referral_name: '',
  business_name: '', business_type: '', tax_code: '',
  representative_name: '', representative_cccd: '',
  bank_name: '', bank_account: '',
  permanent_address: '', contact_address: '',
  business_license_url: '', agency_contract_url: '', ekyc_photos_url: '',
})

const step1Rules = {
  full_name: [{ required: true, message: t('admin.agent.full_name'), trigger: 'blur' }],
  phone: [{ required: true, message: t('admin.agent.phone'), trigger: 'blur' }],
  email: [{ required: true, message: t('admin.agent.email'), trigger: 'blur' }],
  province: [{ required: true, message: t('admin.agent.province'), trigger: 'blur' }],
  parent_id: [{ required: true, message: t('admin.agent.parent_id'), trigger: 'blur' }],
}

const allRequiredDocs = computed(() => form.business_license_url && form.agency_contract_url)

const nextStep = async () => {
  if (activeStep.value === 0 && step1Ref.value) {
    try {
      await step1Ref.value.validate()
    } catch { return }
  }
  activeStep.value++
}

const handleSubmit = async (mode) => {
  submitting.value = true
  try {
    const res = await onboardingAgent({ ...form, mode })
    if (res.code === 0) {
      result.value = res.data
      ElMessage.success(res.data.message || t('admin.onboarding.create_success'))
    } else {
      if (res.data?.details) {
        res.data.details.forEach(d => ElMessage.error(`${d.field}: ${d.message}`))
      } else {
        ElMessage.error(res.msg || t('admin.common.fail'))
      }
    }
  } finally {
    submitting.value = false
  }
}
</script>
