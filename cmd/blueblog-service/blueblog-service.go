// Package blueblog_service
// @Author      : lilinzhen
// @Time        : 2022/5/8 14:47:49
// @Description :
package main

import (
	"blueblog/api/blueblog-service/api"
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
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

func init() {
	env.Active().WithApp(configs.AppNameForService)
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
		//repo.WithRabbitMQ(),
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

	// 初始化 gRPC 服务
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.Get().Server.GrpcPort))
	if err != nil {
		accessLogger.Fatal("failed to listen", zap.Error(err))
	}
	grpcServer := grpc.NewServer()
	api.SetGrpcRegister(grpcServer, accessLogger, rp)
	accessLogger.Info("gRPC server listen on " + fmt.Sprintf(":%d", configs.Get().Server.GrpcPort))

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil && err != grpc.ErrServerStopped {
			accessLogger.Fatal("gRPC server startup err", zap.Error(err))
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

		// 关闭 grpc server
		func() {
			grpcServer.Stop()
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
