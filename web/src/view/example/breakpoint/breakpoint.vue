<template>
  <div class="break-point">
    <div class="gva-table-box">
      <el-divider content-position="left">Large file upload</el-divider>
      <form id="fromCont" method="post">
        <!-- Button container (Flexbox alignment) -->
        <div class="button-container">
          <div class="fileUpload" @click="inputChange">
            <span class="takeFile">Choose file</span>
            <input
              v-show="false"
              id="file"
              ref="FileInput"
              multiple="multiple"
              type="file"
              @change="choseFile"
            />
          </div>
          <el-button
            :disabled="limitFileSize"
            type="primary"
            class="uploadBtn"
            @click="getFile"
          >Upload</el-button>
        </div>
      </form>
      <div class="el-upload__tip">Please upload a file up to 5MB</div>
      <div class="list">
        <transition name="list" tag="p">
          <div v-if="file" class="list-item">
            <el-icon>
              <document />
            </el-icon>
            <span>{{ file.name }}</span>
            <span class="percentage">{{ percentage }}%</span>
            <el-progress
              :show-text="false"
              :text-inside="false"
              :stroke-width="2"
              :percentage="percentage"
            />
          </div>
        </transition>
      </div>
      <div class="tips">
        This is a preview build for testing. Styling and performance are still in progress. Uploaded chunks and merged files are stored under QMPlusserver (breakpointDir / fileDir).
      </div>
    </div>
  </div>
</template>

<script setup>
import SparkMD5 from 'spark-md5'
import {
  findFile,
  breakpointContinueFinish,
  removeChunk,
  breakpointContinue
} from '@/api/breakpoint'
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'

defineOptions({
  name: 'BreakPoint'
})

const file = ref(null)
const fileMd5 = ref('')
const formDataList = ref([])
const waitUpLoad = ref([])
const waitNum = ref(NaN)
const limitFileSize = ref(false)
const percentage = ref(0)
const percentageFlage = ref(true)

// Handle file selection
const choseFile = async (e) => {
  // If user cancels selection, return
  if (!e.target.files.length) {
    return
  }
  const fileR = new FileReader() // read file stream
  const fileInput = e.target.files[0] // current file
  const maxSize = 5 * 1024 * 1024
  file.value = fileInput // keep file in ref for later use
  percentage.value = 0
  if (file.value.size < maxSize) {
    fileR.readAsArrayBuffer(file.value) // read as ArrayBuffer to align with backend stream
    fileR.onload = async (e) => {
      // ArrayBuffer callback; stream is in e.target.result
      const blob = e.target.result
      const spark = new SparkMD5.ArrayBuffer() // build md5 for file identity
      spark.append(blob)
      fileMd5.value = spark.end() // full file md5
      const FileSliceCap = 1 * 1024 * 1024 // chunk size (bytes)
      let start = 0 // chunk start offset
      let end = 0 // chunk end offset
      let i = 0 // chunk index
      formDataList.value = [] // chunk payload pool
      while (end < file.value.size) {
        // slice until end exceeds file size
        start = i * FileSliceCap
        end = (i + 1) * FileSliceCap
        var fileSlice = file.value.slice(start, end) // h5 slice(start,end)
        const formData = new window.FormData()
        formData.append('fileMd5', fileMd5.value)
        formData.append('file', fileSlice)
        formData.append('chunkNumber', i)
        formData.append('fileName', file.value.name)
        formDataList.value.push({ key: i, formData })
        i++
      }
      const params = {
        fileName: file.value.name,
        fileMd5: fileMd5.value,
        chunkTotal: formDataList.value.length
      }
      const res = await findFile(params)
      // Query backend for existing chunks to support resume
      const finishList = res.data.file.ExaFileChunk // uploaded chunks
      const IsFinish = res.data.file.IsFinish // same file md5 but different name => instant upload
      if (!IsFinish) {
        // Resume upload: filter chunks not uploaded yet
        waitUpLoad.value = formDataList.value.filter((all) => {
          return !(
            finishList &&
            finishList.some((fi) => fi.FileChunkNumber === all.key)
          ) // chunks still needed
        })
      } else {
        waitUpLoad.value = [] // instant upload: nothing to upload
        ElMessage.success('Instant upload completed')
      }
      waitNum.value = waitUpLoad.value.length // track for progress display
    }
  } else {
    limitFileSize.value = true
    ElMessage('Please upload a file smaller than 5MB')
  }
}

