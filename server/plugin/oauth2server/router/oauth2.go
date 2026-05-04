package router

import (
	"github.com/huuhoait/gin-vue-admin/server/middleware"

	"github.com/gin-gonic/gin"
)

type oauth2Router struct{}

// Init mounts:
//   - /oauth2/authorize on the PRIVATE group: authorization codes are only
//     issued to admin users authenticated via JWT.
//   - /oauth2/token, /oauth2/introspect, /oauth2/revoke on the PUBLIC group:
//     these endpoints authenticate clients via Basic auth (or body params),
//     not via the admin JWT.
func (r *oauth2Router) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("oauth2").Use(middleware.OperationRecord())
		group.GET("authorize", apiOAuth2Flow.Authorize)
	}
	{
		group := public.Group("oauth2")
		group.POST("token", apiOAuth2Flow.Token)
		group.POST("introspect", apiOAuth2Flow.Introspect)
		group.POST("revoke", apiOAuth2Flow.Revoke)
	}
}
