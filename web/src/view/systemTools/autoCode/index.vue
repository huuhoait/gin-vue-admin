<template>
  <div>
    <warning-bar
      href="https://www.bilibili.com/video/BV1kv4y1g7nT?p=3"
      :title="t('admin.system_tools.auto_code.dev_only_warning')"
    />
    <div class="gva-search-box" v-if="!isAdd">
      <div class="text-lg mb-2 text-gray-600">
        {{ t('admin.system_tools.auto_code.create_with_ai') }}
        <a
          class="text-blue-600 text-sm ml-4"
          href="https://plugin.gin-vue-admin.com/#/layout/userInfo/center"
          target="_blank"
          >{{ t('admin.system_tools.auto_code.get_ai_path') }}</a
        >
      </div>
      <div class="relative">
        <el-input
          v-model="prompt"
          type="textarea"
          :rows="5"
          :maxlength="2000"
          :placeholder="t('admin.system_tools.auto_code.ai_prompt_placeholder')"
          resize="none"
          @focus="handleFocus"
          @blur="handleBlur"
        />

        <div class="flex absolute right-28 bottom-2">
          <el-tooltip effect="light">
            <template #content>
              <div>
                {{ t('admin.system_tools.auto_code.free_go_to') }}
                <a
                  class="text-blue-600"
                  href="https://plugin.gin-vue-admin.com/#/layout/userInfo/center"
                  target="_blank"
              >{{ t('admin.system_tools.auto_code.plugin_market_profile') }}</a
              >
                {{ t('admin.system_tools.auto_code.apply_ai_path_hint') }}
              </div>
            </template>
            <el-button
                :disabled="form.onlyTemplate"
                type="primary"
                @click="eyeFunc()"
            >
              <el-icon size="18">
                <ai-gva />
              </el-icon>
              {{ t('admin.system_tools.auto_code.image_to_text') }}
            </el-button>
          </el-tooltip>
        </div>

        <div class="flex absolute right-2 bottom-2">
          <el-tooltip effect="light">
            <template #content>
              <div>
                {{ t('admin.system_tools.auto_code.free_go_to') }}
                <a
                  class="text-blue-600"
                  href="https://plugin.gin-vue-admin.com/#/layout/userInfo/center"
                  target="_blank"
                  >{{ t('admin.system_tools.auto_code.plugin_market_profile') }}</a
                >
                {{ t('admin.system_tools.auto_code.apply_ai_path_hint') }}
              </div>
            </template>
            <el-button
              :disabled="form.onlyTemplate"
              type="primary"
              @click="llmAutoFunc()"
            >
              <el-icon size="18">
                <ai-gva />
              </el-icon>
              {{ t('admin.common.generate') }}
            </el-button>
          </el-tooltip>
        </div>
      </div>
    </div>
    <!-- Fetch fields directly from DB -->
    <div class="gva-search-box" v-if="!isAdd">
      <div class="text-lg mb-2 text-gray-600">
        {{ t('admin.system_tools.auto_code.create_from_database') }}
      </div>
      <el-form
        ref="getTableForm"
        :inline="true"
        :model="dbform"
        label-width="120px"
      >
        <el-row class="w-full">
          <el-col :span="6">
            <el-form-item
              :label="t('admin.system_tools.auto_code.business_db')"
              prop="selectDBtype"
              class="w-full"
            >
              <template #label>
                <el-tooltip
                  :content="t('admin.system_tools.auto_code.business_db_tooltip')"
                  placement="bottom"
                  effect="light"
                >
                  <div>
                    {{ t('admin.system_tools.auto_code.business_db') }}
                    <el-icon><QuestionFilled /></el-icon>
                  </div>
                </el-tooltip>
              </template>
              <el-select
                v-model="dbform.businessDB"
                clearable
                :placeholder="t('admin.system_tools.auto_code.select_business_db')"
                @change="getDbFunc"
                class="w-full"
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
                    <span
                      style="float: right; color: #8492a6; font-size: 13px"
                      >{{ item.dbName }}</span
                    >
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item
              :label="t('admin.system_tools.auto_code.database')"
              prop="structName"
              class="w-full"
            >
              <el-select
                v-model="dbform.dbName"
                clearable
                filterable
                :placeholder="t('admin.system_tools.auto_code.select_database')"
                class="w-full"
                @change="getTableFunc"
              >
                <el-option
                  v-for="item in dbOptions"
                  :key="item.database"
                  :label="item.database"
                  :value="item.database"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item
              :label="t('admin.system_tools.auto_code.table')"
              prop="structName"
              class="w-full"
            >
              <el-select
                v-model="dbform.tableName"
                :disabled="!dbform.dbName"
                class="w-full"
                filterable
                :placeholder="t('admin.system_tools.auto_code.select_table')"
              >
                <el-option
                  v-for="item in tableOptions"
                  :key="item.tableName"
                  :label="item.tableName"
                  :value="item.tableName"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item class="w-full">
              <div class="flex justify-end w-full">
                <el-button type="primary" @click="getColumnFunc">
                  {{ t('admin.system_tools.auto_code.use_this_table') }}
                </el-button>
              </div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>
    <div class="gva-search-box">
      <!-- Initial version auto-code tool -->
      <div class="text-lg mb-2 text-gray-600">
        {{ t('admin.system_tools.auto_code.auto_code_structure') }}
      </div>
      <el-form
        :disabled="isAdd"
        ref="autoCodeForm"
        :rules="rules"
        :model="form"
        label-width="120px"
        :inline="true"
      >
        <el-row class="w-full">
          <el-col :span="6">
            <el-form-item
              :label="t('admin.system_tools.auto_code.struct_name')"
              prop="structName"
              class="w-full"
            >
              <div class="flex gap-2 w-full">
                <el-input
                  v-model="form.structName"
                  :placeholder="t('admin.system_tools.auto_code.struct_name_placeholder')"
                />
                <el-button
                  :disabled="form.onlyTemplate"
                  type="primary"
                  @click="llmAutoFunc(true)"
                >
                  <el-icon size="18">
                    <ai-gva />
                  </el-icon>
                  {{ t('admin.common.generate') }}
                </el-button>
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="abbreviation" prop="abbreviation" class="w-full">
              <template #label>
                <el-tooltip
                  content="Used as request object name and route group."
                  placement="bottom"
                  effect="light"
                >
                  <div>
                    {{ t('admin.system_tools.auto_code.abbreviation') }}
                    <el-icon><QuestionFilled /></el-icon>
                  </div>
                </el-tooltip>
              </template>
              <el-input
                v-model="form.abbreviation"
                :placeholder="t('admin.system_tools.auto_code.abbreviation_placeholder')"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item
              :label="t('admin.system_tools.auto_code.display_name')"
              prop="description"
              class="w-full"
            >
              <el-input
                v-model="form.description"
                :placeholder="t('admin.system_tools.auto_code.display_name_placeholder')"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="Table name" prop="tableName" class="w-full">
              <el-input
                v-model="form.tableName"
                placeholder="Specify table name (optional)"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row class="w-full">
          <el-col :span="6">
            <el-form-item prop="packageName" class="w-full">
              <template #label>
                <el-tooltip
                  :content="t('admin.system_tools.auto_code.file_name_tooltip')"
                  placement="bottom"
                  effect="light"
                >
                  <div>
                    {{ t('admin.system_tools.auto_code.file_name') }}
                    <el-icon><QuestionFilled /></el-icon>
                  </div>
                </el-tooltip>
              </template>
              <el-input
                v-model="form.packageName"
                :placeholder="t('admin.system_tools.auto_code.file_name_placeholder')"
                @blur="toLowerCaseFunc(form, 'packageName')"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item
              :label="t('admin.system_tools.auto_code.template')"
              prop="package"
              class="w-full relative"
            >
              <el-select v-model="form.package" class="w-full pr-12" filterable>
                <el-option
                  v-for="item in pkgs"
                  :key="item.ID"
                  :value="item.packageName"
                  :label="item.packageName"
                />
              </el-select>
              <span class="absolute right-0">
                <el-icon
                  class="cursor-pointer ml-2 text-gray-600"
                  @click="getPkgs"
                >
                  <refresh />
                </el-icon>
                <el-icon
                  class="cursor-pointer ml-2 text-gray-600"
                  @click="goPkgs"
                >
                  <document-add />
                </el-icon>
              </span>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="Business DB" prop="businessDB" class="w-full">
              <template #label>
                <el-tooltip
                  content="Note: Configure multiple databases in db-list first. If empty, it uses the default GVA DB (global.GVA_DB). If set, it generates code for the specified DB (global.MustGetGlobalDBByDBName(dbname))."
                  placement="bottom"
                  effect="light"
                >
                  <div>
                    Business DB <el-icon><QuestionFilled /></el-icon>
                  </div>
                </el-tooltip>
              </template>
              <el-select
                v-model="form.businessDB"
                clearable
                :placeholder="t('admin.system_tools.auto_code.select_business_db')"
                class="w-full"
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
                    <span
                      style="float: right; color: #8492a6; font-size: 13px"
                      >{{ item.dbName }}</span
                    >
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>
    <div class="gva-search-box">
      <el-collapse class="no-border-collapse">
        <el-collapse-item>
          <template #title>
            <div class="text-lg text-gray-600 font-normal">
              {{ t('admin.system_tools.auto_code.expert_mode') }}
            </div>
          </template>
          <template #icon="{ isActive }">
          <span class="text-lg ml-auto mr-4 font-normal">
            {{ isActive ? t('admin.system_tools.auto_code.collapse') : t('admin.system_tools.auto_code.expand') }}
          </span>
          </template>
          <div class="p-4">
            <!-- Basic settings -->
            <div class="border-b border-gray-200 last:border-0">
              <h3 class="text-lg font-medium mb-4 text-gray-700">
                {{ t('admin.system_tools.auto_code.basic_settings') }}
              </h3>
              <el-row :gutter="20">
                <el-col :span="3">
                  <el-tooltip
                      content="Note: Automatically embeds global.Model fields (primary key and soft delete)."
                      placement="top"
                      effect="light"
                  >
                    <el-form-item label="Use GVA model">
                      <el-checkbox v-model="form.gvaModel" @change="useGva" />
                    </el-form-item>
                  </el-tooltip>
                </el-col>
                <el-col :span="3">
                  <el-tooltip
                      content="Note: Generates page button-permission config. If you don't assign buttons in role management, they won't be visible."
                      placement="top"
                      effect="light"
                  >
                    <el-form-item label="Create button permissions">
                      <el-checkbox :disabled="!form.generateWeb" v-model="form.autoCreateBtnAuth" />
                    </el-form-item>
                  </el-tooltip>
                </el-col>
                <el-col :span="3">
                  <el-form-item :label="t('admin.system_tools.auto_code.generate_frontend')">
                    <el-checkbox v-model="form.generateWeb" />
                  </el-form-item>
                </el-col>
                <el-col :span="3">
                  <el-form-item :label="t('admin.system_tools.auto_code.generate_backend')">
                    <el-checkbox disabled v-model="form.generateServer" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>

            <!-- Automation settings -->
            <div class="border-b border-gray-200 last:border-0">
              <h3 class="text-lg font-medium mb-4 text-gray-700">
                {{ t('admin.system_tools.auto_code.automation_settings') }}
              </h3>
              <el-row :gutter="20">
                <el-col :span="3">
                  <el-tooltip
                      content="Note: Register generated APIs into the database."
                      placement="top"
                      effect="light"
                  >
                    <el-form-item label="Auto-create APIs">
                      <el-checkbox  :disabled="!form.generateServer" v-model="form.autoCreateApiToSql" />
                    </el-form-item>
                  </el-tooltip>
                </el-col>
                <el-col :span="3">
                  <el-tooltip
                      content="Note: Register generated menus into the database."
                      placement="top"
                      effect="light"
                  >
                    <el-form-item label="Auto-create menus">
                      <el-checkbox :disabled="!form.generateWeb" v-model="form.autoCreateMenuToSql" />
                    </el-form-item>
                  </el-tooltip>
                </el-col>
                <el-col :span="3">
                  <el-tooltip
                      content="Note: Auto-migrate database table structure. Disable if not needed."
                      placement="top"
                      effect="light"
                  >
                    <el-form-item label="Auto-migrate tables">
                      <el-checkbox  :disabled="!form.generateServer" v-model="form.autoMigrate" />
                    </el-form-item>
                  </el-tooltip>
                </el-col>
              </el-row>
            </div>

            <!-- Advanced settings -->
            <div class="border-b border-gray-200 last:border-0">
              <h3 class="text-lg font-medium mb-4 text-gray-700">
                {{ t('admin.system_tools.auto_code.advanced_settings') }}
              </h3>
              <el-row :gutter="20">
                <el-col :span="3">
                  <el-tooltip
                      content="Note: Adds created_by/updated_by/deleted_by for resource access control."
                      placement="top"
                      effect="light"
                  >
                    <el-form-item label="Add resource markers">
                      <el-checkbox v-model="form.autoCreateResource" />
                    </el-form-item>
                  </el-tooltip>
                </el-col>
                <el-col :span="3">
                  <el-tooltip
                      content="Note: Base template generates no structs or CRUD. It only configures properties like enter to support custom non-CRUD logic."
                      placement="top"
                      effect="light"
                  >
                    <el-form-item label="Base template">
                      <el-checkbox v-model="form.onlyTemplate" />
                    </el-form-item>
                  </el-tooltip>
                </el-col>
              </el-row>
            </div>

            <!-- Tree structure settings -->
            <div class="last:pb-0">
              <h3 class="text-lg font-medium mb-4 text-gray-700">
                {{ t('admin.system_tools.auto_code.tree_structure_settings') }}
              </h3>
              <el-row :gutter="20" align="middle">
                <el-col :span="24">
                    <el-form-item label="Tree structure">
                      <div class="flex items-center gap-4">
                        <el-tooltip
                            content="Note: Creates parentID for parent/child relations. Only supports int primary keys."
                            placement="top"
                            effect="light"
                        >
                          <el-checkbox v-model="form.isTree" />
                        </el-tooltip>
                        <el-input
                            v-model="form.treeJson"
                            :disabled="!form.isTree"
                            placeholder="Frontend display JSON property"
                            class="flex-1"
                        />
                      </div>
                    </el-form-item>
                </el-col>
              </el-row>
            </div>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
    <!-- Field list -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          @click="editAndAddField()"
          :disabled="form.onlyTemplate"
        >
          Add field
        </el-button>
      </div>
      <div class="draggable">
        <el-table :data="form.fields" row-key="fieldName">
          <el-table-column
            v-if="!isAdd"
            fixed="left"
            align="left"
            type="index"
            width="60"
          >
            <template #default>
              <el-icon class="cursor-grab drag-column">
                <MoreFilled />
              </el-icon>
            </template>
          </el-table-column>
          <el-table-column
            fixed="left"
            align="left"
            type="index"
            label="#"
            width="60"
          />
          <el-table-column
            fixed="left"
            align="left"
            type="index"
            label="PK"
            width="60"
          >
            <template #default="{ row }">
              <el-checkbox :disabled="row.disabled" v-model="row.primaryKey" />
            </template>
          </el-table-column>
          <el-table-column
            fixed="left"
            align="left"
            prop="fieldName"
            label="Field name"
            width="160"
          >
            <template #default="{ row }">
              <el-input disabled v-model="row.fieldName" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="fieldDesc"
            label="Label"
            width="160"
          >
            <template #default="{ row }">
              <el-input :disabled="row.disabled" v-model="row.fieldDesc" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="defaultValue"
            label="Default"
            width="160"
          >
            <template #default="{ row }">
              <el-input :disabled="row.disabled" v-model="row.defaultValue" />
            </template>
          </el-table-column>
          <el-table-column align="left" prop="require" label="Required">
            <template #default="{ row }">
              <el-checkbox :disabled="row.disabled" v-model="row.require" />
            </template>
          </el-table-column>
          <el-table-column align="left" prop="sort" label="Sort">
            <template #default="{ row }">
              <el-checkbox :disabled="row.disabled" v-model="row.sort" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="form"
            width="100"
            label="Create/Edit"
          >
            <template #default="{ row }">
              <el-checkbox :disabled="row.disabled" v-model="row.form" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="table"
            :label="t('admin.system_tools.auto_code.column_table')"
          >
            <template #default="{ row }">
              <el-checkbox :disabled="row.disabled" v-model="row.table" />
            </template>
          </el-table-column>
          <el-table-column align="left" prop="desc" label="Detail">
            <template #default="{ row }">
              <el-checkbox :disabled="row.disabled" v-model="row.desc" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="excel"
            width="100"
            label="Import/Export"
            v-if="!isAdd"
          >
            <template #default="{ row }">
              <el-checkbox v-model="row.excel" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="fieldJson"
            width="160px"
            label="Field JSON"
          >
            <template #default="{ row }">
              <el-input :disabled="row.disabled" v-model="row.fieldJson" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="fieldType"
            label="Field type"
            width="160"
          >
            <template #default="{ row }">
              <el-select
                v-model="row.fieldType"
                style="width: 100%"
                placeholder="Select field type"
                :disabled="row.disabled"
                clearable
              >
                <el-option
                  v-for="item in typeOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="fieldIndexType"
            label="Index type"
            width="160"
          >
            <template #default="{ row }">
              <el-select
                v-model="row.fieldIndexType"
                style="width: 100%"
                placeholder="Select index type"
                :disabled="row.disabled"
                clearable
              >
                <el-option
                  v-for="item in typeIndexOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="dataTypeLong"
            label="Length / enum"
            width="160"
          >
            <template #default="{ row }">
              <el-input :disabled="row.disabled" v-model="row.dataTypeLong" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="columnName"
            label="DB column"
            width="160"
          >
            <template #default="{ row }">
              <el-input :disabled="row.disabled" v-model="row.columnName" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="comment"
            label="DB comment"
            width="160"
          >
            <template #default="{ row }">
              <el-input :disabled="row.disabled" v-model="row.comment" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="fieldSearchType"
            label="Search"
            width="130"
          >
            <template #default="{ row }">
              <el-select
                v-model="row.fieldSearchType"
                style="width: 100%"
                placeholder="Select search condition"
                clearable
                :disabled="row.fieldType === 'json' || row.disabled"
              >
                <el-option
                  v-for="item in typeSearchOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                  :disabled="canSelect(row.fieldType,item.value)"
                />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column align="left" label="Actions" width="300" fixed="right">
            <template #default="scope">
              <el-button
                v-if="!scope.row.disabled"
                type="primary"
                link
                icon="edit"
                @click="editAndAddField(scope.row)"
              >
                Advanced edit
              </el-button>
              <el-button
                v-if="!scope.row.disabled"
                type="primary"
                link
                icon="delete"
                @click="deleteField(scope.$index)"
              >
                Delete
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <!-- Actions -->
      <div class="gva-btn-list justify-end mt-4">
        <el-button type="primary" :disabled="isAdd" @click="exportJson()">
          {{ t('admin.common.export') }} JSON
        </el-button>
        <el-upload
          class="flex items-center"
          :before-upload="importJson"
          :show-file-list="false"
          :headers="{'x-token': token}"
          accept=".json"
        >
          <el-button type="primary" class="mx-2" :disabled="isAdd"
            >{{ t('admin.common.import') }} JSON</el-button
          >
        </el-upload>
        <el-button type="primary" :disabled="isAdd" @click="clearCatch()">
          {{ t('admin.system_tools.auto_code.clear_stash') }}
        </el-button>
        <el-button type="primary" :disabled="isAdd" @click="catchData()">
          {{ t('admin.system_tools.auto_code.stash') }}
        </el-button>
        <el-button type="primary" :disabled="isAdd" @click="enterForm(false)">
          {{ t('admin.system_tools.auto_code.generate_code') }}
        </el-button>
        <el-button type="primary" @click="enterForm(true)">
          {{
            isAdd
              ? t('admin.system_tools.auto_code.view_code')
              : t('admin.system_tools.auto_code.preview_code')
          }}
        </el-button>
      </div>
    </div>
    <!-- Field drawer -->
    <el-drawer v-model="dialogFlag" size="70%" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.system_tools.auto_code.field_editor') }}</span>
          <div>
            <el-button @click="closeDialog">{{ t('admin.common.cancel') }}</el-button>
            <el-button type="primary" @click="enterDialog">{{ t('admin.common.confirm') }}</el-button>
          </div>
        </div>
      </template>

      <FieldDialog
        v-if="dialogFlag"
        ref="fieldDialogNode"
        :dialog-middle="dialogMiddle"
        :type-options="typeOptions"
        :type-search-options="typeSearchOptions"
        :type-index-options="typeIndexOptions"
      />
    </el-drawer>

    <el-drawer v-model="previewFlag" size="80%" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ t('admin.system_tools.auto_code.toolbar') }}</span>
          <div>
            <el-button type="primary" @click="selectText">{{ t('admin.system_tools.auto_code.select_all') }}</el-button>
            <el-button type="primary" @click="copy">{{ t('admin.system_tools.auto_code.copy') }}</el-button>
          </div>
        </div>
      </template>
      <PreviewCodeDialog
        v-if="previewFlag"
        :is-add="isAdd"
        ref="previewNode"
        :preview-code="preViewCode"
      />
    </el-drawer>
  </div>
