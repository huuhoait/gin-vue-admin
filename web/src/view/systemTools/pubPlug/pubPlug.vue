<template>
  <div class="gva-form-box">
    <div class="p-4 bg-white dark:bg-slate-900">
      <WarningBar
        title="Only standard plugins are supported (generated from the plugin template). For non-standard plugins, please package manually."
      />
      <div class="flex items-center gap-3">
        <el-input v-model="plugName" placeholder='Plugin name (from template)' />
      </div>
      <el-card class="mt-2 text-center">
        <WarningBar title="In the transfer box, only select leaf menus." />
        <el-input
          v-model="parentMenu"
          placeholder="Enter menu group name, e.g. Announcement Management"
          class="mb-2"
        ></el-input>
        <el-transfer
          v-model="menus"
          :props="{
            key: 'ID'
          }"
          class="plugin-transfer"
          :data="menusData"
          filterable
          :filter-method="filterMenuMethod"
          filter-placeholder="Search by menu name/path"
          :titles="['Available menus', 'Selected menus']"
          :button-texts="['Remove', 'Select']"
        >
          <template #default="{ option }">
            {{ option.meta.title }} {{ option.component }}
          </template>
        </el-transfer>
        <div class="flex justify-end mt-2">
          <el-button type="primary" @click="fmtInitMenu">
            Generate install menus
          </el-button>
        </div>
      </el-card>
      <el-card class="mt-2 text-center">
        <el-transfer
          v-model="apis"
          :props="{
            key: 'ID'
          }"
          class="plugin-transfer"
          :data="apisData"
          filterable
          :filter-method="filterApiMethod"
          filter-placeholder="Search by API description/path"
          :titles="['Available APIs', 'Selected APIs']"
          :button-texts="['Remove', 'Select']"
        >
          <template #default="{ option }">
            {{ option.description }} {{ option.path }}
          </template>
        </el-transfer>
        <div class="flex justify-end mt-2">
          <el-button type="primary" @click="fmtInitAPI">
            Generate install APIs
          </el-button>
        </div>
      </el-card>
      <el-card class="mt-2 text-center">
        <el-transfer
          v-model="dictionaries"
          :props="{
            key: 'ID'
          }"
          class="plugin-transfer"
          :data="dictionariesData"
          filterable
          :filter-method="filterDictionaryMethod"
          filter-placeholder="Search by dictionary name/type"
          :titles="['Available dictionaries', 'Selected dictionaries']"
          :button-texts="['Remove', 'Select']"
        >
          <template #default="{ option }">
            {{ option.name }} {{ option.type }}
          </template>
        </el-transfer>
        <div class="flex justify-end mt-2">
          <el-button type="primary" @click="fmtInitDictionary">
            Generate install dictionaries
          </el-button>
        </div>
      </el-card>
    </div>
    <div class="flex justify-end">
      <el-button type="primary" @click="pubPlugin"> Package plugin </el-button>
    </div>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { pubPlug, initMenu, initAPI, initDictionary } from '@/api/autoCode.js'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { getAllApis } from '@/api/api'
  import { getMenuList } from '@/api/menu'
  import { getSysDictionaryList } from '@/api/sysDictionary'

  const plugName = ref('')

  const menus = ref([])
  const menusData = ref([])
  const apis = ref([])
  const apisData = ref([])
  const dictionaries = ref([])
  const dictionariesData = ref([])
  const parentMenu = ref('')

  const fmtMenu = (menus) => {
    // If a menu has children, flatten recursively.
    const res = []
    menus.forEach((item) => {
      if (item.children) {
        res.push(...fmtMenu(item.children))
      } else {
        res.push(item)
      }
    })
    return res
  }

  const initData = async () => {
    const menuRes = await getMenuList()
    if (menuRes.code === 0) {
      menusData.value = fmtMenu(menuRes.data)
    }
    const apiRes = await getAllApis()
    if (apiRes.code === 0) {
      apisData.value = apiRes.data.apis
    }
    const dictionaryRes = await getSysDictionaryList({
      page: 1,
      pageSize: 9999
    })
    if (dictionaryRes.code === 0) {
      dictionariesData.value = dictionaryRes.data
    }
  }

  const filterMenuMethod = (query, item) => {
    return (
      item.meta.title.indexOf(query) > -1 || item.component.indexOf(query) > -1
    )
  }

  const filterApiMethod = (query, item) => {
    return item.description.indexOf(query) > -1 || item.path.indexOf(query) > -1
  }

  const filterDictionaryMethod = (query, item) => {
    return item.name.indexOf(query) > -1 || item.type.indexOf(query) > -1
  }

  initData()


  const pubPlugin = async () => {
    ElMessageBox.confirm(
      `Please check whether required exports are enabled in server/plugin/${plugName.value}/plugin.go (initialize.Api(ctx), initialize.Menu(ctx), initialize.Dictionary(ctx)).`,
      'Package',
      {
        confirmButtonText: 'Package',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )
      .then(async () => {
        const res = await pubPlug({ plugName: plugName.value })
        if (res.code === 0) {
          ElMessage.success(res.msg)
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: 'Cancelled'
        })
      })
  }

  const fmtInitMenu = () => {
    if (!parentMenu.value) {
      ElMessage.error('Please enter a menu group name')
      return
    }
    if (menus.value.length === 0) {
      ElMessage.error('Please select at least one menu')
      return
    }
    if (plugName.value === '') {
      ElMessage.error('Please enter a plugin name')
      return
    }
    ElMessageBox.confirm(
      `This will overwrite server/plugin/${plugName.value}/initialize/menu. Continue?`,
      'Generate initial menus',
      {
        confirmButtonText: 'Generate',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )
      .then(async () => {
        const req = {
          plugName: plugName.value,
          parentMenu: parentMenu.value,
          menus: menus.value
        }
        const res = await initMenu(req)
        if (res.code === 0) {
          ElMessage.success('Menus initialized')
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: 'Cancelled'
        })
      })
  }
  const fmtInitAPI = () => {
    if (apis.value.length === 0) {
      ElMessage.error('Please select at least one API')
      return
    }
    if (plugName.value === '') {
      ElMessage.error('Please enter a plugin name')
      return
    }
    ElMessageBox.confirm(
      `This will overwrite server/plugin/${plugName.value}/initialize/api. Continue?`,
      'Generate initial APIs',
      {
        confirmButtonText: 'Generate',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )
      .then(async () => {
        const req = {
          plugName: plugName.value,
          apis: apis.value
        }
        const res = await initAPI(req)
        if (res.code === 0) {
          ElMessage.success('APIs initialized')
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: 'Cancelled'
        })
      })
  }

  const fmtInitDictionary = () => {
    if (dictionaries.value.length === 0) {
      ElMessage.error('Please select at least one dictionary')
      return
    }
    if (plugName.value === '') {
      ElMessage.error('Please enter a plugin name')
      return
    }
    ElMessageBox.confirm(
      `This will overwrite server/plugin/${plugName.value}/initialize/dictionary. Continue?`,
      'Generate initial dictionaries',
      {
        confirmButtonText: 'Generate',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )
      .then(async () => {
        const req = {
          plugName: plugName.value,
          dictionaries: dictionaries.value
        }
        const res = await initDictionary(req)
        if (res.code === 0) {
          ElMessage.success('Dictionaries initialized')
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: 'Cancelled'
        })
      })
  }
</script>

<style lang="scss">
  .plugin-transfer {
    .el-transfer-panel {
      width: 400px !important;
    }
  }
</style>
