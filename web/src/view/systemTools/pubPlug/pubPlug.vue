<template>
  <div class="gva-form-box">
    <div class="p-4 bg-white dark:bg-slate-900">
      <WarningBar
        :title="t('admin.systemtools.pubplug.warning_standard')"
      />
      <div class="flex items-center gap-3">
        <el-input v-model="plugName" :placeholder="t('admin.systemtools.pubplug.plugin_name_placeholder')" />
      </div>
      <el-card class="mt-2 text-center">
        <WarningBar :title="t('admin.systemtools.pubplug.leaf_menu_warning')" />
        <el-input
          v-model="parentMenu"
          :placeholder="t('admin.systemtools.pubplug.menu_group_placeholder')"
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
          :filter-placeholder="t('admin.systemtools.pubplug.search_menu')"
          :titles="[t('admin.systemtools.pubplug.available_menus'), t('admin.systemtools.pubplug.selected_menus')]"
          :button-texts="[t('admin.systemtools.pubplug.remove'), t('admin.systemtools.pubplug.select')]"
        >
          <template #default="{ option }">
            {{ option.meta.title }} {{ option.component }}
          </template>
        </el-transfer>
        <div class="flex justify-end mt-2">
          <el-button type="primary" @click="fmtInitMenu">
            {{ t('admin.systemtools.pubplug.generate_menus') }}
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
          :filter-placeholder="t('admin.systemtools.pubplug.search_api')"
          :titles="[t('admin.systemtools.pubplug.available_apis'), t('admin.systemtools.pubplug.selected_apis')]"
          :button-texts="[t('admin.systemtools.pubplug.remove'), t('admin.systemtools.pubplug.select')]"
        >
          <template #default="{ option }">
            {{ option.description }} {{ option.path }}
          </template>
        </el-transfer>
        <div class="flex justify-end mt-2">
          <el-button type="primary" @click="fmtInitAPI">
            {{ t('admin.systemtools.pubplug.generate_apis') }}
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
          :filter-placeholder="t('admin.systemtools.pubplug.search_dict')"
          :titles="[t('admin.systemtools.pubplug.available_dicts'), t('admin.systemtools.pubplug.selected_dicts')]"
          :button-texts="[t('admin.systemtools.pubplug.remove'), t('admin.systemtools.pubplug.select')]"
        >
          <template #default="{ option }">
            {{ option.name }} {{ option.type }}
          </template>
        </el-transfer>
        <div class="flex justify-end mt-2">
          <el-button type="primary" @click="fmtInitDictionary">
            {{ t('admin.systemtools.pubplug.generate_dicts') }}
          </el-button>
        </div>
      </el-card>
    </div>
    <div class="flex justify-end">
      <el-button type="primary" @click="pubPlugin"> {{ t('admin.systemtools.pubplug.package_plugin') }} </el-button>
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
  import { useI18n } from 'vue-i18n'

  const { t } = useI18n()

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
      t('admin.systemtools.pubplug.package_confirm', { name: plugName.value }),
      t('admin.systemtools.pubplug.package_title'),
      {
        confirmButtonText: t('admin.systemtools.pubplug.package_action'),
        cancelButtonText: t('admin.systemtools.pubplug.cancel'),
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
          message: t('admin.systemtools.pubplug.cancelled')
        })
      })
  }

  const fmtInitMenu = () => {
    if (!parentMenu.value) {
      ElMessage.error(t('admin.systemtools.pubplug.enter_menu_group'))
      return
    }
    if (menus.value.length === 0) {
      ElMessage.error(t('admin.systemtools.pubplug.select_at_least_one_menu'))
      return
    }
    if (plugName.value === '') {
      ElMessage.error(t('admin.systemtools.pubplug.enter_plugin_name'))
      return
    }
    ElMessageBox.confirm(
      t('admin.systemtools.pubplug.menu_overwrite_confirm', { name: plugName.value }),
      t('admin.systemtools.pubplug.generate_menu_title'),
      {
        confirmButtonText: t('admin.systemtools.pubplug.generate_action'),
        cancelButtonText: t('admin.systemtools.pubplug.cancel'),
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
          ElMessage.success(t('admin.systemtools.pubplug.menus_initialized'))
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('admin.systemtools.pubplug.cancelled')
        })
      })
  }
  const fmtInitAPI = () => {
    if (apis.value.length === 0) {
      ElMessage.error(t('admin.systemtools.pubplug.select_at_least_one_api'))
      return
    }
    if (plugName.value === '') {
      ElMessage.error(t('admin.systemtools.pubplug.enter_plugin_name'))
      return
    }
    ElMessageBox.confirm(
      t('admin.systemtools.pubplug.api_overwrite_confirm', { name: plugName.value }),
      t('admin.systemtools.pubplug.generate_api_title'),
      {
        confirmButtonText: t('admin.systemtools.pubplug.generate_action'),
        cancelButtonText: t('admin.systemtools.pubplug.cancel'),
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
          ElMessage.success(t('admin.systemtools.pubplug.apis_initialized'))
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('admin.systemtools.pubplug.cancelled')
        })
      })
  }

  const fmtInitDictionary = () => {
    if (dictionaries.value.length === 0) {
      ElMessage.error(t('admin.systemtools.pubplug.select_at_least_one_dict'))
      return
    }
    if (plugName.value === '') {
      ElMessage.error(t('admin.systemtools.pubplug.enter_plugin_name'))
      return
    }
    ElMessageBox.confirm(
      t('admin.systemtools.pubplug.dict_overwrite_confirm', { name: plugName.value }),
      t('admin.systemtools.pubplug.generate_dict_title'),
      {
        confirmButtonText: t('admin.systemtools.pubplug.generate_action'),
        cancelButtonText: t('admin.systemtools.pubplug.cancel'),
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
          ElMessage.success(t('admin.systemtools.pubplug.dicts_initialized'))
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('admin.systemtools.pubplug.cancelled')
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
