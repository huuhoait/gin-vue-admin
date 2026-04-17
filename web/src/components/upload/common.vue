<template>
  <div>
    <el-upload
      :action="`${getBaseUrl()}/fileUploadAndDownload/upload`"
      :before-upload="checkFile"
      :on-error="uploadError"
      :on-success="uploadSuccess"
      :show-file-list="false"
      :data="{'classId': props.classId}"
      :headers="{'x-token': token}"
      multiple
      class="upload-btn"
    >
      <el-button type="primary" :icon="Upload">{{ t('admin.common.upload') }}</el-button>
    </el-upload>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import { isVideoMime, isImageMime } from '@/utils/image'
  import { getBaseUrl } from '@/utils/format'
  import { Upload } from "@element-plus/icons-vue";
  import { useUserStore } from "@/pinia";

  const { t } = useI18n()

  defineOptions({
    name: 'UploadCommon'
  })

  const userStore = useUserStore()

  const token = userStore.token

  const props = defineProps({
    classId: {
      type: Number,
      default: 0
    }
  })

  const emit = defineEmits(['on-success'])

  const fullscreenLoading = ref(false)

  const checkFile = (file) => {
    fullscreenLoading.value = true
    const isLt500K = file.size / 1024 / 1024 < 0.5 // 500K, should be configurable
    const isLt5M = file.size / 1024 / 1024 < 5 // 5MB, should be configurable
    const isVideo = isVideoMime(file.type)
    const isImage = isImageMime(file.type)
    let pass = true
    if (!isVideo && !isImage) {
      ElMessage.error(t('admin.common.validation.upload_type_limit'))
      fullscreenLoading.value = false
      pass = false
    }
    if (!isLt5M && isVideo) {
      ElMessage.error(t('admin.common.validation.video_size_limit'))
      fullscreenLoading.value = false
      pass = false
    }
    if (!isLt500K && isImage) {
      ElMessage.error(t('admin.common.validation.image_size_limit'))
      fullscreenLoading.value = false
      pass = false
    }

    console.log('upload file check result: ', pass)

    return pass
  }

  const uploadSuccess = (res) => {
    const { data } = res
    if (data.file) {
      emit('on-success', data.file.url)
    }
  }

  const uploadError = () => {
    ElMessage({
      type: 'error',
      message: t('admin.common.messages.upload_failed')
    })
    fullscreenLoading.value = false
  }
</script>
