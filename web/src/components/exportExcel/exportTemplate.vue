<template>
  <el-button type="primary" icon="download" @click="exportTemplateFunc"
    >{{ t('admin.components.export_excel.download_template') }}</el-button
  >
</template>

<script setup>
  import { ElMessage } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import {exportTemplate} from "@/api/exportTemplate";

  const { t } = useI18n()

  const props = defineProps({
    templateId: {
      type: String,
      required: true
    }
  })


  const exportTemplateFunc = async () => {
    if (props.templateId === '') {
      ElMessage.error(t('admin.common.validation.template_id_required'))
      return
    }
    let baseUrl = import.meta.env.VITE_BASE_API
    if (baseUrl === "/"){
      baseUrl = ""
    }

    const res = await exportTemplate({
      templateID: props.templateId
    })

    if(res.code === 0){
      ElMessage.success(t('admin.components.export_excel.export_task_created'))
      const url = `${baseUrl}${res.data}`
      window.open(url, '_blank')
    }

  }
</script>
