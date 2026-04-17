package initialize

import (
	"fmt"
	"{{.Module}}/global"
	"{{.Module}}/plugin/{{ .Package }}/plugin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Viper() {
	err := global.GVA_VP.UnmarshalKey("{{ .Package }}", &plugin.Config)
	if err != nil {
		err = errors.Wrap(err, "Failed to initialize config file!")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
