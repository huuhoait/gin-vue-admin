package main

import (
	"github.com/huuhoait/gin-vue-admin/server/core"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/global/i18n"
	"github.com/huuhoait/gin-vue-admin/server/initialize"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @Tag in this section controls ordering; add endpoints that need sorting using the format below
// swag init resolves @Tag from the entry file only (default: main.go)
// Or use --generalInfo to point at another file
// @Tag.Name        Base
// @Tag.Name        SysUser
// @Tag.Description User

// @title                       Gin-Vue-Admin Swagger API
// @version                     v2.9.1
// @description                 Full-stack rapid development with Gin + Vue
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	// Initialize system
	initializeSystem()
	// Run HTTP server
	core.RunServer()
}

// initializeSystem wires all subsystems
// Split out so reload paths can call it deterministically
func initializeSystem() {
	global.GVA_VP = core.Viper() // init Viper
	initialize.OtherInit()
	// Story 8.3: load i18n message bundles before anything can call the
	// response helpers. Path is relative to the admin server working
	// directory (same convention as resource/rbac_model.conf etc.).
	if err := i18n.Load("resource/i18n"); err != nil {
		// Non-fatal: helpers will echo bundle keys instead of translations.
		zap.L().Warn("i18n bundles failed to load", zap.Error(err))
	}
	global.GVA_LOG = core.Zap() // init zap logger
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // connect database via GORM
	initialize.Timer()
	initialize.DBList()
	initialize.SetupHandlers() // register global handlers
	if global.GVA_DB != nil {
		initialize.RegisterTables() // migrate / register tables
	}
}
