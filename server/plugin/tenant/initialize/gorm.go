package initialize

import (
	"context"
	"fmt"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/tenant/model"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.Tenant),
		new(model.UserTenant),
	)
	if err != nil {
		err = errors.Wrap(err, "tenant: table migration failed")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
