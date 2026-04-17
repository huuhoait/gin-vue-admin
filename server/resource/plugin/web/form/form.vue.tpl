{{- if .IsAdd }}
// Add the following code in the new form
{{- range .Fields}}
          {{- if .Form}}
<el-form-item label="{{.FieldDesc}}:"  prop="{{.FieldJson}}" >
          {{- if .CheckDataSource}}
    <el-select {{if eq .DataSource.Association 2}} multiple {{ end }} v-model="formData.{{.FieldJson}}" placeholder="Please select {{.FieldDesc}}" style="width:100%" :clearable="{{.Clearable}}" >
        <el-option v-for="(item,key) in dataSource.{{.FieldJson}}" :key="key" :label="item.label" :value="item.value" />
    </el-select>
          {{- else }}
          {{- if eq .FieldType "bool" }}
    <el-switch v-model="formData.{{.FieldJson}}" active-color="#13ce66" inactive-color="#ff4949" active-text="Yes" inactive-text="No" clearable ></el-switch>
          {{- end }}
          {{- if eq .FieldType "string" }}
          {{- if .DictType}}
    <el-select {{if eq .FieldType "array"}}multiple {{end}}v-model="formData.{{ .FieldJson }}" placeholder="Please select {{.FieldDesc}}" style="width:100%" :clearable="{{.Clearable}}" >
        <el-option v-for="(item,key) in {{ .DictType }}Options" :key="key" :label="item.label" :value="item.value" />
    </el-select>
          {{- else }}
    <el-input v-model="formData.{{.FieldJson}}" :clearable="{{.Clearable}}"  placeholder="Please enter {{.FieldDesc}}" />
          {{- end }}
          {{- end }}
          {{- if eq .FieldType "richtext" }}
    <RichEdit v-model="formData.{{.FieldJson}}"/>
          {{- end }}
          {{- if eq .FieldType "json" }}
    // This field is a JSON structure; control its rendering and binding on the frontend. Bind it to formData.{{.FieldJson}}; the backend stores/reads it according to its JSON type
    {{"{{"}} formData.{{.FieldJson}} {{"}}"}}
          {{- end }}
           {{- if eq .FieldType "array" }}
    <ArrayCtrl v-model="formData.{{ .FieldJson }}" editable/>
           {{- end }}
          {{- if eq .FieldType "int" }}
    <el-input v-model.number="formData.{{ .FieldJson }}" :clearable="{{.Clearable}}" placeholder="Please enter {{.FieldDesc}}" />
          {{- end }}
          {{- if eq .FieldType "time.Time" }}
    <el-date-picker v-model="formData.{{ .FieldJson }}" type="date" style="width:100%" placeholder="Select date" :clearable="{{.Clearable}}"  />
          {{- end }}
          {{- if eq .FieldType "float64" }}
    <el-input-number v-model="formData.{{ .FieldJson }}"  style="width:100%" :precision="2" :clearable="{{.Clearable}}"  />
          {{- end }}
          {{- if eq .FieldType "enum" }}
    <el-select v-model="formData.{{ .FieldJson }}" placeholder="Please select {{.FieldDesc}}" style="width:100%" :clearable="{{.Clearable}}" >
       <el-option v-for="item in [{{.DataTypeLong}}]" :key="item" :label="item" :value="item" />
    </el-select>
          {{- end }}
          {{- if eq .FieldType "picture" }}
    <SelectImage
     v-model="formData.{{ .FieldJson }}"
     file-type="image"
    />
          {{- end }}
          {{- if eq .FieldType "pictures" }}
    <SelectImage
     multiple
     v-model="formData.{{ .FieldJson }}"
     file-type="image"
     />
          {{- end }}
          {{- if eq .FieldType "video" }}
    <SelectImage
    v-model="formData.{{ .FieldJson }}"
    file-type="video"
    />
           {{- end }}
          {{- if eq .FieldType "file" }}
    <SelectFile v-model="formData.{{ .FieldJson }}" />
          {{- end }}
          {{- end }}
