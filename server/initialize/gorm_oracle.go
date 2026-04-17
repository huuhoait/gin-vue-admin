package initialize

import (
	oracle "github.com/dzwvip/gorm-oracle"
	"github.com/huuhoait/gin-vue-admin/server/config"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/initialize/internal"
	"gorm.io/gorm"
)

// GormOracle initializeoracleDatabase
func GormOracle() *gorm.DB {
	m := global.GVA_CONFIG.Oracle
	return initOracleDatabase(m)
}

// GormOracleByConfig initializeOracledatabase name passed inconfiguration
func GormOracleByConfig(m config.Oracle) *gorm.DB {
	return initOracleDatabase(m)
}

// initOracleDatabase initializeOracledatabase helper functions
func initOracleDatabase(m config.Oracle) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	// database configuration
	general := m.GeneralDB
	if db, err := gorm.Open(oracle.Open(m.Dsn()), internal.Gorm.Config(general)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
