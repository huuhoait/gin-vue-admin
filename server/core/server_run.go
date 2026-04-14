package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

// initServer StartServiceAndImplementExcellentelegantDisable
func initServer(address string, router *gin.Engine, readTimeout, writeTimeout time.Duration) {
	// CreateService
	srv := &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// AtgoroutineInStartService
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
			zap.L().Error("serverStartfailed", zap.Error(err))
			os.Exit(1)
		}
	}()

	// etc.PendingInBreakTrustNumberByExcellentelegantLocationDisableServer
	quit := make(chan os.Signal, 1)
	// kill (NoneParameter) defaultsend syscall.SIGTERM
	// kill -2 send syscall.SIGINT
	// kill -9 send syscall.SIGKILL, ButYesUnableBycatchGain, ofByNotNeedAdd
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("DisableWEBService...")

	// set5Secondoftimeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("WEBServiceDisableException", zap.Error(err))
	}

	zap.L().Info("WEBServiceAlreadyDisable")
}
