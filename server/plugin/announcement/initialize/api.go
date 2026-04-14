package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{
			Path:        "/info/createInfo",
			Description: "createAnnouncement",
			ApiGroup:    "Announcement",
			Method:      "POST",
		},
		{
			Path:        "/info/deleteInfo",
			Description: "delete announcement",
			ApiGroup:    "Announcement",
			Method:      "DELETE",
		},
		{
			Path:        "/info/deleteInfoByIds",
			Description: "batch delete announcements",
			ApiGroup:    "Announcement",
			Method:      "DELETE",
		},
		{
			Path:        "/info/updateInfo",
			Description: "update announcement",
			ApiGroup:    "Announcement",
			Method:      "PUT",
		},
		{
			Path:        "/info/findInfo",
			Description: "get by IDAnnouncement",
			ApiGroup:    "Announcement",
			Method:      "GET",
		},
		{
			Path:        "/info/getInfoList",
			Description: "getAnnouncementList",
			ApiGroup:    "Announcement",
			Method:      "GET",
		},
	}
	utils.RegisterApis(entities...)
}
