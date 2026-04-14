<template>
  <div>
    <warning-bar
      title="Dictionary caching is already implemented in frontend utils/dictionary. See comments in that file for usage."
    />
    <el-splitter class="h-full">
      <el-splitter-panel size="300px" min="200px" max="800px" collapsible>
        <div
          class="flex-none bg-white text-slate-700 dark:text-slate-400 dark:bg-slate-900 rounded p-4"
        >
          <div class="flex justify-between items-center relative">
            <span class="text font-bold">Dictionaries</span>
            <el-input
              class="!absolute top-0 left-0 z-2 ease-in-out animate-slide-left"
              placeholder="Search"
              v-if="showSearchInput"
              v-model="searchName"
              clearable
              :autofocus="showSearchInput"
              @clear="clearSearchInput"
              :prefix-icon="Search"
              v-click-outside="handleCloseSearchInput"
              @keydown="handleInputKeyDown"
            >
              <template #append>
                <el-button
                  :type="searchName ? 'primary' : 'info'"
                  @click="getTableData"
                  >Search</el-button
                >
              </template>
            </el-input>
            <el-button-group class="ml-auto">
              <el-tooltip content="Search" placement="top">
                <el-button
                  :icon="Search"
                  @click="showSearchInputHandler"
                />
              </el-tooltip>
              <el-tooltip content="Import dictionaries" placement="top">
                <el-button
                  type="success"
                  :icon="Upload"
                  @click="openImportDialog"
                />
              </el-tooltip>
              <el-tooltip content="Generate with AI" placement="top">
                <el-button
                  type="warning"
                  @click="openAiDialog"
                >
                  AI
                </el-button>
              </el-tooltip>
              <el-tooltip content="New dictionary" placement="top">
                <el-button
                  type="primary"
                  :icon="Plus"
                  @click="openDrawer"
                />
              </el-tooltip>
            </el-button-group>
          </div>
          <el-scrollbar class="mt-4" style="height: calc(100vh - 300px)">
            <div
              v-for="dictionary in dictionaryData"
              :key="dictionary.ID"
              class="rounded flex justify-between items-center px-2 py-4 cursor-pointer mt-2 hover:bg-blue-50 dark:hover:bg-blue-900 bg-gray-50 dark:bg-gray-800 gap-4"
              :class="[
                selectID === dictionary.ID
                  ? 'text-active'
                  : 'text-slate-700 dark:text-slate-50',
                dictionary.parentID ? 'ml-4 border-l-2 border-blue-200' : ''
              ]"
              @click="toDetail(dictionary)"
            >
              <div class="max-w-[160px] truncate">
                <span
                  v-if="dictionary.parentID"
                  class="text-xs text-gray-400 mr-1"
                  >└─</span
                >
                {{ dictionary.name }}
                <span class="mr-auto text-sm">({{ dictionary.type }})</span>
              </div>

              <div class="min-w-[60px] flex items-center gap-2">
                <el-icon
                  class="!text-green-500"
                  @click.stop="exportDictionary(dictionary)"
                  title="Export dictionary"
                >
                  <Download />
                </el-icon>
                <el-icon
                  class="!text-blue-500"
                  @click.stop="updateSysDictionaryFunc(dictionary)"
                >
                  <Edit />
                </el-icon>
                <el-icon
                  class="!text-red-500"
                  @click="deleteSysDictionaryFunc(dictionary)"
                >
                  <Delete />
                </el-icon>
              </div>
            </div>
          </el-scrollbar>
        </div>
      </el-splitter-panel>
      <el-splitter-panel :min="200">
        <div
          class="flex-1 bg-white text-slate-700 dark:text-slate-400 dark:bg-slate-900"
        >
          <sysDictionaryDetail :sys-dictionary-i-d="selectID" />
        </div>
      </el-splitter-panel>
    </el-splitter>

    <el-drawer
      v-model="drawerFormVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      :before-close="closeDrawer"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{
            type === 'create' ? 'Add dictionary' : 'Edit dictionary'
          }}</span>
          <div>
            <el-button @click="closeDrawer">Cancel</el-button>
            <el-button type="primary" @click="enterDrawer">Confirm</el-button>
          </div>
        </div>
      </template>
      <el-form
        ref="drawerForm"
        :model="formData"
        :rules="rules"
        label-width="110px"
      >
        <el-form-item label="Parent dictionary" prop="parentID">
          <el-select
            v-model="formData.parentID"
            placeholder="Select parent dictionary (optional)"
            clearable
            filterable
            :style="{ width: '100%' }"
          >
            <el-option
              v-for="dict in availableParentDictionaries"
              :key="dict.ID"
              :label="`${dict.name}(${dict.type})`"
              :value="dict.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Name" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="Enter name"
            clearable
            :style="{ width: '100%' }"
          />
        </el-form-item>
        <el-form-item label="Type" prop="type">
          <el-input
            v-model="formData.type"
            placeholder="Enter type"
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
        <el-form-item label="Description" prop="desc">
          <el-input
            v-model="formData.desc"
            placeholder="Enter description"
            clearable
            :style="{ width: '100%' }"
          />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Import drawer -->
    <el-drawer
      v-model="importDrawerVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      :before-close="closeImportDrawer"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Import dictionary JSON</span>
          <div>
            <el-button @click="closeImportDrawer">Cancel</el-button>
            <el-button type="primary" @click="handleImport" :loading="importing">
              Import
            </el-button>
          </div>
        </div>
      </template>
      
      <div class="import-drawer-content">
        <div class="mb-4">
          <el-alert
            title="Paste, edit, or drag a JSON file into the area below"
            type="info"
            :closable="false"
            show-icon
          />
        </div>

        <!-- Drag-and-drop upload -->
        <div
          class="drag-upload-area"
          :class="{ 'is-dragging': isDragging }"
          @drop.prevent="handleDrop"
          @dragover.prevent="handleDragOver"
          @dragleave.prevent="handleDragLeave"
          @click="triggerFileInput"
        >
          <el-icon class="upload-icon"><Upload /></el-icon>
          <div class="upload-text">
            <p>Drag a JSON file here, or click to choose a file</p>
            <p class="upload-hint">You can also edit directly in the textbox below</p>
          </div>
          <input
            ref="fileInputRef"
            type="file"
            accept=".json,application/json"
            style="display: none"
            @change="handleFileSelect"
          />
        </div>

        <div class="json-editor-container mt-4">
          <el-input
            v-model="importJsonText"
            type="textarea"
            :rows="15"
            placeholder='Enter JSON, for example:
{
  "name": "gender",
  "type": "gender",
  "status": true,
  "desc": "gender dictionary",
  "sysDictionaryDetails": [
    {
      "label": "male",
      "value": "1",
      "status": true,
      "sort": 1
    },
    {
      "label": "female",
      "value": "2",
      "status": true,
      "sort": 2
    }
  ]
}'
            class="json-textarea"
          />
        </div>

        <div class="mt-4" v-if="jsonPreviewError">
          <el-alert
            :title="jsonPreviewError"
            type="error"
            :closable="false"
            show-icon
          />
        </div>

    
      </div>
    </el-drawer>

    <!-- AI dialog -->
    <el-dialog
      v-model="aiDialogVisible"
      title="Generate dictionary with AI"
      width="520px"
      :before-close="closeAiDialog"
    >
      <div class="relative">
        <el-input
          v-model="aiPrompt"
          type="textarea"
          :rows="6"
          :maxlength="2000"
          placeholder="Describe the dictionary to generate (e.g. user status: enabled/disabled). You can paste or upload an image to extract text."
          resize="none"
          @keydown.ctrl.enter="handleAiGenerate"
          @paste="handlePaste"
          @focus="handleFocus"
          @blur="handleBlur"
        />

        <input
          ref="imageFileInputRef"
          type="file"
          accept="image/*"
          style="display:none"
          @change="handleImageSelect"
        />

        <div class="flex absolute right-2 bottom-2">
          <el-tooltip effect="light">
            <template #content>
              <div>Paste or upload an image, then extract its contents to generate the dictionary.</div>
            </template>
            <el-button type="primary" @click="eyeFunc">
                <el-icon size="18">
                <ai-gva />
              </el-icon>
              Extract
            </el-button>
          </el-tooltip>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeAiDialog">Cancel</el-button>
          <el-button type="primary" @click="handleAiGenerate" :loading="aiGenerating">
            Confirm
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    createSysDictionary,
    deleteSysDictionary,
    updateSysDictionary,
    findSysDictionary,
    getSysDictionaryList,
    exportSysDictionary,
    importSysDictionary
  } from '@/api/sysDictionary' // replace API module path if needed
  import { llmAuto } from '@/api/autoCode'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { ref, computed, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'

  import sysDictionaryDetail from './sysDictionaryDetail.vue'
  import { Edit, Plus, Search, Download, Upload } from '@element-plus/icons-vue'
  import { useAppStore } from '@/pinia'

  defineOptions({
    name: 'SysDictionary'
  })

  const appStore = useAppStore()

  const selectID = ref(0)

  const formData = ref({
    name: null,
    type: null,
    status: true,
    desc: null,
    parentID: null
  })
  const searchName = ref('')
  const showSearchInput = ref(false)
  const rules = ref({
    name: [
      {
        required: true,
        message: 'Name is required',
        trigger: 'blur'
      }
    ],
    type: [
      {
        required: true,
        message: 'Type is required',
        trigger: 'blur'
      }
    ],
    desc: [
      {
        required: true,
        message: 'Description is required',
        trigger: 'blur'
      }
    ]
  })

  const dictionaryData = ref([])
  const availableParentDictionaries = ref([])

  // Import
  const importDrawerVisible = ref(false)
  const importJsonText = ref('')
  const importing = ref(false)
  const jsonPreviewError = ref('')
  const jsonPreview = ref(null)
  const isDragging = ref(false)
  const fileInputRef = ref(null)

  // AI
  const aiDialogVisible = ref(false)
  const aiPrompt = ref('')
  const aiGenerating = ref(false)

  // Image upload / OCR
  const imageFileInputRef = ref(null)
  const focused = ref(false)

  const handleFocus = () => {
    focused.value = true
  }
  const handleBlur = () => {
    focused.value = false
  }

  // Trigger image selection
  const triggerImageSelect = () => {
    imageFileInputRef.value?.click()
  }

  const handlePaste = (event) => {
    const items = event.clipboardData.items;
    for (let i = 0; i < items.length; i++) {
      if (items[i].type.indexOf('image') !== -1) {
        const file = items[i].getAsFile();
        const reader = new FileReader();
        reader.onload =async (e) => {
          const base64String = e.target.result;
          const res = await llmAuto({ _file_path: base64String, mode:"dictEye" })
          if (res.code === 0) {
            aiPrompt.value = res.data.text
          }
        };
        reader.readAsDataURL(file);
      }
    }
  };

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

          const res = await llmAuto({ _file_path: base64String, mode:"dictEye" })
          if (res.code === 0) {
            aiPrompt.value = res.data.text
          }
        };
        reader.readAsDataURL(file);
      }
    };

    input.click();
  }



  // Watch JSON changes for live preview
  watch(importJsonText, (newVal) => {
    if (!newVal.trim()) {
      jsonPreview.value = null
      jsonPreviewError.value = ''
      return
    }
    try {
      jsonPreview.value = JSON.parse(newVal)
      jsonPreviewError.value = ''
    } catch (e) {
      jsonPreviewError.value = 'Invalid JSON: ' + e.message
      jsonPreview.value = null
    }
  })

  // Format JSON preview
  const jsonPreviewFormatted = computed(() => {
    if (!jsonPreview.value) return ''
    return JSON.stringify(jsonPreview.value, null, 2)
  })


  // Query
  const getTableData = async () => {
    const res = await getSysDictionaryList({
      name: searchName.value.trim()
    })
    if (res.code === 0) {
      dictionaryData.value = res.data
      selectID.value = res.data[0].ID
      // Update parent dictionary options
      updateAvailableParentDictionaries()
    }
  }

  // Update parent dictionary options
  const updateAvailableParentDictionaries = () => {
    // In edit mode, exclude current dictionary and its children
    if (type.value === 'update' && formData.value.ID) {
      availableParentDictionaries.value = dictionaryData.value.filter(
        (dict) => {
          return (
            dict.ID !== formData.value.ID &&
            !isChildDictionary(dict.ID, formData.value.ID)
          )
        }
      )
    } else {
      // In create mode, show all dictionaries
      availableParentDictionaries.value = [...dictionaryData.value]
    }
  }

  // Check child dictionaries (avoid cycles)
  const isChildDictionary = (dictId, parentId) => {
    const dict = dictionaryData.value.find((d) => d.ID === dictId)
    if (!dict || !dict.parentID) return false
    if (dict.parentID === parentId) return true
    return isChildDictionary(dict.parentID, parentId)
  }

  getTableData()

  const toDetail = (row) => {
    selectID.value = row.ID
  }

  const drawerFormVisible = ref(false)
  const type = ref('')
  const updateSysDictionaryFunc = async (row) => {
    const res = await findSysDictionary({ ID: row.ID, status: row.status })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = res.data.resysDictionary
      drawerFormVisible.value = true
      // Update parent dictionary options
      updateAvailableParentDictionaries()
    }
  }
  const closeDrawer = () => {
    drawerFormVisible.value = false
    formData.value = {
      name: null,
      type: null,
      status: true,
      desc: null,
      parentID: null
    }
  }
  const deleteSysDictionaryFunc = async (row) => {
    ElMessageBox.confirm('Delete this dictionary?', 'Confirm', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning'
    }).then(async () => {
      const res = await deleteSysDictionary({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: 'Deleted'
        })
        getTableData()
      }
    })
  }

  const drawerForm = ref(null)
  const enterDrawer = async () => {
    drawerForm.value.validate(async (valid) => {
      if (!valid) return
      let res
      switch (type.value) {
        case 'create':
          res = await createSysDictionary(formData.value)
          break
        case 'update':
          res = await updateSysDictionary(formData.value)
          break
        default:
          res = await createSysDictionary(formData.value)
          break
      }
      if (res.code === 0) {
        ElMessage.success('Success')
        closeDrawer()
        getTableData()
      }
    })
  }
  const openDrawer = () => {
    type.value = 'create'
    drawerForm.value && drawerForm.value.clearValidate()
    drawerFormVisible.value = true
    // Update parent dictionary options
    updateAvailableParentDictionaries()
  }

  const clearSearchInput = () => {
    if (!showSearchInput.value) return
    searchName.value = ''
    showSearchInput.value = false
    getTableData()
  }
  const handleCloseSearchInput = () => {
    if (!showSearchInput.value || searchName.value.trim() != '') return
    showSearchInput.value = false
  }

  const showSearchInputHandler = () => {
    showSearchInput.value = true
  }

  const handleInputKeyDown = (e) => {
    if (e.key === 'Enter' && searchName.value.trim() !== '') {
      getTableData()
    }
  }

  // Export dictionary
  const exportDictionary = async (row) => {
    try {
      const res = await exportSysDictionary({ ID: row.ID })
      if (res.code === 0) {
        // Convert JSON to string and download
        const jsonStr = JSON.stringify(res.data, null, 2)
        const blob = new Blob([jsonStr], { type: 'application/json' })
        const url = URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `${row.type}_${row.name}_dictionary.json`
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        URL.revokeObjectURL(url)
        ElMessage.success('Exported')
      }
    } catch (error) {
      ElMessage.error('Export failed: ' + error.message)
    }
  }

  // Open import drawer
  const openImportDialog = () => {
    importDrawerVisible.value = true
    importJsonText.value = ''
    jsonPreview.value = null
    jsonPreviewError.value = ''
    isDragging.value = false
  }

  // Close import drawer
  const closeImportDrawer = () => {
    importDrawerVisible.value = false
    importJsonText.value = ''
    jsonPreview.value = null
    jsonPreviewError.value = ''
    isDragging.value = false
  }

  // Handle drag enter
  const handleDragOver = (e) => {
    isDragging.value = true
  }

  // Handle drag leave
  const handleDragLeave = (e) => {
    isDragging.value = false
  }
  // Handle file drop
  const handleDrop = (e) => {
    isDragging.value = false
    const files = e.dataTransfer.files
    if (files.length === 0) return

    const file = files[0]
    readJsonFile(file)
  }

  // Trigger file selection
  const triggerFileInput = () => {
    fileInputRef.value?.click()
  }

  // Handle file selection
  const handleFileSelect = (e) => {
    const files = e.target.files
    if (files.length === 0) return

    const file = files[0]
    readJsonFile(file)
    
    // Clear input so the same file can be selected again
    e.target.value = ''
  }

  // Read JSON file
  const readJsonFile = (file) => {
    // Check file type
    if (!file.name.endsWith('.json')) {
      ElMessage.warning('Please upload a JSON file')
      return
    }

    // Read file contents
    const reader = new FileReader()
    reader.onload = (event) => {
      try {
        const content = event.target.result
        // Validate JSON
        JSON.parse(content)
        importJsonText.value = content
        ElMessage.success('File loaded')
      } catch (error) {
        ElMessage.error('File content is not valid JSON')
      }
    }
    reader.onerror = () => {
      ElMessage.error('Failed to read file')
    }
    reader.readAsText(file)
  }

  // Import
  const handleImport = async () => {
    if (!importJsonText.value.trim()) {
      ElMessage.warning('Please enter JSON')
      return
    }

    if (jsonPreviewError.value) {
      ElMessage.error('Invalid JSON, please retry')
      return
    }

    try {
      importing.value = true
      const res = await importSysDictionary({ json: importJsonText.value })
      if (res.code === 0) {
        ElMessage.success('Imported')
        closeImportDrawer()
        getTableData()
      }
    } catch (error) {
      ElMessage.error('Import failed: ' + error.message)
    } finally {
      importing.value = false
    }
  }

  // Open AI dialog
  const openAiDialog = () => {
    aiDialogVisible.value = true
    aiPrompt.value = ''
  }

  // Close AI dialog
  const closeAiDialog = () => {
    aiDialogVisible.value = false
    aiPrompt.value = ''
  }

  // Generate with AI
  const handleAiGenerate = async () => {
    if (!aiPrompt.value.trim()) {
      ElMessage.warning('Please enter a description')
      return
    }
    try {
      aiGenerating.value = true
      const aiRes = await llmAuto({
        prompt: aiPrompt.value,
        mode: 'dict'
      })
      if (aiRes && aiRes.code === 0) {
        ElMessage.success('AI generation succeeded')
        try {
          // Fill the returned data into import textbox (string or object)
          if (typeof aiRes.data === 'string') {
            importJsonText.value = aiRes.data
          } else {
            importJsonText.value = JSON.stringify(aiRes.data, null, 2)
          }
          // Clear parse errors and open import drawer
          jsonPreviewError.value = ''
          importDrawerVisible.value = true
          closeAiDialog()
        } catch (e) {
          ElMessage.error('Failed to process AI result: ' + (e.message || e))
        }
      } 
    } catch (err) {
      ElMessage.error('AI request failed: ' + (err.message || err))
    } finally {
      aiGenerating.value = false
    }
  }
