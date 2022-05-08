// Package api
// @Author      : lilinzhen
// @Time        : 2022/5/8 15:00:32
// @Description :
package api

import (
	"blueblog/api/blueblog-service/grpc_proto"
	"blueblog/internal/app/blueblog-service/service/article_service"
	"blueblog/internal/pkg/repo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func SetGrpcRegister(server *grpc.Server, logger *zap.Logger, rp repo.Repo) {

	grpc_proto.RegisterArticleServiceServer(server, article_service.New(logger, rp))
}
