package core

import (
    "fmt"
    "github.com/huuhoait/gin-vue-admin/server/core/internal"
    "github.com/huuhoait/gin-vue-admin/server/global"
    "github.com/huuhoait/gin-vue-admin/server/utils"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "os"
)

// Zap get zap.Logger
// Author [SliverHorn](https://github.com/SliverHorn)
func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.GVA_CONFIG.Zap.Director); !ok { // JudgeYesNoHaveDirectorFileclip
		fmt.Printf("create %v directory\n", global.GVA_CONFIG.Zap.Director)
		_ = os.Mkdir(global.GVA_CONFIG.Zap.Director, os.ModePerm)
	}
	levels := global.GVA_CONFIG.Zap.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i])
		cores = append(cores, core)
	}
    // build base logger(Error LevelofInboundLogicAlreadyAtCustom ZapCore InHandle)
    logger = zap.New(zapcore.NewTee(cores...))
	// Enable Error AndByUpperLevelofheapStackcapture, ensure entry.Stack CanUse
	opts := []zap.Option{zap.AddStacktrace(zapcore.ErrorLevel)}
	if global.GVA_CONFIG.Zap.ShowLine {
		opts = append(opts, zap.AddCaller())
	}
	logger = logger.WithOptions(opts...)
	return logger
}
