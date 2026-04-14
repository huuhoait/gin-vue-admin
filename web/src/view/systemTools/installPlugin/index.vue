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
      <div class="el-upload__text">Drag files here or <em>click to upload</em></div>
      <template #tip>
        <div class="el-upload__tip">Drop the plugin zip here to upload</div>
      </template>
    </el-upload>

    <!-- Plugin List Table -->
    <div style="margin-top: 20px;">
      <el-table :data="pluginList" style="width: 100%">
        <el-table-column type="expand">
            <template #default="props">
                <div style="padding: 20px;">
                    <h3>API list</h3>
                    <el-table :data="props.row.apis" border>
                        <el-table-column prop="path" label="Path" />
                        <el-table-column prop="method" label="Method" />
                        <el-table-column prop="description" label="Description" />
                        <el-table-column prop="apiGroup" label="APIGROUP" />
                    </el-table>
                    <h3>Menu list</h3>
                    <el-table :data="props.row.menus" row-key="name" :tree-props="{children: 'children', hasChildren: 'hasChildren'}" border>
                        <el-table-column prop="meta.title" label="Title" />
                        <el-table-column prop="name" label="Name" />
                        <el-table-column prop="path" label="Path" />
                    </el-table>
                     <h3>Dictionary list</h3>
                     <el-table :data="props.row.dictionaries" border>
                         <el-table-column prop="name" label="Name" />
                         <el-table-column prop="type" label="Type" />
                         <el-table-column prop="desc" label="Description" />
                     </el-table>
                </div>
            </template>
        </el-table-column>
        <el-table-column prop="pluginName" label="Plugin name" />
        <el-table-column prop="pluginType" label="Plugin type">
          <template #default="scope">
              {{ typeMap[scope.row.pluginType] || 'Unknown' }}
          </template>
        </el-table-column>
        <el-table-column label="Actions">
          <template #default="scope">
            <el-button type="primary" link icon="delete" @click="deletePlugin(scope.row)">Delete</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getBaseUrl } from '@/utils/format'
  import { useUserStore } from "@/pinia";
  import { getPluginList, removePlugin } from '@/api/autoCode'
  import { ElMessageBox } from 'element-plus'

  const userStore = useUserStore()
  const token = userStore.token
  const pluginList = ref([])

  const getTableData = async () => {
    const res = await getPluginList()
    if (res.code === 0) {
      pluginList.value = res.data
    }
  }

  const typeMap = {
    "server": "Backend",
    "web": "Frontend",
    "full": "Fullstack"
  }

  const deletePlugin = (row) => {
    ElMessageBox.confirm(
    'This will permanently delete the plugin and its related API/menu/dictionary data. Continue?',
    'Confirm',
    {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  )
    .then(async () => {
      const res = await removePlugin({ pluginName: row.pluginName, pluginType: row.pluginType })
      if (res.code === 0) {
        ElMessage.success('Deleted')
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
