package audit

import (
	"reflect"
	"sync/atomic"

	"gorm.io/gorm"
)

const (
	colCreatedBy = "CreatedBy"
	colUpdatedBy = "UpdatedBy"
	colDeletedBy = "DeletedBy"
)

// registered guards Register against repeated invocation — the callback
// would otherwise stack and stamp twice on the same write.
var registered atomic.Bool

// Register installs the audit callbacks on the supplied DB. Idempotent
// across multiple calls (atomic CompareAndSwap once-guard). Call once at
// boot, after global.GVA_DB is assigned. Safe to call from a plugin init
// even if the framework also wires it.
func Register(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	if !registered.CompareAndSwap(false, true) {
		return nil
	}
	if err := db.Callback().Create().Before("gorm:before_create").Register("audit:stamp_create", stampCreate); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:before_update").Register("audit:stamp_update", stampUpdate); err != nil {
		return err
	}
	// Soft-delete is implemented in GORM as an UPDATE that sets deleted_at;
	// the same Update callback picks it up. We additionally hook the
	// dedicated Delete chain so hard deletes carry the actor too if the
	// model has a DeletedBy column (mostly for forensics).
	if err := db.Callback().Delete().Before("gorm:before_delete").Register("audit:stamp_delete", stampDelete); err != nil {
		return err
	}
	return nil
}

func hasField(db *gorm.DB, name string) bool {
	if db.Statement.Schema == nil {
		return false
	}
	return db.Statement.Schema.LookUpField(name) != nil
}

// stampCreate fills CreatedBy and UpdatedBy from context. Skips when the
// caller already set a non-zero value (allows backfills / cross-actor
// imports to preserve their explicit origin).
func stampCreate(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	uid, ok := UserIDFromContext(db.Statement.Context)
	if !ok {
		return
	}
	rv := reflect.Indirect(db.Statement.ReflectValue)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			stampCreateElement(db, rv.Index(i), uid)
		}
	case reflect.Struct:
		stampCreateElement(db, rv, uid)
	}
}

// stampUpdate refreshes UpdatedBy on every update. Soft-deletes (which
// GORM dispatches as an UPDATE setting deleted_at) also pass through here
// — we additionally fill DeletedBy when we detect a soft-delete by
// checking whether the statement is updating deleted_at to a non-null
// value.
func stampUpdate(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	uid, ok := UserIDFromContext(db.Statement.Context)
	if !ok {
		return
	}
	if hasField(db, colUpdatedBy) {
		// Always overwrite — UpdatedBy tracks the most recent actor.
		db.Statement.SetColumn(toColumn(db, colUpdatedBy), uid)
	}
	if hasField(db, colDeletedBy) && isSoftDelete(db) {
		db.Statement.SetColumn(toColumn(db, colDeletedBy), uid)
	}
}

// stampDelete catches the explicit Delete callback chain. GORM uses this
// for unscoped/permanent deletes; soft-deletes route through the Update
// chain instead. Setting DeletedBy here is best-effort for forensics —
// the row may be physically gone before anyone reads it.
func stampDelete(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	uid, ok := UserIDFromContext(db.Statement.Context)
	if !ok {
		return
	}
	if hasField(db, colDeletedBy) {
		db.Statement.SetColumn(toColumn(db, colDeletedBy), uid)
	}
}

// isSoftDelete checks the statement's destination value for a non-zero
// DeletedAt — this is how GORM signals "this Update is actually a soft
// delete" to its own callbacks.
func isSoftDelete(db *gorm.DB) bool {
	dest := db.Statement.Dest
	if dest == nil {
		return false
	}
	// Most reliable signal: GORM sets db.Statement.Settings["gorm:soft_delete"].
	// Available in v1.25+. Fall back to scanning the Dest map for a
	// deleted_at key.
	if _, ok := db.Statement.Settings.Load("gorm:soft_delete"); ok {
		return true
	}
	if m, ok := dest.(map[string]interface{}); ok {
		_, has := m["deleted_at"]
		return has
	}
	return false
}

func toColumn(db *gorm.DB, fieldName string) string {
	if db.Statement.Schema == nil {
		return fieldName
	}
	if f := db.Statement.Schema.LookUpField(fieldName); f != nil {
		return f.DBName
	}
	return fieldName
}

// stampCreateElement stamps CreatedBy/UpdatedBy on one row. elem must be a
// struct value or pointer-to-struct (batch Create passes slice elements).
func stampCreateElement(db *gorm.DB, elem reflect.Value, uid uint) {
	elem = reflect.Indirect(elem)
	if elem.Kind() != reflect.Struct {
		return
	}
	if hasField(db, colCreatedBy) {
		setFieldIfZero(db, elem, colCreatedBy, uid)
	}
	if hasField(db, colUpdatedBy) {
		setFieldIfZero(db, elem, colUpdatedBy, uid)
	}
}

func setFieldIfZero(db *gorm.DB, elem reflect.Value, fieldName string, uid uint) {
	f := db.Statement.Schema.LookUpField(fieldName)
	if f == nil || f.Set == nil {
		return
	}
	v, zero := f.ValueOf(db.Statement.Context, elem)
	if !zero && v != nil {
		if u, ok := v.(uint); ok && u != 0 {
			return
		}
		if p, ok := v.(*uint); ok && p != nil && *p != 0 {
			return
		}
	}
	_ = f.Set(db.Statement.Context, elem, uid)
}
