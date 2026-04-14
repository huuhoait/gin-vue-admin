package internal

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	astutil "github.com/flipped-aurora/gin-vue-admin/server/utils/ast"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/stacktrace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

type ZapCore struct {
	level zapcore.Level
	zapcore.Core
}

func NewZapCore(level zapcore.Level) *ZapCore {
	entity := &ZapCore{level: level}
	syncer := entity.WriteSyncer()
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == level
	})
	entity.Core = zapcore.NewCore(global.GVA_CONFIG.Zap.Encoder(), syncer, levelEnabler)
	return entity
}

func (z *ZapCore) WriteSyncer(formats ...string) zapcore.WriteSyncer {
	cutter := NewCutter(
		global.GVA_CONFIG.Zap.Director,
		z.level.String(),
		global.GVA_CONFIG.Zap.RetentionDay,
		CutterWithLayout(time.DateOnly),
		CutterWithFormats(formats...),
	)
	if global.GVA_CONFIG.Zap.LogInConsole {
		multiSyncer := zapcore.NewMultiWriteSyncer(os.Stdout, cutter)
		return zapcore.AddSync(multiSyncer)
	}
	return zapcore.AddSync(cutter)
}

func (z *ZapCore) Enabled(level zapcore.Level) bool {
	return z.level == level
}

func (z *ZapCore) With(fields []zapcore.Field) zapcore.Core {
	return z.Core.With(fields)
}

func (z *ZapCore) Check(entry zapcore.Entry, check *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(entry.Level) {
		return check.AddCore(entry, z)
	}
	return check
}

func (z *ZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	for i := 0; i < len(fields); i++ {
		if fields[i].Key == "business" || fields[i].Key == "folder" || fields[i].Key == "directory" {
			syncer := z.WriteSyncer(fields[i].String)
			z.Core = zapcore.NewCore(global.GVA_CONFIG.Zap.Encoder(), syncer, z.level)
		}
	}
	// FirstWriteInOriginalLogItemTag
	err := z.Core.Write(entry, fields)

	// capture Error AndByUpperLevelLogAndInbound, And CanRaiseFetch zap.Error(err) ofErrorcontent
	if entry.Level >= zapcore.ErrorLevel {
		// AvoidAnd GORM zap WriteInmutualRelatedRecursion:SkipBy gorm logger writer TriggerofLog
		if strings.Contains(entry.Caller.File, "gorm_logger_writer.go") {
			return err
		}

		form := "AfterEnd"
		level := entry.Level.String()
		// GenerateBasicInformation
		info := entry.Message

		// RaiseFetch zap.Error(err) content
		var errStr string
		for i := 0; i < len(fields); i++ {
			f := fields[i]
			if f.Type == zapcore.ErrorType || f.Key == "error" || f.Key == "err" {
				if f.Interface != nil {
					errStr = fmt.Sprintf("%v", f.Interface)
				} else if f.String != "" {
					errStr = f.String
				}
				break
			}
		}
		if errStr != "" {
			info = fmt.Sprintf("%s | Error: %s", info, errStr)
		}

		// AdditionalSourceAndheapStackInformation
		if entry.Caller.File != "" {
			info = fmt.Sprintf("%s \n SourceFile:%s:%d", info, entry.Caller.File, entry.Caller.Line)
		}
		stack := entry.Stack
		if stack != "" {
			info = fmt.Sprintf("%s \n InvokeStack:%s", info, stack)
			// ParseFinalBusinessInvokeSide, AndRaiseFetchItsmethodSourceCode
			if frame, ok := stacktrace.FindFinalCaller(stack); ok {
				fnName, fnSrc, sLine, eLine, exErr := astutil.ExtractFuncSourceByPosition(frame.File, frame.Line)
				if exErr == nil {
					info = fmt.Sprintf("%s \n final callermethod:%s:%d (%s lines %d-%d)\n----- GenerateLogofmethodGenerationCodeSuch AsDown -----\n%s", info, frame.File, frame.Line, fnName, sLine, eLine, fnSrc)
				} else {
					info = fmt.Sprintf("%s \n final callermethod:%s:%d (%s) | extract_err=%v", info, frame.File, frame.Line, fnName, exErr)
				}
			}
		}

		// UseAfterDeviceContext, AvoidDependency gin.Context
		ctx := context.Background()
		_ = service.ServiceGroupApp.SystemServiceGroup.SysErrorService.CreateSysError(ctx, &system.SysError{
			Form:  &form,
			Info:  &info,
			Level: level,
		})
	}
	return err
}

func (z *ZapCore) Sync() error {
	return z.Core.Sync()
}
