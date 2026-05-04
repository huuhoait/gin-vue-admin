package initialize

import (
	"context"
	"fmt"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/model"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.OAuth2Client),
		new(model.OAuth2AuthCode),
		new(model.OAuth2Token),
	)
	if err != nil {
		err = errors.Wrap(err, "oauth2server: table migration failed")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
