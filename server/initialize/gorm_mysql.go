package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize/internal"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormMysql initializeMysqlDatabase
// Author [piexlmax](https://github.com/piexlmax)
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [ByteZhou-2018](https://github.com/ByteZhou-2018)
func GormMysql() *gorm.DB {
	m := global.GVA_CONFIG.Mysql
	return initMysqlDatabase(m)
}

// GormMysqlByConfig ApprovedTransmitInconfigurationinitializeMysqlDatabase
func GormMysqlByConfig(m config.Mysql) *gorm.DB {
	return initMysqlDatabase(m)
}

// initMysqlDatabase initializeMysqldatabase helper functions
func initMysqlDatabase(m config.Mysql) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}

	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string typeFieldofdefaultLength
		SkipInitializeWithVersion: false,   // According toVersionAutomaticconfiguration
	}
	// database configuration
	general := m.GeneralDB
	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(general)); err != nil {
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
