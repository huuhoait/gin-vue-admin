package service

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	sysmodel "github.com/huuhoait/gin-vue-admin/server/model/system"
	sysservice "github.com/huuhoait/gin-vue-admin/server/service/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/onlineusers/model"
	"github.com/huuhoait/gin-vue-admin/server/plugin/onlineusers/model/request"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Redis layout:
//   HASH gva:online:sessions   field=uuid → JSON(SessionInfo)
//   ZSET gva:online:lru        member=uuid, score=lastSeenAt unix nanos
// Two-key design: hash for O(1) lookup + listing, zset for cheap stale-prune
// via ZRANGEBYSCORE.
const (
	redisHashSessions = "gva:online:sessions"
	redisZsetLRU      = "gva:online:lru"
	defaultPruneTTL   = 30 * time.Minute
)

// ErrRedisDisabled is returned when Redis is not configured. Touch swallows
// it (online tracking degrades gracefully); admin endpoints surface it.
var ErrRedisDisabled = errors.New("redis not configured; online-user tracking disabled")

type session struct{}

// Touch upserts the session record for the authenticated user and refreshes
// its LRU position. Preserves LoginAt across calls so the column reflects
// session start, not the last request.
func (s *session) Touch(ctx context.Context, uuid string, userID uint, username, nickName string, authorityID uint, token, ip, userAgent string) error {
	if global.GVA_REDIS == nil {
		return nil
	}
	now := time.Now().UTC()
	loginAt := now
	if existing, err := s.findOne(ctx, uuid); err == nil && existing != nil {
		loginAt = existing.LoginAt
	}
	info := model.SessionInfo{
		UUID:        uuid,
		UserID:      userID,
		Username:    username,
		NickName:    nickName,
		AuthorityId: authorityID,
		IP:          ip,
		UserAgent:   userAgent,
		LoginAt:     loginAt,
		LastSeenAt:  now,
		Token:       token,
	}
	raw, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = global.GVA_REDIS.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.HSet(ctx, redisHashSessions, uuid, raw)
		p.ZAdd(ctx, redisZsetLRU, redis.Z{Score: float64(now.UnixNano()), Member: uuid})
		return nil
	})
	return err
}

func (s *session) findOne(ctx context.Context, uuid string) (*model.SessionInfo, error) {
	raw, err := global.GVA_REDIS.HGet(ctx, redisHashSessions, uuid).Result()
	if err != nil {
		return nil, err
	}
	var info model.SessionInfo
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// List returns paginated, username-filtered sessions sorted by LastSeenAt
// desc. The Token field is wiped before returning.
func (s *session) List(ctx context.Context, req request.ListSessionsReq) (list []model.SessionInfo, total int64, err error) {
	if global.GVA_REDIS == nil {
		return nil, 0, ErrRedisDisabled
	}
	all, err := global.GVA_REDIS.HGetAll(ctx, redisHashSessions).Result()
	if err != nil {
		return nil, 0, err
	}
	sessions := make([]model.SessionInfo, 0, len(all))
	uname := strings.ToLower(strings.TrimSpace(req.Username))
	for _, raw := range all {
		var info model.SessionInfo
		if err := json.Unmarshal([]byte(raw), &info); err != nil {
			continue
		}
		if uname != "" && !strings.Contains(strings.ToLower(info.Username), uname) {
			continue
		}
		info.Token = ""
		sessions = append(sessions, info)
	}
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].LastSeenAt.After(sessions[j].LastSeenAt) })
	total = int64(len(sessions))
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	start := (req.Page - 1) * req.PageSize
	if start >= len(sessions) {
		return []model.SessionInfo{}, total, nil
	}
	end := start + req.PageSize
	if end > len(sessions) {
		end = len(sessions)
	}
	return sessions[start:end], total, nil
}

// Kick blacklists the user's JWT and removes their session record. Idempotent:
// missing sessions return nil (the user was already gone).
func (s *session) Kick(ctx context.Context, uuid string) error {
	if global.GVA_REDIS == nil {
		return ErrRedisDisabled
	}
	info, err := s.findOne(ctx, uuid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	if info == nil {
		return nil
	}
	if info.Token != "" {
		if blErr := sysservice.JwtServiceApp.JsonInBlacklist(sysmodel.JwtBlacklist{Jwt: info.Token}); blErr != nil {
			global.GVA_LOG.Warn("kick: blacklist token failed", zap.Error(blErr))
		}
	}
	_, _ = global.GVA_REDIS.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.HDel(ctx, redisHashSessions, uuid)
		p.ZRem(ctx, redisZsetLRU, uuid)
		return nil
	})
	return nil
}

// Prune drops sessions whose LastSeenAt is older than ttl. Returns the count
// pruned. Called from the plugin's cron task.
func (s *session) Prune(ctx context.Context, ttl time.Duration) (int, error) {
	if global.GVA_REDIS == nil {
		return 0, nil
	}
	if ttl <= 0 {
		ttl = defaultPruneTTL
	}
	threshold := time.Now().Add(-ttl).UnixNano()
	members, err := global.GVA_REDIS.ZRangeByScore(ctx, redisZsetLRU, &redis.ZRangeBy{
		Min: "0",
		Max: strconv.FormatInt(threshold, 10),
	}).Result()
	if err != nil {
		return 0, err
	}
	if len(members) == 0 {
		return 0, nil
	}
	zremArgs := make([]any, len(members))
	for i, m := range members {
		zremArgs[i] = m
	}
	_, err = global.GVA_REDIS.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.HDel(ctx, redisHashSessions, members...)
		p.ZRem(ctx, redisZsetLRU, zremArgs...)
		return nil
	})
	return len(members), err
}
