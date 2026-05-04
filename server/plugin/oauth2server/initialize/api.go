package initialize

import (
	"context"

	model "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(_ context.Context) {
	entities := []model.SysApi{
		{Path: "/oauth2Client/create", Description: "create OAuth2 client", ApiGroup: "OAuth2Client", Method: "POST"},
		{Path: "/oauth2Client/update", Description: "update OAuth2 client", ApiGroup: "OAuth2Client", Method: "PUT"},
		{Path: "/oauth2Client/delete", Description: "delete OAuth2 client", ApiGroup: "OAuth2Client", Method: "DELETE"},
		{Path: "/oauth2Client/find", Description: "find OAuth2 client", ApiGroup: "OAuth2Client", Method: "GET"},
		{Path: "/oauth2Client/list", Description: "list OAuth2 clients", ApiGroup: "OAuth2Client", Method: "GET"},
		{Path: "/oauth2Client/regenerateSecret", Description: "rotate client secret", ApiGroup: "OAuth2Client", Method: "POST"},
		{Path: "/oauth2/authorize", Description: "OAuth2 authorize endpoint", ApiGroup: "OAuth2", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}
