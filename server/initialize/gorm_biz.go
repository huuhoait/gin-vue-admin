package initialize

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	_ "github.com/huuhoait/gin-vue-admin/server/model/example"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}
