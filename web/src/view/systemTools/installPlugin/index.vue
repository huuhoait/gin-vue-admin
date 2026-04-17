<template>
  <div class="gva-form-box">
    <el-upload
      drag
      :action="`${getBaseUrl()}/autoCode/installPlugin`"
      :show-file-list="false"
      :on-success="handleSuccess"
      :on-error="handleSuccess"
      :headers="{'x-token': token}"
      name="plug"
    >
      <el-icon class="el-icon--upload"><upload-filled /></el-icon>
      <div class="el-upload__text">{{ t('admin.systemtools.install_plugin.drag_text_1') }} <em>{{ t('admin.systemtools.install_plugin.drag_text_2') }}</em></div>
      <template #tip>
        <div class="el-upload__tip">{{ t('admin.systemtools.install_plugin.drop_plugin_tip') }}</div>
      </template>
    </el-upload>

    <!-- Plugin List Table -->
    <div style="margin-top: 20px;">
      <el-table :data="pluginList" style="width: 100%">
        <el-table-column type="expand">
            <template #default="props">
                <div style="padding: 20px;">
                    <h3>{{ t('admin.systemtools.install_plugin.api_list') }}</h3>
                    <el-table :data="props.row.apis" border>
                        <el-table-column prop="path" :label="t('admin.systemtools.install_plugin.path')" />
                        <el-table-column prop="method" :label="t('admin.systemtools.install_plugin.method')" />
                        <el-table-column prop="description" :label="t('admin.systemtools.install_plugin.description')" />
                        <el-table-column prop="apiGroup" :label="t('admin.systemtools.install_plugin.api_group')" />
                    </el-table>
                    <h3>{{ t('admin.systemtools.install_plugin.menu_list') }}</h3>
                    <el-table :data="props.row.menus" row-key="name" :tree-props="{children: 'children', hasChildren: 'hasChildren'}" border>
                        <el-table-column prop="meta.title" :label="t('admin.systemtools.install_plugin.title')" />
                        <el-table-column prop="name" :label="t('admin.systemtools.install_plugin.name')" />
                        <el-table-column prop="path" :label="t('admin.systemtools.install_plugin.path')" />
                    </el-table>
                     <h3>{{ t('admin.systemtools.install_plugin.dictionary_list') }}</h3>
                     <el-table :data="props.row.dictionaries" border>
                         <el-table-column prop="name" :label="t('admin.systemtools.install_plugin.name')" />
                         <el-table-column prop="type" :label="t('admin.systemtools.install_plugin.type')" />
                         <el-table-column prop="desc" :label="t('admin.systemtools.install_plugin.description')" />
                     </el-table>
                </div>
            </template>
        </el-table-column>
        <el-table-column prop="pluginName" :label="t('admin.systemtools.install_plugin.plugin_name')" />
        <el-table-column prop="pluginType" :label="t('admin.systemtools.install_plugin.plugin_type')">
          <template #default="scope">
              {{ typeMap[scope.row.pluginType] || t('admin.systemtools.install_plugin.type_unknown') }}
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.systemtools.install_plugin.actions')">
          <template #default="scope">
            <el-button type="primary" link icon="delete" @click="deletePlugin(scope.row)">{{ t('admin.systemtools.install_plugin.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getBaseUrl } from '@/utils/format'
  import { useUserStore } from "@/pinia";
  import { getPluginList, removePlugin } from '@/api/autoCode'
  import { ElMessageBox } from 'element-plus'
  import { useI18n } from 'vue-i18n'

  const { t } = useI18n()

  const userStore = useUserStore()
  const token = userStore.token
  const pluginList = ref([])

  const getTableData = async () => {
    const res = await getPluginList()
    if (res.code === 0) {
      pluginList.value = res.data
    }
  }

  const typeMap = computed(() => ({
    "server": t('admin.systemtools.install_plugin.type_backend'),
    "web": t('admin.systemtools.install_plugin.type_frontend'),
    "full": t('admin.systemtools.install_plugin.type_full')
  }))

  const deletePlugin = (row) => {
    ElMessageBox.confirm(
    t('admin.systemtools.install_plugin.delete_confirm'),
    t('admin.common.confirms.delete_title'),
    {
      confirmButtonText: t('admin.common.confirm'),
      cancelButtonText: t('admin.common.cancel'),
      type: 'warning',
    }
  )
    .then(async () => {
      const res = await removePlugin({ pluginName: row.pluginName, pluginType: row.pluginType })
      if (res.code === 0) {
        ElMessage.success(t('admin.common.messages.deleted'))
        getTableData()
      }
    })
    .catch(() => {
    })
  }

  onMounted(() => {
    getTableData()
  })

  const handleSuccess = (res) => {
    if (res.code === 0) {
      let msg = ``
      res.data &&
        res.data.forEach((item, index) => {
          msg += `${index + 1}.${item.msg}\n`
        })
      alert(msg)
      getTableData() // Refresh list on success
    } else {
      ElMessage.error(res.msg)
    }
  }
</script>
