package service

import (
	"context"

	"github.com/huuhoait/gin-vue-admin/server/global"
	sysreq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TrackOnlineSession is a post-handler hook (registered via
// middleware.RegisterPostHandler in the plugin's init()). It runs AFTER the
// main handler chain, so claims set by JWTAuth are available. Skips silently
// for unauthenticated requests, missing tokens, or unconfigured Redis.
//
// The hook uses context.Background instead of c.Request.Context() because the
// request context is already canceled by the time post-handlers run.
func TrackOnlineSession(c *gin.Context) {
	rawClaims, ok := c.Get("claims")
	if !ok {
		return
	}
	claims, ok := rawClaims.(*sysreq.CustomClaims)
	if !ok || claims == nil {
		return
	}
	token := utils.GetToken(c)
	if token == "" {
		return
	}
	err := Service.Session.Touch(
		context.Background(),
		claims.BaseClaims.UUID.String(),
		claims.BaseClaims.ID,
		claims.BaseClaims.Username,
		claims.BaseClaims.NickName,
		claims.BaseClaims.AuthorityId,
		token,
		c.ClientIP(),
		c.Request.UserAgent(),
	)
	if err != nil {
		global.GVA_LOG.Debug("online-session touch failed", zap.Error(err))
	}
}
