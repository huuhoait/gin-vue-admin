<template>
  <div class="flex justify-between items-center gap-2 w-full">
    <el-cascader
      v-if="pathIsSelect"
      :placeholder="t('admin.superadmin.menu.cascader.select_path')"
      :options="pathOptions"
      v-model="activeComponent"
      filterable
      class="!w-full"
      clearable
      @change="emitChange"
    />
    <el-input
      v-else
      v-model="tempPath"
      :placeholder="t('admin.superadmin.menu.cascader.manual_placeholder')"
      @change="emitChange"
    />
    <el-button @click="togglePathIsSelect"
      >{{ pathIsSelect ? t('admin.superadmin.menu.cascader.manual') : t('admin.superadmin.menu.cascader.quick_select') }}
    </el-button>
  </div>
</template>

<script setup>
  import { onMounted, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import pathInfo from '@/pathInfo.json'

  const { t } = useI18n()

  const props = defineProps({
    component: {
      type: String,
      default: ''
    }
  })

  const emits = defineEmits(['change'])

  const pathOptions = ref([])
  const tempPath = ref('')
  const activeComponent = ref([])
  const pathIsSelect = ref(true)

  const togglePathIsSelect = () => {
    if (pathIsSelect.value) {
      tempPath.value = activeComponent.value?.join('/') || ''
    } else {
      activeComponent.value = tempPath.value?.split('/') || []
    }

    pathIsSelect.value = !pathIsSelect.value
    emitChange()
  }

  function convertToCascaderOptions(data) {
    const result = []

    for (const path in data) {
      const label = data[path]
      const parts = path.split('/').filter(Boolean)

      // when the first segment is 'src', start parsing from the second segment
      const startIndex = parts[0] === 'src' ? 1 : 0

      let currentLevel = result

      for (let i = startIndex; i < parts.length; i++) {
        const part = parts[i]
        let node = currentLevel.find((item) => item.value === part)

        if (!node) {
          node = {
            value: part,
            label: part,
            children: []
          }
          currentLevel.push(node)
        }

        if (i === parts.length - 1) {
          // last path segment: attach label and drop children
          node.label = label
          delete node.children
        }

        currentLevel = node.children || []
      }
    }

    return result
  }

  watch(
    () => props.component,
    (value) => {
      initCascader(value)
    }
  )

  onMounted(() => {
    pathOptions.value = convertToCascaderOptions(pathInfo)
    initCascader(props.component)
  })

  const initCascader = (value) => {
    // create mode
    if (value === '') {
      pathIsSelect.value = true
      return
    }

    // edit mode: pick select vs. input based on known paths
    if (pathInfo[`/src/${value}`]) {
      activeComponent.value = value.split('/').filter(Boolean)
      tempPath.value = ''
      pathIsSelect.value = true
      return
    }
    tempPath.value = value
    activeComponent.value = []
    pathIsSelect.value = false
  }

  const emitChange = () => {
    emits(
      'change',
      pathIsSelect.value ? activeComponent.value?.join('/') : tempPath.value
    )
  }
</script>

<style scoped lang="scss"></style>
