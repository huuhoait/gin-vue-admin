<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list justify-between flex items-center">
        <span class="text font-bold">Dictionary details</span>
        <div class="flex items-center gap-2">
          <el-input
            placeholder="Search display label"
            v-model="searchName"
            clearable
            class="!w-64"
            @clear="clearSearchInput"
            :prefix-icon="Search"
            v-click-outside="handleCloseSearchInput"
            @keydown="handleInputKeyDown"
          >
            <template #append>
              <el-button
                :type="searchName ? 'primary' : 'info'"
                @click="applySearch"
                >Search</el-button
              >
            </template>
          </el-input>
          <el-button type="primary" icon="plus" @click="openDrawer">
            Add item
          </el-button>
        </div>
      </div>
      <!-- Table view -->
      <el-table
        :data="displayTreeData"
        style="width: 100%"
        tooltip-effect="dark"
        :tree-props="{ children: 'children'}"
        row-key="ID"
        default-expand-all
      >
        <el-table-column type="selection" width="55" />

        <el-table-column align="left" label="Label" prop="label" min-width="100"/>

        <el-table-column align="left" label="Value" prop="value" />

        <el-table-column align="left" label="Extend" prop="extend" />

        <el-table-column align="left" label="Level" prop="level" width="80" />

        <el-table-column
          align="left"
          label="Status"
          prop="status"
          width="100"
        >
          <template #default="scope">
            {{ formatBoolean(scope.row.status) }}
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          label="Sort"
          prop="sort"
          width="100"
        />

        <el-table-column
          align="left"
          :label="$t('admin.common.operation')"
          :min-width="appStore.operateMinWith"
          fixed="right"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="plus"
              @click="addChildNode(scope.row)"
            >
              Add child
            </el-button>
            <el-button
              type="primary"
              link
              icon="edit"
              @click="updateSysDictionaryDetailFunc(scope.row)"
            >
              Edit
            </el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteSysDictionaryDetailFunc(scope.row)"
            >
              {{ $t('admin.common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-drawer
      v-model="drawerFormVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      :before-close="closeDrawer"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{
            type === 'create' ? 'Add item' : 'Edit item'
          }}</span>
          <div>
            <el-button @click="closeDrawer">{{ $t('admin.common.cancel') }}</el-button>
            <el-button type="primary" @click="enterDrawer">{{ $t('admin.common.confirm') }}</el-button>
          </div>
        </div>
      </template>
      <el-form
        ref="drawerForm"
        :model="formData"
        :rules="rules"
        label-width="110px"
      >
        <el-form-item label="Parent item" prop="parentID">
          <el-cascader
            v-model="formData.parentID"
            :options="[rootOption,...treeData]"
            :props="cascadeProps"
            placeholder="Select parent item (optional)"
            clearable
            filterable
            :style="{ width: '100%' }"
            @change="handleParentChange"
          />
        </el-form-item>
        <el-form-item label="Label" prop="label">
          <el-input
            v-model="formData.label"
            placeholder="Enter label"
            clearable
            :style="{ width: '100%' }"
          />
        </el-form-item>
        <el-form-item label="Value" prop="value">
          <el-input
            v-model="formData.value"
            placeholder="Enter value"
            clearable
            :style="{ width: '100%' }"
          />
        </el-form-item>
        <el-form-item label="Extend" prop="extend">
          <el-input
            v-model="formData.extend"
            placeholder="Enter extend"
            clearable
            :style="{ width: '100%' }"
          />
        </el-form-item>
        <el-form-item label="Status" prop="status" required>
          <el-switch
            v-model="formData.status"
            :active-text="$t('admin.common.enabled')"
            :inactive-text="$t('admin.common.disabled')"
          />
        </el-form-item>
        <el-form-item label="Sort" prop="sort">
          <el-input-number
            v-model.number="formData.sort"
            placeholder="Sort"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createSysDictionaryDetail,
    deleteSysDictionaryDetail,
    updateSysDictionaryDetail,
    findSysDictionaryDetail,
    getDictionaryTreeList
  } from '@/api/sysDictionaryDetail' // replace API module path if needed
  import { ref, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { formatBoolean, formatDate } from '@/utils/format'
  import { useAppStore } from '@/pinia'
  import { Search } from '@element-plus/icons-vue'

  defineOptions({
    name: 'SysDictionaryDetail'
  })

  const appStore = useAppStore()
  const searchName = ref('')

  const props = defineProps({
    sysDictionaryID: {
      type: Number,
      default: 0
    }
  })

  const formData = ref({
    label: null,
    value: null,
    status: true,
    sort: null,
    parentID: null
  })

  const rules = ref({
    label: [
      {
        required: true,
        message: 'Label is required',
        trigger: 'blur'
      }
    ],
    value: [
      {
        required: true,
        message: 'Value is required',
        trigger: 'blur'
      }
    ],
    sort: [
      {
        required: true,
        message: 'Sort is required',
        trigger: 'blur'
      }
    ]
  })

  const treeData = ref([])
  const displayTreeData = ref([])

  // Cascader config
  const cascadeProps = {
    value: 'ID',
    label: 'label',
    children: 'children',
    checkStrictly: true, // allow selecting any level
    emitPath: false // return only selected node value
  }


  const normalizeSearch = (value) => (value ?? '').toString().toLowerCase()

  const filterTree = (nodes, keyword) => {
    const trimmed = normalizeSearch(keyword).trim()
    if (!trimmed) {
      return nodes
    }
    const walk = (list) => {
      const result = []
      for (const node of list) {
        const label = normalizeSearch(node.label)
        const children = Array.isArray(node.children) ? walk(node.children) : []
        if (label.includes(trimmed) || children.length > 0) {
          result.push({
            ...node,
            children
          })
        }
      }
      return result
    }
    return walk(nodes)
  }

  const applySearch = () => {
    displayTreeData.value = filterTree(treeData.value, searchName.value)
  }

  // Fetch tree data
  const getTreeData = async () => {
    if (!props.sysDictionaryID) return
    try {
      const res = await getDictionaryTreeList({
        sysDictionaryID: props.sysDictionaryID
      })
      if (res.code === 0) {
        treeData.value = res.data.list || []
        applySearch()
      }
    } catch (error) {
      console.error('failed to fetch tree data:', error)
      ElMessage.error('Failed to load tree data')
    }
  }

  const rootOption = {
    ID: null,
    label: 'No parent (root)'
  }


  // Initial load
  getTreeData()

  const type = ref('')
  const drawerFormVisible = ref(false)

  const updateSysDictionaryDetailFunc = async (row) => {
    drawerForm.value && drawerForm.value.clearValidate()
    const res = await findSysDictionaryDetail({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = res.data.reSysDictionaryDetail
      drawerFormVisible.value = true
    }
  }

  // Add child node
  const addChildNode = (parentNode) => {
    console.log(parentNode)
    type.value = 'create'
    formData.value = {
      label: null,
      value: null,
      status: true,
      sort: null,
      parentID: parentNode.ID,
      sysDictionaryID: props.sysDictionaryID
    }
    drawerForm.value && drawerForm.value.clearValidate()
    drawerFormVisible.value = true
  }

  // Handle parent selection change
  const handleParentChange = (value) => {
    formData.value.parentID = value
  }

  const closeDrawer = () => {
    drawerFormVisible.value = false
    formData.value = {
      label: null,
      value: null,
      status: true,
      sort: null,
      parentID: null,
      sysDictionaryID: props.sysDictionaryID
    }
  }

  const deleteSysDictionaryDetailFunc = async (row) => {
    ElMessageBox.confirm('Delete this item?', 'Confirm', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning'
    }).then(async () => {
      const res = await deleteSysDictionaryDetail({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: 'Deleted'
        })
        await getTreeData() // reload
      }
    })
  }

  const drawerForm = ref(null)
  const enterDrawer = async () => {
    drawerForm.value.validate(async (valid) => {
      formData.value.sysDictionaryID = props.sysDictionaryID
      if (!valid) return
      let res
      switch (type.value) {
        case 'create':
          res = await createSysDictionaryDetail(formData.value)
          break
        case 'update':
          res = await updateSysDictionaryDetail(formData.value)
          break
        default:
          res = await createSysDictionaryDetail(formData.value)
          break
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: 'Saved'
        })
        closeDrawer()
        await getTreeData() // reload
      }
    })
  }

  const openDrawer = () => {
    type.value = 'create'
    formData.value.parentID = null
    drawerForm.value && drawerForm.value.clearValidate()
    drawerFormVisible.value = true
  }

  const clearSearchInput = () => {
    searchName.value = ''
    applySearch()
  }

  const handleCloseSearchInput = () => {
    // Handle closing the search input
  }

  const handleInputKeyDown = (e) => {
    if (e.key === 'Enter') {
      applySearch()
    }
  }

  watch(
    () => props.sysDictionaryID,
    () => {
      getTreeData()
    }
  )
</script>

<style scoped>

</style>
