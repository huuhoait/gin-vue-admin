<template>
  <div>
    <warning-bar
      title="Email must be configured first. To prevent spam, this feature is disabled in online demo environments."
    />
    <div class="gva-form-box">
      <el-form
        ref="emailForm"
        label-position="right"
        label-width="80px"
        :model="form"
      >
        <el-form-item label="To">
          <el-input v-model="form.to" />
        </el-form-item>
        <el-form-item label="Subject">
          <el-input v-model="form.subject" />
        </el-form-item>
        <el-form-item label="Body">
          <el-input v-model="form.body" type="textarea" />
        </el-form-item>
        <el-form-item>
          <el-button @click="sendTestEmail">Send test email</el-button>
          <el-button @click="sendEmail">Send email</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { emailTest } from '@/plugin/email/api/email.js'
  import { ElMessage } from 'element-plus'
  import { reactive, ref } from 'vue'

  defineOptions({
    name: 'Email'
  })

  const emailForm = ref(null)
  const form = reactive({
    to: '',
    subject: '',
    body: ''
  })
  const sendTestEmail = async () => {
    const res = await emailTest()
    if (res.code === 0) {
      ElMessage.success('Sent')
    }
  }

  const sendEmail = async () => {
    const res = await emailTest()
    if (res.code === 0) {
      ElMessage.success('Sent, please check your inbox')
    }
  }
</script>
