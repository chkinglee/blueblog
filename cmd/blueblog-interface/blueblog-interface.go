// Package blueblog_interface
// @Author      : lilinzhen
// @Time        : 2022/5/7 23:56:38
// @Description :
package main

import (
	"blueblog/internal/app/blueblog-interface/router"
	"blueblog/internal/pkg/configs"
	"blueblog/pkg/env"
	"blueblog/pkg/logger"
	"blueblog/pkg/shutdown"
	"blueblog/pkg/time_parse"
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func init() {
	env.Active().WithApp(configs.AppNameForInterface)
	configs.Init()
}

func main() {
	// 初始化 access logger
	accessLogger, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		logger.WithTimeLayout(time_parse.CSTLayout),
		logger.WithFilePath(configs.Get().Logger.File),
		logger.WithLevel(configs.Get().Logger.Level),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = accessLogger.Sync()
	}()

	// 初始化 HTTP 服务
	s, err := router.NewHTTPServer(accessLogger)
	if err != nil {
		panic(err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.Get().Server.Port),
		Handler: s.Mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			accessLogger.Fatal("http server startup err", zap.Error(err))
		}
	}()

	// 优雅关闭
	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				accessLogger.Error("http server shutdown err", zap.Error(err))
			}
		},
	)

}
