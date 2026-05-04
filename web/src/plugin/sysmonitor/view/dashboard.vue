<template>
  <div class="sysmonitor-dashboard">
    <div class="gva-btn-list flex justify-between items-center mb-4">
      <span class="text-sm text-gray-500">{{ t('admin.plugin.sysmonitor.auto_refresh_10s') }}</span>
      <el-button icon="refresh-right" @click="loadAll">{{ t('admin.plugin.sysmonitor.refresh_now') }}</el-button>
    </div>

    <el-row :gutter="16">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header><strong>{{ t('admin.plugin.sysmonitor.host') }}</strong></template>
          <el-descriptions v-if="server" :column="2" border size="small">
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.hostname')">{{ server.hostname || '—' }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.os')">{{ server.os.goos }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.cores')">{{ server.os.numCpu }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.goroutines')">{{ server.os.numGoroutine }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.uptime')">{{ formatUptime(server.uptimeSeconds) }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.booted')">{{ server.bootedAt ? formatDate(server.bootedAt) : '—' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card shadow="hover">
          <template #header><strong>{{ t('admin.plugin.sysmonitor.ram') }}</strong></template>
          <div v-if="server">
            <el-progress :percentage="server.ram.usedPercent" :stroke-width="14" />
            <div class="mt-2 text-sm text-gray-500">
              {{ server.ram.usedMb }} MB / {{ server.ram.totalMb }} MB
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="mt-4">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <strong>{{ t('admin.plugin.sysmonitor.cpu') }}</strong>
            {{ t('admin.plugin.sysmonitor.cpu_cores_suffix', { n: server?.cpu?.cores || 0 }) }}
          </template>
          <div v-if="server" class="grid grid-cols-4 gap-3">
            <div v-for="(pct, idx) in server.cpu.cpus" :key="idx" class="text-center">
              <div class="text-xs text-gray-500 mb-1">{{ t('admin.plugin.sysmonitor.core_label', { n: idx }) }}</div>
              <el-progress type="dashboard" :percentage="Math.round(pct)" :width="80" />
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row v-if="server?.disk?.length" :gutter="16" class="mt-4">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header><strong>{{ t('admin.plugin.sysmonitor.disk') }}</strong></template>
          <el-table :data="server.disk" size="small">
            <el-table-column prop="mountPoint" :label="t('admin.plugin.sysmonitor.mount')" width="200" />
            <el-table-column :label="t('admin.plugin.sysmonitor.used')">
              <template #default="scope">
                <el-progress :percentage="scope.row.usedPercent" />
              </template>
            </el-table-column>
            <el-table-column :label="t('admin.plugin.sysmonitor.capacity')" width="200">
              <template #default="scope">
                {{ scope.row.usedGb }} GB / {{ scope.row.totalGb }} GB
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="mt-4">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header><strong>{{ t('admin.plugin.sysmonitor.go_runtime') }}</strong></template>
          <el-descriptions v-if="rt" :column="2" border size="small">
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.version')">{{ rt.goVersion }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.gomaxprocs')">{{ rt.gomaxprocs }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.goroutines')">{{ rt.numGoroutine }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.gc_count')">{{ rt.numGc }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.heap_alloc')">{{ rt.heapAllocMb }} MB</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.heap_sys')">{{ rt.heapSysMb }} MB</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.stack_inuse')">{{ rt.stackInuseMb }} MB</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.next_gc')">{{ rt.nextGcMb }} MB</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card shadow="hover">
          <template #header><strong>{{ t('admin.plugin.sysmonitor.cache_redis') }}</strong></template>
          <div v-if="!cache?.isConnected" class="text-gray-500 text-sm">
            {{ t('admin.plugin.sysmonitor.redis_not_configured') }}
          </div>
          <el-descriptions v-else :column="2" border size="small">
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.version')">{{ cache.version }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.uptime')">{{ formatUptime(cache.uptimeSeconds) }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.memory')">{{ cache.usedMemoryHuman }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.clients')">{{ cache.connectedClients }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.ops_per_sec')">{{ cache.opsPerSec }}</el-descriptions-item>
            <el-descriptions-item :label="t('admin.plugin.sysmonitor.hit_rate')">{{ cache.hitRate }}</el-descriptions-item>
            <el-descriptions-item v-for="(v, k) in cache.dbKeys" :key="k" :label="k">{{ v }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/format'
import {
  getServerStats,
  getRuntimeStats,
  getCacheStats
} from '@/plugin/sysmonitor/api/monitor'

defineOptions({ name: 'SysMonitor' })

const { t } = useI18n()
const server = ref(null)
const rt = ref(null)
const cache = ref(null)

const loadAll = async () => {
  const [srv, rtRes, cacheRes] = await Promise.all([
    getServerStats(),
    getRuntimeStats(),
    getCacheStats()
  ])
  if (srv.code === 0) server.value = srv.data
  else ElMessage.error(srv.msg || t('admin.plugin.sysmonitor.server_failed'))
  if (rtRes.code === 0) rt.value = rtRes.data
  if (cacheRes.code === 0) cache.value = cacheRes.data
}

const formatUptime = (raw) => {
  const sec = Number(raw)
  if (!sec || Number.isNaN(sec)) return '—'
  const d = Math.floor(sec / 86400)
  const h = Math.floor((sec % 86400) / 3600)
  const m = Math.floor((sec % 3600) / 60)
  return `${d}d ${h}h ${m}m`
}

let timer = null
onMounted(() => {
  loadAll()
  timer = setInterval(loadAll, 10_000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.sysmonitor-dashboard :deep(.el-progress-bar__outer) {
  background-color: var(--el-color-info-light-9);
}
</style>