const getFile = () => {
  // Upload action
  if (file.value === null) {
    ElMessage('Please choose a file first')
    return
  }
  // Check progress
  if (percentage.value === 100) {
    ElMessage.success('Upload already completed')
    percentageFlage.value = false
    return // stop if already completed
  }
  // Continue uploading chunks
  sliceFile()
}

const sliceFile = () => {
  waitUpLoad.value &&
  waitUpLoad.value.forEach((item) => {
    // Upload remaining chunks
    item.formData.append('chunkTotal', formDataList.value.length)
    const fileR = new FileReader()
    const fileF = item.formData.get('file')
    fileR.readAsArrayBuffer(fileF)
    fileR.onload = (e) => {
      const spark = new SparkMD5.ArrayBuffer()
      spark.append(e.target.result)
      item.formData.append('chunkMd5', spark.end()) // chunk md5 for integrity
      upLoadFileSlice(item)
    }
  })
}

watch(
  () => waitNum.value,
  () => {
    percentage.value = Math.floor(
      ((formDataList.value.length - waitNum.value) /
        formDataList.value.length) *
      100
    )
  }
)

const upLoadFileSlice = async (item) => {
  // Upload a chunk
  const fileRe = await breakpointContinue(item.formData)
  if (fileRe.code !== 0) {
    return
  }
  waitNum.value--
  if (waitNum.value === 0) {
    // Merge file after all chunks uploaded
    const params = {
      fileName: file.value.name,
      fileMd5: fileMd5.value
    }
    const res = await breakpointContinueFinish(params)
    if (res.code === 0) {
      // Remove cached chunks after merge
      const params = {
        fileName: file.value.name,
        fileMd5: fileMd5.value,
        filePath: res.data.filePath
      }
      ElMessage.success('Upload succeeded')
      await removeChunk(params)
    }
  }
}

const FileInput = ref(null)
const inputChange = () => {
  FileInput.value.dispatchEvent(new MouseEvent('click'))
}
</script>

<style lang="scss" scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
#fromCont {
  display: inline-block;
}

.gva-table-box {
  display: block;
}

.button-container {
  display: flex;
  align-items: center;
}

.fileUpload,
.uploadBtn {
  width: 90px;
  height: 35px;
  line-height: 35px;
  font-size: 14px;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  border-radius: 5px;
  cursor: pointer;
}

.fileUpload {
  padding: 0 15px;
  background-color: #007bff;
  color: #ffffff;
  font-weight: 500;
  transition: all 0.3s ease-in-out;
  margin-right: 5px;
}

.uploadBtn {
  background-color: #007bff;
  color: #fff;
  margin-left: 10px;
}

.fileUpload:hover {
  background-color: #0056b3;
}

.uploadBtn:hover {
  background-color: #0056b3;
}


.fileUpload:active,
.uploadBtn:active {
  transform: translateY(2px);
}

.fileUpload input {
  position: relative;
  font-size: 100px;
  right: 0;
  top: 0;
  opacity: 0;
  cursor: pointer;
  width: 100%;
  height: 100%;
}



.fileName {
  display: inline-block;
  vertical-align: top;
  margin: 6px 15px 0 15px;
}
.tips {
  margin-top: 30px;
  font-size: 14px;
  font-weight: 400;
  color: #606266;
}
.el-divider {
  margin: 0 0 30px 0;
}

.list {
  margin-top: 15px;
}
.list-item {
  display: block;
  margin-right: 10px;
  color: #606266;
  line-height: 25px;
  margin-bottom: 5px;
  width: 40%;
  .percentage {
    float: right;
  }
}
.list-enter-active,
.list-leave-active {
  transition: all 1s;
}
.list-enter, .list-leave-to
  /* .list-leave-active for below version 2.1.8 */ {
  opacity: 0;
  transform: translateY(-30px);
}
</style>