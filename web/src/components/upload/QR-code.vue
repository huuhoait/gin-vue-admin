<template>
  <div>
    <el-button type="primary" icon="iphone" @click="createQrCode"> Scan to upload</el-button>
  </div>

  <el-dialog v-model="dialogVisible" title="Scan to upload" width="320px" :show-close="false" append-to-body :close-on-click-modal="false"
             draggable
  >
    <div class="m-2">
      <vue-qr :logoSrc="logoSrc"
              :size="291"
              :margin="0"
              :autoColor="true"
              :dotScale="1"
              :text="codeUrl"
              colorDark="green"
              colorLight="white"
              ref="qrcode"
      />
    </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="onFinished">Done</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import logoSrc from '@/assets/logo.png'
import vueQr from 'vue-qr/src/packages/vue-qr.vue'
import { ref } from 'vue'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'QRCodeUpload'
})

const emit = defineEmits(['on-success'])

const props = defineProps({
  classId: {
    type: Number,
    default: 0
  }
})

const dialogVisible = ref(false)
const userStore = useUserStore()
const codeUrl = ref('')

const createQrCode = () => {
  const local = window.location
  codeUrl.value = local.protocol + '//' + local.host + '/#/scanUpload?id=' + props.classId + '&token=' + userStore.token + '&t=' + Date.now()
  dialogVisible.value = true
  console.log(codeUrl.value)
}

const onFinished = () => {
  dialogVisible.value = false
  codeUrl.value = ''
  emit('on-success', '')
}
</script>
