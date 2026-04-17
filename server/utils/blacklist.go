package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
)

// Storing the full JWT in Redis bloats memory and logs; a SHA-256 of the token
// is unique enough for a blacklist check and keeps keys fixed-width.
func blacklistKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return "gva:jwt:bl:" + hex.EncodeToString(sum[:])
}

// BlacklistTTL returns how long a blacklist entry should live. It mirrors the
// token's remaining lifetime so expired tokens roll off automatically.
func BlacklistTTL(token string) time.Duration {
	claims := &request.CustomClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(token, claims)
	if err != nil || claims.ExpiresAt == nil {
		dr, perr := ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
		if perr != nil {
			return 24 * time.Hour
		}
		return dr
	}
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		return time.Minute
	}
	return remaining
}

// BlacklistAdd inserts a token into the blacklist. Redis is preferred so that
// the blacklist is shared across instances; when Redis is unavailable we fall
// back to the local cache so single-node deployments keep working.
func BlacklistAdd(token string) error {
	ttl := BlacklistTTL(token)
	if global.GVA_REDIS != nil {
		if err := global.GVA_REDIS.Set(context.Background(), blacklistKey(token), 1, ttl).Err(); err == nil {
			return nil
		}
	}
	global.BlackCache.Set(token, struct{}{}, ttl)
	return nil
}

// BlacklistContains checks Redis first, then falls back to the local cache.
func BlacklistContains(token string) bool {
	if global.GVA_REDIS != nil {
		n, err := global.GVA_REDIS.Exists(context.Background(), blacklistKey(token)).Result()
		if err == nil && n > 0 {
			return true
		}
	}
	if _, ok := global.BlackCache.Get(token); ok {
		return true
	}
	return false
}
