<template>
  <div>
    <WarningBar
      title="This feature provides synchronous table export. For large datasets requiring async export, please contact us for customization."
      href="https://huuhoaitvn.feishu.cn/docx/KwjxdnvatozgwIxGV0rcpkZSn4d"
    />
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="Created at" prop="createdAt">
          <template #label>
            <span>
              Created at
              <el-tooltip
                content="Range is start date (inclusive) to end date (exclusive)."
              >
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.startCreatedAt"
            type="datetime"
            placeholder="Start"
            :disabled-date="
              (time) =>
                searchInfo.endCreatedAt
                  ? time.getTime() > searchInfo.endCreatedAt.getTime()
                  : false
            "
          />
          —
          <el-date-picker
            v-model="searchInfo.endCreatedAt"
            type="datetime"
            placeholder="End"
            :disabled-date="
              (time) =>
                searchInfo.startCreatedAt
                  ? time.getTime() < searchInfo.startCreatedAt.getTime()
                  : false
            "
          />
        </el-form-item>
        <el-form-item label="Template name" prop="name">
          <el-input v-model="searchInfo.name" placeholder="Search" />
        </el-form-item>
        <el-form-item label="Table name" prop="tableName">
          <el-input v-model="searchInfo.tableName" placeholder="Search" />
        </el-form-item>
        <el-form-item label="Template key" prop="templateID">
          <el-input v-model="searchInfo.templateID" placeholder="Search" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit"
            >Search</el-button
          >
          <el-button icon="refresh" @click="onReset">Reset</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog"
          >New</el-button
        >

        <el-button
          icon="delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
          >Delete</el-button
        >
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="Date" width="180">
          <template #default="scope">{{
            formatDate(scope.row.CreatedAt)
          }}</template>
        </el-table-column>
        <el-table-column align="left" label="Database" width="120">
          <template #default="scope">
            <span>{{ scope.row.dbName || 'GVA DB' }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Template key"
          prop="templateID"
          width="120"
        />
        <el-table-column
          align="left"
          label="Template name"
          prop="name"
          width="120"
        />
        <el-table-column
          align="left"
          label="Table name"
          prop="tableName"
          width="120"
        />
        <el-table-column
          align="left"
          label="Template info"
          prop="templateInfo"
          min-width="120"
          show-overflow-tooltip
        />
        <el-table-column align="left" label="Actions" min-width="280">
          <template #default="scope">
            <el-button
                type="primary"
                link
                icon="documentCopy"
                class="table-button"
                @click="copyFunc(scope.row)"
            >Copy</el-button
            >
            <el-button
              type="primary"
              link
              icon="edit-pen"
              class="table-button"
              @click="showCode(scope.row)"
              >Code & SQL preview</el-button
            >
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateSysExportTemplateFunc(scope.row)"
              >Edit</el-button
            >
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
              >Delete</el-button
            >
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
    <el-drawer
      v-model="dialogFormVisible"
      size="60%"
      :before-close="closeDialog"
      :title="type === 'create' ? 'Create' : 'Edit'"
      :show-close="false"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Create' : 'Edit' }}</span>
          <div>
            <el-button @click="closeDialog">Cancel</el-button>
            <el-button type="primary" @click="enterDialog">Confirm</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="elFormRef"
        :model="formData"
        label-position="right"
        :rules="rule"
        label-width="100px"
        v-loading="aiLoading"
        element-loading-text="Thinking..."
      >
        <el-form-item label="Business DB" prop="dbName">
          <template #label>
            <el-tooltip
              content="Note: configure multiple databases in db-list first. Restart the service after configuration. If you cannot select, set disabled:false in config.yaml to choose the target DB."
              placement="bottom"
              effect="light"
            >
              <div>
                Business DB <el-icon><QuestionFilled /></el-icon>
              </div>
            </el-tooltip>
          </template>
          <el-select
            v-model="formData.dbName"
            clearable
            @change="dbNameChange"
            placeholder="Select business DB"
          >
            <el-option
              v-for="item in dbList"
              :key="item.aliasName"
              :value="item.aliasName"
              :label="item.aliasName"
              :disabled="item.disable"
            >
              <div>
                <span>{{ item.aliasName }}</span>
                <span style="float: right; color: #8492a6; font-size: 13px">{{
                  item.dbName
                }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="Tables" prop="tables">
          <el-select
            multiple
            v-model="tables"
            clearable
            placeholder="Select when using AI"
          >
            <el-option
              v-for="item in tableOptions"
              :key="item.tableName"
              :label="item.tableName"
              :value="item.tableName"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="AI:" prop="ai">
          <div class="relative w-full">
            <el-input
              type="textarea"
              v-model="prompt"
              :clearable="true"
              :rows="5"
              placeholder="Describe the export you want. Select the business DB first; if not selected, the default GVA DB will be used."
            />
            <el-button
              class="absolute bottom-2 right-2"
              type="primary"
              @click="autoExport"
              ><el-icon><ai-gva /></el-icon>Generate</el-button
            >
          </div>
        </el-form-item>

        <el-form-item label="Table:" clearable prop="tableName">
          <div class="w-full flex gap-4">
            <el-select
              v-model="formData.tableName"
              class="flex-1"
              filterable
              placeholder="Select a table"
            >
              <el-option
                v-for="item in tableOptions"
                :key="item.tableName"
                :label="item.tableName"
                :value="item.tableName"
              />
            </el-select>
            <el-button
              :disabled="!formData.tableName"
              type="primary"
              @click="getColumnFunc(true)"
              ><el-icon><ai-gva /></el-icon>Auto-fill</el-button
            >
            <el-button
              :disabled="!formData.tableName"
              type="primary"
              @click="getColumnFunc(false)"
              >Auto-generate template</el-button
            >
          </div>
        </el-form-item>

        <el-form-item label="Template name:" prop="name">
          <el-input
            v-model="formData.name"
            :clearable="true"
            placeholder="Enter template name"
          />
        </el-form-item>

        <el-form-item label="Template key:" prop="templateID">
          <el-input
            v-model="formData.templateID"
            :clearable="true"
            placeholder="This key is used by frontend components to reference the template"
          />
        </el-form-item>

        <el-tabs v-model="activeName">
          <el-tab-pane label="Auto build" name="auto" class="pt-2">
            <el-form-item label="Join conditions:">
              <div
                v-for="(join, key) in formData.joinTemplate"
                :key="key"
                class="flex gap-4 w-full mb-2"
              >
                <el-select v-model="join.joins" placeholder="Select join type">
                  <el-option label="LEFT JOIN" value="LEFT JOIN" />
                  <el-option label="INNER JOIN" value="INNER JOIN" />
                  <el-option label="RIGHT JOIN" value="RIGHT JOIN" />
                </el-select>
                <el-input v-model="join.table" placeholder="Enter join table" />
                <el-input
                  v-model="join.on"
                  placeholder="Join condition, e.g. table1.a = table2.b"
                />
                <el-button
                  type="danger"
                  icon="delete"
                  @click="() => formData.joinTemplate.splice(key, 1)"
                  >Delete</el-button
                >
              </div>
              <div class="flex justify-end w-full">
                <el-button type="primary" icon="plus" @click="addJoin"
                  >Add condition</el-button
                >
              </div>
            </el-form-item>

            <el-form-item label="Default export limit:">
              <el-input-number
                v-model="formData.limit"
                :step="1"
                :step-strictly="true"
                :precision="0"
              />
            </el-form-item>
            <el-form-item label="Default sort:">
              <el-input v-model="formData.order" placeholder="e.g. id desc" />
            </el-form-item>
            <el-form-item label="Export conditions:">
              <div
                v-for="(condition, key) in formData.conditions"
                :key="key"
                class="flex gap-4 w-full mb-2"
              >
                <el-input
                  v-model="condition.from"
                  placeholder="JSON key from query conditions"
                />
                <el-input v-model="condition.column" placeholder="Table column" />
                <el-select
                  v-model="condition.operator"
                  placeholder="Select operator"
                >
                  <el-option
                    v-for="item in typeSearchOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
                <el-button
                  type="danger"
                  icon="delete"
                  @click="() => formData.conditions.splice(key, 1)"
                  >Delete</el-button
                >
              </div>
              <div class="flex justify-end w-full">
                <el-button type="primary" icon="plus" @click="addCondition"
                  >Add condition</el-button
                >
              </div>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="Custom SQL" name="sql"  class="pt-2">
            <el-form-item label="Export SQL:" prop="sql">
              <el-input
                v-model="formData.sql"
                type="textarea"
                :rows="10"
                placeholder="Enter export SQL. Supports GORM named params, e.g. SELECT * FROM sys_apis WHERE id = @id"
              />
            </el-form-item>
            <el-form-item label="Import SQL:" prop="importSql">
              <el-input
                v-model="formData.importSql"
                type="textarea"
                :rows="10"
                placeholder="Enter import SQL. Supports GORM named params, e.g. INSERT INTO sys_apis (path, description, api_group, method) VALUES (@path, @description, @api_group, @method). Param names must match keys in template info."
              />
            </el-form-item>
            <el-form-item label="Export conditions:">
              In this mode, conditions must be passed as condition = {key1:"value1", key2:"value2"}. Keys must match the @key placeholders in your SQL.
            </el-form-item>
          </el-tab-pane>
        </el-tabs>

        <el-form-item label="Template info:" prop="templateInfo">
          <el-input
            v-model="formData.templateInfo"
            type="textarea"
            :rows="12"
            :clearable="true"
            :placeholder="templatePlaceholder"
          />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Combined drawer: code template + SQL preview -->
    <el-drawer
      v-model="drawerVisible"
      size="70%"
      :title="'Template & Preview'"
      :show-close="true"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Template & Preview</span>
          <div>
            <el-button @click="drawerVisible = false">Close</el-button>
            <el-button v-if="activeTab === 'sql'" type="primary" @click="runPreview">Generate</el-button>
          </div>
        </div>
      </template>
      <el-tabs v-model="activeTab" type="border-card">
        <el-tab-pane label="Code" name="code">
          <v-ace-editor
            v-model:value="webCode"
            lang="vue"
            theme="github_dark"
            class="w-full h-96"
            :options="{ showPrintMargin: false, fontSize: 14 }"
          />
        </el-tab-pane>
        <el-tab-pane label="SQL preview" name="sql">
          <div class="flex flex-col gap-4">
            <div class="w-full">
              <el-form :model="previewForm" label-width="120px">
                <el-form-item label="Filter deleted">
                  <el-switch v-model="previewForm.filterDeleted" />
                </el-form-item>
                <el-form-item label="Default sort">
                  <el-input v-model="previewForm.order" placeholder="e.g. id desc" />
                </el-form-item>
                <el-form-item label="Limit">
                  <el-input-number v-model="previewForm.limit" :min="0" />
                </el-form-item>
                <el-form-item label="Offset">
                  <el-input-number v-model="previewForm.offset" :min="0" />
                </el-form-item>

                <el-divider content-position="left">Query conditions</el-divider>
                <div v-if="previewConditions.length === 0" class="text-gray">No conditions for this template</div>
                <template v-for="(cond, idx) in previewConditions" :key="idx">
                  <el-form-item :label="cond.column + ' ' + cond.operator">
                    <template v-if="cond.operator === 'BETWEEN'">
                      <div class="flex gap-2 w-full">
                        <el-input v-model="previewForm['start' + cond.from]" placeholder="Start: start{{cond.from}}" />
                        <el-input v-model="previewForm['end' + cond.from]" placeholder="End: end{{cond.from}}" />
                      </div>
                    </template>
                    <template v-else>
                      <el-input v-model="previewForm[cond.from]" :placeholder="'Variable: ' + cond.from" />
                    </template>
                  </el-form-item>
                </template>
              </el-form>
            </div>
            <div class="w-full">
              <v-ace-editor
                v-model:value="previewSQLCode"
                lang="sql"
                theme="github_dark"
                class="w-full h-96"
                :options="aceOptions"
              />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createSysExportTemplate,
    deleteSysExportTemplate,
    deleteSysExportTemplateByIds,
    updateSysExportTemplate,
    findSysExportTemplate,
    getSysExportTemplateList
  } from '@/api/exportTemplate.js'
  import { previewSQL } from '@/api/exportTemplate.js'

  // Import formatting utility (keep if needed)
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, reactive } from 'vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { getDB, getTable, getColumn, llmAuto } from '@/api/autoCode'
  import { getCode } from './code'
  import { VAceEditor } from 'vue3-ace-editor'
  import { useI18n } from 'vue-i18n'

  import 'ace-builds/src-noconflict/mode-vue'
  import 'ace-builds/src-noconflict/theme-github_dark'
  import 'ace-builds/src-noconflict/mode-sql'

  const { t } = useI18n()

  defineOptions({
    name: 'ExportTemplate'
  })

  const templatePlaceholder = t('admin.systemtools.export_template.template_placeholder')

  // Auto-generated dictionary (optional) and fields
  const formData = ref({
    name: '',
    tableName: '',
    dbName: '',
    templateID: '',
    templateInfo: '',
    limit: 0,
    order: '',
    conditions: [],
    joinTemplate: [],
    sql: '',
    importSql: ''
  })

  const activeName = ref('auto')

  const prompt = ref('')
  const tables = ref([])

  const typeSearchOptions = ref([
    {
      label: '=',
      value: '='
    },
    {
      label: '<>',
      value: '<>'
    },
    {
      label: '>',
      value: '>'
    },
    {
      label: '<',
      value: '<'
    },
    {
      label: 'LIKE',
      value: 'LIKE'
    },
    {
      label: 'BETWEEN',
      value: 'BETWEEN'
    },
    {
      label: 'NOT BETWEEN',
      value: 'NOT BETWEEN'
    },
    {
      label: 'IN',
      value: 'IN'
    },
    {
      label: 'NOT IN',
      value: 'NOT IN'
    },
  ])

  const addCondition = () => {
    formData.value.conditions.push({
      from: '',
      column: '',
      operator: ''
    })
  }

  const addJoin = () => {
    formData.value.joinTemplate.push({
      joins: 'LEFT JOIN',
      table: '',
      on: ''
    })
  }

  // Validation rules
  const rule = reactive({
    name: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: 'Whitespace only is not allowed',
        trigger: ['input', 'blur']
      }
    ],
    tableName: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: 'Whitespace only is not allowed',
        trigger: ['input', 'blur']
      }
    ],
    templateID: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: 'Whitespace only is not allowed',
        trigger: ['input', 'blur']
      }
    ],
    templateInfo: [
      {
        required: true,
        message: '',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: 'Whitespace only is not allowed',
        trigger: ['input', 'blur']
      }
    ]
  })

  const searchRule = reactive({
    createdAt: [
      {
        validator: (rule, value, callback) => {
          if (
            searchInfo.value.startCreatedAt &&
            !searchInfo.value.endCreatedAt
          ) {
    callback(new Error('Please provide an end date'))
          } else if (
            !searchInfo.value.startCreatedAt &&
            searchInfo.value.endCreatedAt
          ) {
    callback(new Error('Please provide a start date'))
          } else if (
            searchInfo.value.startCreatedAt &&
            searchInfo.value.endCreatedAt &&
            (searchInfo.value.startCreatedAt.getTime() ===
              searchInfo.value.endCreatedAt.getTime() ||
              searchInfo.value.startCreatedAt.getTime() >
                searchInfo.value.endCreatedAt.getTime())
          ) {
    callback(new Error('Start date must be earlier than end date'))
          } else {
            callback()
          }
        },
        trigger: 'change'
      }
    ]
  })

  const elFormRef = ref()
  const elSearchFormRef = ref()

  // =========== table ===========
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})

  const dbList = ref([])
  const tableOptions = ref([])
  const aiLoading = ref(false)

  const getTablesCloumn = async () => {
    const tablesMap = {}
    const promises = tables.value.map(async (item) => {
      const res = await getColumn({
        businessDB: formData.value.dbName,
        tableName: item
      })
      if (res.code === 0) {
        tablesMap[item] = res.data.columns
      }
    })
    await Promise.all(promises)
    return tablesMap
  }

  const autoExport = async () => {
    if (tables.value.length === 0) {
      ElMessage({
        type: 'error',
        message: 'Please select the tables to export'
      })
      return
    }
    aiLoading.value = true
    const tableMap = await getTablesCloumn()
    const aiRes = await llmAuto({
      prompt: prompt.value,
      tableMap: JSON.stringify(tableMap),
      mode: 'autoExportTemplate'
    })
    aiLoading.value = false
    if (aiRes.code === 0) {
      const aiData = JSON.parse(aiRes.data)
      formData.value.name = aiData.name
      formData.value.tableName = aiData.tableName
      formData.value.templateID = aiData.templateID
      formData.value.templateInfo = JSON.stringify(aiData.templateInfo, null, 2)
      formData.value.joinTemplate = aiData.joinTemplate
    }
  }

  const getDbFunc = async () => {
    const res = await getDB()
    if (res.code === 0) {
      dbList.value = res.data.dbList
    }
  }

  getDbFunc()

  const dbNameChange = () => {
    formData.value.tableName = ''
    formData.value.templateInfo = ''
    tables.value = []
    getTableFunc()
  }

  const getTableFunc = async () => {
    const res = await getTable({ businessDB: formData.value.dbName })
    if (res.code === 0) {
      tableOptions.value = res.data.tables
    }
    formData.value.tableName = ''
  }
  getTableFunc()
  const getColumnFunc = async (aiFLag) => {
    if (!formData.value.tableName) {
      ElMessage({
        type: 'error',
        message: 'Please select a business DB and table before continuing'
      })
      return
    }
    formData.value.templateInfo = ''
    aiLoading.value = true
    const res = await getColumn({
      businessDB: formData.value.dbName,
      tableName: formData.value.tableName
    })
    if (res.code === 0) {
      if (aiFLag) {
        const aiRes = await llmAuto({
          data: JSON.stringify(res.data.columns),
          mode: 'exportCompletion'
        })
        if (aiRes.code === 0) {
          const aiData = JSON.parse(aiRes.data)
          aiLoading.value = false
          formData.value.templateInfo = JSON.stringify(
            aiData.templateInfo,
            null,
            2
          )
          formData.value.name = aiData.name
          formData.value.templateID = aiData.templateID
          return
        }
    ElMessage.warning(t('admin.systemtools.export_template.ai_autofill_failed'))
      }

    // Convert data.columns to a JSON map: columnName -> columnComment
      const templateInfo = {}
      res.data.columns.forEach((item) => {
        templateInfo[item.columnName] = item.columnComment || item.columnName
      })
      formData.value.templateInfo = JSON.stringify(templateInfo, null, 2)
    }
    aiLoading.value = false
  }

    // Reset
  const onReset = () => {
    searchInfo.value = {}
    getTableData()
  }

  // Search
  const onSubmit = () => {
    elSearchFormRef.value?.validate(async (valid) => {
      if (!valid) return
      page.value = 1
      getTableData()
    })
  }

  // Pagination
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  // Page change
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // Query
  const getTableData = async () => {
    const table = await getSysExportTemplateList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  // ============== end table ===============

  // Optional dictionary init (keep if needed)
  const setOptions = async () => {}

  // Optional dictionary init (keep if needed)
  setOptions()

  // Selection
  const multipleSelection = ref([])
  // Selection change
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  // Delete row
  const deleteRow = (row) => {
  ElMessageBox.confirm(t('admin.common.confirms.delete_item'), t('admin.common.confirms.delete_title'), {
    confirmButtonText: t('admin.common.delete'),
    cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(() => {
      deleteSysExportTemplateFunc(row)
    })
  }

  // Bulk delete
  const onDelete = async () => {
  ElMessageBox.confirm(t('admin.systemtools.export_template.delete_selected_confirm'), t('admin.common.confirms.delete_title'), {
    confirmButtonText: t('admin.common.delete'),
    cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const ids = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: t('admin.common.validation.select_at_least_one')
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map((item) => {
          ids.push(item.ID)
        })
      const res = await deleteSysExportTemplateByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: t('admin.common.messages.deleted')
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }

  // Mode flag (create vs edit in the drawer)
  const type = ref('')

  // Copy
  const copyFunc = async (row) => {
    let copyData
    const res = await findSysExportTemplate({ ID: row.ID })
    if (res.code === 0) {
      copyData = JSON.parse(JSON.stringify(res.data.resysExportTemplate))
      if (!copyData.conditions) {
        copyData.conditions = []
      }
      if (!copyData.joinTemplate) {
        copyData.joinTemplate = []
      }
      if (!copyData.sql) {
        copyData.sql = ''
      }
      if (!copyData.importSql) {
        copyData.importSql = ''
      }
      delete copyData.ID
      delete copyData.CreatedAt
      delete copyData.UpdatedAt
      copyData.templateID = copyData.templateID + '_copy'
      copyData.name = copyData.name + '_copy'
      formData.value = copyData
      dialogFormVisible.value = true
    }
  }

  // Update row
  const updateSysExportTemplateFunc = async (row) => {
    const res = await findSysExportTemplate({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = res.data.resysExportTemplate
      if (!formData.value.conditions) {
        formData.value.conditions = []
      }
      if (!formData.value.joinTemplate) {
        formData.value.joinTemplate = []
      }
      if (!formData.value.sql) {
        formData.value.sql = ''
      }
      if (!formData.value.importSql) {
        formData.value.importSql = ''
      }
      if (formData.value.sql || formData.value.importSql) {
        activeName.value = 'sql'
      } else {
        activeName.value = 'auto'
      }
      dialogFormVisible.value = true
    }
  }

  // Delete row
  const deleteSysExportTemplateFunc = async (row) => {
    const res = await deleteSysExportTemplate({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Deleted'
      })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  }
  const drawerVisible = ref(false)
  const activeTab = ref('code')
  // Drawer visibility flag
  const dialogFormVisible = ref(false)

  const webCode = ref('')

  const showCode = (row) => {
    webCode.value = getCode(row.templateID)
    activeTab.value = 'code'
    drawerVisible.value = true
  }

  // Preview SQL
  const previewForm = ref({ filterDeleted: true, order: '', limit: 0, offset: 0 })
  const previewSQLCode = ref('')
  const previewTemplate = ref(null)
  const previewConditions = ref([])
  const aceOptions = { wrap: true, showPrintMargin: false, fontSize: 14 }

  const openPreview = async (row) => {
  // Load full template info to render condition inputs
    const res = await findSysExportTemplate({ ID: row.ID })
    if (res.code === 0) {
      previewTemplate.value = res.data.resysExportTemplate
      previewConditions.value = (previewTemplate.value.conditions || []).map((c) => ({
        from: c.from,
        column: c.column,
        operator: c.operator
      }))
  // Prefill default sort/limit
      previewForm.value.order = previewTemplate.value.order || ''
      previewForm.value.limit = previewTemplate.value.limit || 0
      previewForm.value.offset = 0
      previewSQLCode.value = ''
      activeTab.value = 'sql'
      drawerVisible.value = true
    }
  }

  const runPreview = async () => {
    if (!previewTemplate.value) return
  // Build params consistent with export component
    const paramsCopy = JSON.parse(JSON.stringify(previewForm.value))
  // Encode booleans/numbers using export component rules
    if (paramsCopy.filterDeleted) paramsCopy.filterDeleted = 'true'
    const entries = Object.entries(paramsCopy).filter(([key, v]) => {
      if (v === '' || v === null || v === undefined) return false
      if ((key === 'limit' || key === 'offset') && Number(v) === 0) return false
      return true
    })
    const params = entries
      .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
      .join('&')

    const res = await previewSQL({ templateID: previewTemplate.value.templateID, params })
    if (res.code === 0) {
      previewSQLCode.value = res.data.sql || ''
    }
  }

  // Open drawer
  const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
  }

  // Close drawer
  const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
      name: '',
      tableName: '',
      templateID: '',
      templateInfo: '',
      limit: 0,
      order: '',
      conditions: [],
      joinTemplate: [],
      sql: '',
      importSql: ''
    }
    activeName.value = 'auto'
  }
  // Confirm in drawer
  const enterDialog = async () => {
  // Validate templateInfo JSON; attempt to auto-fix if needed
    try {
      JSON.parse(formData.value.templateInfo)
    } catch (_) {
      ElMessage({
        type: 'error',
        message: 'Invalid template info format. Please check.'
      })
      return
    }

    const reqData = JSON.parse(JSON.stringify(formData.value))
    if (activeName.value === 'sql') {
      reqData.conditions = []
      reqData.joinTemplate = []
      reqData.limit = 0
      reqData.order = ''
    } else {
      reqData.sql = ''
      reqData.importSql = ''
    }

    for (let i = 0; i < reqData.conditions.length; i++) {
      if (
        !reqData.conditions[i].from ||
        !reqData.conditions[i].column ||
        !reqData.conditions[i].operator
      ) {
        ElMessage({
          type: 'error',
        message: 'Please fill all export conditions'
        })
        return
      }
      reqData.conditions[i].templateID = reqData.templateID
    }

    for (let i = 0; i < reqData.joinTemplate.length; i++) {
      if (!reqData.joinTemplate[i].joins || !reqData.joinTemplate[i].on) {
        ElMessage({
          type: 'error',
        message: 'Please complete the join configuration'
        })
        return
      }
      reqData.joinTemplate[i].templateID = reqData.templateID
    }

    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      let res
      switch (type.value) {
        case 'create':
          res = await createSysExportTemplate(reqData)
          break
        case 'update':
          res = await updateSysExportTemplate(reqData)
          break
        default:
          res = await createSysExportTemplate(reqData)
          break
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
        message: 'Saved'
        })
        closeDialog()
        getTableData()
      }
    })
  }
</script>

<style></style>
