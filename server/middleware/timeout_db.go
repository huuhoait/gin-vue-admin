package middleware

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// DBQueryTimeout is a GORM plugin that wraps every statement in a context
// with a 5-second deadline. This prevents runaway queries from holding
// connections indefinitely.
type DBQueryTimeout struct {
	Timeout time.Duration
}

func (t *DBQueryTimeout) Name() string { return "gva:db_query_timeout" }

func (t *DBQueryTimeout) Initialize(db *gorm.DB) error {
	timeout := t.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	cb := db.Callback()
	addTimeout := func(db *gorm.DB) {
		if db.Statement == nil || db.Statement.Context == nil {
			return
		}
		// Only add timeout if the context doesn't already have a deadline
		if _, ok := db.Statement.Context.Deadline(); ok {
			return
		}
		ctx, cancel := context.WithTimeout(db.Statement.Context, timeout)
		// Store cancel in the statement so it's called when the statement finishes.
		db.Statement.Context = ctx
		db.InstanceSet("gva:cancel", cancel)
	}
	cancelAfter := func(db *gorm.DB) {
		if v, ok := db.InstanceGet("gva:cancel"); ok {
			if cancel, ok := v.(context.CancelFunc); ok {
				cancel()
			}
		}
	}

	cb.Create().Before("gorm:create").Register("gva:timeout_before_create", addTimeout)
	cb.Create().After("gorm:create").Register("gva:timeout_after_create", cancelAfter)
	cb.Query().Before("gorm:query").Register("gva:timeout_before_query", addTimeout)
	cb.Query().After("gorm:query").Register("gva:timeout_after_query", cancelAfter)
	cb.Update().Before("gorm:update").Register("gva:timeout_before_update", addTimeout)
	cb.Update().After("gorm:update").Register("gva:timeout_after_update", cancelAfter)
	cb.Delete().Before("gorm:delete").Register("gva:timeout_before_delete", addTimeout)
	cb.Delete().After("gorm:delete").Register("gva:timeout_after_delete", cancelAfter)
	// Row/Raw return a lazy *sql.Row or *sql.Rows; the caller's .Scan()/.Next()
	// runs after the After-callbacks fire. Canceling here would invalidate
	// that still-pending work (this is what breaks AutoMigrate.HasTable). The
	// 5s deadline on the context still fires on its own if the query hangs.
	cb.Row().Before("gorm:row").Register("gva:timeout_before_row", addTimeout)
	cb.Raw().Before("gorm:raw").Register("gva:timeout_before_raw", addTimeout)
	return nil
}
