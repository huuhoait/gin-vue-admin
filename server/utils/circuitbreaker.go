package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

var (
	// DBBreaker protects the primary database connection.
	DBBreaker = newBreaker("gva-db")
	// RedisBreaker protects the Redis connection.
	RedisBreaker = newBreaker("gva-redis")
)

func newBreaker(name string) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 5,            // half-open: allow 5 probes before closing
		Interval:    60 * time.Second, // rolling window for error counting
		Timeout:     30 * time.Second, // open→half-open wait
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip when ≥5 consecutive failures or failure rate > 60% with ≥10 requests
			if counts.ConsecutiveFailures >= 5 {
				return true
			}
			total := counts.Requests
			if total >= 10 && float64(counts.TotalFailures)/float64(total) >= 0.6 {
				return true
			}
			return false
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			if global.GVA_LOG != nil {
				global.GVA_LOG.Warn("circuit breaker state change",
					zap.String("breaker", name),
					zap.String("from", from.String()),
					zap.String("to", to.String()),
				)
			}
		},
	})
}

// DBExec executes fn inside the DB circuit breaker.
// Returns gobreaker.ErrOpenState if the circuit is open.
func DBExec(fn func() error) error {
	_, err := DBBreaker.Execute(func() (any, error) {
		return nil, fn()
	})
	return err
}

// RedisExec executes fn inside the Redis circuit breaker.
func RedisExec(fn func() error) error {
	_, err := RedisBreaker.Execute(func() (any, error) {
		return nil, fn()
	})
	return err
}

// WithDBTimeout returns a context with a 5-second deadline for DB queries.
// This is the standard per-query timeout; adjust per call site if needed.
func WithDBTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, 5*time.Second)
}

// WithDBTimeoutDuration returns a context with a custom deadline for DB queries.
func WithDBTimeoutDuration(parent context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, d)
}

// RetryDB retries fn up to maxAttempts times with exponential backoff.
// Only retries on non-circuit-open errors (circuit open = stop immediately).
func RetryDB(ctx context.Context, maxAttempts int, fn func(ctx context.Context) error) error {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		if i > 0 {
			wait := time.Duration(1<<uint(i-1)) * 100 * time.Millisecond // 100ms, 200ms, 400ms
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(wait):
			}
		}
		lastErr = fn(ctx)
		if lastErr == nil {
			return nil
		}
		// Don't retry if circuit is open
		if lastErr == gobreaker.ErrOpenState || lastErr == gobreaker.ErrTooManyRequests {
			return lastErr
		}
		// Don't retry context errors
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return fmt.Errorf("after %d attempts: %w", maxAttempts, lastErr)
}
