package initialize

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	pluginApi "github.com/huuhoait/gin-vue-admin/server/plugin/skyagent/api"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Readonly opens read-only Postgres connections from config and injects them
// into the dashboard API layer. Call after Viper config is loaded. Does not
// panic — dashboard degrades gracefully if connections are unavailable.
func Readonly() {
	cfg := global.GVA_CONFIG.Proxy

	if cfg.CoreDBDSN != "" {
		db, err := openReadonly(cfg.CoreDBDSN)
		if err != nil {
			global.GVA_LOG.Warn("skyagent: core read-only DB unavailable", zap.Error(err))
		} else {
			pluginApi.DashboardDBs.CoreDB = db
			global.GVA_LOG.Info("skyagent: core read-only DB connected")
		}
	}

	if cfg.OrderDBDSN != "" {
		db, err := openReadonly(cfg.OrderDBDSN)
		if err != nil {
			global.GVA_LOG.Warn("skyagent: order read-only DB unavailable", zap.Error(err))
		} else {
			pluginApi.DashboardDBs.OrderDB = db
			global.GVA_LOG.Info("skyagent: order read-only DB connected")
		}
	}
}

func openReadonly(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(5)
	return db, nil
}
