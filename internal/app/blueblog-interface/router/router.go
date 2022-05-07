// Package router
// @Author      : lilinzhen
// @Time        : 2022/5/8 01:27:25
// @Description :
package router

import (
	"blueblog/internal/pkg/core"
	//"blueblog/internal/pkg/metrics"
	"blueblog/pkg/errors"

	"go.uber.org/zap"
)

type resource struct {
	mux    core.Mux
	logger *zap.Logger
}

type Server struct {
	Mux core.Mux
}

func NewHTTPServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	r := new(resource)
	r.logger = logger

	mux, err := core.New(logger,
		core.WithEnableCors(),
		core.WithEnableRate(),
		//core.WithRecordMetrics(metrics.RecordMetrics),
	)

	if err != nil {
		panic(err)
	}

	r.mux = mux

	// 设置 API 路由
	//setApiRouter(r)

	s := new(Server)
	s.Mux = mux

	return s, nil
}

