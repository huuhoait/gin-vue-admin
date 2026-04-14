package system

import (
	"github.com/gin-gonic/gin"
)

type AutoCodeHistoryRouter struct{}

func (s *AutoCodeRouter) InitAutoCodeHistoryRouter(Router *gin.RouterGroup) {
	autoCodeHistoryRouter := Router.Group("autoCode")
	{
		autoCodeHistoryRouter.POST("getMeta", autocodeHistoryApi.First)         // According toidgetmetaInformation
		autoCodeHistoryRouter.POST("rollback", autocodeHistoryApi.RollBack)     // rollback
		autoCodeHistoryRouter.POST("delSysHistory", autocodeHistoryApi.Delete)  // deleterollbackRecord
		autoCodeHistoryRouter.POST("getSysHistory", autocodeHistoryApi.GetList) // getrollbackRecordPagination
	}
}
