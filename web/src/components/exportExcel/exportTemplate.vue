<template>
  <el-button type="primary" icon="download" @click="exportTemplateFunc"
    >Download template</el-button
  >
</template>

<script setup>
  import { ElMessage } from 'element-plus'
  import {exportTemplate} from "@/api/exportTemplate";

  const props = defineProps({
    templateId: {
      type: String,
      required: true
    }
  })


  const exportTemplateFunc = async () => {
    if (props.templateId === '') {
      ElMessage.error('Template ID is required')
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
      ElMessage.success('Export task created. Downloading...')
      const url = `${baseUrl}${res.data}`
      window.open(url, '_blank')
    }

  }
</script>
