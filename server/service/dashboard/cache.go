package dashboard

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	cacheKey = "dashboard:overview"
	cacheTTL = 30 * time.Second
)

// Cache wraps Redis for dashboard metric caching.
type Cache struct {
	rdb redis.UniversalClient
}

// NewCache creates a dashboard Cache. rdb may be nil — caching is skipped.
func NewCache(rdb redis.UniversalClient) *Cache {
	return &Cache{rdb: rdb}
}

// Get attempts to read cached OverviewMetrics. Returns nil on miss or error.
func (c *Cache) Get(ctx context.Context) *OverviewMetrics {
	if c.rdb == nil {
		return nil
	}
	val, err := c.rdb.Get(ctx, cacheKey).Bytes()
	if err != nil {
		return nil
	}
	var m OverviewMetrics
	if json.Unmarshal(val, &m) != nil {
		return nil
	}
	return &m
}

// Set writes OverviewMetrics to cache with a 30 s TTL.
func (c *Cache) Set(ctx context.Context, m *OverviewMetrics) {
	if c.rdb == nil || m == nil {
		return
	}
	data, err := json.Marshal(m)
	if err != nil {
		return
	}
	c.rdb.Set(ctx, cacheKey, data, cacheTTL)
}
