<template>
  <div>
    <warning-bar
      title="id, created_at, updated_at, deleted_at are auto-generated. Do not create them manually. For search conditions, LIKE only supports strings."
    />
    <el-form
      ref="fieldDialogForm"
      :model="middleDate"
      label-width="120px"
      label-position="right"
      :rules="rules"
      class="grid grid-cols-2"
    >
      <el-form-item label="Field name" prop="fieldName">
        <el-input
          v-model="middleDate.fieldName"
          autocomplete="off"
          style="width: 80%"
        />
        <el-button style="width: 18%; margin-left: 2%" @click="autoFill">
          <span style="font-size: 12px">Auto fill</span>
        </el-button>
      </el-form-item>
      <el-form-item label="Field label (CN)" prop="fieldDesc">
        <el-input v-model="middleDate.fieldDesc" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Field JSON" prop="fieldJson">
        <el-input v-model="middleDate.fieldJson" autocomplete="off" />
      </el-form-item>
      <el-form-item label="DB column name" prop="columnName">
        <el-input v-model="middleDate.columnName" autocomplete="off" />
      </el-form-item>
      <el-form-item label="DB column comment" prop="comment">
        <el-input v-model="middleDate.comment" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Field type" prop="fieldType">
        <el-select
          v-model="middleDate.fieldType"
          style="width: 100%"
          placeholder="Select field type"
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
        :label="middleDate.fieldType === 'enum' ? 'Enum values' : 'Type length'"
        prop="dataTypeLong"
      >
        <el-input
          v-model="middleDate.dataTypeLong"
          :placeholder="
            middleDate.fieldType === 'enum'
              ? `e.g. 'beijing','tianjin'`
              : 'DB type length'
          "
        />
      </el-form-item>
      <el-form-item label="Search operator" prop="fieldSearchType">
        <el-select
          v-model="middleDate.fieldSearchType"
          :disabled="middleDate.fieldType === 'json'"
          style="width: 100%"
          placeholder="Select search operator"
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
      <el-form-item label="Dictionary" prop="dictType">
        <el-select
          v-model="middleDate.dictType"
          style="width: 100%"
          :disabled="middleDate.fieldType !== 'string' && middleDate.fieldType !== 'array'"
          placeholder="Select dictionary"
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
      <el-form-item label="Default value">
        <el-input
          v-model="middleDate.defaultValue"
          placeholder="Enter default value"
        />
      </el-form-item>
      <el-form-item label="Primary key">
        <el-checkbox v-model="middleDate.primaryKey" />
      </el-form-item>
      <el-form-item label="Index type" prop="fieldIndexType">
        <el-select
          v-model="middleDate.fieldIndexType"
          :disabled="middleDate.fieldType === 'json'"
          style="width: 100%"
          placeholder="Select index type"
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
      <el-form-item label="FE create / edit">
        <el-switch v-model="middleDate.form" />
      </el-form-item>
      <el-form-item label="FE table column">
        <el-switch v-model="middleDate.table" />
      </el-form-item>
      <el-form-item label="FE details">
        <el-switch v-model="middleDate.desc" />
      </el-form-item>
      <el-form-item label="Import / export">
        <el-switch v-model="middleDate.excel" />
      </el-form-item>
      <el-form-item label="Sortable">
        <el-switch v-model="middleDate.sort" />
      </el-form-item>
      <el-form-item label="Required">
        <el-switch v-model="middleDate.require" />
      </el-form-item>
      <el-form-item label="Clearable">
        <el-switch v-model="middleDate.clearable" />
      </el-form-item>
      <el-form-item label="Hide search field">
        <el-switch
          v-model="middleDate.fieldSearchHide"
          :disabled="!middleDate.fieldSearchType"
        />
      </el-form-item>
      <el-form-item label="Validation message">
        <el-input v-model="middleDate.errorText" />
      </el-form-item>
    </el-form>
    <el-collapse v-model="activeNames">
      <el-collapse-item
        title="Data source (advanced; incorrect configuration may break generated code)"
        name="1"
      >
        <el-row :gutter="8">
          <el-col :span="4">
            <el-select
              v-model="middleDate.dataSource.dbName"
              placeholder="Database (empty = GVA default)"
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
              placeholder="Association"
              @change="associationChange"
            >
              <el-option label="One-to-one" :value="1" />
              <el-option label="One-to-many" :value="2" />
            </el-select>
          </el-col>
          <el-col :span="5">
            <el-select
              v-model="middleDate.dataSource.table"
              placeholder="Select source table"
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
              placeholder="Select value field"
            >
              <template #label="{ value }">
                <span>Value: </span>
                <span style="font-weight: bold">{{ value }}</span>
              </template>
              <el-option
                v-for="item in dbColumnList"
                :key="item.columnName"
                :value="item.columnName"
              >
                <span style="float: left">
                  <el-tag :type="item.isPrimary ? 'primary' : 'info'">
                    {{ item.isPrimary ? 'PK' : 'Non-PK' }}
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
                  Type: {{ item.type }}
                  <block v-if="item.comment != ''"
                    >, Comment: {{ item.comment }}</block
                  >
                </span>
              </el-option>
            </el-select>
          </el-col>
          <el-col :span="5">
            <el-select
              v-model="middleDate.dataSource.label"
              placeholder="Select label field"
            >
              <template #label="{ value }">
                <span>Label: </span>
                <span style="font-weight: bold">{{ value }}</span>
              </template>
              <el-option
                v-for="item in dbColumnList"
                :key="item.columnName"
                :value="item.columnName"
              >
                <span style="float: left">
                  <el-tag :type="item.isPrimary ? 'primary' : 'info'">
                    {{ item.isPrimary ? 'PK' : 'Non-PK' }}
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
                  Type: {{ item.type }}
                  <span v-if="item.comment != ''"
                    >, Comment: {{ item.comment }}</span
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
      callback(new Error('Invalid enum value'))
    } else {
      callback()
    }
  }

  const rules = ref({
    fieldName: [
      { required: true, message: 'Field name is required', trigger: 'blur' }
    ],
    fieldDesc: [
      { required: true, message: 'Field label (CN) is required', trigger: 'blur' }
    ],
    fieldJson: [
      { required: true, message: 'Field JSON is required', trigger: 'blur' }
    ],
    columnName: [
      { required: true, message: 'DB column name is required', trigger: 'blur' }
    ],
    fieldType: [{ required: true, message: 'Field type is required', trigger: 'blur' }],
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
        'In one-to-many association, data type becomes an array (backend stored as JSON). Continue?',
        'Confirm',
        {
          confirmButtonText: 'Continue',
          cancelButtonText: 'Cancel',
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
