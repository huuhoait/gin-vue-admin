// Package audit provides automatic audit-trail population for models that
// embed global.GVA_MODEL. Three columns are stamped from the request
// context's user id by a GORM callback:
//
//   - CreatedBy on INSERT
//   - UpdatedBy on INSERT and UPDATE
//   - DeletedBy on soft DELETE (UPDATE deleted_at)
//
// Adoption is opt-in by context: handlers that pass `c.Request.Context()`
// to GORM (e.g. `db.WithContext(c.Request.Context()).Find(...)`) get the
// stamping for free, provided the JWT middleware injected the user id.
// Calls without a user-scoped context use 0, the conventional "system
// origin" marker (seeds, background tasks, migrations).
package audit

import "context"

type userIDCtxKey struct{}

// WithUserID returns a child context carrying the authenticated user id.
// Called from the JWT middleware after claims have been parsed.
func WithUserID(ctx context.Context, id uint) context.Context {
	if id == 0 {
		return ctx
	}
	return context.WithValue(ctx, userIDCtxKey{}, id)
}

// UserIDFromContext extracts the authenticated user id; ok=false means the
// request was unauthenticated or the JWT middleware did not run (system /
// public endpoint, captcha, login itself).
func UserIDFromContext(ctx context.Context) (uint, bool) {
	if ctx == nil {
		return 0, false
	}
	v := ctx.Value(userIDCtxKey{})
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint)
	if !ok || id == 0 {
		return 0, false
	}
	return id, true
}
