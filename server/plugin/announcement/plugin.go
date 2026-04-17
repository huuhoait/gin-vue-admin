package announcement

import (
	"context"
	"github.com/huuhoait/gin-vue-admin/server/plugin/announcement/initialize"
	interfaces "github.com/huuhoait/gin-vue-admin/server/utils/plugin/v2"
	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func init() {
	interfaces.Register(Plugin)
}

func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background()
	// if neededconfigurationFile, PleaseToconfig.ConfigInFillFillconfigurationStructure, AndToDownSideIssueInFillInItsAtconfig.yamlInkey
	// initialize.Viper()
	// auto-registered when installing pluginapidata please refer to belowmethod.Apimethodimplemented in
	initialize.Api(ctx)
	// auto-registered when installing pluginMenudata please refer to belowmethod.Menumethodimplemented in
	initialize.Menu(ctx)
	// auto-registered when installing pluginDictionarydata please refer to belowmethod.Dictionarymethodimplemented in
	initialize.Dictionary(ctx)
	initialize.Gorm(ctx)
	initialize.Router(group)
}
