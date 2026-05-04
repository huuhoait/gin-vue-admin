package initialize

import (
	"context"

	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(_ context.Context) {
	entities := []model.SysApi{
		{
			Path:        "/onlineUsers/list",
			Description: "list online sessions",
			ApiGroup:    "OnlineUsers",
			Method:      "GET",
		},
		{
			Path:        "/onlineUsers/kick",
			Description: "kick an online session",
			ApiGroup:    "OnlineUsers",
			Method:      "POST",
		},
	}
	utils.RegisterApis(entities...)
}
