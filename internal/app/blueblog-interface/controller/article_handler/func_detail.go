// Package article_handler
// @Author      : lilinzhen
// @Time        : 2022/5/8 10:45:13
// @Description :
package article_handler

import (
	"blueblog/internal/pkg/code"
	"blueblog/internal/pkg/configs"
	"blueblog/internal/pkg/core"
	"blueblog/pkg/errno"
	"errors"
	"net/http"
)

type detailRequest struct {
	Uid string `uri:"uid"` // 用户ID
	Id  string `uri:"id"`  // 文章ID
}

type detailResponse struct {
	Id                string `json:"id"`                  // 文章ID
	UserName          string `json:"user_name"`           // 用户名
	ArticleTitle      string `json:"article_title"`       // 文章标题
	ArticleContent    string `json:"article_content"`     // 文章内容
	ArticleCreateTime string `json:"article_create_time"` // 创建时间
}

// Detail 详情
// @Summary 详情
// @Description 详情
// @Tags API.article
// @Accept json
// @Produce json
// @Param Request body detailRequest true "请求信息"
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Router /api/article/{uid}/{id} [get]
func (h *handler) Detail() core.HandlerFunc {
	return func(c core.Context) {
		req := new(detailRequest)
		res := new(detailResponse)
		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(errno.NewError(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithErr(err),
			)
			return
		}

		/* TODO 此处可以包含一些鉴权操作，例如请求「用户微服务」验证uid是否存在，这样的验证也可以放到统一入口去做
			暂时假定从「用户微服务」中得到的username是admin
		*/
		username := "admin"

		info, err := h.articleService.Detail(c, req.Uid, req.Id)
		if err != nil {
			if errors.As(err, &configs.ErrNotExist) {
				c.AbortWithError(errno.NewError(
					http.StatusBadRequest,
					code.ArticleNotExist,
					code.Text(code.ArticleNotExist)).WithErr(err),
				)
				return
			}
			c.AbortWithError(errno.NewError(
				http.StatusBadRequest,
				code.ArticleDetailError,
				code.Text(code.ArticleDetailError)).WithErr(err),
			)
			return
		}

		res.Id = req.Id
		res.ArticleTitle = info.Title
		res.ArticleContent = info.Content
		res.ArticleCreateTime = info.CreateTime
		res.UserName = username
		c.Payload(res)
	}
}
