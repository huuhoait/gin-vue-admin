<template>
  <div>
    <warning-bar :title="t('admin.superadmin.user.warning_switch_role')" />
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.superadmin.user.username')">
          <el-input v-model="searchInfo.username" :placeholder="t('admin.superadmin.user.username')" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.nickname')">
          <el-input v-model="searchInfo.nickname" :placeholder="t('admin.superadmin.user.nickname')" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.phone')">
          <el-input v-model="searchInfo.phone" :placeholder="t('admin.superadmin.user.phone')" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.email')">
          <el-input v-model="searchInfo.email" :placeholder="t('admin.superadmin.user.email')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            {{ t('admin.common.search') }}
          </el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addUser"
          >{{ t('admin.superadmin.user.add') }}</el-button
        >
      </div>
      <el-table :data="tableData" row-key="ID" :default-sort="{ prop: 'ID', order: 'descending' }" @sort-change="sortChange">
        <el-table-column align="left" :label="t('admin.superadmin.user.avatar')" min-width="75">
          <template #default="scope">
            <CustomPic style="margin-top: 8px" :pic-src="scope.row.headerImg" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="ID" min-width="50" prop="ID" sortable="custom" />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.user.username')"
          min-width="150"
          prop="userName"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.user.nickname')"
          min-width="150"
          prop="nickName"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.user.phone')"
          min-width="180"
          prop="phone"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.user.email')"
          min-width="180"
          prop="email"
        />
        <el-table-column align="left" :label="t('admin.superadmin.user.role')" min-width="200">
          <template #default="scope">
            <el-cascader
              v-model="scope.row.authorityIds"
              :options="authOptions"
              :show-all-levels="false"
              collapse-tags
              :props="{
                multiple: true,
                checkStrictly: true,
                label: 'authorityName',
                value: 'authorityId',
                disabled: 'disabled',
                emitPath: false
              }"
              :clearable="false"
              @visible-change="
                (flag) => {
                  changeAuthority(scope.row, flag, 0)
                }
              "
              @remove-tag="
                (removeAuth) => {
                  changeAuthority(scope.row, false, removeAuth)
                }
              "
            />
          </template>
        </el-table-column>
        <el-table-column align="left" :label="t('admin.superadmin.user.enabled')" min-width="150">
          <template #default="scope">
            <el-switch
              v-model="scope.row.enable"
              inline-prompt
              :active-value="1"
              :inactive-value="2"
              @change="
                () => {
                  switchEnable(scope.row)
                }
              "
            />
          </template>
        </el-table-column>

        <el-table-column :label="t('admin.common.operation')" :min-width="appStore.operateMinWith" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteUserFunc(scope.row)"
              >{{ t('admin.common.delete') }}</el-button
            >
            <el-button
              type="primary"
              link
              icon="edit"
              @click="openEdit(scope.row)"
              >{{ t('admin.common.edit') }}</el-button
            >
            <el-button
              type="primary"
              link
              icon="magic-stick"
              @click="resetPasswordFunc(scope.row)"
              >{{ t('admin.superadmin.user.reset_password') }}</el-button
            >
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
    <!-- password reset dialog -->
    <el-dialog
      v-model="resetPwdDialog"
      :title="t('admin.superadmin.user.reset_password')"
      width="500px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-form :model="resetPwdInfo" ref="resetPwdForm" label-width="100px">
        <el-form-item :label="t('admin.superadmin.user.username')">
          <el-input v-model="resetPwdInfo.userName" disabled />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.nickname')">
          <el-input v-model="resetPwdInfo.nickName" disabled />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.new_password')">
          <div class="flex w-full">
            <el-input class="flex-1" v-model="resetPwdInfo.password" :placeholder="t('admin.superadmin.user.new_password_placeholder')" show-password />
            <el-button type="primary" @click="generateRandomPassword" style="margin-left: 10px">
              {{ t('admin.superadmin.user.generate_random') }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeResetPwdDialog">{{ t('admin.common.cancel') }}</el-button>
          <el-button type="primary" @click="confirmResetPassword">{{ t('admin.common.confirm') }}</el-button>
        </div>
      </template>
    </el-dialog>

    <el-drawer
      v-model="addUserDialog"
      :size="appStore.drawerSize"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.superadmin.user.title') }}</span>
          <div>
            <el-button @click="closeAddUserDialog">{{ t('admin.common.cancel') }}</el-button>
            <el-button type="primary" @click="enterAddUserDialog"
              >{{ t('admin.common.confirm') }}</el-button
            >
          </div>
        </div>
      </template>

      <el-form
        ref="userForm"
        :rules="rules"
        :model="userInfo"
        label-width="80px"
      >
        <el-form-item
          v-if="dialogFlag === 'add'"
          :label="t('admin.superadmin.user.username')"
          prop="userName"
        >
          <el-input v-model="userInfo.userName" />
        </el-form-item>
        <el-form-item v-if="dialogFlag === 'add'" :label="t('admin.superadmin.user.password')" prop="password">
          <el-input v-model="userInfo.password" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.nickname')" prop="nickName">
          <el-input v-model="userInfo.nickName" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.phone')" prop="phone">
          <el-input v-model="userInfo.phone" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.email')" prop="email">
          <el-input v-model="userInfo.email" />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.role')" prop="authorityId">
          <el-cascader
            v-model="userInfo.authorityIds"
            style="width: 100%"
            :options="authOptions"
            :show-all-levels="false"
            :props="{
              multiple: true,
              checkStrictly: true,
              label: 'authorityName',
              value: 'authorityId',
              disabled: 'disabled',
              emitPath: false
            }"
            :clearable="false"
          />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.enabled')" prop="disabled">
          <el-switch
            v-model="userInfo.enable"
            inline-prompt
            :active-value="1"
            :inactive-value="2"
          />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.user.avatar')" label-width="80px">
          <SelectImage v-model="userInfo.headerImg" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    getUserList,
    setUserAuthorities,
    register,
    deleteUser
  } from '@/api/user'

  import { getAuthorityList } from '@/api/authority'
  import CustomPic from '@/components/customPic/index.vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { setUserInfo, resetPassword } from '@/api/user.js'

  import { nextTick, ref, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import SelectImage from '@/components/selectImage/selectImage.vue'
  import { useAppStore } from "@/pinia";
  import { toSQLLine } from '@/utils/stringFun'

  defineOptions({
    name: 'User'
  })

  const { t } = useI18n()
  const appStore = useAppStore()

  const searchInfo = ref({
    username: '',
    nickname: '',
    phone: '',
    email: ''
  })

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const onReset = () => {
    searchInfo.value = {
      username: '',
      nickname: '',
      phone: '',
      email: ''
    }
    orderKey.value = 'id'
    desc.value = true
    getTableData()
  }
  // initialisation helpers
  const setAuthorityOptions = (AuthorityData, optionsData) => {
    AuthorityData &&
      AuthorityData.forEach((item) => {
        if (item.children && item.children.length) {
          const option = {
            authorityId: item.authorityId,
            authorityName: item.authorityName,
            children: []
          }
          setAuthorityOptions(item.children, option.children)
          optionsData.push(option)
        } else {
          const option = {
            authorityId: item.authorityId,
            authorityName: item.authorityName
          }
          optionsData.push(option)
        }
      })
  }

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const orderKey = ref('id')
  const desc = ref(true)

  const sortChange = ({ prop, order }) => {
    if (prop) {
      orderKey.value = prop === 'ID' ? 'id' : toSQLLine(prop)
      desc.value = order === 'descending'
    }
    getTableData()
  }
  // pagination
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // fetch table data
  const getTableData = async () => {
    const table = await getUserList({
      page: page.value,
      pageSize: pageSize.value,
      orderKey: orderKey.value,
      desc: desc.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  watch(
    () => tableData.value,
    () => {
      setAuthorityIds()
    }
  )

  const authOptions = ref([])
  const setOptions = (authData) => {
    authOptions.value = []
    setAuthorityOptions(authData, authOptions.value)
  }

  const initPage = async () => {
    getTableData()
    const res = await getAuthorityList()
    setOptions(res.data)
  }

  initPage()

  // reset password dialog state
  const resetPwdDialog = ref(false)
  const resetPwdForm = ref(null)
  const resetPwdInfo = ref({
    ID: '',
    userName: '',
    nickName: '',
    password: ''
  })

  // generate a random password
  const generateRandomPassword = () => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
    let password = ''
    for (let i = 0; i < 12; i++) {
      password += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    resetPwdInfo.value.password = password
    // copy to clipboard
    navigator.clipboard.writeText(password).then(() => {
      ElMessage({
        type: 'success',
        message: t('admin.superadmin.user.password_copied')
      })
    }).catch(() => {
      ElMessage({
        type: 'error',
        message: t('admin.superadmin.user.password_copy_failed')
      })
    })
  }

  // open reset password dialog
  const resetPasswordFunc = (row) => {
    resetPwdInfo.value.ID = row.ID
    resetPwdInfo.value.userName = row.userName
    resetPwdInfo.value.nickName = row.nickName
    resetPwdInfo.value.password = ''
    resetPwdDialog.value = true
  }

  // confirm reset password
  const confirmResetPassword = async () => {
    if (!resetPwdInfo.value.password) {
      ElMessage({
        type: 'warning',
        message: t('admin.superadmin.user.password_required')
      })
      return
    }

    const res = await resetPassword({
      ID: resetPwdInfo.value.ID,
      password: resetPwdInfo.value.password
    })

    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg || t('admin.superadmin.user.reset_password_success')
      })
      resetPwdDialog.value = false
    } else {
      ElMessage({
        type: 'error',
        message: res.msg || t('admin.superadmin.user.reset_password_failed')
      })
    }
  }

  // close reset password dialog
  const closeResetPwdDialog = () => {
    resetPwdInfo.value.password = ''
    resetPwdDialog.value = false
  }
  const setAuthorityIds = () => {
    tableData.value &&
      tableData.value.forEach((user) => {
        user.authorityIds =
          user.authorities &&
          user.authorities.map((i) => {
            return i.authorityId
          })
      })
  }

  const deleteUserFunc = async (row) => {
    ElMessageBox.confirm(t('admin.common.delete_confirm'), t('admin.common.are_you_sure'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    }).then(async () => {
      const res = await deleteUser({ id: row.ID })
      if (res.code === 0) {
        ElMessage.success(t('admin.superadmin.user.delete_success'))
        await getTableData()
      }
    })
  }

  // dialog state
  const userInfo = ref({
    userName: '',
    password: '',
    nickName: '',
    headerImg: '',
    authorityId: '',
    authorityIds: [],
    enable: 1
  })

  // Validator messages call t() per invocation so locale switches are reactive.
  const rules = ref({
    userName: [
      { required: true, message: () => t('admin.superadmin.user.rules.username_required'), trigger: 'blur' },
      { min: 5, message: () => t('admin.superadmin.user.rules.min5'), trigger: 'blur' }
    ],
    password: [
      { required: true, message: () => t('admin.superadmin.user.rules.password_required'), trigger: 'blur' },
      { min: 6, message: () => t('admin.superadmin.user.rules.min6'), trigger: 'blur' }
    ],
    nickName: [{ required: true, message: () => t('admin.superadmin.user.rules.nickname_required'), trigger: 'blur' }],
    phone: [
      {
        pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/,
        message: () => t('admin.superadmin.user.rules.phone_invalid'),
        trigger: 'blur'
      }
    ],
    email: [
      {
        pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g,
        message: () => t('admin.superadmin.user.rules.email_invalid'),
        trigger: 'blur'
      }
    ],
    authorityId: [
      { required: true, message: () => t('admin.superadmin.user.rules.role_required'), trigger: 'blur' }
    ]
  })
  const userForm = ref(null)
  const enterAddUserDialog = async () => {
    userInfo.value.authorityId = userInfo.value.authorityIds[0]
    userForm.value.validate(async (valid) => {
      if (valid) {
        const req = {
          ...userInfo.value
        }
        if (dialogFlag.value === 'add') {
          const res = await register(req)
          if (res.code === 0) {
            ElMessage({ type: 'success', message: t('admin.superadmin.user.create_success') })
            await getTableData()
            closeAddUserDialog()
          }
        }
        if (dialogFlag.value === 'edit') {
          const res = await setUserInfo(req)
          if (res.code === 0) {
            ElMessage({ type: 'success', message: t('admin.superadmin.user.edit_success') })
            await getTableData()
            closeAddUserDialog()
          }
        }
      }
    })
  }

  const addUserDialog = ref(false)
  const closeAddUserDialog = () => {
    userForm.value.resetFields()
    userInfo.value.headerImg = ''
    userInfo.value.authorityIds = []
    addUserDialog.value = false
  }

  const dialogFlag = ref('add')

  const addUser = () => {
    dialogFlag.value = 'add'
    addUserDialog.value = true
  }

  const tempAuth = {}
  const changeAuthority = async (row, flag, removeAuth) => {
    if (flag) {
      if (!removeAuth) {
        tempAuth[row.ID] = [...row.authorityIds]
      }
      return
    }
    await nextTick()
    const res = await setUserAuthorities({
      ID: row.ID,
      authorityIds: row.authorityIds
    })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: t('admin.superadmin.user.role_update_success') })
    } else {
      if (!removeAuth) {
        row.authorityIds = [...tempAuth[row.ID]]
        delete tempAuth[row.ID]
      } else {
        row.authorityIds = [removeAuth, ...row.authorityIds]
      }
    }
  }

  const openEdit = (row) => {
    dialogFlag.value = 'edit'
    userInfo.value = JSON.parse(JSON.stringify(row))
    addUserDialog.value = true
  }

  const switchEnable = async (row) => {
    userInfo.value = JSON.parse(JSON.stringify(row))
    await nextTick()
    const req = {
      ...userInfo.value
    }
    const res = await setUserInfo(req)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: req.enable === 2
          ? t('admin.superadmin.user.disable_success')
          : t('admin.superadmin.user.enable_success')
      })
      await getTableData()
      userInfo.value.headerImg = ''
      userInfo.value.authorityIds = []
    }
  }
</script>

<style lang="scss">
  .header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
  }
</style>
