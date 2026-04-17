package {{ .Package }}

import (
	"context"
	"{{.Module}}/plugin/{{ .Package }}/initialize"
	interfaces "{{.Module}}/utils/plugin/v2"
	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func init() {
	interfaces.Register(Plugin)
}


// If a config file is needed, populate the config struct in config.Config and wire up the key from config.yaml in the method below, adding:
// initialize.Viper()
// To auto-register API data when installing the plugin, implement it in the .Api method and add:
// initialize.Api(ctx)
// To auto-register menus when installing the plugin, implement it in the .Menu method and add:
// initialize.Menu(ctx)
// To auto-register dictionaries when installing the plugin, implement it in the .Dictionary method and add:
// initialize.Dictionary(ctx)
func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background() 
	initialize.Gorm(ctx)
	initialize.Router(group)
}
