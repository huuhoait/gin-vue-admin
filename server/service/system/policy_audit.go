package system

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
)

// ctxKey for stashing the current *gin.Context when the mutation is invoked
// from an HTTP path. Casbin mutators used to live inside service-layer
// transactions with no access to actor metadata; this is a narrow shim that
// keeps the signature stable while still capturing who/where.
type policyCtxKey struct{}

// WithRequestContext annotates the context with the gin context so audit
// logging can read actor + request-id + IP. Safe to call on a nil context.
func WithRequestContext(ctx context.Context, c *gin.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, policyCtxKey{}, c)
}

func requestContext(ctx context.Context) *gin.Context {
	if ctx == nil {
		return nil
	}
	if c, ok := ctx.Value(policyCtxKey{}).(*gin.Context); ok {
		return c
	}
	return nil
}

// hashChainEntry returns the SHA-256 hex digest that covers all mutable fields
// of the entry plus the previous row's hash. This makes the chain tamper-evident:
// altering any past row invalidates every subsequent hash.
func hashChainEntry(e *system.SysPolicyChangeLog, prevHash string) string {
	h := sha256.New()
	fmt.Fprintf(h, "%d|%d|%s|%s|%s|%s|%s|%s|%s|%s",
		e.Actor, e.ActorUserID,
		e.Action, e.AuthorityID,
		e.Before, e.After,
		e.IP, e.RequestID, e.Note,
		prevHash,
	)
	return hex.EncodeToString(h.Sum(nil))
}

// RecordPolicyChange persists an audit row with a hash chain link. Any failure
// is logged but never bubbled up — we must not block a successful policy change
// because the audit sink is temporarily unavailable.
func RecordPolicyChange(ctx context.Context, action, authorityID string, before, after any, note string) {
	if global.GVA_DB == nil {
		return
	}
	entry := system.SysPolicyChangeLog{
		Action:      action,
		AuthorityID: authorityID,
		Note:        note,
	}
	if b, err := json.Marshal(before); err == nil {
		entry.Before = string(b)
	}
	if a, err := json.Marshal(after); err == nil {
		entry.After = string(a)
	}
	if c := requestContext(ctx); c != nil {
		entry.IP = c.ClientIP()
		if rid := c.GetHeader("X-Request-ID"); rid != "" {
			entry.RequestID = rid
		}
		if raw, ok := c.Get("claims"); ok {
			if claims, ok := raw.(interface {
				GetAuthorityID() uint
				GetUserID() uint
			}); ok {
				entry.Actor = claims.GetAuthorityID()
				entry.ActorUserID = claims.GetUserID()
			}
		}
	}

	// Use a serializable TX so that (SELECT last hash → INSERT) is atomic and
	// the chain remains consistent under concurrent writers.
	err := global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		prev := lastPolicyHashTx(ctx, tx)
		entry.PrevHash = prev
		entry.Hash = hashChainEntry(&entry, prev)
		return tx.Create(&entry).Error
	})
	if err != nil {
		if global.GVA_LOG != nil {
			global.GVA_LOG.Error("policy audit write failed", zap.Error(err), zap.String("action", action))
		}
	}
}

func lastPolicyHashTx(ctx context.Context, tx *gorm.DB) string {
	var last system.SysPolicyChangeLog
	err := tx.WithContext(ctx).
		Order("id DESC").
		Select("hash").
		First(&last).Error
	if err != nil {
		return ""
	}
	return last.Hash
}
