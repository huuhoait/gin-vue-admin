package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize/internal"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// GormSqlite initializeSqliteDatabase
func GormSqlite() *gorm.DB {
	s := global.GVA_CONFIG.Sqlite
	return initSqliteDatabase(s)
}

// GormSqliteByConfig initializeSqlitedatabase name passed inconfiguration
func GormSqliteByConfig(s config.Sqlite) *gorm.DB {
	return initSqliteDatabase(s)
}

// initSqliteDatabase initializeSqliteDatabaseAssisthelpFunctionNumber
func initSqliteDatabase(s config.Sqlite) *gorm.DB {
	if s.Dbname == "" {
		return nil
	}

	// database configuration
	general := s.GeneralDB
	if db, err := gorm.Open(sqlite.Open(s.Dsn()), internal.Gorm.Config(general)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
		sqlDB.SetMaxOpenConns(s.MaxOpenConns)
		return db
	}
}
