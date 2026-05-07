// Package gorm wires automatic tenant_id scoping into GORM via callbacks.
//
// The plugin replaces the older opt-in WithTenantScope/WithStrictTenantScope
// helpers with strict-by-default behaviour: any model whose schema exposes a
// "TenantID" field (the convention populated by tenantmodel.TenantModel)
// receives an automatic "tenant_id = ?" predicate on SELECT/UPDATE/DELETE and
// is auto-stamped on INSERT, sourced from the request context.
//
// Bypass paths (in priority order):
//
//  1. service.IsTenantIgnored(ctx) — explicit opt-out for audit / support /
//     backfill code that needs cross-tenant access.
//  2. service.IsSystemTenant(id) — the reserved tenant_id=0 used by super-admin
//     contexts. The callback treats this as "do not inject" so a super-admin
//     can read across the fleet.
//  3. The caller already added a "tenant_id" predicate. We log at debug level
//     and defer to the caller's value rather than ANDing a second one (which
//     would either be a no-op or accidentally fail-closed when values disagree).
//
// We deliberately avoid raw-SQL string parsing: detection walks GORM's
// structured WHERE expressions only. Anything more obscure (e.g. tenant_id
// inside a hand-written sub-select) gets the auto-predicate added in addition
// to the caller's, which is the safe default.
package gorm

import (
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/huuhoait/gin-vue-admin/server/global"
	tenantservice "github.com/huuhoait/gin-vue-admin/server/plugin/tenant/service"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	tenantFieldName  = "TenantID"
	tenantColumnName = "tenant_id"

	cbBeforeQuery  = "tenant:before_query"
	cbBeforeUpdate = "tenant:before_update"
	cbBeforeDelete = "tenant:before_delete"
	cbBeforeCreate = "tenant:before_create"
)

// registered guards Register against double-installation. We allow Register to
// be called more than once (e.g. plugin reload during tests) and only attach
// callbacks the first time per process.
var registered atomic.Bool

// Register hooks the tenant-scoping callbacks onto db. Subsequent calls are
// no-ops so the plugin's init path can call this defensively.
func Register(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	if !registered.CompareAndSwap(false, true) {
		return nil
	}

	if err := db.Callback().Query().Before("gorm:query").Register(cbBeforeQuery, injectTenantQuery); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").Register(cbBeforeUpdate, injectTenantUpdate); err != nil {
		return err
	}
	if err := db.Callback().Delete().Before("gorm:delete").Register(cbBeforeDelete, injectTenantDelete); err != nil {
		return err
	}
	if err := db.Callback().Create().Before("gorm:before_create").Register(cbBeforeCreate, stampTenantCreate); err != nil {
		return err
	}
	return nil
}

// resetForTests rolls registration state back so the test suite can install
// the callbacks against a fresh in-memory DB. Lower-cased so it is unexported.
func resetForTests() { registered.Store(false) }

// shouldInject returns the active tenant id when the callback should attach a
// predicate. It encapsulates the "skip" decisions so the three read/write
// callbacks stay short.
func shouldInject(db *gorm.DB) (uint, bool) {
	if db.Statement == nil || db.Statement.Schema == nil {
		return 0, false
	}
	// Plain raw queries (db.Raw / db.Exec) have no schema attached; nothing to
	// scope. Same for migrations.
	if db.Statement.Schema.LookUpField(tenantFieldName) == nil {
		return 0, false
	}
	ctx := db.Statement.Context
	if tenantservice.IsTenantIgnored(ctx) {
		return 0, false
	}
	id, ok := tenantservice.FromContext(ctx)
	if !ok {
		return 0, false
	}
	if tenantservice.IsSystemTenant(id) {
		return 0, false
	}
	return id, true
}

// hasExplicitTenantPredicate walks the structured WHERE expressions for any
// reference to tenant_id. We handle the forms GORM produces from typical
// caller idioms:
//
//   - db.Where("tenant_id = ?", 1)        → clause.Expr SQL "tenant_id = ?"
//   - db.Where(map[string]any{"tenant_id":1}) → clause.Eq with column "tenant_id"
//   - db.Where(&Model{TenantID: 1})       → clause.Eq via reflection
//
// Anything outside these forms is treated as "no explicit predicate"; we'll
// add ours, which is the safe default.
func hasExplicitTenantPredicate(stmt *gorm.Statement) bool {
	c, ok := stmt.Clauses["WHERE"]
	if !ok {
		return false
	}
	where, ok := c.Expression.(clause.Where)
	if !ok {
		return false
	}
	for _, expr := range where.Exprs {
		if exprMentionsTenant(expr) {
			return true
		}
	}
	return false
}

