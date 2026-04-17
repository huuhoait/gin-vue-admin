package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

// TryDistLock acquires a best-effort distributed mutex via Redis SET NX PX.
// Returns a release function that MUST be called (typically via defer) once
// the critical section completes. If Redis is unavailable the call acts as
// a no-op success (single-instance behavior preserved) so local dev does
// not require Redis to schedule cron tasks.
//
// ttl should be comfortably longer than the longest expected runtime of
// the protected section — on crash the key expires automatically and the
// next scheduler run can take over, so a fresh instance never deadlocks.
func TryDistLock(ctx context.Context, name string, ttl time.Duration) (release func(), acquired bool, err error) {
	if global.GVA_REDIS == nil {
		// No shared coordinator available; return a no-op release so the
		// caller can proceed. In multi-instance deployments Redis MUST be
		// configured — secrets.go already warns operators in release mode.
		return func() {}, true, nil
	}
	tokBytes := make([]byte, 16)
	if _, rErr := rand.Read(tokBytes); rErr != nil {
		return nil, false, rErr
	}
	token := hex.EncodeToString(tokBytes)
	key := "gva:lock:" + name

	ok, sErr := global.GVA_REDIS.SetNX(ctx, key, token, ttl).Result()
	if sErr != nil {
		return nil, false, sErr
	}
	if !ok {
		return func() {}, false, nil
	}

	// Release must only delete the key if the stored token still matches,
	// otherwise a slow job whose lock already expired could delete the
	// lock of a subsequent holder. Lua keeps the check+delete atomic.
	const releaseScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
else
  return 0
end`
	release = func() {
		_ = global.GVA_REDIS.Eval(context.Background(), releaseScript, []string{key}, token).Err()
	}
	return release, true, nil
}

// ErrLockNotAcquired is returned by callers that want to treat missed locks
// as errors instead of silent skips.
var ErrLockNotAcquired = errors.New("distributed lock not acquired")
