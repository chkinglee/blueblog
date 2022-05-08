// Package article_service
// @Author      : lilinzhen
// @Time        : 2022/5/8 10:26:24
// @Description :
package article_service

import (
	"blueblog/internal/pkg/core"
)

var _ Service = (*service)(nil)

type CreateData struct {
}

type UpdateData struct {
}

type SearchData struct {
}

type ArticleData struct {
	Title      string
	Content    string
	CreateTime string
}

type Service interface {
	i()
	Detail(ctx core.Context, uid string, id string) (info *ArticleData, err error)
}

type service struct {
}

func (s *service) Detail(ctx core.Context, uid string, id string) (info *ArticleData, err error) {
	// TODO 请求blueblog-service查询文章详情
	ctx.Logger().Info("send request to blueblog-service")
	if err != nil {
		return
	}
	return
}

func New() Service {
	return &service{}
}

func (s *service) i() {}