</el-form-item>
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
            {{- if eq .FieldType "bool" }}
{{.FieldJson}}: false,
            {{- end }}
            {{- if eq .FieldType "string" }}
{{.FieldJson}}: '',
            {{- end }}
            {{- if eq .FieldType "richtext" }}
{{.FieldJson}}: '',
            {{- end }}
            {{- if eq .FieldType "int" }}
{{.FieldJson}}: {{- if or .DataSource}} undefined{{ else }} 0{{- end }},
            {{- end }}
            {{- if eq .FieldType "time.Time" }}
{{.FieldJson}}: new Date(),
            {{- end }}
            {{- if eq .FieldType "float64" }}
{{.FieldJson}}: 0,
            {{- end }}
            {{- if eq .FieldType "picture" }}
{{.FieldJson}}: "",
            {{- end }}
            {{- if eq .FieldType "video" }}
{{.FieldJson}}: "",
            {{- end }}
            {{- if eq .FieldType "pictures" }}
{{.FieldJson}}: [],
            {{- end }}
            {{- if eq .FieldType "file" }}
{{.FieldJson}}: [],
            {{- end }}
            {{- if eq .FieldType "json" }}
{{.FieldJson}}: {},
            {{- end }}
            {{- if eq .FieldType "array" }}
{{.FieldJson}}: [],
            {{- end }}
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
        <el-form-item label="{{.FieldDesc}}:" prop="{{.FieldJson}}">
       {{- if .CheckDataSource}}
        <el-select {{if eq .DataSource.Association 2}} multiple {{ end }} v-model="formData.{{.FieldJson}}" placeholder="Please select {{.FieldDesc}}" style="width:100%" :clearable="{{.Clearable}}" >
          <el-option v-for="(item,key) in dataSource.{{.FieldJson}}" :key="key" :label="item.label" :value="item.value" />
        </el-select>
       {{- else }}
      {{- if eq .FieldType "bool" }}
          <el-switch v-model="formData.{{.FieldJson}}" active-color="#13ce66" inactive-color="#ff4949" active-text="Yes" inactive-text="No" clearable ></el-switch>
      {{- end }}
      {{- if eq .FieldType "string" }}
      {{- if .DictType}}
           <el-select {{if eq .FieldType "array"}}multiple {{end}}v-model="formData.{{ .FieldJson }}" placeholder="Please select {{.FieldDesc}}" style="width:100%" :clearable="{{.Clearable}}" >
              <el-option v-for="(item,key) in {{ .DictType }}Options" :key="key" :label="item.label" :value="item.value" />
           </el-select>
      {{- else }}
          <el-input v-model="formData.{{.FieldJson}}" :clearable="{{.Clearable}}"  placeholder="Please enter {{.FieldDesc}}" />
      {{- end }}
      {{- end }}
      {{- if eq .FieldType "richtext" }}
          <RichEdit v-model="formData.{{.FieldJson}}"/>
      {{- end }}
      {{- if eq .FieldType "int" }}
          <el-input v-model.number="formData.{{ .FieldJson }}" :clearable="{{.Clearable}}" placeholder="Please enter" />
      {{- end }}
      {{- if eq .FieldType "time.Time" }}
          <el-date-picker v-model="formData.{{ .FieldJson }}" type="date" placeholder="Select date" :clearable="{{.Clearable}}"></el-date-picker>
      {{- end }}
      {{- if eq .FieldType "float64" }}
          <el-input-number v-model="formData.{{ .FieldJson }}" :precision="2" :clearable="{{.Clearable}}"></el-input-number>
      {{- end }}
      {{- if eq .FieldType "enum" }}
        <el-select v-model="formData.{{ .FieldJson }}" placeholder="Please select" style="width:100%" :clearable="{{.Clearable}}">
          <el-option v-for="item in [{{ .DataTypeLong }}]" :key="item" :label="item" :value="item" />
        </el-select>
      {{- end }}
       {{- if eq .FieldType "picture" }}
          <SelectImage v-model="formData.{{ .FieldJson }}" file-type="image"/>
       {{- end }}
       {{- if eq .FieldType "video" }}
          <SelectImage v-model="formData.{{ .FieldJson }}" file-type="video"/>
       {{- end }}
       {{- if eq .FieldType "pictures" }}
           <SelectImage v-model="formData.{{ .FieldJson }}" multiple file-type="image"/>
       {{- end }}
       {{- if eq .FieldType "file" }}
          <SelectFile v-model="formData.{{ .FieldJson }}" />
       {{- end }}
       {{- if eq .FieldType "json" }}
          // This field is a JSON structure; control its rendering and binding on the frontend. Bind it to formData.{{.FieldJson}}; the backend stores/reads it according to its JSON type
          {{"{{"}} formData.{{.FieldJson}} {{"}}"}}
       {{- end }}
       {{- if eq .FieldType "array" }}
          <ArrayCtrl v-model="formData.{{ .FieldJson }}" editable/>
       {{- end }}
       {{- end }}
       </el-form-item>
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
} from '@/plugin/{{.Package}}/api/{{.PackageName}}'

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
            {{- if eq .FieldType "bool" }}
            {{.FieldJson}}: false,
            {{- end }}
            {{- if eq .FieldType "string" }}
            {{.FieldJson}}: '',
            {{- end }}
            {{- if eq .FieldType "richtext" }}
            {{.FieldJson}}: '',
            {{- end }}
            {{- if eq .FieldType "int" }}
            {{.FieldJson}}: {{- if or .DataSource }} undefined{{ else }} 0{{- end }},
            {{- end }}
            {{- if eq .FieldType "time.Time" }}
            {{.FieldJson}}: new Date(),
            {{- end }}
            {{- if eq .FieldType "float64" }}
            {{.FieldJson}}: 0,
            {{- end }}
            {{- if eq .FieldType "picture" }}
            {{.FieldJson}}: "",
            {{- end }}
            {{- if eq .FieldType "video" }}
            {{.FieldJson}}: "",
            {{- end }}
            {{- if eq .FieldType "pictures" }}
            {{.FieldJson}}: [],
            {{- end }}
            {{- if eq .FieldType "file" }}
            {{.FieldJson}}: [],
            {{- end }}
            {{- if eq .FieldType "json" }}
            {{.FieldJson}}: {},
            {{- end }}
            {{- if eq .FieldType "array" }}
            {{.FieldJson}}: [],
            {{- end }}
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
