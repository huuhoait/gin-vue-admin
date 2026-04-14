<template>
  <el-dialog :model-value="modelValue" :title="t('admin.agent.edit')" width="500px" destroy-on-close @close="emit('update:modelValue', false)">
    <el-form ref="formRef" :model="form" label-width="120px">
      <el-form-item :label="t('admin.agent.full_name')" prop="full_name">
        <el-input v-model="form.full_name" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.partner_id')">
        <el-input v-model="form.partner_id" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.user_id')">
        <el-input v-model="form.user_id" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:modelValue', false)">{{ t('admin.common.cancel') }}</el-button>
      <el-button type="primary" :loading="submitting" @click="onSubmit">{{ t('admin.common.save') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { updateAgent } from '@/api/skyagent/agent'

const { t } = useI18n()
const props = defineProps({ modelValue: Boolean, agent: Object })
const emit = defineEmits(['update:modelValue', 'success'])

const formRef = ref(null)
const submitting = ref(false)
const form = reactive({ full_name: '', partner_id: '', user_id: '' })

watch(() => props.agent, (val) => {
  if (val) {
    form.full_name = val.full_name || ''
    form.partner_id = val.partner_id || ''
    form.user_id = val.user_id || ''
  }
}, { immediate: true })

const onSubmit = async () => {
  submitting.value = true
  try {
    const res = await updateAgent(props.agent.id, form)
    if (res.code === 0) {
      ElMessage.success(t('admin.agent.update_success'))
      emit('update:modelValue', false)
      emit('success')
    } else {
      if (res.data?.details) {
        res.data.details.forEach((d) => {
          const field = formRef.value.fields?.find((f) => f.prop === d.field)
          if (field) field.validateMessage = d.message
        })
      }
      ElMessage.error(res.msg || t('admin.common.fail'))
    }
  } finally {
    submitting.value = false
  }
}
</script>
