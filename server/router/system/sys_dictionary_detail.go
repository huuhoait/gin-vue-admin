package system

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DictionaryDetailRouter struct{}

func (s *DictionaryDetailRouter) InitSysDictionaryDetailRouter(Router *gin.RouterGroup) {
	dictionaryDetailRouter := Router.Group("sysDictionaryDetail").Use(middleware.OperationRecord())
	dictionaryDetailRouterWithoutRecord := Router.Group("sysDictionaryDetail")
	{
		dictionaryDetailRouter.POST("createSysDictionaryDetail", dictionaryDetailApi.CreateSysDictionaryDetail)   // createSysDictionaryDetail
		dictionaryDetailRouter.DELETE("deleteSysDictionaryDetail", dictionaryDetailApi.DeleteSysDictionaryDetail) // deleteSysDictionaryDetail
		dictionaryDetailRouter.PUT("updateSysDictionaryDetail", dictionaryDetailApi.UpdateSysDictionaryDetail)    // updateSysDictionaryDetail
	}
	{
		dictionaryDetailRouterWithoutRecord.GET("findSysDictionaryDetail", dictionaryDetailApi.FindSysDictionaryDetail)           // get by IDSysDictionaryDetail
		dictionaryDetailRouterWithoutRecord.GET("getSysDictionaryDetailList", dictionaryDetailApi.GetSysDictionaryDetailList)     // getSysDictionaryDetailList
		dictionaryDetailRouterWithoutRecord.GET("getDictionaryTreeList", dictionaryDetailApi.GetDictionaryTreeList)               // get dictionaryDetailstree structure
		dictionaryDetailRouterWithoutRecord.GET("getDictionaryTreeListByType", dictionaryDetailApi.GetDictionaryTreeListByType)   // According todictionary typeget dictionaryDetailstree structure
		dictionaryDetailRouterWithoutRecord.GET("getDictionaryDetailsByParent", dictionaryDetailApi.GetDictionaryDetailsByParent) // according to parentIDget dictionaryDetails
		dictionaryDetailRouterWithoutRecord.GET("getDictionaryPath", dictionaryDetailApi.GetDictionaryPath)                       // get dictionaryfull detailpath
	}
}
