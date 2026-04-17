{{- $global := . }}
{{- $templateID := printf "%s_%s" .Package .StructName }}
{{- if .IsAdd }}

// Add the following code to the search conditions
{{- range .Fields}}
    {{- if .FieldSearchType}}
{{ GenerateSearchFormItem .}}
    {{ end }}
{{ end }}


// Add the following columns to the table

{{- range .Fields}}
    {{- if .Table}}
       {{ GenerateTableColumn . }}
    {{- end }}
{{- end }}

// Add the following code to the new form
{{- range .Fields}}
   {{- if .Form}}
     {{ GenerateFormItem . }}
   {{- end }}
{{- end }}

// Add the following code to the detail drawer

{{- range .Fields}}
              {{- if .Desc }}
    {{ GenerateDescriptionItem . }}
              {{- end }}
            {{- end }}

// Add the following dictionary code
    {{- range $index, $element := .DictTypes}}
const {{ $element }}Options = ref([])
    {{- end }}

// Add the following calls inside the setOptions method

{{- range $index, $element := .DictTypes }}
    {{ $element }}Options.value = await getDictFunc('{{$element}}')
{{- end }}

// Add the following fields to the base formData struct (in the declaration and in the close-form reset)
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
const dataSource = ref({})
const getDataSourceFunc = async()=>{
  const res = await get{{.StructName}}DataSource()
  if (res.code === 0) {
    dataSource.value = res.data
  }
}
getDataSourceFunc()
{{- end }}

{{- else }}

{{- if not .OnlyTemplate}}
<template>
  <div>
  {{- if not .IsTree }}
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      {{- if .GvaModel }}
      <el-form-item label="Created at" prop="createdAtRange">
      <template #label>
        <span>
          Created at
          <el-tooltip content="The search range is inclusive of the start date and exclusive of the end date">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>

      <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="to"
            start-placeholder="Start time"
            end-placeholder="End time"
          />
       </el-form-item>
      {{ end -}}
           {{- range .Fields}}  {{- if .FieldSearchType}} {{- if not .FieldSearchHide }}
            {{ GenerateSearchFormItem .}}
            {{ end }}{{ end }}{{ end }}

        <template v-if="showAllQuery">
          <!-- Place any search fields that should be toggleable here -->
          {{- range .Fields}}  {{- if .FieldSearchType}} {{- if .FieldSearchHide }}
          {{ GenerateSearchFormItem .}}
          {{ end }}{{ end }}{{ end }}
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">Search</el-button>
          <el-button icon="refresh" @click="onReset">Reset</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">Expand</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>Collapse</el-button>
        </el-form-item>
      </el-form>
    </div>
  {{- end }}
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button {{ if $global.AutoCreateBtnAuth }}v-auth="btnAuth.add"{{ end }} type="primary" icon="plus" @click="openDialog()">Add</el-button>
            <el-button {{ if $global.AutoCreateBtnAuth }}v-auth="btnAuth.batchDelete"{{ end }} icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">Delete</el-button>
            {{ if .HasExcel -}}
            <ExportTemplate {{ if $global.AutoCreateBtnAuth }}v-auth="btnAuth.exportTemplate"{{ end }} template-id="{{$templateID}}" />
            <ExportExcel {{ if $global.AutoCreateBtnAuth }}v-auth="btnAuth.exportExcel"{{ end }} template-id="{{$templateID}}" filterDeleted/>
            <ImportExcel {{ if $global.AutoCreateBtnAuth }}v-auth="btnAuth.importExcel"{{ end }} template-id="{{$templateID}}" @on-success="getTableData" />
            {{- end }}
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="{{.PrimaryField.FieldJson}}"
        @selection-change="handleSelectionChange"
        {{- if .NeedSort}}
        @sort-change="sortChange"
        {{- end}}
        >
        <el-table-column type="selection" width="55" />
        {{ if .GvaModel }}
        <el-table-column sortable align="left" label="Date" prop="CreatedAt" {{ if .IsTree -}} min-{{- end -}}width="180">
            <template #default="scope">{{ "{{ formatDate(scope.row.CreatedAt) }}" }}</template>
        </el-table-column>
        {{ end }}
        {{- range .Fields}}
        {{- if .Table}}
            {{ GenerateTableColumn . }}
        {{- end }}
        {{- end }}
        <el-table-column align="left" label="Actions" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            {{- if .IsTree }}
            <el-button {{ if $global.AutoCreateBtnAuth }}v-auth="btnAuth.add"{{ end }} type="primary" link class="table-button" @click="openDialog(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>Add child</el-button>
            {{- end }}
            <el-button {{ if $global.AutoCreateBtnAuth }}v-auth="btnAuth.info"{{ end }} type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>View</el-button>
            <el-button {{ if $global.AutoCreateBtnAuth }}v-auth="btnAuth.edit"{{ end }} type="primary" link icon="edit" class="table-button" @click="update{{.StructName}}Func(scope.row)">Edit</el-button>
            <el-button {{ if .IsTree }}v-if="!scope.row.children?.length" {{ end }} {{if $global.AutoCreateBtnAuth }}v-auth="btnAuth.delete"{{ end }} type="primary" link icon="delete" @click="deleteRow(scope.row)">Delete</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{"{{"}}type==='create'?'Add':'Edit'{{"}}"}}</span>
                <div>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">Confirm</el-button>
                  <el-button @click="closeDialog">Cancel</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
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
          {{- if .Form}}
            {{ GenerateFormItem . }}
          {{- end }}
          {{- end }}
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="View">
            <el-descriptions :column="1" border>
            {{- if .IsTree }}
            <el-descriptions-item label="Parent node">
                <el-tree-select
                  v-model="detailForm.parentID"
                  :data="[rootNode,...tableData]"
                  check-strictly
                  disabled
                  :render-after-expand="false"
                  :props="defaultProps"
                  clearable
                  style="width: 240px"
                  placeholder="Root node"
                />
            </el-descriptions-item>
            {{- end }}
            {{- range .Fields}}
              {{- if .Desc }}
                    {{ GenerateDescriptionItem . }}
              {{- end }}
            {{- end }}
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  {{- if .HasDataSource }}
    get{{.StructName}}DataSource,
  {{- end }}
  create{{.StructName}},
  delete{{.StructName}},
  delete{{.StructName}}ByIds,
  update{{.StructName}},
  find{{.StructName}},
  get{{.StructName}}List
} from '@/api/{{.Package}}/{{.PackageName}}'

