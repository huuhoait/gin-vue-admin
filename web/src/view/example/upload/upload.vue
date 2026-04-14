<template>
  <div v-loading.fullscreen.lock="fullscreenLoading">
    <div class="flex gap-4 pt-2">
      <div
        class="flex-none w-64 bg-white text-slate-700 dark:text-slate-400 dark:bg-slate-900 rounded p-4"
      >
        <el-scrollbar style="height: calc(100vh - 300px)">
          <el-tree
            :data="categories"
            node-key="id"
            :props="defaultProps"
            @node-click="handleNodeClick"
            default-expand-all
          >
            <template #default="{ node, data }">
              <div
                class="w-36"
                :class="
                  search.classId === data.ID ? 'text-blue-500 font-bold' : ''
                "
              >
                {{ data.name }}
              </div>
              <el-dropdown>
                <el-icon class="ml-3 text-right" v-if="data.ID > 0"
                  ><MoreFilled
                /></el-icon>
                <el-icon class="ml-3 text-right mt-1" v-else><Plus /></el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="addCategoryFun(data)"
                      >Add category</el-dropdown-item
                    >
                    <el-dropdown-item
                      @click="editCategory(data)"
                      v-if="data.ID > 0"
                      >Edit category</el-dropdown-item
                    >
                    <el-dropdown-item
                      @click="deleteCategoryFun(data.ID)"
                      v-if="data.ID > 0"
                      >Delete category</el-dropdown-item
                    >
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-tree>
        </el-scrollbar>
      </div>
      <div
        class="flex-1 bg-white text-slate-700 dark:text-slate-400 dark:bg-slate-900"
      >
        <div class="gva-table-box mt-0 mb-0">
          <warning-bar
            title="Click the file name to edit. The selected category will be used for uploads."
          />
          <div class="gva-btn-list gap-3">
            <upload-common
              :image-common="imageCommon"
              :classId="search.classId"
              @on-success="onSuccess"
            />
            <cropper-image :classId="search.classId" @on-success="onSuccess" />
            <QRCodeUpload :classId="search.classId" @on-success="onSuccess" />
            <upload-image
              :image-url="imageUrl"
              :file-size="512"
              :max-w-h="1080"
              :classId="search.classId"
              @on-success="onSuccess"
            />
            <el-button type="primary" icon="upload" @click="importUrlFunc">
              Import URL
            </el-button>
            <el-input
              v-model="search.keyword"
              class="w-72"
              placeholder="Enter file name or note"
            />
            <el-button type="primary" icon="search" @click="onSubmit"
              >Search
            </el-button>
          </div>

          <el-table :data="tableData">
            <el-table-column align="left" label="Preview" width="100">
              <template #default="scope">
                <CustomPic pic-type="file" :pic-src="scope.row.url" preview />
              </template>
            </el-table-column>
            <el-table-column
              align="left"
              label="Updated at"
              prop="UpdatedAt"
              width="180"
            >
              <template #default="scope">
                <div>{{ formatDate(scope.row.UpdatedAt) }}</div>
              </template>
            </el-table-column>
            <el-table-column
              align="left"
              label="Name / note"
              prop="name"
              width="180"
            >
              <template #default="scope">
                <div
                  class="cursor-pointer"
                  @click="editFileNameFunc(scope.row)"
                >
                  {{ scope.row.name }}
                </div>
              </template>
            </el-table-column>
            <el-table-column
              align="left"
              label="URL"
              prop="url"
              min-width="300"
            />
            <el-table-column align="left" label="Tag" prop="tag" width="100">
              <template #default="scope">
                <el-tag
                  :type="
                    scope.row.tag?.toLowerCase() === 'jpg' ? 'info' : 'success'
                  "
                  disable-transitions
                  >{{ scope.row.tag }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column align="left" label="Actions" width="160">
              <template #default="scope">
                <el-button
                  icon="download"
                  type="primary"
                  link
                  @click="downloadFile(scope.row)"
                  >Download
                </el-button>
                <el-button
                  icon="delete"
                  type="primary"
                  link
                  @click="deleteFileFunc(scope.row)"
                  >Delete
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
              :current-page="page"
              :page-size="pageSize"
              :page-sizes="[10, 30, 50, 100]"
              :style="{ float: 'right', padding: '20px' }"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="handleCurrentChange"
              @size-change="handleSizeChange"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Add category dialog -->
    <el-dialog
      v-model="categoryDialogVisible"
      @close="closeAddCategoryDialog"
      width="520"
      :title="(categoryFormData.ID === 0 ? 'Add' : 'Edit') + ' category'"
      draggable
    >
      <el-form
        ref="categoryForm"
        :rules="rules"
        :model="categoryFormData"
        label-width="80px"
      >
        <el-form-item label="Parent category">
          <el-tree-select
            v-model="categoryFormData.pid"
            :data="categories"
            check-strictly
            :props="defaultProps"
            :render-after-expand="false"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item label="Category name" prop="name">
          <el-input
            v-model.trim="categoryFormData.name"
            placeholder="Category name"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeAddCategoryDialog">Cancel</el-button>
        <el-button type="primary" @click="confirmAddCategory">Confirm</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    getFileList,
    deleteFile,
    editFileName,
    importURL
  } from '@/api/fileUploadAndDownload'
  import { downloadImage } from '@/utils/downloadImg'
  import CustomPic from '@/components/customPic/index.vue'
  import UploadImage from '@/components/upload/image.vue'
  import UploadCommon from '@/components/upload/common.vue'
  import { CreateUUID, formatDate } from '@/utils/format'
  import WarningBar from '@/components/warningBar/warningBar.vue'

  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import {
    addCategory,
    deleteCategory,
    getCategoryList
  } from '@/api/attachmentCategory'
  import CropperImage from '@/components/upload/cropper.vue'
  import QRCodeUpload from '@/components/upload/QR-code.vue'

  defineOptions({
    name: 'Upload'
  })

  const fullscreenLoading = ref(false)
  const path = ref(import.meta.env.VITE_BASE_API)

  const imageUrl = ref('')
  const imageCommon = ref('')

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const search = ref({
    keyword: null,
    classId: 0
  })
  const tableData = ref([])

  // Pagination
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  const onSubmit = () => {
    search.value.classId = 0
    page.value = 1
    getTableData()
  }

  // Query
  const getTableData = async () => {
    const table = await getFileList({
      page: page.value,
      pageSize: pageSize.value,
      ...search.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }
  getTableData()

  const deleteFileFunc = async (row) => {
  ElMessageBox.confirm('This will permanently delete the file. Continue?', 'Confirm', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning'
    })
      .then(async () => {
        const res = await deleteFile(row)
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: 'Deleted'
          })
          if (tableData.value.length === 1 && page.value > 1) {
            page.value--
          }
          await getTableData()
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: 'Cancelled'
        })
      })
  }

  const downloadFile = (row) => {
    if (row.url.indexOf('http://') > -1 || row.url.indexOf('https://') > -1) {
      downloadImage(row.url, row.name)
    } else {
      downloadImage(path.value + '/' + row.url, row.name)
    }
  }

  /**
   * Edit file name or note
   * @param row
   * @returns {Promise<void>}
   */
  const editFileNameFunc = async (row) => {
    ElMessageBox.prompt('Enter file name or note', 'Edit', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      inputPattern: /\S/,
      inputErrorMessage: 'Required',
      inputValue: row.name
    })
      .then(async ({ value }) => {
        row.name = value
        // console.log(row)
        const res = await editFileName(row)
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: 'Saved'
          })
          await getTableData()
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: 'Cancelled'
        })
      })
  }

  /**
   * Import URL
   */
  const importUrlFunc = () => {
    ElMessageBox.prompt('Format: name|url or url only.', 'Import', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      inputType: 'textarea',
      inputPlaceholder:
        'my-image|https://my-oss.com/my.png\nhttps://my-oss.com/my_1.png',
      inputPattern: /\S/,
      inputErrorMessage: 'Required'
    })
      .then(async ({ value }) => {
        let data = value.split('\n')
        let importData = []
        data.forEach((item) => {
          let oneData = item.trim().split('|')
          let url, name
          if (oneData.length > 1) {
            name = oneData[0].trim()
            url = oneData[1]
          } else {
            url = oneData[0].trim()
            let str = url.substring(url.lastIndexOf('/') + 1)
            name = str.substring(0, str.lastIndexOf('.'))
          }
          if (url) {
            importData.push({
              name: name,
              url: url,
              classId: search.value.classId,
              tag: url.substring(url.lastIndexOf('.') + 1),
              key: CreateUUID()
            })
          }
        })

        const res = await importURL(importData)
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: 'Imported'
          })
          await getTableData()
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: 'Cancelled'
        })
      })
  }

  const onSuccess = () => {
    search.value.keyword = null
    page.value = 1
    getTableData()
  }

  const defaultProps = {
    children: 'children',
    label: 'name',
    value: 'ID'
  }

  const categories = ref([])
  const fetchCategories = async () => {
    const res = await getCategoryList()
    let data = {
      name: 'All categories',
      ID: 0,
      pid: 0,
      children: []
    }
    if (res.code === 0) {
      categories.value = res.data || []
      categories.value.unshift(data)
    }
  }

  const handleNodeClick = (node) => {
    search.value.keyword = null
    search.value.classId = node.ID
    page.value = 1
    getTableData()
  }

  const categoryDialogVisible = ref(false)
  const categoryFormData = ref({
    ID: 0,
    pid: 0,
    name: ''
  })

  const categoryForm = ref(null)
  const rules = ref({
    name: [
      { required: true, message: 'Category name is required', trigger: 'blur' },
      { max: 20, message: 'Max 20 characters', trigger: 'blur' }
    ]
  })

  const addCategoryFun = (category) => {
    categoryDialogVisible.value = true
    categoryFormData.value.ID = 0
    categoryFormData.value.pid = category.ID
  }

  const editCategory = (category) => {
    categoryFormData.value = {
      ID: category.ID,
      pid: category.pid,
      name: category.name
    }
    categoryDialogVisible.value = true
  }

  const deleteCategoryFun = async (id) => {
    const res = await deleteCategory({ id: id })
    if (res.code === 0) {
      ElMessage.success({ type: 'success', message: 'Deleted' })
      await fetchCategories()
    }
  }

  const confirmAddCategory = async () => {
    categoryForm.value.validate(async (valid) => {
      if (valid) {
        const res = await addCategory(categoryFormData.value)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: 'Success' })
          await fetchCategories()
          closeAddCategoryDialog()
        }
      }
    })
  }

  const closeAddCategoryDialog = () => {
    categoryDialogVisible.value = false
    categoryFormData.value = {
      ID: 0,
      pid: 0,
      name: ''
    }
  }

  fetchCategories()
</script>
