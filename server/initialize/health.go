package initialize

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

// RegisterHealthRoutes adds liveness and readiness probes.
// /health/live  — process is alive (always 200 if responding)
// /health/ready — dependencies (DB + Redis) are reachable (503 if not)
func RegisterHealthRoutes(r *gin.Engine, prefix string) {
	r.GET(prefix+"/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET(prefix+"/health/ready", func(c *gin.Context) {
		report := gin.H{}
		allOK := true

		// DB check
		if global.GVA_DB != nil {
			sqlDB, err := global.GVA_DB.DB()
			if err == nil {
				ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
				defer cancel()
				err = sqlDB.PingContext(ctx)
			}
			if err != nil {
				allOK = false
				report["db"] = "unreachable: " + err.Error()
				global.GVA_LOG.Warn("readiness: db ping failed", zap.Error(err))
			} else {
				report["db"] = "ok"
			}
		} else {
			report["db"] = "not configured"
		}

		// Redis check
		if global.GVA_REDIS != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()
			if err := global.GVA_REDIS.Ping(ctx).Err(); err != nil {
				allOK = false
				report["redis"] = "unreachable: " + err.Error()
				global.GVA_LOG.Warn("readiness: redis ping failed", zap.Error(err))
			} else {
				report["redis"] = "ok"
			}
		} else {
			report["redis"] = "not configured"
		}

		status := http.StatusOK
		if !allOK {
			status = http.StatusServiceUnavailable
		}
		report["status"] = map[bool]string{true: "ready", false: "degraded"}[allOK]
		c.JSON(status, report)
	})
}
