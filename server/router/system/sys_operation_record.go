package system

import (
	"github.com/gin-gonic/gin"
)

type OperationRecordRouter struct{}

func (s *OperationRecordRouter) InitSysOperationRecordRouter(Router *gin.RouterGroup) {
	operationRecordRouter := Router.Group("sysOperationRecord")
	{
		operationRecordRouter.DELETE("deleteSysOperationRecord", operationRecordApi.DeleteSysOperationRecord)           // delete operation record
		operationRecordRouter.DELETE("deleteSysOperationRecordByIds", operationRecordApi.DeleteSysOperationRecordByIds) // batch delete operation records
		operationRecordRouter.GET("findSysOperationRecord", operationRecordApi.FindSysOperationRecord)                  // get by IDSysOperationRecord
		operationRecordRouter.GET("getSysOperationRecordList", operationRecordApi.GetSysOperationRecordList)            // get operation record list

	}
}
