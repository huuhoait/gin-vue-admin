package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/huuhoait/gin-vue-admin/server/model/common"
	"github.com/huuhoait/gin-vue-admin/server/utils"
)

//@author: [songzhibin97](https://github.com/songzhibin97)
//@function: ClearTable
//@description: CleanDatabaseTableData
//@param: db(DatabaseObject) *gorm.DB, tableName(table name) string, compareField(CompareField) string, interval(IntervalSeparate) string
//@return: error

// clearTableWhitelist pins the set of tables this job is allowed to touch.
// fmt.Sprintf() builds the DELETE statement, so an attacker-controlled table
// name would be an SQL injection vector; the whitelist keeps that closed.
var clearTableWhitelist = map[string]bool{
	"sys_operation_records": true,
	"jwt_blacklists":        true,
}

func ClearTable(db *gorm.DB) error {
	if db == nil {
		return errors.New("db Cannot be empty")
	}

	// When multiple replicas run this cron, racing DELETEs waste locks and
	// bloat replication traffic. A 10-minute Redis lease is plenty for a
	// one-shot sweep of log tables and expires quickly enough that a
	// restarted leader can take over on the next tick.
	release, acquired, err := utils.TryDistLock(context.Background(), "cron:clear-table", 10*time.Minute)
	if err != nil {
		return err
	}
	if !acquired {
		return nil // another instance is running the sweep; that is expected.
	}
	defer release()

	var ClearTableDetail []common.ClearDB

	ClearTableDetail = append(ClearTableDetail, common.ClearDB{
		TableName:    "sys_operation_records",
		CompareField: "created_at",
		Interval:     "2160h",
	})

	ClearTableDetail = append(ClearTableDetail, common.ClearDB{
		TableName:    "jwt_blacklists",
		CompareField: "created_at",
		Interval:     "168h",
	})

	for _, detail := range ClearTableDetail {
		if !clearTableWhitelist[detail.TableName] {
			return fmt.Errorf("table %q is not allowed for clear-table cron", detail.TableName)
		}
		duration, err := time.ParseDuration(detail.Interval)
		if err != nil {
			return err
		}
		if duration < 0 {
			return errors.New("parse duration < 0")
		}
		err = db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s < ?", detail.TableName, detail.CompareField), time.Now().Add(-duration)).Error
		if err != nil {
			return err
		}
	}
	return nil
}
