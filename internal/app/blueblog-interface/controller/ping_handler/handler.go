// Package ping_handler
// @Author      : lilinzhen
// @Time        : 2022/5/8 10:13:15
// @Description :
package ping_handler

import (
	"blueblog/internal/pkg/core"
	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	Ping() core.HandlerFunc
}

type handler struct {
	logger *zap.Logger
}

func (h handler) i() {
}

func (h handler) Ping() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload("PONG")
	}
}

func New(logger *zap.Logger) Handler {
	return &handler{
		logger: logger,
	}
}
