package system

import (
	"time"
)

type SysAuthority struct {
	CreatedAt       time.Time       // created at
	UpdatedAt       time.Time       // updated at
	DeletedAt       *time.Time      `sql:"index"`
	AuthorityId     uint            `json:"authorityId" gorm:"not null;unique;primary_key;comment:role ID;size:90"` // role ID
	AuthorityName   string          `json:"authorityName" gorm:"comment:role name"`                                    // role name
	ParentId        *uint           `json:"parentId" gorm:"comment:parent role ID"`                                       // parent role ID
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;"`
	Children        []SysAuthority  `json:"children" gorm:"-"`
	SysBaseMenus    []SysBaseMenu   `json:"menus" gorm:"many2many:sys_authority_menus;"`
	Users           []SysUser       `json:"-" gorm:"many2many:sys_user_authority;"`
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:default menu;default:dashboard"` // default menu(defaultdashboard)
}

func (SysAuthority) TableName() string {
	return "sys_authorities"
}
