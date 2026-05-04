package request

import "github.com/huuhoait/gin-vue-admin/server/model/common/request"

type CreateClientReq struct {
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description"`
	RedirectURIs []string `json:"redirectUris" binding:"required,min=1,dive,url"`
	GrantTypes   []string `json:"grantTypes" binding:"required,min=1"`
	Scopes       []string `json:"scopes"`
}

type UpdateClientReq struct {
	ID           uint     `json:"id" binding:"required"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	RedirectURIs []string `json:"redirectUris"`
	GrantTypes   []string `json:"grantTypes"`
	Scopes       []string `json:"scopes"`
	Enabled      *bool    `json:"enabled"`
}

type ClientListReq struct {
	request.PageInfo
}

type IdReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}
