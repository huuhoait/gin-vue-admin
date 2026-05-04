package service

import (
	"context"
	"runtime"
	"strings"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/shirou/gopsutil/v3/host"
)

type monitor struct{}

// ServerStats reuses utils.InitOS/CPU/RAM/Disk and adds host uptime.
type ServerStats struct {
	Os       utils.Os     `json:"os"`
	Cpu      utils.Cpu    `json:"cpu"`
	Ram      utils.Ram    `json:"ram"`
	Disk     []utils.Disk `json:"disk"`
	Hostname string       `json:"hostname"`
	UptimeS  uint64       `json:"uptimeSeconds"`
	BootedAt time.Time    `json:"bootedAt"`
}

func (m *monitor) Server() (ServerStats, error) {
	var s ServerStats
	s.Os = utils.InitOS()
	cpuStats, err := utils.InitCPU()
	if err != nil {
		return s, err
	}
	s.Cpu = cpuStats
	ramStats, err := utils.InitRAM()
	if err != nil {
		return s, err
	}
	s.Ram = ramStats
	// Disk depends on config.DiskList — empty slice when unconfigured is fine.
	if diskStats, err := utils.InitDisk(); err == nil {
		s.Disk = diskStats
	}
	if info, err := host.Info(); err == nil {
		s.Hostname = info.Hostname
		s.UptimeS = info.Uptime
		s.BootedAt = time.Unix(int64(info.BootTime), 0).UTC()
	}
	return s, nil
}

// RuntimeStats captures the Go process state.
type RuntimeStats struct {
	GoVersion    string `json:"goVersion"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	NumCPU       int    `json:"numCpu"`
	NumGoroutine int    `json:"numGoroutine"`
	GOMAXPROCS   int    `json:"gomaxprocs"`
	NumGC        uint32 `json:"numGc"`
	HeapAllocMB  uint64 `json:"heapAllocMb"`
	HeapSysMB    uint64 `json:"heapSysMb"`
	StackInuseMB uint64 `json:"stackInuseMb"`
	NextGCMB     uint64 `json:"nextGcMb"`
	LastGC       uint64 `json:"lastGcUnixNanos"`
}

func (m *monitor) Runtime() RuntimeStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	const mb = 1024 * 1024
	return RuntimeStats{
		GoVersion:    runtime.Version(),
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		GOMAXPROCS:   runtime.GOMAXPROCS(0),
		NumGC:        ms.NumGC,
		HeapAllocMB:  ms.HeapAlloc / mb,
		HeapSysMB:    ms.HeapSys / mb,
		StackInuseMB: ms.StackInuse / mb,
		NextGCMB:     ms.NextGC / mb,
		LastGC:       ms.LastGC,
	}
}

// CacheStats summarises Redis INFO. Returns IsConnected=false when Redis is
// unconfigured so the dashboard can show a placeholder rather than an error.
type CacheStats struct {
	IsConnected     bool              `json:"isConnected"`
	Version         string            `json:"version"`
	UptimeSeconds   string            `json:"uptimeSeconds"`
	UsedMemoryHuman string            `json:"usedMemoryHuman"`
	ConnectedClient string            `json:"connectedClients"`
	OpsPerSec       string            `json:"opsPerSec"`
	HitRate         string            `json:"hitRate"`
	DBKeys          map[string]string `json:"dbKeys"`
	Raw             map[string]string `json:"raw,omitempty"`
}

func (m *monitor) Cache(ctx context.Context) (CacheStats, error) {
	var c CacheStats
	if global.GVA_REDIS == nil {
		return c, nil
	}
	raw, err := global.GVA_REDIS.Info(ctx).Result()
	if err != nil {
		return c, err
	}
	parsed := parseRedisInfo(raw)
	c.IsConnected = true
	c.Version = parsed["redis_version"]
	c.UptimeSeconds = parsed["uptime_in_seconds"]
	c.UsedMemoryHuman = parsed["used_memory_human"]
	c.ConnectedClient = parsed["connected_clients"]
	c.OpsPerSec = parsed["instantaneous_ops_per_sec"]
	hits := parsed["keyspace_hits"]
	misses := parsed["keyspace_misses"]
	c.HitRate = computeHitRate(hits, misses)
	c.DBKeys = make(map[string]string)
	for k, v := range parsed {
		if strings.HasPrefix(k, "db") {
			c.DBKeys[k] = v
		}
	}
	return c, nil
}

func parseRedisInfo(raw string) map[string]string {
	out := make(map[string]string, 64)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.IndexByte(line, ':'); i > 0 {
			out[line[:i]] = line[i+1:]
		}
	}
	return out
}

func computeHitRate(hitsStr, missesStr string) string {
	var hits, misses uint64
	for _, c := range hitsStr {
		if c < '0' || c > '9' {
			return ""
		}
		hits = hits*10 + uint64(c-'0')
	}
	for _, c := range missesStr {
		if c < '0' || c > '9' {
			return ""
		}
		misses = misses*10 + uint64(c-'0')
	}
	total := hits + misses
	if total == 0 {
		return "n/a"
	}
	pct := float64(hits) / float64(total) * 100
	// Format manually to avoid an fmt import for one call.
	whole := int(pct)
	frac := int((pct - float64(whole)) * 10)
	return itoa(whole) + "." + itoa(frac) + "%"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
