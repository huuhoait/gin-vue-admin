<template>
  <el-dialog :model-value="modelValue" :title="t('admin.agent.create')" width="500px" destroy-on-close @close="emit('update:modelValue', false)">
    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
      <el-form-item :label="t('admin.agent.full_name')" prop="full_name">
        <el-input v-model="form.full_name" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.phone')" prop="phone">
        <el-input v-model="form.phone" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.email')" prop="email">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item :label="t('admin.agent.parent_id')" prop="parent_id">
        <el-input v-model="form.parent_id" :placeholder="t('admin.agent.parent_id')" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:modelValue', false)">{{ t('admin.common.cancel') }}</el-button>
      <el-button type="primary" :loading="submitting" @click="onSubmit">{{ t('admin.common.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { createAgent } from '@/plugin/skyagent/api/agent'

const { t } = useI18n()
const props = defineProps({ modelValue: Boolean })
const emit = defineEmits(['update:modelValue', 'success'])

const formRef = ref(null)
const submitting = ref(false)
const form = reactive({ full_name: '', phone: '', email: '', parent_id: '' })
const rules = {
  full_name: [{ required: true, message: t('admin.agent.full_name'), trigger: 'blur' }],
  phone: [{ required: true, message: t('admin.agent.phone'), trigger: 'blur' }],
  email: [{ required: true, message: t('admin.agent.email'), trigger: 'blur' }],
}

const onSubmit = async () => {
  await formRef.value.validate()
  submitting.value = true
  try {
    const res = await createAgent(form)
    if (res.code === 0) {
      ElMessage.success(t('admin.agent.create_success'))
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
