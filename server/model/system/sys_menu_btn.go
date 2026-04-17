package system

import "github.com/huuhoait/gin-vue-admin/server/global"

type SysBaseMenuBtn struct {
	global.GVA_MODEL
	Name          string `json:"name" gorm:"comment:ButtonCloseKeykey"`
	Desc          string `json:"desc" gorm:"ButtonRemark"`
	SysBaseMenuID uint   `json:"sysBaseMenuID" gorm:"comment:menu ID"`
}
