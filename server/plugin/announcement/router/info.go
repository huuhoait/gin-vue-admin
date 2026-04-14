package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var Info = new(info)

type info struct{}

// Init initialize Announcement route information
func (r *info) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("info").Use(middleware.OperationRecord())
		group.POST("createInfo", apiInfo.CreateInfo)             // createAnnouncement
		group.DELETE("deleteInfo", apiInfo.DeleteInfo)           // delete announcement
		group.DELETE("deleteInfoByIds", apiInfo.DeleteInfoByIds) // batch delete announcements
		group.PUT("updateInfo", apiInfo.UpdateInfo)              // update announcement
	}
	{
		group := private.Group("info")
		group.GET("findInfo", apiInfo.FindInfo)       // get by IDAnnouncement
		group.GET("getInfoList", apiInfo.GetInfoList) // getAnnouncementList
	}
	{
		group := public.Group("info")
		group.GET("getInfoDataSource", apiInfo.GetInfoDataSource) // getAnnouncementdata source
		group.GET("getInfoPublic", apiInfo.GetInfoPublic)         // getAnnouncementList
	}
}
