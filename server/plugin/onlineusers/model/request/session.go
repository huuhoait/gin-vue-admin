package request

import "github.com/huuhoait/gin-vue-admin/server/model/common/request"

type ListSessionsReq struct {
	request.PageInfo
	Username string `json:"username" form:"username"`
}

type KickSessionReq struct {
	UUID string `json:"uuid" form:"uuid" binding:"required"`
}
