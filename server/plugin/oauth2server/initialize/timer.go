package initialize

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/service"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Timer prunes expired auth codes + fully expired tokens hourly. Revoked
// tokens that haven't reached refresh expiry are kept for the introspection
// endpoint so callers can still see {active:false}.
func Timer() {
	go func() {
		_, err := global.GVA_Timer.AddTaskByFunc(
			"OAuth2Prune",
			"@hourly",
			func() {
				n, err := service.Service.OAuth2.PruneExpired()
				if err != nil {
					global.GVA_LOG.Warn("oauth2 prune failed", zap.Error(err))
					return
				}
				if n > 0 {
					global.GVA_LOG.Info("oauth2 pruned expired records", zap.Int64("rows", n))
				}
			},
			"Drop expired auth codes and tokens",
			cron.WithSeconds(),
		)
		if err != nil {
			global.GVA_LOG.Warn("register OAuth2Prune task failed", zap.Error(err))
		}
	}()
}
