<template>
  <div class="attachment-upload">
    <el-upload
      :action="uploadUrl"
      :headers="uploadHeaders"
      :data="uploadFormData"
      :show-file-list="false"
      :accept="accept"
      :before-upload="handleBeforeUpload"
      :on-success="handleSuccess"
      :on-error="handleError"
      name="file"
    >
      <el-button :loading="uploading" size="small">
        {{ modelValue ? t('admin.onboarding.replace_file') : t('admin.onboarding.choose_file') }}
      </el-button>
    </el-upload>
    <div v-if="modelValue" class="mt-1 flex items-center gap-2">
      <el-link :href="modelValue" target="_blank" type="primary" :underline="false">
        {{ fileName || t('admin.onboarding.uploaded') }}
      </el-link>
      <el-button size="small" text type="danger" @click="clear">
        {{ t('admin.onboarding.remove') }}
      </el-button>
    </div>
    <div v-if="hint" class="mt-1 text-xs text-gray-400">{{ hint }}</div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'

const props = defineProps({
  modelValue: { type: String, default: '' },
  accept: { type: String, default: 'image/jpeg,image/png,image/heic,application/pdf' },
  maxSizeMB: { type: Number, default: 10 },
  hint: { type: String, default: '' }
})
const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const userStore = useUserStore()

const uploading = ref(false)
const fileName = ref('')

// Native GVA upload endpoint: POST {VITE_BASE_API}/fileUploadAndDownload/upload
// returns { code, data: { file: { url, name, key, ... } }, msg }.
const uploadUrl = computed(
  () => `${import.meta.env.VITE_BASE_API}/fileUploadAndDownload/upload`
)
const uploadHeaders = computed(() => ({
  'x-token': userStore.token,
  'x-user-id': userStore.userInfo.ID
}))
const uploadFormData = { noSave: '0' }

const handleBeforeUpload = (file) => {
  const tooBig = file.size > props.maxSizeMB * 1024 * 1024
  if (tooBig) {
    ElMessage.error(t('admin.onboarding.file_too_large', { mb: props.maxSizeMB }))
    return false
  }
  uploading.value = true
  return true
}

const handleSuccess = (res, file) => {
  uploading.value = false
  if (res?.code === 0 && res?.data?.file?.url) {
    fileName.value = res.data.file.name || file.name
    emit('update:modelValue', res.data.file.url)
    ElMessage.success(t('admin.onboarding.upload_success'))
  } else {
    ElMessage.error(res?.msg || t('admin.onboarding.upload_failed'))
  }
}

const handleError = () => {
  uploading.value = false
  ElMessage.error(t('admin.onboarding.upload_failed'))
}

const clear = () => {
  fileName.value = ''
  emit('update:modelValue', '')
}
</script>
