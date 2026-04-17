package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
)

type LoginLogService struct{}

var LoginLogServiceApp = new(LoginLogService)

func (loginLogService *LoginLogService) CreateLoginLog(loginLog system.SysLoginLog) (err error) {
	err = global.GVA_DB.Create(&loginLog).Error
	return err
}

func (loginLogService *LoginLogService) DeleteLoginLogByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysLoginLog{}, "id in (?)", ids.Ids).Error
	return err
}

func (loginLogService *LoginLogService) DeleteLoginLog(loginLog system.SysLoginLog) (err error) {
	err = global.GVA_DB.Delete(&loginLog).Error
	return err
}

func (loginLogService *LoginLogService) GetLoginLog(id uint) (loginLog system.SysLoginLog, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&loginLog).Error
	return
}

func (loginLogService *LoginLogService) GetLoginLogInfoList(info systemReq.SysLoginLogSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// create db
	db := global.GVA_DB.Model(&system.SysLoginLog{})
	var loginLogs []system.SysLoginLog
	// if search conditions exist, search queries will be automatically built below
	if info.Username != "" {
		db = db.Where("username LIKE ?", "%"+info.Username+"%")
	}
	if info.Status != false {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Preload("User").Find(&loginLogs).Error
	return loginLogs, total, err
}
