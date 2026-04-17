package request

import (
	common "github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
)

// Register User register structure
type Register struct {
	Username     string `json:"userName" example:"username"`
	Password     string `json:"passWord" example:"password"`
	NickName     string `json:"nickName" example:"nickname"`
	HeaderImg    string `json:"headerImg" example:"avatar URL"`
	AuthorityId  uint   `json:"authorityId" swaggertype:"string" example:"int Roleid"`
	Enable       int    `json:"enable" swaggertype:"string" example:"int YesNoEnable"`
	AuthorityIds []uint `json:"authorityIds" swaggertype:"string" example:"[]uint Roleid"`
	Phone        string `json:"phone" example:"PhoneNumberCode"`
	Email        string `json:"email" example:"electricitySubEmail"`
}

// Login User login structure
type Login struct {
	Username  string `json:"username"`  // username
	Password  string `json:"password"`  // password
	Captcha   string `json:"captcha"`   // captcha
	CaptchaId string `json:"captchaId"` // captchaID
}

// ChangePasswordReq Modify password structure
type ChangePasswordReq struct {
	ID          uint   `json:"-"`           // From JWT InRaiseFetch user id, Avoidunauthorized access
	Password    string `json:"password"`    // password
	NewPassword string `json:"newPassword"` // Newpassword
}

type ResetPassword struct {
	ID       uint   `json:"ID" form:"ID"`
	Password string `json:"password" form:"password" gorm:"comment:user login password"` // user login password
}

// SetUserAuth Modify user's auth structure
type SetUserAuth struct {
	AuthorityId uint `json:"authorityId"` // role ID
}

// SetUserAuthorities Modify user's auth structure
type SetUserAuthorities struct {
	ID           uint
	AuthorityIds []uint `json:"authorityIds"` // role ID
}

type ChangeUserInfo struct {
	ID           uint                  `gorm:"primarykey"`                                                                           // primary key ID
	NickName     string                `json:"nickName" gorm:"default:System User;comment:user nickname"`                                            // user nickname
	Phone        string                `json:"phone"  gorm:"comment:user phone number"`                                                          // user phone number
	AuthorityIds []uint                `json:"authorityIds" gorm:"-"`                                                                // role ID
	Email        string                `json:"email"  gorm:"comment:user email"`                                                           // user email
	HeaderImg    string                `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:user avatar"` // user avatar
	Enable       int                   `json:"enable" gorm:"comment:FreezeUser"`                                                           //FreezeUser
	Authorities  []system.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
}

type GetUserList struct {
	common.PageInfo
	Username string `json:"username" form:"username"`
	NickName string `json:"nickName" form:"nickName"`
	Phone    string `json:"phone" form:"phone"`
	Email    string `json:"email" form:"email"`
	OrderKey string `json:"orderKey" form:"orderKey"` // sort
	Desc     bool   `json:"desc" form:"desc"`         // sortMethod:Ascendingfalse(default)|Descendingtrue
}

// SetRoleUsers Approvedrole IDfull overwriteassociated userList
type SetRoleUsers struct {
	AuthorityId uint   `json:"authorityId" form:"authorityId"` // role ID
	UserIds     []uint `json:"userIds" form:"userIds"`         // User IDList
}
