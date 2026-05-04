package oauth2server

import (
	"context"

	"github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/initialize"
	interfaces "github.com/huuhoait/gin-vue-admin/server/utils/plugin/v2"

	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func init() { interfaces.Register(Plugin) }

func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background()
	initialize.Api(ctx)
	initialize.Menu(ctx)
	initialize.Gorm(ctx)
	initialize.Router(group)
	initialize.Timer()
}
