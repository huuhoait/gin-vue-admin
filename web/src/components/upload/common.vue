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
      <el-button type="primary" :icon="Upload">Upload</el-button>
    </el-upload>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { isVideoMime, isImageMime } from '@/utils/image'
  import { getBaseUrl } from '@/utils/format'
  import { Upload } from "@element-plus/icons-vue";
  import { useUserStore } from "@/pinia";

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
      ElMessage.error(
        'Images must be jpg/png/svg/webp. Videos must be mp4/webm.'
      )
      fullscreenLoading.value = false
      pass = false
    }
    if (!isLt5M && isVideo) {
      ElMessage.error('Video size must be <= 5MB')
      fullscreenLoading.value = false
      pass = false
    }
    if (!isLt500K && isImage) {
      ElMessage.error('Uncompressed image size must be <= 500KB. Please use compressed upload.')
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
      message: 'Upload failed'
    })
    fullscreenLoading.value = false
  }
</script>
