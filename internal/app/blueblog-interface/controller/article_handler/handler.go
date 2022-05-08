// Package article_handler
// @Author      : lilinzhen
// @Time        : 2022/5/8 10:26:18
// @Description :
package article_handler

import (
	"blueblog/internal/app/blueblog-interface/service/article_service"
	"blueblog/internal/pkg/core"
	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Create 创建/修改
	// @Tags API.article
	// @Router /api/article/{uid} [post]
	Create() core.HandlerFunc

	// Detail 详情
	// @Tags API.article
	// @Router /api/article/{uid}/{id} [get]
	Detail() core.HandlerFunc

	// Delete 删除
	// @Tags API.article
	// @Router /api/article/{uid}/{id} [delete]
	Delete() core.HandlerFunc

	// List 列表
	// @Tags API.article
	// @Router /api/article/{uid} [get]
	List() core.HandlerFunc
}

type handler struct {
	logger         *zap.Logger
	articleService article_service.Service
}

func (h *handler) Create() core.HandlerFunc {
	panic("implement me")
}

func (h *handler) Delete() core.HandlerFunc {
	panic("implement me")
}

func (h *handler) List() core.HandlerFunc {
	panic("implement me")
}

func New(logger *zap.Logger) Handler {
	return &handler{
		logger: logger,
		articleService: article_service.New(),
	}
}

func (h *handler) i() {}