{{- if or .HasPic .HasFile}}
import { getUrl } from '@/utils/image'
{{- end }}
{{- if .HasPic }}
// Image selector component
import SelectImage from '@/components/selectImage/selectImage.vue'
{{- end }}

{{- if .HasRichText }}
// Rich text component
import RichEdit from '@/components/richtext/rich-edit.vue'
import RichView from '@/components/richtext/rich-view.vue'
{{- end }}

{{- if .HasFile }}
// File selector component
import SelectFile from '@/components/selectFile/selectFile.vue'
{{- end }}

{{- if .HasArray}}
// Array control component
import ArrayCtrl from '@/components/arrayCtrl/arrayCtrl.vue'
{{- end }}

// Full import of format utilities — remove any you don't need
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
{{- if .AutoCreateBtnAuth }}
// Import button-auth helper
import { useBtnAuth } from '@/utils/btnAuth'
{{- end }}
import { useAppStore } from "@/pinia"

{{if .HasExcel -}}
// Export component
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
// Import component
import ImportExcel from '@/components/exportExcel/importExcel.vue'
// Export-template component
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'
{{- end}}


defineOptions({
    name: '{{.StructName}}'
})

{{- if .AutoCreateBtnAuth }}
// Button auth instance
    const btnAuth = useBtnAuth()
{{- end }}

// Submit button loading state
const btnLoading = ref(false)
const appStore = useAppStore()

// Controls the show/hide state of the advanced search fields
const showAllQuery = ref(false)

// Auto-generated dictionaries (may be empty) and fields
    {{- range $index, $element := .DictTypes}}
const {{ $element }}Options = ref([])
    {{- end }}
const formData = ref({
        {{- if .IsTree }}
            parentID:undefined,
        {{- end }}
        {{- range .Fields}}
          {{- if .Form}}
            {{ GenerateDefaultFormValue . }}
          {{- end }}
        {{- end }}
        })

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



