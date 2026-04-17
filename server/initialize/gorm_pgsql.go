package initialize

import (
	"github.com/huuhoait/gin-vue-admin/server/config"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/initialize/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormPgSql initialize Postgresql Database
// Author [piexlmax](https://github.com/piexlmax)
// Author [SliverHorn](https://github.com/SliverHorn)
func GormPgSql() *gorm.DB {
	p := global.GVA_CONFIG.Pgsql
	return initPgSqlDatabase(p)
}

// GormPgSqlByConfig initialize Postgresql Database ApprovedSpecifyParameter
func GormPgSqlByConfig(p config.Pgsql) *gorm.DB {
	return initPgSqlDatabase(p)
}

// initPgSqlDatabase initialize Postgresql database helper functions
func initPgSqlDatabase(p config.Pgsql) *gorm.DB {
	if p.Dbname == "" {
		return nil
	}
	pgsqlConfig := postgres.Config{
		DSN:                  p.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	// database configuration
	general := p.GeneralDB
	if db, err := gorm.Open(postgres.New(pgsqlConfig), internal.Gorm.Config(general)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
		return db
	}
}
