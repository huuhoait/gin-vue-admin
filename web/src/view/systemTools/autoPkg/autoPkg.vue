<template>
  <div>
    <warning-bar
      href="https://www.bilibili.com/video/BV1kv4y1g7nT?p=3"
      :title="t('admin.systemtools.auto_pkg.dev_only_warning')"
    />
    <div class="gva-table-box">
      <div class="gva-btn-list gap-3 flex items-center">
        <el-button type="primary" icon="plus" @click="openDialog('addApi')">
          {{ t('admin.systemtools.auto_pkg.add') }}
        </el-button>
      </div>
      <el-table :data="tableData">
        <el-table-column align="left" :label="t('admin.systemtools.auto_pkg.id')" width="120" prop="ID" />
        <el-table-column
          align="left"
          :label="t('admin.systemtools.auto_pkg.package')"
          width="150"
          prop="packageName"
        />
        <el-table-column
          align="left"
          :label="t('admin.systemtools.auto_pkg.template')"
          width="150"
          prop="template"
        />
        <el-table-column align="left" :label="t('admin.systemtools.auto_pkg.label')" width="150" prop="label" />
        <el-table-column
          align="left"
          :label="t('admin.systemtools.auto_pkg.description')"
          min-width="150"
          prop="desc"
        />

        <el-table-column align="left" :label="t('admin.systemtools.auto_pkg.actions')" width="200">
          <template #default="scope">
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteApiFunc(scope.row)"
            >
              {{ t('admin.systemtools.auto_pkg.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-drawer v-model="dialogFormVisible" size="40%" :show-close="false">
      <warning-bar
        :title="t('admin.systemtools.auto_pkg.warn_template_types')"
      />
      <el-form ref="pkgForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item :label="t('admin.systemtools.auto_pkg.package')" prop="packageName">
          <el-input v-model="form.packageName" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="t('admin.systemtools.auto_pkg.template')" prop="template">
          <el-select v-model="form.template">
            <el-option
              v-for="template in templatesOptions"
              :label="template"
              :value="template"
              :key="template"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('admin.systemtools.auto_pkg.label')" prop="label">
          <el-input v-model="form.label" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="t('admin.systemtools.auto_pkg.description')" prop="desc">
          <el-input v-model="form.desc" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.systemtools.auto_pkg.create_package') }}</span>
          <div>
            <el-button @click="closeDialog"> {{ t('admin.systemtools.auto_pkg.cancel') }} </el-button>
            <el-button type="primary" @click="enterDialog"> {{ t('admin.systemtools.auto_pkg.confirm') }} </el-button>
          </div>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createPackageApi,
    getPackageApi,
    deletePackageApi,
    getTemplatesApi
  } from '@/api/autoCode'
  import { ref } from 'vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useI18n } from 'vue-i18n'

  const { t } = useI18n()

  defineOptions({
    name: 'AutoPkg'
  })

  const form = ref({
    packageName: '',
    template: '',
    label: '',
    desc: ''
  })
  const templatesOptions = ref([])

  const getTemplates = async () => {
    const res = await getTemplatesApi()
    if (res.code === 0) {
      templatesOptions.value = res.data
    }
  }

  getTemplates()

  const validateData = (rule, value, callback) => {
    if (/[\u4E00-\u9FA5]/g.test(value)) {
      callback(new Error(t('admin.systemtools.auto_pkg.no_chinese')))
    } else if (/^\d+$/.test(value[0])) {
      callback(new Error(t('admin.systemtools.auto_pkg.no_start_number')))
    } else if (!/^[a-zA-Z0-9_]+$/.test(value)) {
      callback(new Error(t('admin.systemtools.auto_pkg.only_letters')))
    } else {
      callback()
    }
  }

  const rules = ref({
    packageName: [
      { required: true, message: t('admin.systemtools.auto_pkg.package_required'), trigger: 'blur' },
      { validator: validateData, trigger: 'blur' }
    ],
    template: [
      { required: true, message: t('admin.systemtools.auto_pkg.template_required'), trigger: 'change' },
      { validator: validateData, trigger: 'blur' }
    ]
  })

  const dialogFormVisible = ref(false)
  const openDialog = () => {
    dialogFormVisible.value = true
  }

  const closeDialog = () => {
    dialogFormVisible.value = false
    form.value = {
      packageName: '',
      template: '',
      label: '',
      desc: ''
    }
  }

  const pkgForm = ref(null)
  const enterDialog = async () => {
    pkgForm.value.validate(async (valid) => {
      if (valid) {
        const res = await createPackageApi(form.value)
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: t('admin.systemtools.auto_pkg.added'),
            showClose: true
          })
        }
        getTableData()
        closeDialog()
      }
    })
  }

  const tableData = ref([])
  const getTableData = async () => {
    const table = await getPackageApi()
    if (table.code === 0) {
      tableData.value = table.data.pkgs
    }
  }

  const deleteApiFunc = async (row) => {
    ElMessageBox.confirm(
      t('admin.systemtools.auto_pkg.delete_confirm_msg'),
      t('admin.common.confirms.delete_title'),
      {
        confirmButtonText: t('admin.common.confirm'),
        cancelButtonText: t('admin.common.cancel'),
        type: 'warning'
      }
    ).then(async () => {
      const res = await deletePackageApi(row)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: t('admin.common.messages.deleted')
        })
        getTableData()
      }
    })
  }

  getTableData()
</script>