// Validation rules
const rule = reactive({
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
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== Table control section ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

{{- if .NeedSort}}
// Sort
const sortChange = ({ prop, order }) => {
  const sortMap = {
    CreatedAt:"created_at",
    ID:"id",
    {{- range .Fields}}
     {{- if .Table}}
      {{- if and .Sort}}
        {{- if not (eq .ColumnName "")}}
            {{.FieldJson}}: '{{.ColumnName}}',
        {{- end}}
      {{- end}}
     {{- end}}
    {{- end}}
  }

  let sort = sortMap[prop]
  if(!sort){
   sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }

  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
}
{{- end}}

{{- if not .IsTree }}
// Reset
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// Search
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    {{- range .Fields}}{{- if eq .FieldType "bool" }}
    if (searchInfo.value.{{.FieldJson}} === ""){
        searchInfo.value.{{.FieldJson}}=null
    }{{ end }}{{ end }}
    getTableData()
  })
}

// Page size
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// Change current page
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// Query
const getTableData = async() => {
  const table = await get{{.StructName}}List({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}
{{- else }}
// Tree selector config
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

// Query
const getTableData = async() => {
  const table = await get{{.StructName}}List()
  if (table.code === 0) {
    tableData.value = table.data || []
  }
}
{{- end }}

getTableData()

// ============== End of table control section ===============

// Fetch the required dictionaries (may be empty; remove what you don't need)
const setOptions = async () =>{
{{- range $index, $element := .DictTypes }}
    {{ $element }}Options.value = await getDictFunc('{{$element}}')
{{- end }}
}

// Fetch the required dictionaries (may be empty; remove what you don't need)
setOptions()


// Multi-selection data
const multipleSelection = ref([])
// Multi-select handler
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// Delete a row
const deleteRow = (row) => {
    ElMessageBox.confirm('Are you sure you want to delete?', 'Notice', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        type: 'warning'
    }).then(() => {
            delete{{.StructName}}Func(row)
        })
    }

// Bulk delete
const onDelete = async() => {
  ElMessageBox.confirm('Are you sure you want to delete?', 'Notice', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    type: 'warning'
  }).then(async() => {
      const {{.PrimaryField.FieldJson}}s = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: 'Please select rows to delete'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          {{.PrimaryField.FieldJson}}s.push(item.{{.PrimaryField.FieldJson}})
        })
      const res = await delete{{.StructName}}ByIds({ {{.PrimaryField.FieldJson}}s })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: 'Deleted successfully'
        })
        if (tableData.value.length === {{.PrimaryField.FieldJson}}s.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
      })
    }

// Tracks whether the dialog is in create or update mode
const type = ref('')

// Update a row
const update{{.StructName}}Func = async(row) => {
    const res = await find{{.StructName}}({ {{.PrimaryField.FieldJson}}: row.{{.PrimaryField.FieldJson}} })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// Delete a row
const delete{{.StructName}}Func = async (row) => {
    const res = await delete{{.StructName}}({ {{.PrimaryField.FieldJson}}: row.{{.PrimaryField.FieldJson}} })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: 'Deleted successfully'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// Dialog visibility state
const dialogFormVisible = ref(false)

// Open the dialog
const openDialog = ({{- if .IsTree -}}row{{- end -}}) => {
    type.value = 'create'
    {{- if .IsTree }}
    formData.value.parentID = row ? row.{{.PrimaryField.FieldJson}} : undefined
    {{- end }}
    dialogFormVisible.value = true
}

// Close the dialog
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
    {{- range .Fields}}
      {{- if .Form}}
        {{ GenerateDefaultFormValue . }}
      {{- end }}
    {{- end }}
        }
}
// Confirm the dialog
const enterDialog = async () => {
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
                closeDialog()
                getTableData()
              }
      })
}

const detailForm = ref({})

// Detail drawer visibility state
const detailShow = ref(false)


// Open the detail drawer
const openDetailShow = () => {
  detailShow.value = true
}


// Load details
const getDetails = async (row) => {
  // Open the drawer
  const res = await find{{.StructName}}({ {{.PrimaryField.FieldJson}}: row.{{.PrimaryField.FieldJson}} })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// Close the detail drawer
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}


</script>

<style>
{{if .HasFile }}
.file-list{
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.fileBtn{
  margin-bottom: 10px;
}

.fileBtn:last-child{
  margin-bottom: 0;
}
{{end}}
</style>
{{- else}}
<template>
<div>form</div>
</template>
<script setup>
defineOptions({
  name: '{{.StructName}}'
})
</script>
<style>
</style>
{{- end }}

{{- end }}
