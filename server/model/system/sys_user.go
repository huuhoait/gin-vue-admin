package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common"
	"github.com/google/uuid"
)

type Login interface {
	GetUsername() string
	GetNickname() string
	GetUUID() uuid.UUID
	GetUserId() uint
	GetAuthorityId() uint
	GetUserInfo() any
}

var _ Login = new(SysUser)

type SysUser struct {
	global.GVA_MODEL
	UUID          uuid.UUID      `json:"uuid" gorm:"index;comment:user UUID"`                                                                   // user UUID
	Username      string         `json:"userName" gorm:"index;comment:username"`                                                                // username
	Password      string         `json:"-"  gorm:"comment:user login password"`                                                                           // user login password
	NickName      string         `json:"nickName" gorm:"default:System User;comment:user nickname"`                                                          // user nickname
	HeaderImg     string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:user avatar"`               // user avatar
	AuthorityId   uint           `json:"authorityId" gorm:"default:888;comment:user role ID"`                                                      // user role ID
	Authority     SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:UserRole"`                        // UserRole
	Authorities   []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`                                                   // MultipleUserRole
	Phone         string         `json:"phone"  gorm:"comment:user phone number"`                                                                        // user phone number
	Email         string         `json:"email"  gorm:"comment:user email"`                                                                         // user email
	Enable        int            `json:"enable" gorm:"default:1;comment:UserYesNofrozen 1Normal 2Freeze"`                                                    //UserYesNofrozen 1Normal 2Freeze
	OriginSetting common.JSONMap `json:"originSetting" form:"originSetting" gorm:"type:text;default:null;column:origin_setting;comment:configuration;"` //configuration
	OAuthProvider string         `json:"oauth_provider,omitempty" gorm:"size:32;comment:oidc provider name e.g. okta"`
	OAuthSub      string         `json:"oauth_sub,omitempty"      gorm:"size:256;uniqueIndex:idx_oauth;comment:oidc subject identifier"`
}

func (SysUser) TableName() string {
	return "sys_users"
}

func (s *SysUser) GetUsername() string {
	return s.Username
}

func (s *SysUser) GetNickname() string {
	return s.NickName
}

func (s *SysUser) GetUUID() uuid.UUID {
	return s.UUID
}

func (s *SysUser) GetUserId() uint {
	return s.ID
}

func (s *SysUser) GetAuthorityId() uint {
	return s.AuthorityId
}

func (s *SysUser) GetUserInfo() any {
	return *s
}
