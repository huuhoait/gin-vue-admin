package system

type SysAuthorityBtn struct {
	AuthorityId      uint           `gorm:"comment:role ID;index;index:idx_authority_menu_btn,priority:1"`
	SysMenuID        uint           `gorm:"comment:menu ID;index:idx_authority_menu_btn,priority:2"`
	SysBaseMenuBtnID uint           `gorm:"comment:menu buttonsID"`
	SysBaseMenuBtn   SysBaseMenuBtn ` gorm:"comment:ButtonDetails"`
}