</template>

<script setup>
  import FieldDialog from '@/view/systemTools/autoCode/component/fieldDialog.vue'
  import PreviewCodeDialog from '@/view/systemTools/autoCode/component/previewCodeDialog.vue'
  import {
    toUpperCase,
    toHump,
    toSQLLine,
    toLowerCase
  } from '@/utils/stringFun'
  import {
    createTemp,
    getDB,
    getTable,
    getColumn,
    preview,
    getMeta,
    getPackageApi,
    llmAuto
  } from '@/api/autoCode'
  import { getDict } from '@/utils/dictionary'
  import { ref, watch, toRaw, onMounted, nextTick } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import Sortable from 'sortablejs'
  import { useUserStore } from "@/pinia";

  const userStore = useUserStore()
  const { t } = useI18n()

  const token = userStore.token

  const handleFocus = () => {
    document.addEventListener('keydown', handleKeydown);
    document.addEventListener('paste', handlePaste);
  }

  const handleBlur = () => {
    document.removeEventListener('keydown', handleKeydown);
    document.removeEventListener('paste', handlePaste);
  }

  const handleKeydown = (event) => {
    if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
      llmAutoFunc()
    }
  }

  const handlePaste = (event) => {
    const items = event.clipboardData.items;
    for (let i = 0; i < items.length; i++) {
      if (items[i].type.indexOf('image') !== -1) {
        const file = items[i].getAsFile();
        const reader = new FileReader();
        reader.onload =async (e) => {
          const base64String = e.target.result;
          const res = await llmAuto({ _file_path: base64String,mode:"eye" })
          if (res.code === 0) {
            prompt.value = res.data.text
            llmAutoFunc()
          }
        };
        reader.readAsDataURL(file);
      }
    }
  };

  const getOnlyNumber = () => {
    let randomNumber = ''
    while (randomNumber.length < 16) {
      randomNumber += Math.random().toString(16).substring(2)
    }
    return randomNumber.substring(0, 16)
  }

  const prompt = ref('')

  const eyeFunc = async () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';

    input.onchange = (event) => {
      const file = event.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = async (e) => {
          const base64String = e.target.result;

          const res = await llmAuto({ _file_path: base64String,mode:'eye' })
          if (res.code === 0) {
            prompt.value = res.data.text
            llmAutoFunc()
          }
        };
        reader.readAsDataURL(file);
      }
    };

    input.click();
  }


  const llmAutoFunc = async (flag) => {
    if (flag && !form.value.structName) {
      ElMessage.error(t('admin.system_tools.auto_code.msg_enter_struct_name'))
      return
    }
    if (!flag && !prompt.value) {
      ElMessage.error(t('admin.system_tools.auto_code.msg_enter_description'))
      return
    }

    if (form.value.fields.length > 0) {
      const res = await ElMessageBox.confirm(
        t('admin.system_tools.auto_code.confirm_ai_overwrite_message'),
        t('admin.common.confirm'),
        {
          confirmButtonText: t('admin.common.confirm'),
          cancelButtonText: t('admin.common.cancel'),
          type: 'warning'
        }
      )
      if (res !== 'confirm') {
        return
      }
    }

    const res = await llmAuto({
      prompt: flag ? 'Struct name: ' + form.value.structName : prompt.value,
      mode: "ai"
    })
    if (res.code === 0) {
      form.value.fields = []
      const json = JSON.parse(res.data.text)
      json.fields?.forEach((item) => {
        item.fieldName = toUpperCase(item.fieldName)
      })

      for (let key in json) {
        form.value[key] = json[key]
      }

      form.value.generateServer = true
      form.value.generateWeb = true

    }
  }

  const isAdd = ref(false)

  // Row drag-and-drop
  const rowDrop = () => {
    // Parent container for draggable rows
    const tbody = document.querySelector(
      '.draggable .el-table__body-wrapper tbody'
    )
    Sortable.create(tbody, {
      // Draggable row selector
      draggable: '.draggable .el-table__row',
      handle: '.drag-column',
      onEnd: async ({ newIndex, oldIndex }) => {
        await nextTick()
        const currRow = form.value.fields.splice(oldIndex, 1)[0]
        form.value.fields.splice(newIndex, 0, currRow)
      }
    })
  }

  onMounted(() => {
    rowDrop()
  })

  defineOptions({
    name: 'AutoCode'
  })
  const gormModelList = ['id', 'created_at', 'updated_at', 'deleted_at']

  const dataModelList = ['created_by', 'updated_by', 'deleted_by']

  const typeOptions = ref([
    {
      label: 'String',
      value: 'string'
    },
    {
      label: 'Rich text',
      value: 'richtext'
    },
    {
      label: 'Integer',
      value: 'int'
    },
    {
      label: 'Boolean',
      value: 'bool'
    },
    {
      label: 'Float',
      value: 'float64'
    },
    {
      label: 'Time',
      value: 'time.Time'
    },
    {
      label: 'Enum',
      value: 'enum'
    },
    {
      label: 'Single image (string)',
      value: 'picture'
    },
    {
      label: 'Multiple images (JSON string)',
      value: 'pictures'
    },
    {
      label: 'Video (string)',
      value: 'video'
    },
    {
      label: 'File (JSON string)',
      value: 'file'
    },
    {
      label: 'JSON',
      value: 'json'
    },
    {
      label: 'Array',
      value: 'array'
    }
  ])

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
    }
  ])

  const typeIndexOptions = ref([
    {
      label: 'index',
      value: 'index'
    },
    {
      label: 'uniqueIndex',
      value: 'uniqueIndex'
    }
  ])

  const fieldTemplate = {
    fieldName: '',
    fieldDesc: '',
    fieldType: '',
    dataType: '',
    fieldJson: '',
    columnName: '',
    dataTypeLong: '',
    comment: '',
    defaultValue: '',
    require: false,
    sort: false,
    form: true,
    desc: true,
    table: true,
    excel: false,
    errorText: '',
    primaryKey: false,
    clearable: true,
    fieldSearchType: '',
    fieldIndexType: '',
    dictType: '',
    dataSource: {
      dbName: '',
      association: 1,
      table: '',
      label: '',
      value: '',
      hasDeletedAt: false
    }
  }
  const route = useRoute()
  const router = useRouter()
  const preViewCode = ref({})
  const dbform = ref({
    businessDB: '',
    dbName: '',
    tableName: ''
  })
  const tableOptions = ref([])
  const addFlag = ref('')
  const fdMap = ref({})
  const form = ref({
    structName: '',
    tableName: '',
    packageName: '',
    package: '',
    abbreviation: '',
    description: '',
    businessDB: '',
    autoCreateApiToSql: true,
    autoCreateMenuToSql: true,
    autoCreateBtnAuth: false,
    autoMigrate: true,
    gvaModel: true,
    autoCreateResource: false,
    onlyTemplate: false,
    isTree: false,
    generateWeb:true,
    generateServer:true,
    treeJson: "",
    fields: []
  })
  const rules = ref({
    structName: [
      { required: true, message: 'Please enter struct name', trigger: 'blur' }
    ],
    abbreviation: [
      { required: true, message: 'Please enter struct abbreviation', trigger: 'blur' }
    ],
    description: [
      { required: true, message: 'Please enter struct description', trigger: 'blur' }
    ],
    packageName: [
      {
        required: true,
        message: 'File name: sysXxxxXxxx',
        trigger: 'blur'
      }
    ],
    package: [{ required: true, message: 'Please select a template', trigger: 'blur' }]
  })
  const dialogMiddle = ref({})
  const bk = ref({})
  const dialogFlag = ref(false)
  const previewFlag = ref(false)

  const useGva = (e) => {
    if (e && form.value.fields.length) {
      ElMessageBox.confirm(
        'Enabling the default GVA model will auto-add ID/CreatedAt/UpdatedAt/DeletedAt fields. This will remove any duplicate fields you created below. Continue?',
        'Warning',
        {
          confirmButtonText: 'Continue',
          cancelButtonText: 'Cancel',
          type: 'warning'
        }
      )
        .then(() => {
          form.value.fields = form.value.fields.filter(
            (item) =>
              !gormModelList.some((gormfd) => gormfd === item.columnName)
          )
        })
        .catch(() => {
          form.value.gvaModel = false
        })
    }
  }

  const toLowerCaseFunc = (form, key) => {
    form[key] = toLowerCase(form[key])
  }
  const previewNode = ref(null)
  const selectText = () => {
    previewNode.value.selectText()
  }
  const copy = () => {
    previewNode.value.copy()
  }
  const editAndAddField = (item) => {
    dialogFlag.value = true
    if (item) {
      addFlag.value = 'edit'
      if (!item.dataSource) {
        item.dataSource = {
          dbName: '',
          association: 1,
          table: '',
          label: '',
          value: '',
          hasDeletedAt: false
        }
      }
      bk.value = JSON.parse(JSON.stringify(item))
      dialogMiddle.value = item
    } else {
      addFlag.value = 'add'
      fieldTemplate.onlyNumber = getOnlyNumber()
      dialogMiddle.value = JSON.parse(JSON.stringify(fieldTemplate))
    }
  }

  const fieldDialogNode = ref(null)
  const enterDialog = () => {
    fieldDialogNode.value.fieldDialogForm.validate((valid) => {
      if (valid) {
        dialogMiddle.value.fieldName = toUpperCase(dialogMiddle.value.fieldName)
        if (addFlag.value === 'add') {
          form.value.fields.push(dialogMiddle.value)
        }
        dialogFlag.value = false
      } else {
        return false
      }
    })
  }
  const closeDialog = () => {
    if (addFlag.value === 'edit') {
      dialogMiddle.value = bk.value
    }
    dialogFlag.value = false
  }
  const deleteField = (index) => {
    ElMessageBox.confirm('Delete this field?', 'Confirmation', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning'
    }).then(async () => {
      form.value.fields.splice(index, 1)
    })
  }
  const autoCodeForm = ref(null)
  const enterForm = async (isPreview) => {
    if (form.value.isTree && !form.value.treeJson){
      ElMessage({
        type: 'error',
        message: 'Please fill in the frontend display JSON property for tree structure'
      })
      return false
    }
    if(!form.value.generateWeb && !form.value.generateServer){
      ElMessage({
        type: 'error',
        message: 'Please select at least one generation target'
      })
      return false
    }
    if (!form.value.onlyTemplate) {
      if (form.value.fields.length <= 0) {
        ElMessage({
          type: 'error',
          message: 'Please add at least one field'
        })
        return false
      }

      if (
        !form.value.gvaModel &&
        form.value.fields.every((item) => !item.primaryKey)
      ) {
        ElMessage({
          type: 'error',
          message: 'You need at least one primary key to generate code reliably'
        })
        return false
      }

      if (
        form.value.fields.some(
          (item) => item.fieldName === form.value.structName
        )
      ) {
        ElMessage({
          type: 'error',
          message: 'A field name conflicts with the struct name'
        })
        return false
      }

      if (
        form.value.fields.some((item) => item.fieldJson === form.value.package)
      ) {
        ElMessage({
          type: 'error',
          message: 'A field JSON key conflicts with the template name'
        })
        return false
      }

      if (form.value.fields.some((item) => !item.fieldType)) {
        ElMessage({
          type: 'error',
          message: 'Please select field types for all fields before submitting'
        })
        return false
      }

      if (form.value.package === form.value.abbreviation) {
        ElMessage({
          type: 'error',
          message: 'Template and struct abbreviation must not be the same'
        })
        return false
      }
    }

    autoCodeForm.value.validate(async (valid) => {
      if (valid) {
        for (const key in form.value) {
          if (typeof form.value[key] === 'string') {
            form.value[key] = form.value[key].trim()
          }
        }
        form.value.structName = toUpperCase(form.value.structName)
        form.value.tableName = form.value.tableName.replace(' ', '')
        if (!form.value.tableName) {
          form.value.tableName = toSQLLine(toLowerCase(form.value.structName))
        }
        if (form.value.structName === form.value.abbreviation) {
          ElMessage({
            type: 'error',
            message: 'Struct name and abbreviation must not be the same'
          })
          return false
        }
        form.value.humpPackageName = toSQLLine(form.value.packageName)

        form.value.fields?.forEach((item) => {
          item.fieldName = toUpperCase(item.fieldName)
          if (item.fieldType === 'enum') {
            // Ensure enum values are single-quoted
            item.dataTypeLong = item.dataTypeLong.replace(/[\[\]{}()]/g, '')
            const arr = item.dataTypeLong.split(',')
            arr.forEach((ele, index) => {
              if (ele.indexOf("'") === -1) {
                arr[index] = `'${ele}'`
              }
            })
            item.dataTypeLong = arr.join(',')
          }
        })

        delete form.value.primaryField
        if (isPreview) {
          const res = await preview({
            ...form.value,
            isAdd: !!isAdd.value,
            fields: form.value.fields.filter((item) => !item.disabled)
          })
          if(res.code !== 0){
            return
          }
          preViewCode.value = res.data.autoCode
          previewFlag.value = true
        } else {
          const res = await createTemp(form.value)
          if (res.code !== 0) {
            return
          }
          ElMessage({
            type: 'success',
            message: 'Auto-code created and moved successfully'
          })
          clearCatch()
        }
      }
    })
  }

  const dbList = ref([])
  const dbOptions = ref([])

  const getDbFunc = async () => {
    dbform.value.dbName = ''
    dbform.value.tableName = ''
    const res = await getDB({ businessDB: dbform.value.businessDB })
    if (res.code === 0) {
      dbOptions.value = res.data.dbs
      dbList.value = res.data.dbList
    }
  }
  const getTableFunc = async () => {
    const res = await getTable({
      businessDB: dbform.value.businessDB,
      dbName: dbform.value.dbName
    })
    if (res.code === 0) {
      tableOptions.value = res.data.tables
    }
    dbform.value.tableName = ''
  }

  const getColumnFunc = async () => {
    const res = await getColumn(dbform.value)
    if (res.code === 0) {
      let dbtype = ''
      if (dbform.value.businessDB !== '') {
        const dbtmp = dbList.value.find(
          (item) => item.aliasName === dbform.value.businessDB
        )
        const dbraw = toRaw(dbtmp)
        dbtype = dbraw.dbtype
      }
      form.value.gvaModel = false
      const tbHump = toHump(dbform.value.tableName)
      form.value.structName = toUpperCase(tbHump)
      form.value.tableName = dbform.value.tableName
      form.value.packageName = toLowerCase(tbHump)
      form.value.abbreviation = toLowerCase(tbHump)
      form.value.description = `${tbHump} table`
      form.value.autoCreateApiToSql = true
      form.value.generateServer = true
      form.value.generateWeb = true
      form.value.fields = []
      res.data.columns &&
        res.data.columns.forEach((item) => {
          if (needAppend(item)) {
            const fbHump = toHump(item.columnName)
            form.value.fields.push({
              onlyNumber: getOnlyNumber(),
              fieldName: toUpperCase(fbHump),
              fieldDesc: item.columnComment || `${fbHump} field`,
              fieldType: fdMap.value[item.dataType],
              dataType: item.dataType,
              fieldJson: fbHump,
              primaryKey: item.primaryKey,
              dataTypeLong:
                item.dataTypeLong && item.dataTypeLong.split(',')[0],
              columnName:
                dbtype === 'oracle'
                  ? item.columnName.toUpperCase()
                  : item.columnName,
              comment: item.columnComment,
              require: false,
              errorText: '',
              clearable: true,
              fieldSearchType: '',
              fieldIndexType: '',
              dictType: '',
              form: true,
              table: true,
              excel: false,
              desc: true,
              dataSource: {
                dbName: '',
                association: 1,
                table: '',
                label: '',
                value: '',
                hasDeletedAt: false
              }
            })
          }
        })
    }
  }

  const needAppend = (item) => {
    let isAppend = true
    if (
      form.value.gvaModel &&
      gormModelList.some((gormfd) => gormfd === item.columnName)
    ) {
      isAppend = false
    }
    if (
      form.value.autoCreateResource &&
      dataModelList.some((datafd) => datafd === item.columnName)
    ) {
      isAppend = false
    }
    return isAppend
  }

  const setFdMap = async () => {
    const fdTypes = ['string', 'int', 'bool', 'float64', 'time.Time']
    fdTypes.forEach(async (fdtype) => {
      const res = await getDict(fdtype)
      res &&
        res.forEach((item) => {
          fdMap.value[item.label] = fdtype
        })
    })
  }
  const getAutoCodeJson = async (id) => {
    const res = await getMeta({ id: Number(id) })
    if (res.code === 0) {
      const add = route.query.isAdd
      isAdd.value = add
      form.value = JSON.parse(res.data.meta)
      if (isAdd.value) {
        form.value.fields.forEach((item) => {
          item.disabled = true
        })
      }
    }
  }

  const pkgs = ref([])
  const getPkgs = async () => {
    const res = await getPackageApi()
    if (res.code === 0) {
      pkgs.value = res.data.pkgs
    }
  }

  const goPkgs = () => {
    router.push({ name: 'autoPkg' })
  }

  const init = () => {
    getDbFunc()
    setFdMap()
    getPkgs()
    const id = route.params.id
    if (id) {
      getAutoCodeJson(id)
    }
  }
  init()

  watch(()=>form.value.generateServer,()=>{
    if(!form.value.generateServer){
      form.value.autoCreateApiToSql = false
      form.value.autoMigrate = false
    }
  })

  watch(()=>form.value.generateWeb,()=>{
    if(!form.value.generateWeb){
      form.value.autoCreateMenuToSql = false
      form.value.autoCreateBtnAuth = false
    }
  })

  const catchData = () => {
    window.sessionStorage.setItem('autoCode', JSON.stringify(form.value))
    ElMessage.success('Stashed')
  }

  const getCatch = () => {
    const data = window.sessionStorage.getItem('autoCode')
    if (data) {
      form.value = JSON.parse(data)
    }
  }

  const clearCatch = async () => {
    form.value = {
      structName: '',
      tableName: '',
      packageName: '',
      package: '',
      abbreviation: '',
      description: '',
      businessDB: '',
      autoCreateApiToSql: true,
      autoCreateMenuToSql: true,
      autoCreateBtnAuth: false,
      autoMigrate: true,
      gvaModel: true,
      autoCreateResource: false,
      onlyTemplate: false,
      isTree: false,
      treeJson: "",
      fields: []
    }
    await nextTick()
    window.sessionStorage.removeItem('autoCode')
  }

  getCatch()

  const exportJson = () => {
    const dataStr = JSON.stringify(form.value, null, 2)
    const blob = new Blob([dataStr], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'form_data.json'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  const importJson = (file) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      try {
        form.value = JSON.parse(e.target.result)
        form.value.generateServer = true
        form.value.generateWeb = true
        ElMessage.success('JSON imported')
      } catch (_) {
        ElMessage.error('Invalid JSON file')
      }
    }
    reader.readAsText(file)
    return false
  }

  watch(
    () => form.value.onlyTemplate,
    (val) => {
      if (val) {
        ElMessageBox.confirm(
          'Base template will not generate any structs or CRUD. It only configures properties like enter to support custom non-CRUD logic.',
          'Warning',
          {
            confirmButtonText: 'Continue',
            cancelButtonText: 'Cancel',
            type: 'warning'
          }
        )
          .then(() => {
            form.value.fields = []
          })
          .catch(() => {
            form.value.onlyTemplate = false
          })
      }
    }
  )

  const canSelect = (fieldType,item) => {
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
</script>

<style>
.no-border-collapse{
  @apply border-none;
  .el-collapse-item__header{
    @apply border-none;
  }
  .el-collapse-item__wrap{
    @apply border-none;
  }
  .el-collapse-item__content{
    @apply pb-0;
  }
}
</style>
