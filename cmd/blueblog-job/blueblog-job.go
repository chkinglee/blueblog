// Package blueblog_job
// @Author      : lilinzhen
// @Time        : 2022/5/8 20:12:43
// @Description :
package main

import (
	"blueblog/internal/app/blueblog-job/job/article-job"
	"blueblog/internal/pkg/configs"
	"blueblog/internal/pkg/core"
	"blueblog/internal/pkg/repo"
	"blueblog/internal/pkg/router"
	"blueblog/pkg/env"
	"blueblog/pkg/errors"
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
	env.Active().WithApp(configs.AppNameForJob)
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

	// 初始化 repo 连接
	rp := repo.NewRepo(accessLogger,
		repo.WithDb(),
		repo.WithCache(),
		repo.WithRabbitMQ(),
	)

	// 初始化 HTTP 服务
	s, err := NewHTTPServer(accessLogger)
	if err != nil {
		panic(err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.Get().Server.HttpPort),
		Handler: s.Mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			accessLogger.Fatal("http server startup err", zap.Error(err))
		}
	}()

	go article_job.ReadArticle(accessLogger, rp)

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

		// 关闭 repo 连接
		func() {
			rp.Close(accessLogger)
		},
	)

}

func NewHTTPServer(logger *zap.Logger) (*router.Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	r := new(router.Resource)
	r.Logger = logger

	mux, err := core.New(logger,
		core.WithEnableCors(),
		core.WithEnableRate(),
		//core.WithRecordMetrics(metrics.RecordMetrics),
	)

	if err != nil {
		panic(err)
	}

	r.Mux = mux

	// 设置 API 路由
	//api.SetApiRouter(r)

	s := new(router.Server)
	s.Mux = mux

	return s, nil
}
