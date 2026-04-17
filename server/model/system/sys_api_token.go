package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"time"
)

type SysApiToken struct {
	global.GVA_MODEL
	UserID      uint      `json:"userId" gorm:"comment:User ID"`
	User        SysUser   `json:"user" gorm:"foreignKey:UserID;"`
	AuthorityID uint      `json:"authorityId" gorm:"comment:role ID"`
	Token       string    `json:"token" gorm:"type:text;comment:Token"`
	Status      bool      `json:"status" gorm:"default:true;comment:status"` // trueHaveEffect falseNoneEffect
	ExpiresAt   time.Time `json:"expiresAt" gorm:"comment:expiration time"`
	Remark      string    `json:"remark" gorm:"comment:Remark"`
}
