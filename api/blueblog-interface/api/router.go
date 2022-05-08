// Package api
// @Author      : lilinzhen
// @Time        : 2022/5/8 10:07:20
// @Description :
package api

import (
	"blueblog/internal/app/blueblog-interface/controller/article_handler"
	"blueblog/internal/app/blueblog-interface/controller/ping_handler"
	"blueblog/internal/pkg/router"
)

func SetApiRouter(r *router.Resource) {

	// interface作为对外流量入口，在router应当加入用户解析、鉴权校验等操作，即作为rpc client请求外部rpc server。在此省略

	// agent
	agent := r.Mux.Group("/agent")
	{
		pingHandler := ping_handler.New(r.Logger)
		agent.GET("/ping", pingHandler.Ping())
	}

	// article
	article := r.Mux.Group("/article")
	{
		articleHandler := article_handler.New(r.Logger)
		article.GET("/:uid/:id", articleHandler.Detail())
	}
}

