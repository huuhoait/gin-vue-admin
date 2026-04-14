<template>
  <div class="authority">
    <warning-bar :title="t('admin.superadmin.authority.warning_switch_role')" />
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addAuthority(0)"
          >{{ t('admin.superadmin.authority.add') }}</el-button
        >
      </div>
      <el-table
        :data="tableData"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        row-key="authorityId"
        style="width: 100%"
      >
        <el-table-column :label="t('admin.superadmin.authority.role_id')" min-width="180" prop="authorityId" />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.authority.role_name')"
          min-width="180"
          prop="authorityName"
        />
        <el-table-column align="left" :label="t('admin.common.operation')" width="560">
          <template #default="scope">
            <el-button
              icon="setting"
              type="primary"
              link
              @click="openDrawer(scope.row)"
              >{{ t('admin.superadmin.authority.set_permissions') }}</el-button
            >
            <el-button
              icon="user"
              type="primary"
              link
              @click="openAssignDrawer(scope.row)"
              >{{ t('admin.superadmin.authority.assign_to_users') }}</el-button
            >
            <el-button
              icon="plus"
              type="primary"
              link
              @click="addAuthority(scope.row.authorityId)"
              >{{ t('admin.superadmin.authority.add_child') }}</el-button
            >
            <el-button
              icon="copy-document"
              type="primary"
              link
              @click="copyAuthorityFunc(scope.row)"
              >{{ t('admin.superadmin.authority.copy') }}</el-button
            >
            <el-button
              icon="edit"
              type="primary"
              link
              @click="editAuthority(scope.row)"
              >{{ t('admin.common.edit') }}</el-button
            >
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteAuth(scope.row)"
              >{{ t('admin.common.delete') }}</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </div>
    <!-- role create / edit drawer -->
    <el-drawer v-model="authorityFormVisible" :size="appStore.drawerSize" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ authorityTitleForm }}</span>
          <div>
            <el-button @click="closeAuthorityForm">{{ t('admin.common.cancel') }}</el-button>
            <el-button type="primary" @click="submitAuthorityForm"
              >{{ t('admin.common.confirm') }}</el-button
            >
          </div>
        </div>
      </template>
      <el-form
        ref="authorityForm"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item :label="t('admin.superadmin.authority.parent_role')" prop="parentId">
          <el-cascader
            v-model="form.parentId"
            style="width: 100%"
            :disabled="dialogType === 'add'"
            :options="AuthorityOption"
            :props="{
              checkStrictly: true,
              label: 'authorityName',
              value: 'authorityId',
              disabled: 'disabled',
              emitPath: false
            }"
            :show-all-levels="false"
            filterable
          />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.authority.role_id')" prop="authorityId">
          <el-input
            v-model="form.authorityId"
            :disabled="dialogType === 'edit'"
            autocomplete="off"
            maxlength="15"
          />
        </el-form-item>
        <el-form-item :label="t('admin.superadmin.authority.role_name')" prop="authorityName">
          <el-input v-model="form.authorityName" autocomplete="off" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      v-if="drawer"
      v-model="drawer"
      :size="appStore.drawerSize"
      :title="t('admin.superadmin.authority.title')"
    >
      <el-tabs :before-leave="autoEnter" type="border-card">
        <el-tab-pane :label="t('admin.superadmin.authority.tabs.menus')">
          <Menus ref="menus" :row="activeRow" @changeRow="changeRow" />
        </el-tab-pane>
        <el-tab-pane :label="t('admin.superadmin.authority.tabs.apis')">
          <Apis ref="apis" :row="activeRow" @changeRow="changeRow" />
        </el-tab-pane>
        <el-tab-pane :label="t('admin.superadmin.authority.tabs.data')">
          <Datas
            ref="datas"
            :authority="tableData"
            :row="activeRow"
            @changeRow="changeRow"
          />
        </el-tab-pane>
      </el-tabs>
    </el-drawer>

    <!-- assign-to-users drawer -->
    <el-drawer
      v-model="assignDrawerVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.superadmin.authority.assign_title', { name: assignRow.authorityName }) }}</span>
          <div>
            <el-button @click="assignDrawerVisible = false">{{ t('admin.common.cancel') }}</el-button>
            <el-button type="primary" :loading="assignSubmitting" @click="confirmAssign">{{ t('admin.common.confirm') }}</el-button>
          </div>
        </div>
      </template>
      <warning-bar :title="t('admin.superadmin.authority.assign_warning')" />
      <div class="gva-search-box">
        <el-form :inline="true" :model="userSearchInfo">
          <el-form-item :label="t('admin.superadmin.user.username')">
            <el-input v-model="userSearchInfo.username" :placeholder="t('admin.superadmin.user.username')" />
          </el-form-item>
          <el-form-item :label="t('admin.superadmin.user.nickname')">
            <el-input v-model="userSearchInfo.nickName" :placeholder="t('admin.superadmin.user.nickname')" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="search" @click="searchUserData">{{ t('admin.common.search') }}</el-button>
            <el-button icon="refresh" @click="resetUserSearch">{{ t('admin.common.reset') }}</el-button>
          </el-form-item>
        </el-form>
      </div>
      <el-table
        ref="userTableRef"
        v-loading="assignLoading"
        :data="userTableData"
        row-key="ID"
        :default-sort="{ prop: 'ID', order: 'descending' }"
        @sort-change="sortChange"
        @select="handleSelect"
        @select-all="handleSelectAll"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="ID" prop="ID" width="80" sortable="custom" />
        <el-table-column :label="t('admin.superadmin.user.username')" prop="userName" min-width="120" />
        <el-table-column :label="t('admin.superadmin.user.nickname')" prop="nickName" min-width="120" />
      </el-table>
      <div class="flex justify-center mt-4">
        <el-pagination
          :current-page="userSearchInfo.page"
          :page-size="userSearchInfo.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="userTotal"
          layout="total, sizes, prev, pager, next"
          @current-change="handleUserPageChange"
          @size-change="handleUserSizeChange"
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    getAuthorityList,
    deleteAuthority,
    createAuthority,
    updateAuthority,
    copyAuthority,
    getUsersByAuthorityId,
    setRoleUsers
  } from '@/api/authority'
  import { getUserList } from '@/api/user'

  import Menus from '@/view/superAdmin/authority/components/menus.vue'
  import Apis from '@/view/superAdmin/authority/components/apis.vue'
  import Datas from '@/view/superAdmin/authority/components/datas.vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'

  import { ref, nextTick } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import { useAppStore } from "@/pinia"
  import { toSQLLine } from '@/utils/stringFun'

  defineOptions({
    name: 'Authority'
  })

  const { t } = useI18n()

  const mustUint = (rule, value, callback) => {
    if (!/^[0-9]*[1-9][0-9]*$/.test(value)) {
      return callback(new Error(t('admin.superadmin.authority.rules.must_positive_int')))
    }
    return callback()
  }

  const AuthorityOption = ref([
    {
      authorityId: 0,
      authorityName: t('admin.superadmin.authority.root_role_strict')
    }
  ])
  const drawer = ref(false)
  const dialogType = ref('add')
  const activeRow = ref({})
  const appStore = useAppStore()

  const authorityTitleForm = ref(t('admin.superadmin.authority.add'))
  const authorityFormVisible = ref(false)
  const apiDialogFlag = ref(false)
  const copyForm = ref({})

  const form = ref({
    authorityId: 0,
    authorityName: '',
    parentId: 0
  })
  const rules = ref({
    authorityId: [
      { required: true, message: () => t('admin.superadmin.authority.rules.role_id_required'), trigger: 'blur' },
      { validator: mustUint, trigger: 'blur', message: () => t('admin.superadmin.authority.rules.must_positive_int') }
    ],
    authorityName: [
      { required: true, message: () => t('admin.superadmin.authority.rules.role_name_required'), trigger: 'blur' }
    ],
    parentId: [{ required: true, message: () => t('admin.superadmin.authority.rules.parent_required'), trigger: 'blur' }]
  })

  const tableData = ref([])

  // fetch
  const getTableData = async () => {
    const table = await getAuthorityList()
    if (table.code === 0) {
      tableData.value = table.data
    }
  }

  getTableData()

  const changeRow = (key, value) => {
    activeRow.value[key] = value
  }
  const menus = ref(null)
  const apis = ref(null)
  const datas = ref(null)
  const autoEnter = (activeName, oldActiveName) => {
    const paneArr = [menus, apis, datas]
    if (oldActiveName) {
      if (paneArr[oldActiveName].value.needConfirm) {
        paneArr[oldActiveName].value.enterAndNext()
        paneArr[oldActiveName].value.needConfirm = false
      }
    }
  }
  // copy role
  const copyAuthorityFunc = (row) => {
    setOptions()
    authorityTitleForm.value = t('admin.superadmin.authority.copy')
    dialogType.value = 'copy'
    for (const k in form.value) {
      form.value[k] = row[k]
    }
    copyForm.value = row
    authorityFormVisible.value = true
  }
  const openDrawer = (row) => {
    drawer.value = true
    activeRow.value = row
  }
  // delete role
  const deleteAuth = (row) => {
    ElMessageBox.confirm(t('admin.superadmin.authority.delete_confirm_message'), t('admin.common.are_you_sure'), {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning'
    })
      .then(async () => {
        const res = await deleteAuthority({ authorityId: row.authorityId })
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: t('admin.superadmin.authority.delete_success')
          })

          getTableData()
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('admin.superadmin.authority.delete_cancelled')
        })
      })
  }
  // reset form state
  const authorityForm = ref(null)
  const initForm = () => {
    if (authorityForm.value) {
      authorityForm.value.resetFields()
    }
    form.value = {
      authorityId: 0,
      authorityName: '',
      parentId: 0
    }
  }
  // close drawer
  const closeAuthorityForm = () => {
    initForm()
    authorityFormVisible.value = false
    apiDialogFlag.value = false
  }
  // submit drawer

  const submitAuthorityForm = () => {
    authorityForm.value.validate(async (valid) => {
      if (valid) {
        form.value.authorityId = Number(form.value.authorityId)
        switch (dialogType.value) {
          case 'add':
            {
              const res = await createAuthority(form.value)
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: t('admin.superadmin.authority.add_success')
                })
                getTableData()
                closeAuthorityForm()
              }
            }
            break
          case 'edit':
            {
              const res = await updateAuthority(form.value)
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: t('admin.superadmin.authority.add_success')
                })
                getTableData()
                closeAuthorityForm()
              }
            }
            break
          case 'copy': {
            const data = {
              authority: {
                authorityId: 0,
                authorityName: '',
                datauthorityId: [],
                parentId: 0
              },
              oldAuthorityId: 0
            }
            data.authority.authorityId = form.value.authorityId
            data.authority.authorityName = form.value.authorityName
            data.authority.parentId = form.value.parentId
            data.authority.dataAuthorityId = copyForm.value.dataAuthorityId
            data.oldAuthorityId = copyForm.value.authorityId
            const res = await copyAuthority(data)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: t('admin.superadmin.authority.copy_success')
              })
              getTableData()
            }
          }
        }

        initForm()
        authorityFormVisible.value = false
      }
    })
  }
  const setOptions = () => {
    AuthorityOption.value = [
      {
        authorityId: 0,
        authorityName: t('admin.superadmin.authority.root_role_user')
      }
    ]
    setAuthorityOptions(tableData.value, AuthorityOption.value, false)
  }
  const setAuthorityOptions = (AuthorityData, optionsData, disabled) => {
    AuthorityData &&
      AuthorityData.forEach((item) => {
        if (item.children && item.children.length) {
          const option = {
            authorityId: item.authorityId,
            authorityName: item.authorityName,
            disabled: disabled || item.authorityId === form.value.authorityId,
            children: []
          }
          setAuthorityOptions(
            item.children,
            option.children,
            disabled || item.authorityId === form.value.authorityId
          )
          optionsData.push(option)
        } else {
          const option = {
            authorityId: item.authorityId,
            authorityName: item.authorityName,
            disabled: disabled || item.authorityId === form.value.authorityId
          }
          optionsData.push(option)
        }
      })
  }
  // add role
  const addAuthority = (parentId) => {
    initForm()
    authorityTitleForm.value = t('admin.superadmin.authority.add')
    dialogType.value = 'add'
    form.value.parentId = parentId
    setOptions()
    authorityFormVisible.value = true
  }
  // edit role
  const editAuthority = (row) => {
    setOptions()
    authorityTitleForm.value = t('admin.superadmin.authority.edit')
    dialogType.value = 'edit'
    for (const key in form.value) {
      form.value[key] = row[key]
    }
    setOptions()
    authorityForm.value && authorityForm.value.clearValidate()
    authorityFormVisible.value = true
  }

  // assign role to users
  const assignDrawerVisible = ref(false)
  const assignRow = ref({})
  const userTableData = ref([])
  const userTotal = ref(0)
  const userSearchInfo = ref({ page: 1, pageSize: 10, username: '', nickName: '', orderKey: 'id', desc: true })
  const assignLoading = ref(false)
  const assignSubmitting = ref(false)
  const userTableRef = ref(null)

  const selectedUserIds = ref(new Set())

  const openAssignDrawer = async (row) => {
    assignRow.value = row
    userSearchInfo.value = { page: 1, pageSize: 10, username: '', nickName: '' }
    selectedUserIds.value = new Set()
    assignDrawerVisible.value = true
    const res = await getUsersByAuthorityId(row.authorityId)
    if (res.code === 0 && res.data) {
      selectedUserIds.value = new Set(res.data)
    }
    getUserData()
  }

  const getUserData = async () => {
    assignLoading.value = true
    const res = await getUserList(userSearchInfo.value)
    if (res.code === 0) {
      userTableData.value = res.data.list
      userTotal.value = res.data.total
      await nextTick()
      userTableData.value.forEach((user) => {
        userTableRef.value && userTableRef.value.toggleRowSelection(user, selectedUserIds.value.has(user.ID))
      })
    }
    assignLoading.value = false
  }

  const handleSelect = (selection, row) => {
    if (selection.some(u => u.ID === row.ID)) {
      selectedUserIds.value.add(row.ID)
    } else {
      selectedUserIds.value.delete(row.ID)
    }
  }

  const handleSelectAll = (selection) => {
    const selectedIds = new Set(selection.map(u => u.ID))
    userTableData.value.forEach((user) => {
      if (selectedIds.has(user.ID)) {
        selectedUserIds.value.add(user.ID)
      } else {
        selectedUserIds.value.delete(user.ID)
      }
    })
  }

  const sortChange = ({ prop, order }) => {
    if (prop) {
      userSearchInfo.value.orderKey = prop === 'ID' ? 'id' : toSQLLine(prop)
      userSearchInfo.value.desc = order === 'descending'
    }
    getUserData()
  }

  const searchUserData = () => {
    userSearchInfo.value.page = 1
    getUserData()
  }

  const resetUserSearch = () => {
    userSearchInfo.value = { page: 1, pageSize: 10, username: '', nickName: '' }
    getUserData()
  }

  const handleUserPageChange = (page) => {
    userSearchInfo.value.page = page
    getUserData()
  }

  const handleUserSizeChange = (size) => {
    userSearchInfo.value.pageSize = size
    userSearchInfo.value.page = 1
    getUserData()
  }

  const confirmAssign = async () => {
    assignSubmitting.value = true
    try {
      const res = await setRoleUsers({
        authorityId: assignRow.value.authorityId,
        userIds: [...selectedUserIds.value]
      })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: t('admin.superadmin.authority.assign_success') })
        assignDrawerVisible.value = false
      }
    } catch {
      ElMessage({ type: 'error', message: t('admin.superadmin.authority.assign_failed') })
    }
    assignSubmitting.value = false
  }
</script>

<style lang="scss">
  .authority {
    .el-input-number {
      margin-left: 15px;
      span {
        display: none;
      }
    }
  }
  .tree-content {
    margin-top: 10px;
    height: calc(100vh - 158px);
    overflow: auto;
  }
</style>
