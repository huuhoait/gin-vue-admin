package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type ExaCustomer struct {
	global.GVA_MODEL
	CustomerName       string         `json:"customerName" form:"customerName" gorm:"comment:customer name"`                // customer name
	CustomerPhoneData  string         `json:"customerPhoneData" form:"customerPhoneData" gorm:"comment:customer phone number"`    // customer phone number
	SysUserID          uint           `json:"sysUserId" form:"sysUserId" gorm:"comment:admin ID"`                     // admin ID
	SysUserAuthorityID uint           `json:"sysUserAuthorityID" form:"sysUserAuthorityID" gorm:"comment:admin role ID"` // admin role ID
	SysUser            system.SysUser `json:"sysUser" form:"sysUser" gorm:"comment:admin details"`                         // admin details
}
