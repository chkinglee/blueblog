// Package article_service
// @Author      : lilinzhen
// @Time        : 2022/5/8 10:26:24
// @Description :
package article_service

import (
	"blueblog/api/blueblog-service/grpc_proto"
	"blueblog/internal/pkg/code"
	"blueblog/internal/pkg/configs"
	"blueblog/internal/pkg/core"
	"blueblog/pkg/grpclient"
	"blueblog/pkg/logger"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
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
	BlueblogServiceAddr string
}


func (s *service) Detail(ctx core.Context, uid string, id string) (info *ArticleData, err error) {
	// 请求blueblog-service查询文章详情
	gRpcConn, err := grpclient.New(s.BlueblogServiceAddr, grpclient.WithDialTimeout(time.Second*10))
	if err != nil {
		logger.ErrorT(ctx.RequestContext(), fmt.Sprintf(code.Text(code.GrpcConnectError), s.BlueblogServiceAddr), zap.Error(err))
		return
	}
	logger.DebugT(ctx, fmt.Sprintf(code.Text(code.GrpcConnectSuccess), s.BlueblogServiceAddr))
	defer func(gRpcConn *grpc.ClientConn) {
		_ = gRpcConn.Close()
	}(gRpcConn)

	gRpcClient := grpc_proto.NewArticleServiceClient(gRpcConn)
	req := new(grpc_proto.DetailRequest)
	req.TraceId = ctx.Trace().ID()
	req.Uid = uid
	req.Id = id
	res, err := gRpcClient.Detail(context.Background(), req)
	if err != nil {
		logger.ErrorT(ctx.RequestContext(), fmt.Sprintf("send request to blueblog-service err. [%s]", s.BlueblogServiceAddr), zap.Error(err))
		return nil, err
	}
	logger.InfoT(ctx.RequestContext(), "send request to blueblog-service success")
	return &ArticleData{
		Title:      res.Article.Title,
		Content:    res.Article.Content,
		CreateTime: res.Article.CreateTime,
	}, nil
}

func New() Service {
	return &service{
		BlueblogServiceAddr: configs.Get().Blues.Blueblog.Service.Addr,
	}
}

func (s *service) i() {}
