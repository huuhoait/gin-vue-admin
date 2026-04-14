<template>
  <div class="gva-form-box">
    <el-form :model="form" ref="formRef" label-width="100px" :rules="rules">
      <el-form-item label="Tool name" prop="name">
        <el-input v-model="form.name" placeholder="e.g. CurrentTime" />
      </el-form-item>
      <el-form-item label="Description" prop="description">
        <el-input type="textarea" v-model="form.description" placeholder="Enter tool description" />
      </el-form-item>
      <el-form-item label="Parameters">
        <el-table :data="form.params"  style="width: 100%">
          <el-table-column prop="name" label="Name" width="120">
            <template #default="scope">
              <el-input v-model="scope.row.name" placeholder="Name" />
            </template>
          </el-table-column>
          <el-table-column prop="description" label="Description" min-width="180">
            <template #default="scope">
              <el-input v-model="scope.row.description" placeholder="Description" />
            </template>
          </el-table-column>
          <el-table-column prop="type" label="Type" width="120">
            <template #default="scope">
              <el-select v-model="scope.row.type" placeholder="Type">
                <el-option label="string" value="string" />
                <el-option label="number" value="number" />
                <el-option label="boolean" value="boolean" />
                <el-option label="object" value="object" />
                <el-option label="array" value="array" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="Default" width="300">
            <template #default="scope">
              <el-input :disabled="scope.row.type === 'object'" v-model="scope.row.default" />
            </template>
          </el-table-column>
          <el-table-column prop="required" label="Required" width="80">
            <template #default="scope">
              <el-checkbox v-model="scope.row.required" />
            </template>
          </el-table-column>
          <el-table-column label="Actions" width="80">
            <template #default="scope">
              <el-button type="text" @click="removeParam(scope.$index)">Delete</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <div class="flex justify-end">
        <el-button type="primary" icon="plus" @click="addParam" style="margin-top: 10px;">Add parameter</el-button>
      </div>
      <el-form-item label="Responses">
        <el-table :data="form.response" style="width: 100%">
          <el-table-column prop="type" label="Type" min-width="120">
            <template #default="scope">
              <el-select v-model="scope.row.type" placeholder="Type">
                <el-option label="text" value="text" />
                <el-option label="image" value="image" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="Actions" width="80">
            <template #default="scope">
              <el-button type="text" @click="removeResponse(scope.$index)">Delete</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <div class="flex justify-end">
        <el-button type="primary" icon="plus" @click="addResponse" style="margin-top: 10px;">Add response</el-button>
      </div>

      <div class="flex justify-end mt-8">
        <el-button type="primary" @click="submit">Generate</el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { mcp } from '@/api/autoCode'

defineOptions({
  name: 'MCP'
})

const formRef = ref(null)
const form = reactive({
  name: '',
  description: '',
  type: '',
  params: [],
  response: []
})

const rules = {
  name: [{ required: true, message: 'Tool name is required', trigger: 'blur' }],
  description: [{ required: true, message: 'Description is required', trigger: 'blur' }],
  type: [{ required: true, message: 'Type is required', trigger: 'change' }]
}

function addParam() {
  form.params.push({
    name: '',
    description: '',
    type: '',
    required: false
  })
}

function removeParam(index) {
  form.params.splice(index, 1)
}

function addResponse() {
  form.response.push({
    type: ''
  })
}

function removeResponse(index) {
  form.response.splice(index, 1)
}

function submit() {
  formRef.value.validate(async (valid) => {
    if (!valid) return
    // Basic param validation
    for (const p of form.params) {
      if (!p.name || !p.description || !p.type) {
        ElMessage.error('Please complete all parameter fields')
        return
      }
    }
    // Validate response types
    for (const r of form.response) {
      if (!r.type) {
        ElMessage.error('Please select response type')
        return
      }
    }
      const res = await mcp(form)
      if (res.code === 0) {
        ElMessage.success(res.msg)
      }
  })
}
</script>