func exprMentionsTenant(expr clause.Expression) bool {
	switch v := expr.(type) {
	case clause.Eq:
		return columnMatches(v.Column)
	case clause.Neq:
		return columnMatches(v.Column)
	case clause.IN:
		return columnMatches(v.Column)
	case clause.Expr:
		// Best-effort: if the literal SQL fragment names tenant_id we honour the
		// caller's predicate. Lower-case compare so "Tenant_ID" still matches.
		return strings.Contains(strings.ToLower(v.SQL), tenantColumnName)
	case clause.AndConditions:
		for _, sub := range v.Exprs {
			if exprMentionsTenant(sub) {
				return true
			}
		}
	case clause.OrConditions:
		for _, sub := range v.Exprs {
			if exprMentionsTenant(sub) {
				return true
			}
		}
	}
	return false
}

func columnMatches(col any) bool {
	switch v := col.(type) {
	case string:
		return strings.EqualFold(v, tenantColumnName)
	case clause.Column:
		return strings.EqualFold(v.Name, tenantColumnName)
	}
	return false
}

func injectTenantQuery(db *gorm.DB) { injectTenantWhere(db) }
func injectTenantUpdate(db *gorm.DB) { injectTenantWhere(db) }
func injectTenantDelete(db *gorm.DB) { injectTenantWhere(db) }

// injectTenantWhere is the shared body for SELECT/UPDATE/DELETE. We attach a
// table-qualified predicate so JOINs stay unambiguous.
func injectTenantWhere(db *gorm.DB) {
	id, ok := shouldInject(db)
	if !ok {
		return
	}
	if hasExplicitTenantPredicate(db.Statement) {
		// Caller knows what they're doing (e.g. cross-tenant support tooling
		// that opted in by passing an explicit predicate). Honour it rather
		// than silently AND-ing a second one — that would either be a no-op
		// or surprise the caller by zeroing out their result set.
		if logger := global.GVA_LOG; logger != nil {
			logger.Debug("tenant: skipping auto-predicate, caller supplied explicit tenant_id",
				zap.String("table", db.Statement.Table))
		}
		return
	}
	db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{
		clause.Eq{
			Column: clause.Column{Table: db.Statement.Table, Name: tenantColumnName},
			Value:  id,
		},
	}})
}

// stampTenantCreate fills tenant_id on inserts when the caller left it zero.
// We never override a non-zero value: that lets seeders and migration code
// pin a specific tenant when needed.
func stampTenantCreate(db *gorm.DB) {
	id, ok := shouldInject(db)
	if !ok {
		return
	}
	// shouldInject already verified the field exists; LookUpField is internally
	// cached by GORM so this call is O(1). Pass the field's index path down so
	// stampElement uses FieldByIndex (direct pointer arithmetic) instead of
	// FieldByName (string lookup) on every row.
	field := db.Statement.Schema.LookUpField(tenantFieldName)
	if field == nil {
		return
	}
	indexPath := field.StructField.Index

	// ReflectValue is set during create; it can be a single struct or a slice
	// (batch create). Handle both: stamp every element whose TenantID is zero.
	rv := db.Statement.ReflectValue
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			stampElement(db, rv.Index(i), id, indexPath)
		}
	case reflect.Struct:
		stampElement(db, rv, id, indexPath)
	}
}

// stampElement writes id into the TenantID field on a single row when the
// existing value is zero. We poke the struct directly (rather than via
// schema.Field.Set) because the value also needs to be visible to the rest of
// the create chain — calling SetColumn on the Statement keeps the column
// cache and the struct in sync.
func stampElement(db *gorm.DB, elem reflect.Value, id uint, indexPath []int) {
	if elem.Kind() == reflect.Ptr {
		if elem.IsNil() {
			return
		}
		elem = elem.Elem()
	}
	fv := elem.FieldByIndex(indexPath)
	if !fv.IsValid() || !fv.CanSet() {
		return
	}
	// Only stamp uint zero — any explicit value (including a deliberate non-
	// matching tenant in a backfill job) is preserved.
	if fv.Kind() != reflect.Uint || fv.Uint() != 0 {
		return
	}
	fv.SetUint(uint64(id))
	db.Statement.SetColumn(tenantColumnName, id)
}
