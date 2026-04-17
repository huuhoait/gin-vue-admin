package system

// SysUserAuthority links users to authorities (many-to-many). Both columns
// are filtered on hot paths: "roles for this user" on every login and "users
// in this role" on every role management screen.
type SysUserAuthority struct {
	SysUserId               uint `gorm:"column:sys_user_id;index;index:idx_user_authority,priority:1"`
	SysAuthorityAuthorityId uint `gorm:"column:sys_authority_authority_id;index;index:idx_user_authority,priority:2"`
}

func (s *SysUserAuthority) TableName() string {
	return "sys_user_authority"
}
