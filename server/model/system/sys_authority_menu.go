package system

type SysMenu struct {
	SysBaseMenu
	MenuId      uint                   `json:"menuId" gorm:"comment:menu ID"`
	AuthorityId uint                   `json:"-" gorm:"comment:role ID"`
	Children    []SysMenu              `json:"children" gorm:"-"`
	Parameters  []SysBaseMenuParameter `json:"parameters" gorm:"foreignKey:SysBaseMenuID;references:MenuId"`
	Btns        map[string]uint        `json:"btns" gorm:"-"`
}

// Lookups on this join table happen on every authorized request (middleware
// checks "which menus does this role see"). Without indexes GORM's default
// schema is a plain table, so the check scales O(rows). The composite index
// covers the common (authority_id, menu_id) filter; the single-column index
// on menu_id keeps the reverse lookup ("which roles link to this menu") fast
// for menu deletion.
type SysAuthorityMenu struct {
	MenuId      string `json:"menuId" gorm:"comment:menu ID;column:sys_base_menu_id;index;index:idx_authority_menu,priority:2"`
	AuthorityId string `json:"-" gorm:"comment:role ID;column:sys_authority_authority_id;index:idx_authority_menu,priority:1"`
}

func (s SysAuthorityMenu) TableName() string {
	return "sys_authority_menus"
}
