package system

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/utils"
)

// ChainVerifyResult is returned by VerifyAuditChain.
type ChainVerifyResult struct {
	OK          bool      `json:"ok"`
	RowsChecked int       `json:"rows_checked"`
	BrokenAt    *uint     `json:"broken_at,omitempty"`   // ID of first bad row
	BrokenHash  string    `json:"broken_hash,omitempty"` // stored hash that failed
	Expected    string    `json:"expected,omitempty"`    // recomputed hash
	CheckedAt   time.Time `json:"checked_at"`
}

// VerifyAuditChain walks the entire sys_policy_change_logs table in insertion
// order and re-derives each SHA-256 hash from the row's content + the
// previous row's hash. It returns on the first broken link so the caller
// knows exactly which row was tampered with.
//
// The full table scan is intentional: integrity requires checking every link.
// For large deployments, run this as an overnight job rather than
// on every request — the route is admin-only and rate-limited.
func VerifyAuditChain(ctx context.Context) (ChainVerifyResult, error) {
	result := ChainVerifyResult{CheckedAt: time.Now()}
	if global.GVA_DB == nil {
		return result, fmt.Errorf("database not initialised")
	}

	const batchSize = 500
	var offset int
	prevHash := ""

	for {
		var batch []system.SysPolicyChangeLog
		err := global.GVA_DB.WithContext(ctx).
			Order("id ASC").
			Limit(batchSize).
			Offset(offset).
			Find(&batch).Error
		if err != nil {
			return result, fmt.Errorf("db scan at offset %d: %w", offset, err)
		}
		if len(batch) == 0 {
			break
		}

		for i := range batch {
			row := &batch[i]
			expected := hashChainEntry(row, prevHash)
			result.RowsChecked++

			if row.Hash != expected {
				id := row.ID
				result.OK = false
				result.BrokenAt = &id
				result.BrokenHash = row.Hash
				result.Expected = expected
				global.GVA_LOG.Error("audit chain broken",
					zap.Uint("id", id),
					zap.String("stored", row.Hash),
					zap.String("expected", expected),
				)
				return result, nil
			}
			prevHash = row.Hash
		}
		offset += batchSize
	}

	result.OK = true
	return result, nil
}

// ScheduledAuditChainVerify is meant to be registered with the timer/cron
// subsystem. It acquires a distributed lock so only one replica runs the
// check, logs the result, and returns an error if the chain is broken.
func ScheduledAuditChainVerify() error {
	release, acquired, err := utils.TryDistLock(context.Background(), "cron:audit-chain-verify", 30*time.Minute)
	if err != nil {
		return err
	}
	if !acquired {
		return nil
	}
	defer release()

	result, err := VerifyAuditChain(context.Background())
	if err != nil {
		global.GVA_LOG.Error("audit chain verification error", zap.Error(err))
		return err
	}
	if !result.OK {
		global.GVA_LOG.Error("AUDIT CHAIN INTEGRITY FAILURE — possible tampering",
			zap.Uintp("broken_at_id", result.BrokenAt),
			zap.String("stored_hash", result.BrokenHash),
			zap.String("expected_hash", result.Expected),
			zap.Int("rows_checked", result.RowsChecked),
		)
		return fmt.Errorf("audit chain broken at id=%d", *result.BrokenAt)
	}
	global.GVA_LOG.Info("audit chain verified OK", zap.Int("rows_checked", result.RowsChecked))
	return nil
}
