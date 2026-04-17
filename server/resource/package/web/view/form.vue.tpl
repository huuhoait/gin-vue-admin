{{- if .IsAdd }}
// Add the following code in the new form
{{- range .Fields}}
     {{- if .Form}}
        {{ GenerateFormItem . }}
     {{- end }}
{{- end }}

// Add the following dictionary code
    {{- range $index, $element := .DictTypes}}
const {{ $element }}Options = ref([])
    {{- end }}

// Add the following calls inside the init method

{{- range $index, $element := .DictTypes }}
    {{ $element }}Options.value = await getDictFunc('{{$element}}')
{{- end }}

// Add the following fields to the base formData struct
{{- range .Fields}}
          {{- if .Form}}
            {{ GenerateDefaultFormValue . }}
          {{- end }}
        {{- end }}
// Add the following fields to the validation rules

{{- range .Fields }}
        {{- if .Form }}
            {{- if eq .Require true }}
{{.FieldJson }} : [{
    required: true,
    message: '{{ .ErrorText }}',
    trigger: ['input','blur'],
},
               {{- if eq .FieldType "string" }}
{
    whitespace: true,
    message: 'Whitespace only is not allowed',
    trigger: ['input', 'blur'],
}
              {{- end }}
],
            {{- end }}
        {{- end }}
    {{- end }}

{{- if .HasDataSource }}
// Please import
get{{.StructName}}DataSource,

// Fetch the data source
const dataSource = ref([])
const getDataSourceFunc = async()=>{
  const res = await get{{.StructName}}DataSource()
  if (res.code === 0) {
    dataSource.value = res.data
  }
}
getDataSourceFunc()
{{- end }}
{{- else }}
{{- if not .OnlyTemplate }}
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        {{- if .IsTree }}
          <el-form-item label="Parent node:" prop="parentID" >
              <el-tree-select
                  v-model="formData.parentID"
                  :data="[rootNode,...tableData]"
                  check-strictly
                  :render-after-expand="false"
                  :props="defaultProps"
                  clearable
                  style="width: 240px"
                  placeholder="Root node"
              />
          </el-form-item>
        {{- end }}
      {{- range .Fields}}
      {{- if .Form }}
        {{ GenerateFormItem . }}
      {{- end }}
      {{- end }}
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">Save</el-button>
          <el-button type="primary" @click="back">Back</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  {{- if .HasDataSource }}
    get{{.StructName}}DataSource,
  {{- end }}
  {{- if .IsTree }}
    get{{.StructName}}List,
  {{- end }}
  create{{.StructName}},
  update{{.StructName}},
  find{{.StructName}}
} from '@/api/{{.Package}}/{{.PackageName}}'

defineOptions({
    name: '{{.StructName}}Form'
})

// Auto-fetch dictionaries
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
{{- if .HasPic }}
// Image selector component
import SelectImage from '@/components/selectImage/selectImage.vue'
{{- end }}

{{- if .HasFile }}
// File selector component
import SelectFile from '@/components/selectFile/selectFile.vue'
{{- end }}

{{- if .HasRichText }}
// Rich text component
import RichEdit from '@/components/richtext/rich-edit.vue'
{{- end }}

{{- if .HasArray}}
// Array control component
import ArrayCtrl from '@/components/arrayCtrl/arrayCtrl.vue'
{{- end }}


const route = useRoute()
const router = useRouter()

{{- if .IsTree }}
const tableData = ref([])

const defaultProps = {
  children: "children",
  label: "{{ .TreeJson }}",
  value: "{{ .PrimaryField.FieldJson }}"
}

const rootNode = {
  {{ .PrimaryField.FieldJson }}: 0,
  {{ .TreeJson }}: 'Root node',
  children: []
}

const getTableData = async() => {
  const table = await get{{.StructName}}List()
  if (table.code === 0) {
    tableData.value = table.data || []
  }
}

getTableData()

{{- end }}

// Submit button loading state
const btnLoading = ref(false)

const type = ref('')
    {{- range $index, $element := .DictTypes}}
const {{ $element }}Options = ref([])
    {{- end }}
const formData = ref({
        {{- if .IsTree }}
            parentID: undefined,
        {{- end }}
        {{- range .Fields}}
          {{- if .Form }}
            {{ GenerateDefaultFormValue . }}
          {{- end }}
        {{- end }}
        })
// Validation rules
const rule = reactive({
    {{- range .Fields }}
            {{- if eq .Require true }}
               {{.FieldJson }} : [{
                   required: true,
                   message: '{{ .ErrorText }}',
                   trigger: ['input','blur'],
               }],
            {{- end }}
    {{- end }}
})

const elFormRef = ref()

{{- if .HasDataSource }}
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await get{{.StructName}}DataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()
{{- end }}

// Initialization
const init = async () => {
 // Suggestion: read the target ID from the URL query and call find to decide whether this page is a create or update. The example below reads id from the URL
    if (route.query.id) {
      const res = await find{{.StructName}}({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    {{- range $index, $element := .DictTypes }}
    {{ $element }}Options.value = await getDictFunc('{{$element}}')
    {{- end }}
}

init()
// Save button
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await create{{.StructName}}(formData.value)
               break
             case 'update':
               res = await update{{.StructName}}(formData.value)
               break
             default:
               res = await create{{.StructName}}(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: 'Created/updated successfully'
             })
           }
       })
}

// Back button
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
{{- else }}
<template>
<div>form</div>
</template>
<script setup>
</script>
<style>
</style>
{{- end }}
{{- end }}
