package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/onlineusers/service"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const sessionTTL = 30 * time.Minute

// Timer registers a 1-minute prune task that drops sessions whose LastSeenAt
// is older than sessionTTL. Idempotent: safe to run on every plugin reload.
func Timer() {
	go func() {
		_, err := global.GVA_Timer.AddTaskByFunc(
			"OnlineUsersPrune",
			"@every 1m",
			func() {
				n, err := service.Service.Session.Prune(context.Background(), sessionTTL)
				if err != nil {
					global.GVA_LOG.Warn("online-users prune failed", zap.Error(err))
					return
				}
				if n > 0 {
					global.GVA_LOG.Debug("online-users pruned", zap.Int("count", n))
				}
			},
			"Drop stale online-user sessions older than 30m",
			cron.WithSeconds(),
		)
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("register OnlineUsersPrune task failed: %v", err))
		}
	}()
}
