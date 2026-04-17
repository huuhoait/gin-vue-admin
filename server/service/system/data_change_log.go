package system

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
)

// sensitiveDataField matches JSON string values whose key name suggests
// sensitive content. The value is replaced with [REDACTED] before storage.
var sensitiveDataField = regexp.MustCompile(
	`("(password|passwd|pwd|token|secret|authorization|x[-_]?token|signingKey|refreshToken|accessToken|apiKey)"\s*:\s*")[^"]*(")`+
		`|("(password|passwd|pwd|token|secret|authorization|x[-_]?token|signingKey|refreshToken|accessToken|apiKey)"\s*:\s*)(\d+|true|false|null)`,
)

func scrubDataJSON(s string) string {
	return sensitiveDataField.ReplaceAllStringFunc(s, func(m string) string {
		// preserve the key prefix and closing quote; replace only the value
		return sensitiveDataField.ReplaceAllString(m, `${1}[REDACTED]${3}`)
	})
}

// DataChangeCtx carries operator metadata extracted from the gin context.
type DataChangeCtx struct {
	OperatorID   uint
	OperatorName string
	IP           string
	RequestID    string
}

func extractDataChangeCtx(ctx context.Context) DataChangeCtx {
	gc := requestContext(ctx) // reuse policyCtxKey helper from policy_audit.go
	if gc == nil {
		return DataChangeCtx{}
	}
	var d DataChangeCtx
	d.IP = gc.ClientIP()
	d.RequestID = gc.GetHeader("X-Request-ID")
	if raw, ok := gc.Get("claims"); ok {
		type claimsIface interface {
			GetUserID() uint
			GetUsername() string
		}
		if c, ok := raw.(claimsIface); ok {
			d.OperatorID = c.GetUserID()
			d.OperatorName = c.GetUsername()
		}
	}
	return d
}

// hashDataChangeEntry returns the SHA-256 hex digest covering all mutable fields
// of the entry plus the previous row's hash. This makes the chain tamper-evident:
// altering any past row invalidates every subsequent hash.
func hashDataChangeEntry(e *system.SysDataChangeLog, prevHash string) string {
	h := sha256.New()
	fmt.Fprintf(h, "%d|%s|%s|%s|%s|%s|%s|%s|%s|%s",
		e.OperatorID, e.OperatorName,
		e.TargetType, e.TargetID, e.Action,
		e.Before, e.After,
		e.IP, e.RequestID,
		prevHash,
	)
	return hex.EncodeToString(h.Sum(nil))
}

// lastDataChangeHashTx reads the most recent hash from the table within an
// existing transaction, returning "" if the table is empty.
func lastDataChangeHashTx(ctx context.Context, tx *gorm.DB) string {
	var last system.SysDataChangeLog
	err := tx.WithContext(ctx).
		Order("id DESC").
		Select("hash").
		First(&last).Error
	if err != nil {
		return ""
	}
	return last.Hash
}

// RecordDataChange persists a before/after audit row for any admin mutation.
//
//   - targetType: entity name, e.g. "SysUser", "SystemConfig"
//   - targetID:   primary key or identifier of the entity as string
//   - action:     verb, e.g. "update", "delete", "reset_password", "set_authority"
//   - before/after: old and new state; pass nil for creates/deletes where one side is absent
//
// Failures are logged but never returned — the caller's mutation already
// succeeded and blocking on the audit sink would degrade UX.
func RecordDataChange(ctx context.Context, targetType, targetID, action string, before, after any) {
	if global.GVA_DB == nil {
		return
	}
	op := extractDataChangeCtx(ctx)
	entry := system.SysDataChangeLog{
		OperatorID:   op.OperatorID,
		OperatorName: op.OperatorName,
		TargetType:   targetType,
		TargetID:     targetID,
		Action:       action,
		IP:           op.IP,
		RequestID:    op.RequestID,
	}
	if before != nil {
		if b, err := json.Marshal(before); err == nil {
			entry.Before = scrubDataJSON(string(b))
		}
	}
	if after != nil {
		if a, err := json.Marshal(after); err == nil {
			entry.After = scrubDataJSON(string(a))
		}
	}

	// Use a serializable TX so that (SELECT last hash → INSERT) is atomic and
	// the chain remains consistent under concurrent writers.
	err := global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		prev := lastDataChangeHashTx(ctx, tx)
		entry.PrevHash = prev
		entry.Hash = hashDataChangeEntry(&entry, prev)
		return tx.Create(&entry).Error
	})
	if err != nil {
		if global.GVA_LOG != nil {
			global.GVA_LOG.Error("data change log write failed",
				zap.Error(err),
				zap.String("target_type", targetType),
				zap.String("target_id", targetID),
				zap.String("action", action),
			)
		}
	}
}

// fmtID converts any numeric or string id to a stable string key.
func fmtID(id any) string { return fmt.Sprintf("%v", id) }
