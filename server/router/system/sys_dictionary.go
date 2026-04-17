package system

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DictionaryRouter struct{}

func (s *DictionaryRouter) InitSysDictionaryRouter(Router *gin.RouterGroup) {
	sysDictionaryRouter := Router.Group("sysDictionary").Use(middleware.OperationRecord())
	sysDictionaryRouterWithoutRecord := Router.Group("sysDictionary")
	{
		sysDictionaryRouter.POST("createSysDictionary", dictionaryApi.CreateSysDictionary)   // createSysDictionary
		sysDictionaryRouter.DELETE("deleteSysDictionary", dictionaryApi.DeleteSysDictionary) // deleteSysDictionary
		sysDictionaryRouter.PUT("updateSysDictionary", dictionaryApi.UpdateSysDictionary)    // updateSysDictionary
		sysDictionaryRouter.POST("importSysDictionary", dictionaryApi.ImportSysDictionary)   // importSysDictionary
		sysDictionaryRouter.GET("exportSysDictionary", dictionaryApi.ExportSysDictionary)    // exportSysDictionary
	}
	{
		sysDictionaryRouterWithoutRecord.GET("findSysDictionary", dictionaryApi.FindSysDictionary)       // get by IDSysDictionary
		sysDictionaryRouterWithoutRecord.GET("getSysDictionaryList", dictionaryApi.GetSysDictionaryList) // getSysDictionaryList
	}
}
