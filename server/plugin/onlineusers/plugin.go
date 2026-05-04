package onlineusers

import (
	"context"

	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/huuhoait/gin-vue-admin/server/plugin/onlineusers/initialize"
	"github.com/huuhoait/gin-vue-admin/server/plugin/onlineusers/service"
	interfaces "github.com/huuhoait/gin-vue-admin/server/utils/plugin/v2"

	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func init() {
	interfaces.Register(Plugin)
	// Register the post-handler hook at package init so it is in place before
	// initialize.Routers() installs middleware.PostHandlerChain.
	middleware.RegisterPostHandler(service.TrackOnlineSession)
}

func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background()
	initialize.Api(ctx)
	initialize.Menu(ctx)
	initialize.Router(group)
	initialize.Timer()
}
