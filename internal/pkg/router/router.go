// Package router
// @Author      : lilinzhen
// @Time        : 2022/5/8 10:15:20
// @Description :
package router

import (
	"blueblog/internal/pkg/core"
	"go.uber.org/zap"
)

type Resource struct {
	Mux    core.Mux
	Logger *zap.Logger
}

type Server struct {
	Mux core.Mux
}

