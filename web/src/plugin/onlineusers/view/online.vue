<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item :label="t('admin.plugin.online_users.username_label')">
          <el-input
            v-model="searchInfo.username"
            :placeholder="t('admin.plugin.online_users.username_placeholder')"
            clearable
            @keyup.enter="onSubmit"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">{{ t('admin.plugin.online_users.search') }}</el-button>
          <el-button icon="refresh" @click="onReset">{{ t('admin.plugin.online_users.reset') }}</el-button>
          <el-button icon="refresh-right" @click="loadList">{{ t('admin.plugin.online_users.refresh') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <span class="text-sm text-gray-500">
          {{ t('admin.plugin.online_users.active_count', { n: total }) }}
        </span>
      </div>

      <el-table
        :data="tableData"
        style="width: 100%"
        row-key="uuid"
        tooltip-effect="dark"
      >
        <el-table-column :label="t('admin.plugin.online_users.col_user')" min-width="180">
          <template #default="scope">
            <div>{{ scope.row.username }}<span v-if="scope.row.nickName"> ({{ scope.row.nickName }})</span></div>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.online_users.col_ip')" prop="ip" width="160" />
        <el-table-column :label="t('admin.plugin.online_users.col_user_agent')" prop="userAgent" min-width="240" show-overflow-tooltip />
        <el-table-column :label="t('admin.plugin.online_users.col_logged_in')" width="180">
          <template #default="scope">{{ formatDate(scope.row.loginAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.online_users.col_last_seen')" width="180">
          <template #default="scope">
            <el-tooltip :content="formatDate(scope.row.lastSeenAt)">
              <span>{{ relativeTime(scope.row.lastSeenAt) }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column :label="t('admin.plugin.online_users.col_actions')" width="120" fixed="right">
          <template #default="scope">
            <el-button
              link
              type="danger"
              icon="circle-close"
              :disabled="scope.row.uuid === currentUuid"
              @click="onKick(scope.row)"
            >
              {{ t('admin.plugin.online_users.kick') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :page-sizes="[10, 30, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/format'
import { useUserStore } from '@/pinia/modules/user'
import { listOnlineSessions, kickOnlineSession } from '@/plugin/onlineusers/api/session'

defineOptions({ name: 'OnlineUsers' })

const { t } = useI18n()
const userStore = useUserStore()
const { userInfo } = storeToRefs(userStore)
const currentUuid = ref(userInfo.value?.uuid || '')

const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchInfo = ref({ username: '' })

const loadList = async () => {
  const res = await listOnlineSessions({
    page: page.value,
    pageSize: pageSize.value,
    username: searchInfo.value.username
  })
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
  } else {
    ElMessage.error(res.msg || t('admin.plugin.online_users.load_failed'))
  }
}

const onSubmit = () => {
  page.value = 1
  loadList()
}

const onReset = () => {
  searchInfo.value.username = ''
  page.value = 1
  loadList()
}

const handleCurrentChange = (val) => {
  page.value = val
  loadList()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  loadList()
}

const onKick = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('admin.plugin.online_users.kick_confirm_body', { username: row.username }),
      t('admin.plugin.online_users.kick_confirm_title'),
      {
        confirmButtonText: t('admin.plugin.online_users.kick_confirm_ok'),
        cancelButtonText: t('admin.common.cancel'),
        type: 'warning'
      }
    )
  } catch {
    return
  }
  const res = await kickOnlineSession({ uuid: row.uuid })
  if (res.code === 0) {
    ElMessage.success(t('admin.plugin.online_users.kick_success'))
    loadList()
  } else {
    ElMessage.error(res.msg || t('admin.plugin.online_users.kick_failed'))
  }
}

const relativeTime = (iso) => {
  if (!iso) return ''
  const diffMs = Date.now() - new Date(iso).getTime()
  const sec = Math.floor(diffMs / 1000)
  if (sec < 60) return t('admin.plugin.online_users.ago_seconds', { n: sec })
  const min = Math.floor(sec / 60)
  if (min < 60) return t('admin.plugin.online_users.ago_minutes', { n: min })
  const hr = Math.floor(min / 60)
  if (hr < 24) return t('admin.plugin.online_users.ago_hours', { n: hr })
  return t('admin.plugin.online_users.ago_days', { n: Math.floor(hr / 24) })
}

let timer = null
onMounted(() => {
  loadList()
  timer = setInterval(loadList, 30_000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>
