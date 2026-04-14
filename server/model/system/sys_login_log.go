package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysLoginLog struct {
	global.GVA_MODEL
	Username      string  `json:"username" gorm:"column:username;comment:username"`
	Ip            string  `json:"ip" gorm:"column:ip;comment:request IP"`
	Status        bool    `json:"status" gorm:"column:status;comment:Loginstatus"`
	ErrorMessage  string  `json:"errorMessage" gorm:"column:error_message;comment:error message"`
	Agent         string  `json:"agent" gorm:"column:agent;comment:agent"`
	UserID        uint    `json:"userId" gorm:"column:user_id;comment:user ID"`
	User          SysUser `json:"user" gorm:"foreignKey:UserID"`
}