</script>

<style scoped>
  .dict-box {
    height: calc(100vh - 240px);
  }

  .active {
    background-color: var(--el-color-primary) !important;
    color: #fff;
  }

  .import-drawer-content {
    padding: 0 4px;
  }

  /* Drag-and-drop upload */
  .drag-upload-area {
    border: 2px dashed #dcdfe6;
    border-radius: 8px;
    padding: 40px 20px;
    text-align: center;
    background-color: #fafafa;
    transition: all 0.3s ease;
    cursor: pointer;
  }

  .drag-upload-area:hover {
    border-color: #409eff;
    background-color: #ecf5ff;
  }

  .drag-upload-area.is-dragging {
    border-color: #409eff;
    background-color: #ecf5ff;
    transform: scale(1.02);
  }

  .upload-icon {
    font-size: 48px;
    color: #8c939d;
    margin-bottom: 16px;
  }

  .drag-upload-area.is-dragging .upload-icon {
    color: #409eff;
  }

  .upload-text {
    color: #606266;
  }

  .upload-text p {
    margin: 4px 0;
  }

  .upload-hint {
    font-size: 12px;
    color: #909399;
  }

  .json-editor-container {
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    overflow: hidden;
  }

  .json-textarea :deep(.el-textarea__inner) {
    font-family: 'Courier New', Courier, monospace;
    font-size: 13px;
    line-height: 1.5;
  }

  .json-preview {
    background-color: #f5f7fa;
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    padding: 16px;
    max-height: 400px;
    overflow: auto;
  }

  .json-preview pre {
    margin: 0;
    font-family: 'Courier New', Courier, monospace;
    font-size: 13px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-wrap: break-word;
  }

  .dark .drag-upload-area {
    background-color: #1d1e1f;
    border-color: #414243;
  }

  .dark .drag-upload-area:hover,
  .dark .drag-upload-area.is-dragging {
    background-color: #1a3a52;
    border-color: #409eff;
  }

  .dark .json-preview {
    background-color: #1d1e1f;
    border-color: #414243;
  }
</style>
