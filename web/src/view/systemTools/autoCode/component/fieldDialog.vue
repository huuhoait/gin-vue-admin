<template>
  <div>
    <warning-bar
      :title="t('admin.system_tools.field_dialog.warning')"
    />
    <el-form
      ref="fieldDialogForm"
      :model="middleDate"
      label-width="120px"
      label-position="right"
      :rules="rules"
      class="grid grid-cols-2"
    >
      <el-form-item :label="t('admin.system_tools.field_dialog.field_name')" prop="fieldName">
        <el-input
          v-model="middleDate.fieldName"
          autocomplete="off"
          style="width: 80%"
        />
        <el-button style="width: 18%; margin-left: 2%" @click="autoFill">
          <span style="font-size: 12px">{{ t('admin.system_tools.field_dialog.auto_fill') }}</span>
        </el-button>
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.field_desc')" prop="fieldDesc">
        <el-input v-model="middleDate.fieldDesc" autocomplete="off" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.field_json')" prop="fieldJson">
        <el-input v-model="middleDate.fieldJson" autocomplete="off" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.db_column_name')" prop="columnName">
        <el-input v-model="middleDate.columnName" autocomplete="off" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.db_column_comment')" prop="comment">
        <el-input v-model="middleDate.comment" autocomplete="off" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.field_type')" prop="fieldType">
        <el-select
          v-model="middleDate.fieldType"
          style="width: 100%"
          :placeholder="t('admin.system_tools.field_dialog.select_field_type')"
          clearable
          @change="clearOther"
        >
          <el-option
            v-for="item in typeOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled"
          />
        </el-select>
      </el-form-item>
      <el-form-item
        :label="middleDate.fieldType === 'enum' ? t('admin.system_tools.field_dialog.enum_values') : t('admin.system_tools.field_dialog.type_length')"
        prop="dataTypeLong"
      >
        <el-input
          v-model="middleDate.dataTypeLong"
          :placeholder="
            middleDate.fieldType === 'enum'
              ? t('admin.system_tools.field_dialog.enum_placeholder')
              : t('admin.system_tools.field_dialog.db_type_length')
          "
        />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.search_operator')" prop="fieldSearchType">
        <el-select
          v-model="middleDate.fieldSearchType"
          :disabled="middleDate.fieldType === 'json'"
          style="width: 100%"
          :placeholder="t('admin.system_tools.field_dialog.select_search_operator')"
          clearable
        >
          <el-option
            v-for="item in typeSearchOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="canSelect(item.value)"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.dictionary')" prop="dictType">
        <el-select
          v-model="middleDate.dictType"
          style="width: 100%"
          :disabled="middleDate.fieldType !== 'string' && middleDate.fieldType !== 'array'"
          :placeholder="t('admin.system_tools.field_dialog.select_dictionary')"
          clearable
        >
          <el-option
            v-for="item in dictOptions"
            :key="item.type"
            :label="`${item.type}(${item.name})`"
            :value="item.type"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.default_value')">
        <el-input
          v-model="middleDate.defaultValue"
          :placeholder="t('admin.system_tools.field_dialog.default_value_placeholder')"
        />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.primary_key')">
        <el-checkbox v-model="middleDate.primaryKey" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.index_type')" prop="fieldIndexType">
        <el-select
          v-model="middleDate.fieldIndexType"
          :disabled="middleDate.fieldType === 'json'"
          style="width: 100%"
          :placeholder="t('admin.system_tools.field_dialog.select_index_type')"
          clearable
        >
          <el-option
            v-for="item in typeIndexOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="canSelect(item.value)"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.fe_form')">
        <el-switch v-model="middleDate.form" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.fe_table')">
        <el-switch v-model="middleDate.table" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.fe_details')">
        <el-switch v-model="middleDate.desc" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.import_export')">
        <el-switch v-model="middleDate.excel" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.sortable')">
        <el-switch v-model="middleDate.sort" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.required')">
        <el-switch v-model="middleDate.require" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.clearable')">
        <el-switch v-model="middleDate.clearable" />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.hide_search')">
        <el-switch
          v-model="middleDate.fieldSearchHide"
          :disabled="!middleDate.fieldSearchType"
        />
      </el-form-item>
      <el-form-item :label="t('admin.system_tools.field_dialog.validation_message')">
        <el-input v-model="middleDate.errorText" />
      </el-form-item>
    </el-form>
    <el-collapse v-model="activeNames">
      <el-collapse-item
        :title="t('admin.system_tools.field_dialog.data_source')"
        name="1"
      >
        <el-row :gutter="8">
          <el-col :span="4">
            <el-select
              v-model="middleDate.dataSource.dbName"
              :placeholder="t('admin.system_tools.field_dialog.database_default')"
              @change="dbNameChange"
              clearable
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
          </el-col>
          <el-col :span="4">
            <el-select
              v-model="middleDate.dataSource.association"
              :placeholder="t('admin.system_tools.field_dialog.association')"
              @change="associationChange"
            >
              <el-option :label="t('admin.system_tools.field_dialog.one_to_one')" :value="1" />
              <el-option :label="t('admin.system_tools.field_dialog.one_to_many')" :value="2" />
            </el-select>
          </el-col>
          <el-col :span="5">
            <el-select
              v-model="middleDate.dataSource.table"
              :placeholder="t('admin.system_tools.field_dialog.source_table')"
              filterable
              allow-create
              clearable
              @focus="getDBTableList"
              @change="selectDB"
              @clear="clearAccress"
            >
              <el-option
                v-for="item in dbTableList"
                :key="item.tableName"
                :label="item.tableName"
                :value="item.tableName"
              />
            </el-select>
          </el-col>
          <el-col :span="5">
            <el-select
              v-model="middleDate.dataSource.value"
              :placeholder="t('admin.system_tools.field_dialog.value_field')"
            >
              <template #label="{ value }">
                <span>{{ t('admin.system_tools.field_dialog.value_label') }}</span>
                <span style="font-weight: bold">{{ value }}</span>
              </template>
              <el-option
                v-for="item in dbColumnList"
                :key="item.columnName"
                :value="item.columnName"
              >
                <span style="float: left">
                  <el-tag :type="item.isPrimary ? 'primary' : 'info'">
                    {{ item.isPrimary ? t('admin.system_tools.field_dialog.pk_short') : t('admin.system_tools.field_dialog.non_pk') }}
                  </el-tag>
                  {{ item.columnName }}</span
                >
                <span
                  style="
                    float: right;
                    margin-left: 5px;
                    color: var(--el-text-color-secondary);
                    font-size: 13px;
                  "
                >
                  {{ t('admin.system_tools.field_dialog.type_col') }}{{ item.type }}
                  <block v-if="item.comment != ''"
                    >, {{ t('admin.system_tools.field_dialog.comment_col') }}{{ item.comment }}</block
                  >
                </span>
              </el-option>
            </el-select>
          </el-col>
          <el-col :span="5">
            <el-select
              v-model="middleDate.dataSource.label"
              :placeholder="t('admin.system_tools.field_dialog.label_field')"
            >
              <template #label="{ value }">
                <span>{{ t('admin.system_tools.field_dialog.label_label') }}</span>
                <span style="font-weight: bold">{{ value }}</span>
              </template>
              <el-option
                v-for="item in dbColumnList"
                :key="item.columnName"
                :value="item.columnName"
              >
                <span style="float: left">
                  <el-tag :type="item.isPrimary ? 'primary' : 'info'">
                    {{ item.isPrimary ? t('admin.system_tools.field_dialog.pk_short') : t('admin.system_tools.field_dialog.non_pk') }}
                  </el-tag>
                  {{ item.columnName }}</span
                >
                <span
                  style="
                    float: right;
                    margin-left: 5px;
                    color: var(--el-text-color-secondary);
                    font-size: 13px;
                  "
                >
                  {{ t('admin.system_tools.field_dialog.type_col') }}{{ item.type }}
                  <span v-if="item.comment != ''"
                    >, {{ t('admin.system_tools.field_dialog.comment_col') }}{{ item.comment }}</span
                  >
                </span>
              </el-option>
            </el-select>
            <!-- <el-input v-model="middleDate.dataSource.label" placeholder="Label field" /> -->
          </el-col>
        </el-row>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script setup>
  import { toLowerCase, toSQLLine } from '@/utils/stringFun'
  import { getSysDictionaryList } from '@/api/sysDictionary'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { ref, onMounted } from 'vue'
  import { ElMessageBox } from 'element-plus'
  import { getColumn, getDB, getTable } from '@/api/autoCode'
  import { useI18n } from 'vue-i18n'

  const { t } = useI18n()

  defineOptions({
    name: 'FieldDialog'
  })

  const props = defineProps({
    dialogMiddle: {
      type: Object,
      default: function () {
        return {}
      }
    },
    typeOptions: {
      type: Array,
      default: function () {
        return []
      }
    },
    typeSearchOptions: {
      type: Array,
      default: function () {
        return []
      }
    },
    typeIndexOptions: {
      type: Array,
      default: function () {
        return []
      }
    }
  })

  const activeNames = ref([])

  const middleDate = ref({})
  const dictOptions = ref([])

  const dbList = ref([])

  const getDbFunc = async () => {
    const res = await getDB()
    if (res.code === 0) {
      dbList.value = res.data.dbList
    }
  }

  const validateDataTypeLong = (rule, value, callback) => {
    const regex = /^('([^']*)'(?:,'([^']+)'*)*)$/
    if (middleDate.value.fieldType == 'enum' && !regex.test(value)) {
      callback(new Error(t('admin.system_tools.field_dialog.invalid_enum')))
    } else {
      callback()
    }
  }

  const rules = ref({
    fieldName: [
      { required: true, message: t('admin.system_tools.field_dialog.field_name_required'), trigger: 'blur' }
    ],
    fieldDesc: [
      { required: true, message: t('admin.system_tools.field_dialog.field_desc_required'), trigger: 'blur' }
    ],
    fieldJson: [
      { required: true, message: t('admin.system_tools.field_dialog.field_json_required'), trigger: 'blur' }
    ],
    columnName: [
      { required: true, message: t('admin.system_tools.field_dialog.db_column_required'), trigger: 'blur' }
    ],
    fieldType: [{ required: true, message: t('admin.system_tools.field_dialog.field_type_required'), trigger: 'blur' }],
    dataTypeLong: [{ validator: validateDataTypeLong, trigger: 'blur' }]
  })

  const init = async () => {
    middleDate.value = props.dialogMiddle
    const dictRes = await getSysDictionaryList({
      page: 1,
      pageSize: 999999
    })

    dictOptions.value = dictRes.data
  }
  init()

  const autoFill = () => {
    middleDate.value.fieldJson = toLowerCase(middleDate.value.fieldName)
    middleDate.value.columnName = toSQLLine(middleDate.value.fieldJson)
  }

  const canSelect = (item) => {
    const fieldType = middleDate.value.fieldType;

    if (fieldType === 'richtext') {
      return item !== 'LIKE';
    }

    if (fieldType !== 'string' && item === 'LIKE') {
      return true;
    }

    const nonNumericTypes = ['int', 'time.Time', 'float64'];
    if (!nonNumericTypes.includes(fieldType) && ['BETWEEN', 'NOT BETWEEN'].includes(item)) {
      return true;
    }

    return false;
  }

  const clearOther = () => {
    middleDate.value.fieldSearchType = ''
    middleDate.value.dictType = ''
  }

  const associationChange = (val) => {
    if (val === 2) {
      ElMessageBox.confirm(
        t('admin.system_tools.field_dialog.one_to_many_confirm'),
        t('admin.system_tools.field_dialog.confirm_title'),
        {
          confirmButtonText: t('admin.system_tools.field_dialog.continue'),
          cancelButtonText: t('admin.common.cancel'),
          type: 'warning'
        }
      )
        .then(() => {
          middleDate.value.fieldType = 'array'
        })
        .catch(() => {
          middleDate.value.dataSource.association = 1
        })
    }
  }

  const clearAccress = () => {
    middleDate.value.dataSource.value = ''
    middleDate.value.dataSource.label = ''
  }

  const clearDataSourceTable = () => {
    middleDate.value.dataSource.table = ''
  }

  const dbNameChange = () => {
    getDBTableList()
    clearDataSourceTable()
    clearAccress()
  }

  const dbTableList = ref([])

  const getDBTableList = async () => {
    const res = await getTable({
      businessDB: middleDate.value.dataSource.dbName
    })
    if (res.code === 0) {
      let list = res.data.tables // ensure we read tables array
      dbTableList.value = list.map((item) => ({
        tableName: item.tableName,
        value: item.tableName // assuming value is tableName; adjust if different
      }))
    }
    clearAccress()
  }

  const dbColumnList = ref([])
  const selectDB = async (val, isInit) => {
    middleDate.value.dataSource.hasDeletedAt = false
    middleDate.value.dataSource.table = val
    const res = await getColumn({
      businessDB: middleDate.value.dataSource.dbName,
      tableName: val
    })

    if (res.code === 0) {
      let list = res.data.columns // ensure we read columns array
      dbColumnList.value = list.map((item) => {
        if (item.columnName === 'deleted_at') {
          middleDate.value.dataSource.hasDeletedAt = true
        }
        return {
          columnName: item.columnName,
          value: item.columnName,
          type: item.dataType,
          isPrimary: item.primaryKey,
          comment: item.columnComment
        }
      })
      if (dbColumnList.value.length > 0 && !isInit) {
        middleDate.value.dataSource.label = dbColumnList.value[0].columnName
        middleDate.value.dataSource.value = dbColumnList.value[0].columnName
      }
    }
  }

  const fieldDialogForm = ref(null)
  defineExpose({ fieldDialogForm })

  onMounted(() => {
    getDbFunc()
    if (middleDate.value.dataSource.table) {
      selectDB(middleDate.value.dataSource.table, true)
    }
  })
</script>
