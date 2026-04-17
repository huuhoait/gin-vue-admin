package system

import (
	"context"
	"fmt"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
)

type SysErrorService struct{}

// CreateSysError creates an error log record
// Author [yourname](https://github.com/yourname)
func (sysErrorService *SysErrorService) CreateSysError(ctx context.Context, sysError *system.SysError) (err error) {
	if global.GVA_DB == nil {
		return nil
	}
	err = global.GVA_DB.Create(sysError).Error
	return err
}

// DeleteSysError deletes an error log record
// Author [yourname](https://github.com/yourname)
func (sysErrorService *SysErrorService) DeleteSysError(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&system.SysError{}, "id = ?", ID).Error
	return err
}

// DeleteSysErrorByIds batch deletes error log records
// Author [yourname](https://github.com/yourname)
func (sysErrorService *SysErrorService) DeleteSysErrorByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysError{}, "id in ?", IDs).Error
	return err
}

// UpdateSysError updates an error log record
// Author [yourname](https://github.com/yourname)
func (sysErrorService *SysErrorService) UpdateSysError(ctx context.Context, sysError system.SysError) (err error) {
	err = global.GVA_DB.Model(&system.SysError{}).Where("id = ?", sysError.ID).Updates(&sysError).Error
	return err
}

// GetSysError gets an error log record by ID
// Author [yourname](https://github.com/yourname)
func (sysErrorService *SysErrorService) GetSysError(ctx context.Context, ID string) (sysError system.SysError, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sysError).Error
	return
}

// GetSysErrorInfoList gets error log records with pagination
// Author [yourname](https://github.com/yourname)
func (sysErrorService *SysErrorService) GetSysErrorInfoList(ctx context.Context, info systemReq.SysErrorSearch) (list []system.SysError, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// create db
	db := global.GVA_DB.Model(&system.SysError{}).Order("created_at desc")
	var sysErrors []system.SysError
	// if there are search conditions, the search statement will be automatically created below
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.Form != nil && *info.Form != "" {
		db = db.Where("form = ?", *info.Form)
	}
	if info.Info != nil && *info.Info != "" {
		db = db.Where("info LIKE ?", "%"+*info.Info+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&sysErrors).Error
	return sysErrors, total, err
}

// GetSysErrorSolution asynchronously processes an error
// Author [yourname](https://github.com/yourname)
func (sysErrorService *SysErrorService) GetSysErrorSolution(ctx context.Context, ID string) (err error) {
	// immediately update status to processing
	err = global.GVA_DB.WithContext(ctx).Model(&system.SysError{}).Where("id = ?", ID).Update("status", "processing").Error
	if err != nil {
		return err
	}

	// async goroutine to process and update status
	go func(id string) {
		// query current error info for generating solution
		var se system.SysError
		_ = global.GVA_DB.Model(&system.SysError{}).Where("id = ?", id).First(&se).Error

		// construct LLM request parameters, use butler mode to generate solution based on error info
		var form, info string
		if se.Form != nil {
			form = *se.Form
		}
		if se.Info != nil {
			info = *se.Info
		}

		llmReq := common.JSONMap{
			"mode": "solution",
			"info": info,
			"form": form,
		}

		// call service layer LLMAuto, ignore errors but try to write the solution
		var solution string
		if data, err := (&AutoCodeService{}).LLMAuto(context.Background(), llmReq); err == nil {
			solution = fmt.Sprintf("%v", data.(map[string]interface{})["text"])
			_ = global.GVA_DB.Model(&system.SysError{}).Where("id = ?", id).Updates(map[string]interface{}{"status": "completed", "solution": solution}).Error
		} else {
			// even if generation fails, mark as failed to avoid task getting stuck
			_ = global.GVA_DB.Model(&system.SysError{}).Where("id = ?", id).Update("status", "failed").Error
		}
	}(ID)

	return nil
}
