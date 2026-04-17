package utils

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	enforcerMu           sync.Mutex
)

// GetCasbin returns the Casbin enforcer. If a previous init attempt failed
// (e.g. DB was not ready), subsequent calls retry under a mutex rather than
// permanently returning nil — callers still must nil-check.
func GetCasbin() *casbin.SyncedCachedEnforcer {
	enforcerMu.Lock()
	defer enforcerMu.Unlock()
	if syncedCachedEnforcer != nil {
		return syncedCachedEnforcer
	}
	if global.GVA_DB == nil {
		zap.L().Error("casbin init skipped: GVA_DB is nil")
		return nil
	}
	a, err := gormadapter.NewAdapterByDB(global.GVA_DB)
	if err != nil {
		zap.L().Error("AdaptDatabasefailedPleaseCheckcasbinTableYesNoForInnoDBReferenceengine!", zap.Error(err))
		return nil
	}
	text := `
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[role_definition]
		g = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		zap.L().Error("StringLoadModelfailed!", zap.Error(err))
		return nil
	}
	e, err := casbin.NewSyncedCachedEnforcer(m, a)
	if err != nil {
		zap.L().Error("casbin NewSyncedCachedEnforcer failed", zap.Error(err))
		return nil
	}
	e.SetExpireTime(60 * 60)
	if err := e.LoadPolicy(); err != nil {
		zap.L().Error("casbin LoadPolicy failed", zap.Error(err))
		return nil
	}
	syncedCachedEnforcer = e
	return syncedCachedEnforcer
}
