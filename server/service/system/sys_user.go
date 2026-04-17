package system

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/model/common"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct{}

var UserServiceApp = new(UserService)

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) {
		return userInter, errors.New("username already registered")
	}
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.New()
	err = global.GVA_DB.Create(&u).Error
	if err == nil {
		after := u
		after.Password = ""
		RecordDataChange(context.Background(), "SysUser", fmtID(u.ID), "create", nil, after)
	}
	return u, err
}

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.GVA_DB {
		return nil, fmt.Errorf("db not init")
	}
	var user system.SysUser
	err = global.GVA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("incorrect password")
		}
		MenuServiceApp.UserAuthorityDefaultRouter(&user)
	}
	return &user, err
}

func (userService *UserService) ChangePassword(ctx context.Context, u *system.SysUser, newPassword string) (err error) {
	var user system.SysUser
	err = global.GVA_DB.Select("id, password").Where("id = ?", u.ID).First(&user).Error
	if err != nil {
		return err
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return errors.New("original password is incorrect")
	}
	pwd := utils.BcryptHash(newPassword)
	err = global.GVA_DB.Model(&user).Update("password", pwd).Error
	if err == nil {
		RecordDataChange(ctx, "SysUser", fmtID(u.ID), "change_password", nil, nil)
	}
	return err
}

func (userService *UserService) GetUserInfoList(info systemReq.GetUserList) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysUser{})
	var userList []system.SysUser

	if info.NickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+info.NickName+"%")
	}
	if info.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+info.Phone+"%")
	}
	if info.Username != "" {
		db = db.Where("username LIKE ?", "%"+info.Username+"%")
	}
	if info.Email != "" {
		db = db.Where("email LIKE ?", "%"+info.Email+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	orderStr := "id desc"
	if info.OrderKey != "" {
		allowedOrders := map[string]bool{
			"id":        true,
			"username":  true,
			"nick_name": true,
			"phone":     true,
			"email":     true,
		}
		if allowedOrders[info.OrderKey] {
			orderStr = info.OrderKey
			if info.Desc {
				orderStr = info.OrderKey + " desc"
			}
		}
	}

	err = db.Limit(limit).Offset(offset).Order(orderStr).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

func (userService *UserService) SetUserAuthority(ctx context.Context, id uint, authorityId uint) (err error) {
	assignErr := global.GVA_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("user does not have this role")
	}

	var authority system.SysAuthority
	err = global.GVA_DB.Where("authority_id = ?", authorityId).First(&authority).Error
	if err != nil {
		return err
	}
	var authorityMenu []system.SysAuthorityMenu
	var authorityMenuIDs []string
	err = global.GVA_DB.Where("sys_authority_authority_id = ?", authorityId).Find(&authorityMenu).Error
	if err != nil {
		return err
	}

	for i := range authorityMenu {
		authorityMenuIDs = append(authorityMenuIDs, authorityMenu[i].MenuId)
	}

	var authorityMenus []system.SysBaseMenu
	err = global.GVA_DB.Preload("Parameters").Where("id in (?)", authorityMenuIDs).Find(&authorityMenus).Error
	if err != nil {
		return err
	}
	hasMenu := false
	for i := range authorityMenus {
		if authorityMenus[i].Name == authority.DefaultRouter {
			hasMenu = true
			break
		}
	}
	if !hasMenu {
		return errors.New("default route not found, cannot switch to this role")
	}

	// Capture old authority before update
	var oldUser system.SysUser
	_ = global.GVA_DB.Select("id, authority_id").Where("id = ?", id).First(&oldUser).Error

	err = global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", id).Update("authority_id", authorityId).Error
	if err == nil {
		RecordDataChange(ctx, "SysUser", fmtID(id), "set_authority",
			map[string]any{"authority_id": oldUser.AuthorityId},
			map[string]any{"authority_id": authorityId},
		)
	}
	return err
}

func (userService *UserService) SetUserAuthorities(ctx context.Context, adminAuthorityID, id uint, authorityIds []uint) (err error) {
	// Capture old authorities for the audit log
	var oldAuthorities []system.SysUserAuthority
	_ = global.GVA_DB.Where("sys_user_id = ?", id).Find(&oldAuthorities).Error

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var user system.SysUser
		TxErr := tx.Where("id = ?", id).First(&user).Error
		if TxErr != nil {
			global.GVA_LOG.Debug(TxErr.Error())
			return errors.New("failed to query user data")
		}
		TxErr = tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []system.SysUserAuthority
		for _, v := range authorityIds {
			e := AuthorityServiceApp.CheckAuthorityIDAuth(adminAuthorityID, v)
			if e != nil {
				return e
			}
			useAuthority = append(useAuthority, system.SysUserAuthority{
				SysUserId: id, SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = tx.Model(&user).Update("authority_id", authorityIds[0]).Error
		if TxErr != nil {
			return TxErr
		}
		return nil
	})
	if err == nil {
		var oldIDs []uint
		for _, a := range oldAuthorities {
			oldIDs = append(oldIDs, a.SysAuthorityAuthorityId)
		}
		RecordDataChange(ctx, "SysUser", fmtID(id), "set_authorities",
			map[string]any{"authority_ids": oldIDs},
			map[string]any{"authority_ids": authorityIds},
		)
	}
	return err
}

func (userService *UserService) DeleteUser(ctx context.Context, id int) (err error) {
	// Snapshot the user before deletion
	var before system.SysUser
	_ = global.GVA_DB.Select("id, username, nick_name, email, phone, enable").Where("id = ?", id).First(&before).Error

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&system.SysUser{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		RecordDataChange(ctx, "SysUser", fmtID(id), "delete", before, nil)
	}
	return err
}

func (userService *UserService) SetUserInfo(ctx context.Context, req system.SysUser) error {
	// Snapshot before
	var before system.SysUser
	_ = global.GVA_DB.Select("id, nick_name, header_img, phone, email, enable").Where("id = ?", req.ID).First(&before).Error

	err := global.GVA_DB.Model(&system.SysUser{}).
		Select("updated_at", "nick_name", "header_img", "phone", "email", "enable").
		Where("id=?", req.ID).
		Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"nick_name":  req.NickName,
			"header_img": req.HeaderImg,
			"phone":      req.Phone,
			"email":      req.Email,
			"enable":     req.Enable,
		}).Error
	if err == nil {
		after := map[string]any{
			"nick_name":  req.NickName,
			"header_img": req.HeaderImg,
			"phone":      req.Phone,
			"email":      req.Email,
			"enable":     req.Enable,
		}
		RecordDataChange(ctx, "SysUser", fmtID(req.ID), "update", before, after)
	}
	return err
}

