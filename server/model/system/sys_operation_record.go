// auto-generate templateSysOperationRecord
package system

import (
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

// If it contains time.Time please import the time package yourself
type SysOperationRecord struct {
	global.GVA_MODEL
	Ip           string        `json:"ip" form:"ip" gorm:"column:ip;comment:request IP"`                                   // request IP
	Method       string        `json:"method" form:"method" gorm:"column:method;comment:request method"`                       // request method
	Path         string        `json:"path" form:"path" gorm:"column:path;comment:request path"`                             // request path
	Status       int           `json:"status" form:"status" gorm:"column:status;comment:request status"`                       // request status
	Latency      time.Duration `json:"latency" form:"latency" gorm:"column:latency;comment:latency" swaggertype:"string"` // latency
	Agent        string        `json:"agent" form:"agent" gorm:"type:text;column:agent;comment:agent"`                  // agent
	ErrorMessage string        `json:"error_message" form:"error_message" gorm:"column:error_message;comment:error message"`  // error message
	Body         string        `json:"body" form:"body" gorm:"type:text;column:body;comment:request body"`                 // request body
	Resp         string        `json:"resp" form:"resp" gorm:"type:text;column:resp;comment:response body"`                 // response body
	UserID       int           `json:"user_id" form:"user_id" gorm:"column:user_id;comment:user ID"`                    // user ID
	RequestID    string        `json:"request_id" form:"request_id" gorm:"column:request_id;size:64;index;comment:correlation id"`
	User         SysUser       `json:"user"`
}
