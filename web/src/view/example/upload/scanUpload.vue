<template>
  <div class="flex justify-center w-full pt-2">
    <el-upload
        ref="uploadRef"
        class="h5-uploader"
        :action="`${getBaseUrl()}/fileUploadAndDownload/upload`"
        accept="image/*"
        :show-file-list="false"
        :auto-upload="false"
        :headers="{ 'x-token': token }"
        :data="{'classId': classId}"
        :on-success="handleImageSuccess"
        :on-change="handleFileChange"
    >
      <el-icon class="h5-uploader-icon"><Plus /></el-icon>
    </el-upload>
  </div>

  <div class="flex flex-col w-full h-auto p-0 pt-4">
    <!-- Left editor -->
    <div class="flex-1 min-h-[60vh]">
      <div class="w-screen h-[calc(100vh-175px)] rounded">
        <template v-if="isCrop">
          <VueCropper
              ref="cropperRef"
              :img="imgSrc"
              mode="contain"
              outputType="jpeg"
              :autoCrop="true"
              :autoCropWidth="cropWidth"
              :autoCropHeight="cropHeight"
              :fixedBox="false"
              :fixed="fixedRatio"
              :fixedNumber="fixedNumber"
              :centerBox="true"
              :canMoveBox="true"
              :full="false"
              :maxImgSize="windowWidth"
              :original="true"
          ></VueCropper>
        </template>
        <template v-else>
          <div class="flex justify-center items-center w-full h-[calc(100vh-175px)]">
            <el-image v-if="imgSrc" :src="imgSrc" class="max-w-full max-h-full" mode="cover" />
          </div>
        </template>
      </div>
    </div>
  </div>
  <!-- Toolbar -->
  <div class="toolbar">
    <el-button-group v-if="isCrop">
      <el-tooltip content="Rotate left">
        <el-button @click="rotate(-90)" :icon="RefreshLeft" />
      </el-tooltip>
      <el-tooltip content="Rotate right">
        <el-button @click="rotate(90)" :icon="RefreshRight" />
      </el-tooltip>
      <el-button :icon="Plus" @click="changeScale(1)"></el-button>
      <el-button :icon="Minus" @click="changeScale(-1)"></el-button>
    </el-button-group>


    <el-switch
        size="large"
        v-model="isCrop"
        inline-prompt
        active-text="Crop"
        inactive-text="Crop"
    />

    <el-button type="primary" @click="handleUpload" :loading="uploading"> {{ uploading ? 'Uploading...' : 'Upload' }}
    </el-button>
  </div>


</template>

<script setup>
import { ref, getCurrentInstance, onMounted } from 'vue'
import { ElLoading, ElMessage } from 'element-plus'
import { RefreshLeft, RefreshRight, Plus, Minus } from '@element-plus/icons-vue'
import 'vue-cropper/dist/index.css'
import { VueCropper } from 'vue-cropper'
import { getBaseUrl } from '@/utils/format'
import { useRouter } from 'vue-router'
import { useUserStore } from "@/pinia";

defineOptions({
  name: 'scanUpload'
})

const classId = ref(0)
const token = ref('')
const isCrop = ref(false)

const windowWidth = ref(300)

// Get screen width
const getWindowResize = function() {
  windowWidth.value = window.innerWidth
}

// Lifecycle
onMounted(() => {
  getWindowResize()
  window.addEventListener('resize', getWindowResize)
})

const router = useRouter()
router.isReady().then(() => {
  let query = router.currentRoute.value.query
  //console.log(query)
  classId.value = query.id
  token.value = query.token
}).catch((err) => {
  console.log(err)
})

const uploadRef = ref(null)
// Reactive state
const imgSrc = ref('')
const cropperRef = ref(null)
const { proxy } = getCurrentInstance()
const previews = ref({})
const uploading = ref(false)

// Zoom controls
const changeScale = (value) => {
  proxy.$refs.cropperRef.changeScale(value)
}

const fixedNumber = ref([1, 1])
const cropWidth = ref(300)
const cropHeight = ref(300)

const fixedRatio = ref(false)

// File handling
const handleFileChange = (file) => {
  const isImage = file.raw.type.includes('image')
  if (!isImage) {
    ElMessage.error('Please select an image file')
    return
  }

  if (file.raw.size / 1024 / 1024 > 8) {
    ElMessage.error('File size must be <= 8MB')
    return false
  }

  const loading = ElLoading.service({
    lock: true,
    text: 'Please wait',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  const reader = new FileReader()
  reader.onload = (e) => {
    imgSrc.value = e.target.result
    loading.close()
  }
  reader.readAsDataURL(file.raw)
}

// Rotate controls
const rotate = (degree) => {
  if (degree === -90) {
    proxy.$refs.cropperRef.rotateLeft()
  } else {
    proxy.$refs.cropperRef.rotateRight()
  }
}

// Upload handler
const handleUpload = () => {
  uploading.value = true
  if(isCrop.value === false){
    uploadRef.value.submit()
    return true
  }
  proxy.$refs.cropperRef.getCropBlob((blob) => {
    try {
      const file = new File([blob], `${Date.now()}.jpg`, { type: 'image/jpeg' })
      uploadRef.value.clearFiles()
      uploadRef.value.handleStart(file)
      uploadRef.value.submit()

    } catch (error) {
      uploading.value = false
      ElMessage.error('Upload failed: ' + error.message)
    }
  })
}

const handleImageSuccess = (res) => {
  const { data } = res
  if (data) {
    imgSrc.value = null
    uploading.value = false
    previews.value = {}
    ElMessage.success('Uploaded')
  }
}

</script>

<style scoped>

/* Toolbar (fixed at bottom) */
.toolbar {
  @apply fixed bottom-0 m-0 rounded-none p-2.5 shadow-[0_-2px_10px_rgba(0,0,0,0.1)] z-[1000] flex justify-between w-screen bg-slate-900;

  /* Button group layout */
  .el-button-group {
    @apply flex gap-2;

    .el-button {
      @apply p-2 w-10;
    }
  }
}

:deep(.vue-cropper) {
  @apply bg-transparent;
}

</style>

<style>
.h5-uploader .el-upload {
  @apply rounded cursor-pointer relative overflow-hidden;
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  transition: var(--el-transition-duration-fast);
}

.h5-uploader .el-upload:hover {
  border-color: var(--el-color-primary);
}

.el-icon.h5-uploader-icon {
  @apply text-2xl text-gray-500 w-32 h-32 text-center;
}
</style>
