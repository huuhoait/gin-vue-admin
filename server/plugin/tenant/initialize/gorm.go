package initialize

import (
	"context"
	"fmt"

	"github.com/huuhoait/gin-vue-admin/server/global"
	tenantgorm "github.com/huuhoait/gin-vue-admin/server/plugin/tenant/gorm"
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

	// Wire automatic tenant_id scoping. Idempotent — safe even if the plugin
	// init runs more than once (tests, hot-reload).
	if regErr := tenantgorm.Register(global.GVA_DB); regErr != nil {
		zap.L().Error("tenant: failed to register GORM scoping callbacks", zap.Error(regErr))
	}
}
