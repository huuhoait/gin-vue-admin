<template>
  <div>
    <warning-bar
      :title="t('admin.plugin.email.warning')"
    />
    <div class="gva-form-box">
      <el-form
        ref="emailForm"
        label-position="right"
        label-width="80px"
        :model="form"
      >
        <el-form-item :label="t('admin.plugin.email.to')">
          <el-input v-model="form.to" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.email.subject')">
          <el-input v-model="form.subject" />
        </el-form-item>
        <el-form-item :label="t('admin.plugin.email.body')">
          <el-input v-model="form.body" type="textarea" />
        </el-form-item>
        <el-form-item>
          <el-button @click="sendTestEmail">{{ t('admin.plugin.email.send_test') }}</el-button>
          <el-button @click="sendEmail">{{ t('admin.plugin.email.send') }}</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { emailTest } from '@/plugin/email/api/email.js'
  import { ElMessage } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import { reactive, ref } from 'vue'

  defineOptions({
    name: 'Email'
  })

  const { t } = useI18n()

  const emailForm = ref(null)
  const form = reactive({
    to: '',
    subject: '',
    body: ''
  })
  const sendTestEmail = async () => {
    const res = await emailTest()
    if (res.code === 0) {
      ElMessage.success(t('admin.plugin.email.sent'))
    }
  }

  const sendEmail = async () => {
    const res = await emailTest()
    if (res.code === 0) {
      ElMessage.success(t('admin.plugin.email.sent_check_inbox'))
    }
  }
</script>
