<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="goAutoCode(null)">
          {{ t('admin.systemtools.autocode_admin.new') }}
        </el-button>
      </div>
      <el-table :data="tableData">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" :label="t('admin.systemtools.autocode_admin.id')" width="60" prop="ID" />
        <el-table-column align="left" :label="t('admin.systemtools.autocode_admin.date')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          :label="t('admin.systemtools.autocode_admin.struct_name')"
          min-width="150"
          prop="structName"
        />
        <el-table-column
          align="left"
          :label="t('admin.systemtools.autocode_admin.struct_description')"
          min-width="150"
          prop="description"
        />
        <el-table-column
          align="left"
          :label="t('admin.systemtools.autocode_admin.table_name')"
          min-width="150"
          prop="tableName"
        />
        <el-table-column
          align="left"
          :label="t('admin.systemtools.autocode_admin.rollback_flag')"
          min-width="150"
          prop="flag"
        >
          <template #default="scope">
            <el-tag v-if="scope.row.flag" type="danger" effect="dark">
              {{ t('admin.systemtools.autocode_admin.rolled_back') }}
            </el-tag>
            <el-tag v-else type="success" effect="dark"> {{ t('admin.systemtools.autocode_admin.not_rolled_back') }} </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.systemtools.autocode_admin.actions')" min-width="240">
          <template #default="scope">
            <div>
              <el-button
                type="primary"
                link
                :disabled="scope.row.flag === 1"
                @click="addFuncBtn(scope.row)"
              >
                {{ t('admin.systemtools.autocode_admin.add_method') }}
              </el-button>
              <el-button type="primary" link @click="goAutoCode(scope.row, 1)">
                {{ t('admin.systemtools.autocode_admin.add_field') }}
              </el-button>
              <el-button
                type="primary"
                link
                :disabled="scope.row.flag === 1"
                @click="openDialog(scope.row)"
              >
                {{ t('admin.systemtools.autocode_admin.rollback') }}
              </el-button>
              <el-button type="primary" link @click="goAutoCode(scope.row)">
                {{ t('admin.systemtools.autocode_admin.reuse') }}
              </el-button>
              <el-button type="primary" link @click="deleteRow(scope.row)">
                {{ t('admin.systemtools.autocode_admin.delete') }}
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-dialog
      v-model="dialogFormVisible"
      :title="dialogFormTitle"
      :before-close="closeDialog"
      width="600px"
    >
      <el-form :inline="true" :model="formData" label-width="80px">
        <el-form-item :label="t('admin.systemtools.autocode_admin.options')">
          <el-checkbox v-model="formData.deleteApi" :label="t('admin.systemtools.autocode_admin.delete_apis')" />
          <el-checkbox v-model="formData.deleteMenu" :label="t('admin.systemtools.autocode_admin.delete_menus')" />
          <el-checkbox
            v-model="formData.deleteTable"
            :label="t('admin.systemtools.autocode_admin.delete_table')"
            @change="deleteTableCheck"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog"> {{ t('admin.systemtools.autocode_admin.cancel') }} </el-button>
          <el-popconfirm
            :title="t('admin.systemtools.autocode_admin.rollback_confirm')"
            @confirm="enterDialog"
          >
            <template #reference>
              <el-button type="primary"> {{ t('admin.common.confirm') }} </el-button>
            </template>
          </el-popconfirm>
        </div>
      </template>
    </el-dialog>

    <el-drawer
      v-model="funcFlag"
      size="60%"
      :show-close="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.systemtools.autocode_admin.actions') }}</span>
          <div>
            <el-button type="primary" @click="runFunc" :loading="aiLoading">
              {{ t('admin.systemtools.autocode_admin.generate') }}
            </el-button>
            <el-button type="primary" @click="closeFunc" :loading="aiLoading">
              {{ t('admin.systemtools.autocode_admin.cancel') }}
            </el-button>
          </div>
        </div>
      </template>
      <div class="">
        <el-form
          v-loading="aiLoading"
          label-position="top"
          :element-loading-text="t('admin.systemtools.autocode_admin.thinking_wait')"
          :model="autoFunc"
          label-width="80px"
        >
          <el-row :gutter="12">
            <el-col :span="8">
              <el-form-item :label="t('admin.systemtools.autocode_admin.package')">
                <el-input
                    v-model="autoFunc.package"
                    :placeholder="t('admin.systemtools.autocode_admin.enter_package')"
                    disabled
                />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item :label="t('admin.systemtools.autocode_admin.struct_name_field')">
                <el-input
                    v-model="autoFunc.structName"
                    :placeholder="t('admin.systemtools.autocode_admin.enter_struct_name')"
                    disabled
                />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item :label="t('admin.systemtools.autocode_admin.frontend_filename')">
                <el-input
                    v-model="autoFunc.packageName"
                    :placeholder="t('admin.systemtools.autocode_admin.enter_filename')"
                    disabled
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="12">
            <el-col :span="8">
              <el-form-item :label="t('admin.systemtools.autocode_admin.backend_filename')">
                <el-input
                    v-model="autoFunc.humpPackageName"
                    :placeholder="t('admin.systemtools.autocode_admin.enter_filename')"
                    disabled
                />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item :label="t('admin.systemtools.autocode_admin.description')">
                <el-input
                    v-model="autoFunc.description"
                    :placeholder="t('admin.systemtools.autocode_admin.enter_description')"
                    disabled
                />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item :label="t('admin.systemtools.autocode_admin.abbreviation')">
                <el-input
                    v-model="autoFunc.abbreviation"
                    :placeholder="t('admin.systemtools.autocode_admin.enter_abbreviation')"
                    disabled
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item :label="t('admin.systemtools.autocode_admin.use_ai')">
            <el-switch v-model="autoFunc.isAi" />
            <span class="text-sm text-red-600 p-2"
              >{{ t('admin.systemtools.autocode_admin.ai_warning') }}</span
            >
          </el-form-item>
          <template v-if="autoFunc.isAi">
            <el-form-item :label="t('admin.systemtools.autocode_admin.ai_prompt')">
              <div class="relative w-full">
                <el-input
                  type="textarea"
                  :placeholder="t('admin.systemtools.autocode_admin.ai_prompt_placeholder')"
                  v-model="autoFunc.prompt"
                  :rows="5"
                  @input="autoFunc.router = autoFunc.router.replace(/\//g, '')"
                />
                <el-button
                  @click="aiAddFunc"
                  type="primary"
                  class="absolute right-2 bottom-2"
                  ><ai-gva />{{ t('admin.systemtools.autocode_admin.generate') }}</el-button
                >
              </div>
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.autocode_admin.api_method')">
              <v-ace-editor
                v-model:value="autoFunc.apiFunc"
                lang="golang"
                theme="github_dark"
                class="h-80 w-full"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.autocode_admin.server_method')">
              <v-ace-editor
                v-model:value="autoFunc.serverFunc"
                lang="golang"
                theme="github_dark"
                class="h-80 w-full"
              />
            </el-form-item>
            <el-form-item :label="t('admin.systemtools.autocode_admin.frontend_js_api')">
              <v-ace-editor
                v-model:value="autoFunc.jsFunc"
                lang="javascript"
                theme="github_dark"
                class="h-80 w-full"
              />
            </el-form-item>
          </template>

          <el-form-item :label="t('admin.systemtools.autocode_admin.method_summary')">
            <div class="flex w-full gap-2">
              <el-input
                class="flex-1"
                v-model="autoFunc.funcDesc"
                :placeholder="t('admin.systemtools.autocode_admin.enter_method_summary')"
              />
              <el-button type="primary" @click="autoComplete"
                ><ai-gva />{{ t('admin.systemtools.autocode_admin.auto_fill') }}</el-button
              >
            </div>
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.autocode_admin.method_name')">
            <el-input
              @blur="autoFunc.funcName = toUpperCase(autoFunc.funcName)"
              v-model="autoFunc.funcName"
              :placeholder="t('admin.systemtools.autocode_admin.enter_method_name')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.autocode_admin.http_method')">
            <el-select v-model="autoFunc.method" :placeholder="t('admin.systemtools.autocode_admin.select_method')">
              <el-option
                v-for="item in ['GET', 'POST', 'PUT', 'DELETE']"
                :key="item"
                :label="item"
                :value="item"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.autocode_admin.require_auth')">
            <el-switch
              v-model="autoFunc.isAuth"
              :active-text="t('admin.systemtools.autocode_admin.yes')"
              :inactive-text="t('admin.systemtools.autocode_admin.no')"
            />
          </el-form-item>
          <el-form-item :label="t('admin.systemtools.autocode_admin.route_path')">
            <el-input
              v-model="autoFunc.router"
              :placeholder="t('admin.systemtools.autocode_admin.route_placeholder')"
              @input="autoFunc.router = autoFunc.router.replace(/\//g, '')"
            />
            <div>
              {{ t('admin.systemtools.autocode_admin.api_path') }}: [{{ autoFunc.method }}] /{{ autoFunc.abbreviation }}/{{
                autoFunc.router
              }}
            </div>
          </el-form-item>
        </el-form>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    getSysHistory,
    rollback,
    delSysHistory,
    addFunc,
    llmAuto
  } from '@/api/autoCode.js'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref } from 'vue'
  import { formatDate } from '@/utils/format'
  import { toUpperCase } from '@/utils/stringFun'
  import { useI18n } from 'vue-i18n'

  import { VAceEditor } from 'vue3-ace-editor'
  import 'ace-builds/src-noconflict/mode-javascript'
  import 'ace-builds/src-noconflict/mode-golang'
  import 'ace-builds/src-noconflict/theme-github_dark'

  const { t } = useI18n()

  defineOptions({
    name: 'AutoCodeAdmin'
  })

  const aiLoading = ref(false)

  const formData = ref({
    id: undefined,
    deleteApi: true,
    deleteMenu: true,
    deleteTable: false
  })

  const router = useRouter()
  const dialogFormVisible = ref(false)
  const dialogFormTitle = ref('')

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])

  const activeInfo = ref('')

  const autoFunc = ref({
    package: '',
    funcName: '',
    structName: '',
    packageName: '',
    description: '',
    abbreviation: '',
    humpPackageName: '',
    businessDB: '',
    method: '',
    funcDesc: '',
    isAuth: false,
    isAi: false,
    apiFunc: '',
    serverFunc: '',
    jsFunc: ''
  })

  const addFuncBtn = (row) => {
    const req = JSON.parse(row.request)
    activeInfo.value = row.request
    autoFunc.value.package = req.package
    autoFunc.value.structName = req.structName
    autoFunc.value.packageName = req.packageName
    autoFunc.value.description = req.description
    autoFunc.value.abbreviation = req.abbreviation
    autoFunc.value.humpPackageName = req.humpPackageName
    autoFunc.value.businessDB = req.businessDB
    autoFunc.value.method = ''
    autoFunc.value.funcName = ''
    autoFunc.value.router = ''
    autoFunc.value.funcDesc = ''
    autoFunc.value.isAuth = false
    autoFunc.value.isAi = false
    autoFunc.value.apiFunc = ''
    autoFunc.value.serverFunc = ''
    autoFunc.value.jsFunc = ''
    funcFlag.value = true
  }

  const funcFlag = ref(false)

  const closeFunc = () => {
    funcFlag.value = false
  }

  const runFunc = async () => {
    // Auto-uppercase first letter
    autoFunc.value.funcName = toUpperCase(autoFunc.value.funcName)

    if (!autoFunc.value.funcName) {
      ElMessage.error(t('admin.systemtools.autocode_admin.enter_method_name_err'))
      return
    }
    if (!autoFunc.value.method) {
      ElMessage.error(t('admin.systemtools.autocode_admin.select_method_err'))
      return
    }
    if (!autoFunc.value.router) {
      ElMessage.error(t('admin.systemtools.autocode_admin.enter_route_err'))
      return
    }
    if (!autoFunc.value.funcDesc) {
      ElMessage.error(t('admin.systemtools.autocode_admin.enter_method_summary_err'))
      return
    }

    if (autoFunc.value.isAi) {
      if (
        !autoFunc.value.apiFunc ||
        !autoFunc.value.serverFunc ||
        !autoFunc.value.jsFunc
      ) {
        ElMessage.error(t('admin.systemtools.autocode_admin.generate_ai_first'))
        return
      }
    }

    const res = await addFunc(autoFunc.value)
    if (res.code === 0) {
      ElMessage.success(t('admin.systemtools.autocode_admin.method_added'))
      closeFunc()
    }
  }

  // Pagination
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // Query
  const getTableData = async () => {
    const table = await getSysHistory({
      page: page.value,
      pageSize: pageSize.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  const deleteRow = async (row) => {
    ElMessageBox.confirm(t('admin.systemtools.autocode_admin.delete_history_confirm'), t('admin.common.confirm'), {
      confirmButtonText: t('admin.systemtools.autocode_admin.delete'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const res = await delSysHistory({ id: Number(row.ID) })
      if (res.code === 0) {
        ElMessage.success(t('admin.common.messages.deleted'))
        getTableData()
      }
    })
  }

  // Open dialog
  const openDialog = (row) => {
    dialogFormTitle.value = t('admin.systemtools.autocode_admin.rollback_title', { name: row.structName })
    formData.value.id = row.ID
    dialogFormVisible.value = true
  }

  // Close dialog
  const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
      id: undefined,
      deleteApi: true,
      deleteMenu: true,
      deleteTable: false
    }
  }

  // Confirm delete table
  const deleteTableCheck = (flag) => {
    if (flag) {
      ElMessageBox.confirm(
        t('admin.systemtools.autocode_admin.drop_table_warn'),
        t('admin.common.confirms.warning'),
        {
          closeOnClickModal: false,
          distinguishCancelAndClose: true,
          confirmButtonText: t('admin.common.confirm'),
          cancelButtonText: t('admin.common.cancel'),
          type: 'warning'
        }
      )
        .then(() => {
          ElMessageBox.confirm(
            t('admin.systemtools.autocode_admin.drop_table_warn_again'),
            t('admin.systemtools.autocode_admin.drop_table_title'),
            {
              closeOnClickModal: false,
              distinguishCancelAndClose: true,
              confirmButtonText: t('admin.common.confirm'),
              cancelButtonText: t('admin.common.cancel'),
              type: 'warning'
            }
          ).catch(() => {
            formData.value.deleteTable = false
          })
        })
        .catch(() => {
          formData.value.deleteTable = false
        })
    }
  }

  const enterDialog = async () => {
    const res = await rollback(formData.value)
    if (res.code === 0) {
      ElMessage.success(t('admin.common.messages.rolled_back'))
      getTableData()
    }
  }

  const goAutoCode = (row, isAdd) => {
    if (row) {
      router.push({
        name: 'autoCodeEdit',
        params: {
          id: row.ID
        },
        query: {
          isAdd: isAdd
        }
      })
    } else {
      router.push({ name: 'autoCode' })
    }
  }

  const aiAddFunc = async () => {
    aiLoading.value = true
    autoFunc.value.apiFunc = ''
    autoFunc.value.serverFunc = ''
    autoFunc.value.jsFunc = ''

    if (!autoFunc.value.prompt) {
      ElMessage.error(t('admin.systemtools.autocode_admin.enter_prompt'))
      return
    }

    const res = await addFunc({ ...autoFunc.value, isPreview: true })
    if (res.code !== 0) {
      aiLoading.value = false
      ElMessage.error(res.msg)
      return
    }

    const aiRes = await llmAuto({
      structInfo: activeInfo.value,
      template: JSON.stringify(res.data),
      prompt: autoFunc.value.prompt,
      mode: 'addFunc'
    })
    aiLoading.value = false
    if (aiRes.code === 0) {
      try {
        const aiData = JSON.parse(aiRes.data)
        autoFunc.value.apiFunc = aiData.api
        autoFunc.value.serverFunc = aiData.server
        autoFunc.value.jsFunc = aiData.js
        autoFunc.value.method = aiData.method
        autoFunc.value.funcName = aiData.funcName
        const routerArr = aiData.router.split('/')
        autoFunc.value.router = routerArr[routerArr.length - 1]
        autoFunc.value.funcDesc = autoFunc.value.prompt
      } catch (_) {
        ElMessage.error(t('admin.systemtools.autocode_admin.ai_busy_retry'))
      }
    }
  }

  const autoComplete = async () => {
    aiLoading.value = true
    const aiRes = await llmAuto({
      prompt: autoFunc.value.funcDesc,
      mode: 'autoCompleteFunc'
    })
    aiLoading.value = false
    if (aiRes.code === 0) {
      try {
        const aiData = JSON.parse(aiRes.data)
        autoFunc.value.method = aiData.method
        autoFunc.value.funcName = aiData.funcName
        autoFunc.value.router = aiData.router
        autoFunc.value.prompt = autoFunc.value.funcDesc
      } catch (_) {
        ElMessage.error(t('admin.systemtools.autocode_admin.ai_failed_retry'))
      }
    }
  }
</script>
