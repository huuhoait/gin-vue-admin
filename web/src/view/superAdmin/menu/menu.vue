<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addMenu(0)">
          {{ t('admin.superadmin.menu.add_root') }}
        </el-button>
      </div>

      <!-- Menu tree matches sidebar 1:1; pageSize defaults to 999 -->
      <el-table :data="tableData" row-key="ID">
        <el-table-column align="left" label="ID" min-width="100" prop="ID" />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.menu.columns.display_name')"
          min-width="120"
          prop="authorityName"
        >
          <template #default="scope">
            <span>{{ scope.row.meta.title }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          :label="t('admin.superadmin.menu.columns.icon')"
          min-width="140"
          prop="authorityName"
        >
          <template #default="scope">
            <div v-if="scope.row.meta.icon" class="icon-column">
              <el-icon>
                <component :is="scope.row.meta.icon" />
              </el-icon>
              <span>{{ scope.row.meta.icon }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          :label="t('admin.superadmin.menu.columns.route_name')"
          show-overflow-tooltip
          min-width="160"
          prop="name"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.menu.columns.route_path')"
          show-overflow-tooltip
          min-width="160"
          prop="path"
        />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.menu.columns.hidden')"
          min-width="100"
          prop="hidden"
        >
          <template #default="scope">
            <span>{{ scope.row.hidden ? t('admin.superadmin.menu.hidden.hidden') : t('admin.superadmin.menu.hidden.visible') }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          :label="t('admin.superadmin.menu.columns.parent')"
          min-width="90"
          prop="parentId"
        />
        <el-table-column align="left" :label="t('admin.superadmin.menu.columns.sort')" min-width="70" prop="sort" />
        <el-table-column
          align="left"
          :label="t('admin.superadmin.menu.columns.component_path')"
          min-width="360"
          prop="component"
        />
        <el-table-column align="left" fixed="right" :label="t('admin.common.operation')" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="plus"
              @click="addMenu(scope.row.ID)"
            >
              {{ t('admin.superadmin.menu.actions.add_child') }}
            </el-button>
            <el-button
              type="primary"
              link
              icon="edit"
              @click="editMenu(scope.row.ID)"
            >
              {{ t('admin.common.edit') }}
            </el-button>
            <el-button
              type="primary"
              link
              icon="user"
              @click="openAssignRoleDrawer(scope.row)"
            >
              {{ t('admin.superadmin.menu.actions.assign_roles') }}
            </el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteMenu(scope.row.ID)"
            >
              {{ t('admin.common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <el-drawer
      v-model="dialogFormVisible"
      :size="appStore.drawerSize"
      :before-close="handleClose"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle }}</span>
          <div>
            <el-button @click="closeDialog">{{ t('admin.common.cancel') }}</el-button>
            <el-button type="primary" @click="enterDialog">{{ t('admin.common.confirm') }}</el-button>
          </div>
        </div>
      </template>

      <warning-bar :title="t('admin.superadmin.menu.warning_add_menu')" />
      
      <!-- Basic info -->
      <div class="border-b border-gray-200">
        <h3 class="font-semibold text-gray-700 mb-4">{{ t('admin.superadmin.menu.sections.basic_info') }}</h3>
        <el-form
          v-if="dialogFormVisible"
          ref="menuForm"
          :inline="true"
          :model="form"
          :rules="rules"
          label-position="top"
        >
          <el-row class="w-full">
            <el-col :span="24">
              <el-form-item :label="t('admin.superadmin.menu.fields.component_path')" prop="component">
                <components-cascader
                  :component="form.component"
                  @change="fmtComponent"
                />
                <div class="form-tip">
                  <el-icon><InfoFilled /></el-icon>
                  <span>{{ t('admin.superadmin.menu.tips.component_router_view') }}</span>
                  <el-button
                    size="small"
                    type="text"
                    @click="form.component = 'view/routerHolder.vue'"
                  >
                    {{ t('admin.superadmin.menu.tips.click_to_set') }}
                  </el-button>
                </div>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row class="w-full">
            <el-col :span="12">
              <el-form-item :label="t('admin.superadmin.menu.fields.display_name')" prop="meta.title">
                <el-input 
                  v-model="form.meta.title" 
                  autocomplete="off" 
                  :placeholder="t('admin.superadmin.menu.placeholders.display_name')"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="t('admin.superadmin.menu.fields.route_name')" prop="path">
                <el-input
                  v-model="form.name"
                  autocomplete="off"
                  :placeholder="t('admin.superadmin.menu.placeholders.unique_name')"
                  @change="changeName"
                />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </div>
       
      <!-- Route config -->
      <div class="border-b border-gray-200">
        <h3 class="font-semibold text-gray-700 mb-4">{{ t('admin.superadmin.menu.sections.route_config') }}</h3>
        <el-form
          :inline="true"
          :model="form"
          :rules="rules"
          label-position="top"
        >
           <el-row class="w-full">
             <el-col :span="12">
               <el-form-item :label="t('admin.superadmin.menu.fields.parent_id')">
                 <el-cascader
                   v-model="form.parentId"
                   style="width: 100%"
                   :disabled="!isEdit"
                   :options="menuOption"
                   :props="{
                     checkStrictly: true,
                     label: 'title',
                     value: 'ID',
                     disabled: 'disabled',
                     emitPath: false
                   }"
                   :show-all-levels="false"
                   filterable
                   :placeholder="t('admin.superadmin.menu.placeholders.select_parent')"
                 />
               </el-form-item>
             </el-col>
             <el-col :span="12">
               <el-form-item prop="path">
                 <template #label>
                  <div class="inline-flex items-center h-4">
                    <span>{{ t('admin.superadmin.menu.fields.route_path') }}</span>
                     <el-checkbox
                       class="ml-2"
                       v-model="checkFlag"
                      >{{ t('admin.superadmin.menu.actions.add_params') }}</el-checkbox
                     >
                    </div>
                 </template>
                 <el-input
                   v-model="form.path"
                   :disabled="!checkFlag"
                   autocomplete="off"
                  :placeholder="t('admin.superadmin.menu.placeholders.path_params_hint')"
                 />
               </el-form-item>
             </el-col>
           </el-row>
        </el-form>
      </div>
       
      <!-- Display settings -->
      <div class="border-b border-gray-200">
        <h3 class="font-semibold text-gray-700 mb-4">{{ t('admin.superadmin.menu.sections.display_settings') }}</h3>
        <el-form
          :inline="true"
          :model="form"
          :rules="rules"
          label-position="top"
        >
           <el-row class="w-full">
              <el-col :span="8">
                <el-form-item :label="t('admin.superadmin.menu.fields.icon')" prop="meta.icon">
                  <icon v-model="form.meta.icon" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="t('admin.superadmin.menu.fields.sort')" prop="sort">
                  <el-input 
                    v-model.number="form.sort" 
                    autocomplete="off" 
                    :placeholder="t('admin.superadmin.menu.placeholders.sort')"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="t('admin.superadmin.menu.fields.hidden')">
                  <el-select
                    v-model="form.hidden"
                    style="width: 100%"
                    :placeholder="t('admin.superadmin.menu.placeholders.hidden')"
                  >
                    <el-option :value="false" :label="t('admin.common.no')" />
                    <el-option :value="true" :label="t('admin.common.yes')" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
        </el-form>
      </div>
        
      <!-- Advanced config -->
      <div class="border-b border-gray-200">
        <h3 class="font-semibold text-gray-700 mb-4">{{ t('admin.superadmin.menu.sections.advanced_config') }}</h3>
        <el-form
          :inline="true"
          :model="form"
          :rules="rules"
          label-position="top"
        >
            <el-row class="w-full">
              <el-col :span="12">
                <el-form-item prop="meta.activeName">
                  <template #label>
                    <div class="label-with-tooltip">
                      <span>{{ t('admin.superadmin.menu.fields.active_menu') }}</span>
                      <el-tooltip
                        :content="t('admin.superadmin.menu.tips.active_menu_help')"
                        placement="top"
                        effect="light"
                      >
                        <el-icon><QuestionFilled /></el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                  <el-input
                    v-model="form.meta.activeName"
                    :placeholder="form.name || t('admin.superadmin.menu.placeholders.active_menu')"
                    autocomplete="off"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="KeepAlive" prop="meta.keepAlive">
                  <el-select
                    v-model="form.meta.keepAlive"
                    style="width: 100%"
                    :placeholder="t('admin.superadmin.menu.placeholders.keep_alive')"
                  >
                    <el-option :value="false" :label="t('admin.common.no')" />
                    <el-option :value="true" :label="t('admin.common.yes')" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
             <el-row class="w-full">
               <el-col :span="8">
                 <el-form-item label="CloseTab" prop="meta.closeTab">
                   <el-select
                     v-model="form.meta.closeTab"
                     style="width: 100%"
                    :placeholder="t('admin.superadmin.menu.placeholders.close_tab')"
                   >
                    <el-option :value="false" :label="t('admin.common.no')" />
                    <el-option :value="true" :label="t('admin.common.yes')" />
                   </el-select>
                 </el-form-item>
               </el-col>
               <el-col :span="8">
                 <el-form-item>
                   <template #label>
                     <div class="label-with-tooltip">
                      <span>{{ t('admin.superadmin.menu.fields.base_view') }}</span>
                       <el-tooltip
                        :content="t('admin.superadmin.menu.tips.base_view_help')"
                         placement="top"
                         effect="light"
                       >
                         <el-icon><QuestionFilled /></el-icon>
                       </el-tooltip>
                     </div>
                   </template>
                   <el-select
                     v-model="form.meta.defaultMenu"
                     style="width: 100%"
                     :placeholder="t('admin.superadmin.menu.placeholders.base_view')"
                   >
                     <el-option :value="false" :label="t('admin.common.no')" />
                     <el-option :value="true" :label="t('admin.common.yes')" />
                   </el-select>
                 </el-form-item>
               </el-col>
               <el-col :span="8">
                 <el-form-item>
                   <template #label>
                     <div class="label-with-tooltip">
                      <span>{{ t('admin.superadmin.menu.fields.transition') }}</span>
                       <el-tooltip
                        :content="t('admin.superadmin.menu.tips.transition_help')"
                         placement="top"
                         effect="light"
                       >
                         <el-icon><QuestionFilled /></el-icon>
                       </el-tooltip>
                     </div>
                   </template>
                   <el-select
                     v-model="form.meta.transitionType"
                     style="width: 100%"
                     :placeholder="t('admin.superadmin.menu.placeholders.follow_global')"
                     clearable
                   >
                     <el-option value="fade" :label="t('admin.superadmin.menu.transitions.fade')" />
                     <el-option value="slide" :label="t('admin.superadmin.menu.transitions.slide')" />
                     <el-option value="zoom" :label="t('admin.superadmin.menu.transitions.zoom')" />
                     <el-option value="none" :label="t('admin.superadmin.menu.transitions.none')" />
                   </el-select>
                 </el-form-item>
               </el-col>
             </el-row>
        </el-form>
      </div>
          
      <!-- Menu params -->
      <div class="border-b border-gray-200">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-semibold text-gray-700">{{ t('admin.superadmin.menu.sections.params_config') }}</h3>
          <el-button type="primary" size="small" @click="addParameter(form)">
            {{ t('admin.superadmin.menu.actions.add_param') }}
          </el-button>
        </div>
            <el-table 
              :data="form.parameters" 
              style="width: 100%"
              class="parameter-table"
            >
              <el-table-column
                align="center"
                prop="type"
                :label="t('admin.superadmin.menu.params.type')"
                width="150"
              >
                <template #default="scope">
                  <el-select 
                    v-model="scope.row.type" 
                    :placeholder="t('admin.common.select')"
                    size="small"
                  >
                    <el-option key="query" value="query" label="query" />
                    <el-option key="params" value="params" label="params" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column align="center" prop="key" :label="t('admin.superadmin.menu.params.key')" width="150">
                <template #default="scope">
                  <el-input 
                    v-model="scope.row.key" 
                    size="small"
                    :placeholder="t('admin.superadmin.menu.placeholders.param_key')"
                  />
                </template>
              </el-table-column>
              <el-table-column align="center" prop="value" :label="t('admin.superadmin.menu.params.value')">
                <template #default="scope">
                  <el-input 
                    v-model="scope.row.value" 
                    size="small"
                    :placeholder="t('admin.superadmin.menu.placeholders.param_value')"
                  />
                </template>
              </el-table-column>
              <el-table-column align="center" :label="t('admin.common.operation')" width="100">
                <template #default="scope">
                  <el-button
                    type="danger"
                    size="small"
                    @click="deleteParameter(form.parameters, scope.$index)"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
      </div>
           
      <!-- Controllable buttons -->
      <div class="mb-2 mt-2">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-semibold text-gray-700">{{ t('admin.superadmin.menu.sections.buttons_config') }}</h3>
          <div class="flex items-center gap-2">
            <el-button type="primary" size="small" @click="addBtn(form)">
              {{ t('admin.superadmin.menu.actions.add_button') }}
            </el-button>
            <el-tooltip
              :content="t('admin.superadmin.menu.tips.button_docs')"
              placement="top"
              effect="light"
            >
              <el-icon
                class="cursor-pointer text-blue-500 hover:text-blue-700"
                @click="toDoc('https://www.gin-vue-admin.com/guide/web/button-auth.html')"
              >
                <QuestionFilled />
              </el-icon>
            </el-tooltip>
          </div>
        </div>
             <el-table 
               :data="form.menuBtn" 
               style="width: 100%"
               class="button-table"
             >
               <el-table-column
                 align="center"
                 prop="name"
                 :label="t('admin.superadmin.menu.buttons.name')"
                 width="150"
               >
                 <template #default="scope">
                   <el-input 
                     v-model="scope.row.name" 
                     size="small"
                     :placeholder="t('admin.superadmin.menu.placeholders.button_name')"
                   />
                 </template>
               </el-table-column>
              <el-table-column align="center" prop="desc" :label="t('admin.superadmin.menu.buttons.remark')">
                 <template #default="scope">
                   <el-input 
                     v-model="scope.row.desc" 
                     size="small"
                    :placeholder="t('admin.superadmin.menu.placeholders.button_remark')"
                   />
                 </template>
               </el-table-column>
              <el-table-column align="center" :label="t('admin.common.operation')" width="100">
                 <template #default="scope">
                   <el-button
                     type="danger"
                     size="small"
                     @click="deleteBtn(form.menuBtn, scope.$index)"
                   >
                     <el-icon><Delete /></el-icon>
                   </el-button>
                 </template>
               </el-table-column>
             </el-table>
       </div>
    </el-drawer>

    <!-- Assign to roles drawer -->
    <el-drawer
      v-model="assignRoleDrawerVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.superadmin.menu.assign.title', { name: assignMenuRow.meta?.title }) }}</span>
          <div>
            <el-button @click="assignRoleDrawerVisible = false">{{ t('admin.common.cancel') }}</el-button>
            <el-button type="primary" :loading="assignRoleSubmitting" @click="confirmAssignRole">{{ t('admin.common.confirm') }}</el-button>
          </div>
        </div>
      </template>
      <warning-bar :title="t('admin.superadmin.menu.assign.warning')" />
      <el-tree
        ref="roleTreeRef"
        v-loading="assignRoleLoading"
        :data="authorityTreeData"
        :props="{ label: 'authorityName', children: 'children', disabled: isRoleDisabled }"
        node-key="authorityId"
        show-checkbox
        check-strictly
        default-expand-all
      />
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    updateBaseMenu,
    getMenuList,
    addBaseMenu,
    deleteBaseMenu,
    getBaseMenuById,
    getMenuRoles,
    setMenuRoles
  } from '@/api/menu'
  import { getAuthorityList } from '@/api/authority'
  import icon from '@/view/superAdmin/menu/icon.vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { canRemoveAuthorityBtnApi } from '@/api/authorityBtn'
  import { reactive, ref, nextTick } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { QuestionFilled, InfoFilled, Delete } from '@element-plus/icons-vue'
  import { toDoc } from '@/utils/doc'
  import { toLowerCase } from '@/utils/stringFun'
  import ComponentsCascader from '@/view/superAdmin/menu/components/components-cascader.vue'
  import { useI18n } from 'vue-i18n'

  import pathInfo from '@/pathInfo.json'
  import { useAppStore } from "@/pinia";

  defineOptions({
    name: 'Menus'
  })

  const appStore = useAppStore()
  const { t } = useI18n()

  const rules = reactive({
    path: [{ required: true, message: t('admin.superadmin.menu.validation.route_name_required'), trigger: 'blur' }],
    component: [{ required: true, message: t('admin.superadmin.menu.validation.component_required'), trigger: 'blur' }],
    'meta.title': [
      { required: true, message: t('admin.superadmin.menu.validation.display_name_required'), trigger: 'blur' }
    ]
  })

  const tableData = ref([])
  // Query
  const getTableData = async () => {
    const table = await getMenuList()
    if (table.code === 0) {
      tableData.value = table.data
    }
  }

  getTableData()

  // Add parameter
  const addParameter = (form) => {
    if (!form.parameters) {
      form.parameters = []
    }
    form.parameters.push({
      type: 'query',
      key: '',
      value: ''
    })
  }

  const fmtComponent = (component) => {
    form.value.component = component.replace(/\\/g, '/')
    form.value.name = toLowerCase(pathInfo['/src/' + component])
    form.value.path = form.value.name
  }

  // Delete parameter
  const deleteParameter = (parameters, index) => {
    parameters.splice(index, 1)
  }

  // Add controllable button
  const addBtn = (form) => {
    if (!form.menuBtn) {
      form.menuBtn = []
    }
    form.menuBtn.push({
      name: '',
      desc: ''
    })
  }
  // Delete controllable button
  const deleteBtn = async (btns, index) => {
    const btn = btns[index]
    if (btn.ID === 0) {
      btns.splice(index, 1)
      return
    }
    const res = await canRemoveAuthorityBtnApi({ id: btn.ID })
    if (res.code === 0) {
      btns.splice(index, 1)
    }
  }

  const form = ref({
    ID: 0,
    path: '',
    name: '',
    hidden: false,
    parentId: 0,
    component: '',
    meta: {
      activeName: '',
      title: '',
      icon: '',
      defaultMenu: false,
      closeTab: false,
      keepAlive: false
    },
    parameters: [],
    menuBtn: []
  })
  const changeName = () => {
    form.value.path = form.value.name
  }

  const handleClose = (done) => {
    initForm()
    done()
  }
  // Delete menu
  const deleteMenu = (ID) => {
    ElMessageBox.confirm(
      t('admin.superadmin.menu.messages.delete_confirm_detail'),
      t('admin.common.confirm'),
      {
        confirmButtonText: t('admin.common.confirm'),
        cancelButtonText: t('admin.common.cancel'),
        type: 'warning'
      }
    )
      .then(async () => {
        const res = await deleteBaseMenu({ ID })
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: t('admin.superadmin.menu.messages.delete_success')
          })

          getTableData()
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('admin.superadmin.menu.messages.delete_cancelled')
        })
      })
  }
  // Drawer form init
  const menuForm = ref(null)
  const checkFlag = ref(false)
  const initForm = () => {
    checkFlag.value = false
    menuForm.value.resetFields()
    form.value = {
      ID: 0,
      path: '',
      name: '',
      hidden: false,
      parentId: 0,
      component: '',
      meta: {
        title: '',
        icon: '',
        defaultMenu: false,
        closeTab: false,
        keepAlive: false
      }
    }
  }
  // Close drawer

  const dialogFormVisible = ref(false)
  const closeDialog = () => {
    initForm()
    dialogFormVisible.value = false
  }
  // Create/update menu
  const enterDialog = async () => {
    menuForm.value.validate(async (valid) => {
      if (valid) {
        let res
        if (isEdit.value) {
          res = await updateBaseMenu(form.value)
        } else {
          res = await addBaseMenu(form.value)
        }
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: isEdit.value ? t('admin.superadmin.menu.messages.edit_success') : t('admin.superadmin.menu.messages.add_success_assign_hint')
          })
          getTableData()
        }
        initForm()
        dialogFormVisible.value = false
      }
    })
  }

  const menuOption = ref([
    {
      ID: '0',
      title: t('admin.superadmin.menu.options.root_menu')
    }
  ])
  const setOptions = () => {
    menuOption.value = [
      {
        ID: 0,
        title: t('admin.superadmin.menu.options.root_dir')
      }
    ]
    setMenuOptions(tableData.value, menuOption.value, false)
  }
  const setMenuOptions = (menuData, optionsData, disabled) => {
    menuData &&
      menuData.forEach((item) => {
        if (item.children && item.children.length) {
          const option = {
            title: item.meta.title,
            ID: item.ID,
            disabled: disabled || item.ID === form.value.ID,
            children: []
          }
          setMenuOptions(
            item.children,
            option.children,
            disabled || item.ID === form.value.ID
          )
          optionsData.push(option)
        } else {
          const option = {
            title: item.meta.title,
            ID: item.ID,
            disabled: disabled || item.ID === form.value.ID
          }
          optionsData.push(option)
        }
      })
  }

  // Add menu; id=0 means add root menu
  const isEdit = ref(false)
  const dialogTitle = ref(t('admin.superadmin.menu.titles.add_menu'))
  const addMenu = (id) => {
    dialogTitle.value = t('admin.superadmin.menu.titles.add_menu')
    form.value.parentId = id
    isEdit.value = false
    setOptions()
    dialogFormVisible.value = true
  }
  // Edit menu
  const editMenu = async (id) => {
    dialogTitle.value = t('admin.superadmin.menu.titles.edit_menu')
    const res = await getBaseMenuById({ id })
    form.value = res.data.menu
    isEdit.value = true
    setOptions()
    dialogFormVisible.value = true
  }

  // Assign roles
  const assignRoleDrawerVisible = ref(false)
  const assignMenuRow = ref({})
  const authorityTreeData = ref([])
  const assignRoleLoading = ref(false)
  const assignRoleSubmitting = ref(false)
  const roleTreeRef = ref(null)
  const defaultRouterAuthorityIds = ref(new Set())

  const isRoleDisabled = (data) => {
    return defaultRouterAuthorityIds.value.has(data.authorityId)
  }

  const openAssignRoleDrawer = async (row) => {
    assignMenuRow.value = row
    defaultRouterAuthorityIds.value = new Set()
    assignRoleDrawerVisible.value = true
    assignRoleLoading.value = true
    // Load role tree and current assignment in parallel
    const [authRes, rolesRes] = await Promise.all([
      getAuthorityList(),
      getMenuRoles(row.ID)
    ])
    if (authRes.code === 0) {
      authorityTreeData.value = authRes.data
    }
    if (rolesRes.code === 0 && rolesRes.data) {
      if (rolesRes.data.defaultRouterAuthorityIds) {
        defaultRouterAuthorityIds.value = new Set(rolesRes.data.defaultRouterAuthorityIds)
      }
      nextTick(() => {
        roleTreeRef.value?.setCheckedKeys(rolesRes.data.authorityIds || [])
      })
    }
    assignRoleLoading.value = false
  }

  const confirmAssignRole = async () => {
    assignRoleSubmitting.value = true
    try {
      const checkedKeys = roleTreeRef.value?.getCheckedKeys(false) || []
      const halfCheckedKeys = roleTreeRef.value?.getHalfCheckedKeys() || []
      const authorityIds = [...checkedKeys, ...halfCheckedKeys]
      const res = await setMenuRoles({
        menuId: assignMenuRow.value.ID,
        authorityIds
      })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: t('admin.superadmin.menu.assign.success') })
        assignRoleDrawerVisible.value = false
      }
    } catch {
      ElMessage({ type: 'error', message: t('admin.superadmin.menu.assign.failed') })
    }
    assignRoleSubmitting.value = false
  }
</script>

<style scoped lang="scss">
  .warning {
    color: #dc143c;
  }
  .icon-column {
    display: flex;
    align-items: center;
    .el-icon {
      margin-right: 8px;
    }
  }


  
  .form-tip {
    margin-top: 8px;
    font-size: 12px;
    color: #909399;
    display: flex;
    align-items: center;
    gap: 8px;
    
    .el-icon {
      color: #409eff;
    }
  }
  
  .label-with-tooltip {
    display: flex;
    align-items: center;
    gap: 6px;
    
    .el-icon {
      color: #909399;
      cursor: help;
      
      &:hover {
        color: #409eff;
      }
    }
  }
  
  .parameter-table,
  .button-table {
    border: 1px solid #ebeef5;
    border-radius: 6px;
    
    :deep(.el-table__header) {
      background-color: #fafafa;
    }
    
    :deep(.el-table__body) {
      .el-table__row {
        &:hover {
          background-color: #f5f7fa;
        }
      }
    }
  }
</style>
