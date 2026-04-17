package system

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
)

type GDPRService struct{}

// GDPRExport contains all personal data held for a user (GDPR Art. 20 data portability).
type GDPRExport struct {
	ExportedAt    time.Time                    `json:"exported_at"`
	UserID        uint                         `json:"user_id"`
	Profile       map[string]any               `json:"profile"`
	LoginLogs     []system.SysLoginLog         `json:"login_logs"`
	OperationLogs []system.SysOperationRecord  `json:"operation_logs"`
	DataChanges   []system.SysDataChangeLog    `json:"data_changes"`
}

// ExportUserData assembles all personal data for userID (GDPR Art. 20).
func (s *GDPRService) ExportUserData(ctx context.Context, userID uint) (*GDPRExport, error) {
	if global.GVA_DB == nil {
		return nil, fmt.Errorf("database not initialised")
	}

	var user system.SysUser
	if err := global.GVA_DB.WithContext(ctx).
		Select("id, uuid, username, nick_name, header_img, phone, email, enable, created_at, updated_at").
		Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	export := &GDPRExport{
		ExportedAt: time.Now().UTC(),
		UserID:     userID,
		Profile: map[string]any{
			"username":   user.Username,
			"nick_name":  user.NickName,
			"email":      user.Email,
			"phone":      user.Phone,
			"header_img": user.HeaderImg,
			"enable":     user.Enable,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	}

	// Login audit logs
	global.GVA_DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(1000).
		Find(&export.LoginLogs)

	// Operation records (UserID is int in SysOperationRecord)
	global.GVA_DB.WithContext(ctx).
		Where("user_id = ?", int(userID)).
		Order("created_at DESC").
		Limit(1000).
		Find(&export.OperationLogs)

	// Data change logs where this user was the subject
	global.GVA_DB.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", "SysUser", fmt.Sprintf("%d", userID)).
		Order("created_at DESC").
		Limit(1000).
		Find(&export.DataChanges)

	RecordDataChange(ctx, "SysUser", fmtID(userID), "gdpr_export", nil, map[string]any{"exported_at": export.ExportedAt})
	return export, nil
}

// EraseUserData anonymizes all PII fields for userID (GDPR Art. 17 right to erasure).
// The user record is soft-deleted (deleted_at set) and all PII fields are nullified.
// Audit logs referencing the user are NOT deleted (legal obligation to retain audit trails).
func (s *GDPRService) EraseUserData(ctx context.Context, userID uint) error {
	if global.GVA_DB == nil {
		return fmt.Errorf("database not initialised")
	}

	var user system.SysUser
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	before := map[string]any{
		"username":   user.Username,
		"nick_name":  user.NickName,
		"email":      user.Email,
		"phone":      user.Phone,
		"header_img": user.HeaderImg,
	}

	// Use a transaction: anonymize + soft-delete atomically.
	err := global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		anonymizedUsername := fmt.Sprintf("deleted-%s", uuid.New().String()[:8])
		updates := map[string]any{
			"username":   anonymizedUsername,
			"nick_name":  "[deleted]",
			"email":      "",
			"phone":      "",
			"header_img": "",
			"password":   "",
		}
		if err := tx.Model(&system.SysUser{}).
			Where("id = ?", userID).
			Updates(updates).Error; err != nil {
			return err
		}
		// Soft delete
		return tx.Delete(&system.SysUser{}, userID).Error
	})
	if err != nil {
		return err
	}

	RecordDataChange(ctx, "SysUser", fmtID(userID), "gdpr_erase", before, map[string]any{"status": "anonymized"})
	if global.GVA_LOG != nil {
		global.GVA_LOG.Info("GDPR erasure completed", zap.Uint("user_id", userID))
	}
	return nil
}