func (userService *UserService) SetSelfInfo(ctx context.Context, req system.SysUser) error {
	var before system.SysUser
	_ = global.GVA_DB.Select("id, nick_name, header_img, phone, email").Where("id = ?", req.ID).First(&before).Error

	err := global.GVA_DB.Model(&system.SysUser{}).Where("id=?", req.ID).Updates(req).Error
	if err == nil {
		after := req
		after.Password = "" // never log password field
		RecordDataChange(ctx, "SysUser", fmtID(req.ID), "self_update", before, after)
	}
	return err
}

func (userService *UserService) SetSelfSetting(ctx context.Context, req common.JSONMap, uid uint) error {
	err := global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", uid).Update("origin_setting", req).Error
	if err == nil {
		RecordDataChange(ctx, "SysUser", fmtID(uid), "update_setting", nil, req)
	}
	return err
}

func (userService *UserService) GetUserInfo(uuid uuid.UUID) (user system.SysUser, err error) {
	var reqUser system.SysUser
	err = global.GVA_DB.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	MenuServiceApp.UserAuthorityDefaultRouter(&reqUser)
	return reqUser, err
}

func (userService *UserService) FindUserById(id int) (user *system.SysUser, err error) {
	var u system.SysUser
	err = global.GVA_DB.Where("id = ?", id).First(&u).Error
	return &u, err
}

func (userService *UserService) FindUserByUuid(uuid string) (user *system.SysUser, err error) {
	var u system.SysUser
	if err = global.GVA_DB.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("user does not exist")
	}
	return &u, nil
}

func (userService *UserService) ResetPassword(ctx context.Context, ID uint, password string) (err error) {
	err = global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash(password)).Error
	if err == nil {
		RecordDataChange(ctx, "SysUser", fmtID(ID), "reset_password", nil, nil)
	}
	return err
}
