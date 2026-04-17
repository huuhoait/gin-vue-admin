package initialize

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// Reload ExcellentelegantLocationRe-LoadSystem Configuration
func Reload() error {
	global.GVA_LOG.Info("PositiveAtRe-LoadSystem Configuration...")

	// Re-LoadconfigurationFile
	if err := global.GVA_VP.ReadInConfig(); err != nil {
		global.GVA_LOG.Error("Re-ReadconfigurationFilefailed!", zap.Error(err))
		return err
	}

	// Re-initialize databaseConnection
	if global.GVA_DB != nil {
		db, _ := global.GVA_DB.DB()
		err := db.Close()
		if err != nil {
			global.GVA_LOG.Error("DisableOriginalDatabaseConnectionfailed!", zap.Error(err))
			return err
		}
	}

	// RepeatcreatestandDatabaseConnection
	global.GVA_DB = Gorm()

	// Re-initializeOthersconfiguration
	OtherInit()
	DBList()

	if global.GVA_DB != nil {
		// ensureDatabaseTableStructureYesLatestof
		RegisterTables()
	}

	// Re-initializescheduled task
	Timer()

	global.GVA_LOG.Info("System ConfigurationRe-LoadComplete")
	return nil
}
