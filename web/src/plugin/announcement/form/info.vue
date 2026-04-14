<template>
  <div>
    <div class="gva-form-box">
      <el-form
        :model="formData"
        ref="elFormRef"
        label-position="right"
        :rules="rule"
        label-width="80px"
      >
        <el-form-item label="Title:" prop="title">
          <el-input
            v-model="formData.title"
            :clearable="true"
            placeholder="Enter title"
          />
        </el-form-item>
        <el-form-item label="Content:" prop="content">
          <RichEdit v-model="formData.content" />
        </el-form-item>
        <el-form-item label="Author:" prop="userID">
          <el-select
            v-model="formData.userID"
            placeholder="Select author"
            style="width: 100%"
            :clearable="true"
          >
            <el-option
              v-for="(item, key) in dataSource.userID"
              :key="key"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Attachments:" prop="attachments">
          <SelectFile v-model="formData.attachments" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">Save</el-button>
          <el-button type="primary" @click="back">Back</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
  import {
    getInfoDataSource,
    createInfo,
    updateInfo,
    findInfo
  } from '@/plugin/announcement/api/info'

  defineOptions({
    name: 'InfoForm'
  })

  // Auto-load dictionary (if needed)
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { ref, reactive } from 'vue'
  import SelectFile from '@/components/selectFile/selectFile.vue'
  // Rich text editor
  import RichEdit from '@/components/richtext/rich-edit.vue'

  const route = useRoute()
  const router = useRouter()

  const type = ref('')
  const formData = ref({
    title: '',
    content: '',
    userID: undefined,
    attachments: []
  })
  // Validation rules
  const rule = reactive({})

  const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async () => {
    const res = await getInfoDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

  // Initialize
  const init = async () => {
    // Use URL param id to decide create vs update; call find() to load existing data
    if (route.query.id) {
      const res = await findInfo({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
  }

  init()
  // Save
  const save = async () => {
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      let res
      switch (type.value) {
        case 'create':
          res = await createInfo(formData.value)
          break
        case 'update':
          res = await updateInfo(formData.value)
          break
        default:
          res = await createInfo(formData.value)
          break
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: 'Saved'
        })
      }
    })
  }

  // Back
  const back = () => {
    router.go(-1)
  }
</script>

<style></style>
