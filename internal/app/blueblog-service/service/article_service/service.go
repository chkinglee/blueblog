// Package article_service
// @Author      : lilinzhen
// @Time        : 2022/5/8 15:14:05
// @Description :
package article_service

import (
	"blueblog/api/blueblog-service/grpc_proto"
	"blueblog/internal/app/blueblog-service/repository/article_repo"
	"blueblog/internal/pkg/configs"
	"blueblog/internal/pkg/core"
	"blueblog/internal/pkg/db"
	"blueblog/internal/pkg/repo"
	"blueblog/pkg/logger"
	"blueblog/pkg/time_parse"
	"blueblog/pkg/trace"
	"context"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type Service interface {
	Detail(ctx context.Context, req *grpc_proto.DetailRequest) (res *grpc_proto.DetailResponse, err error)
}

type service struct {
	logger *zap.Logger
	rp     repo.Repo
}

type Article struct {
	Exist      interface{}
}

func (a *Article) IsExist() bool {
	return a.Exist != nil && strings.ToLower(a.Exist.(string)) == "true"
}

func (s *service) Detail(ctx context.Context, req *grpc_proto.DetailRequest) (res *grpc_proto.DetailResponse, err error) {
	res = new(grpc_proto.DetailResponse)
	stdCtx := core.StdContext{
		Context: ctx,
		Trace: &trace.Trace{
			Identifier: req.TraceId,
		},
		Logger: s.logger,
	}
	logger.InfoT(stdCtx, "receive request", zap.Any("req", req))

	// 先从cache中查
	key := fmt.Sprintf(configs.CacheKey4Article, req.Uid, req.Id)
	ok := s.rp.GetCache().Exists(key)
	if ok {
		// 如果key存在，判断article是否存在，为了防止缓存击穿，Job也会将不存在的article缓存一下
		values, err := s.rp.GetCache().HMGet(key, []string{"exist", "title", "content", "create_time"})
		if err != nil {
			logger.ErrorT(stdCtx, "get values from cache err.", zap.Error(err))
		}
		logger.InfoT(stdCtx, "get values from cache", zap.Any("values", values))
		article := Article{
			Exist: values[0],
		}
		if article.IsExist() {
			res.Article = &grpc_proto.Article{
				Title:      values[1].(string),
				Content:    values[2].(string),
				CreateTime: values[3].(string),
			}
			return res, nil
		} else {
			// TODO grpc跨网络传输error，客户端无法处理
			return nil, configs.ErrNotExist
		}
	} else {
		// 通知blueblog-job缓存该数据
		go func() {
			data := map[string]interface{}{
				"uid": req.Uid,
				"id": req.Id,
			}

			err = s.rp.GetRabbitMQ().Publish(
				configs.MQExchange,
				fmt.Sprintf(configs.MQRoutingKey4Article, "read"),
				data,
			)
			if err != nil {
				logger.ErrorT(stdCtx, "send message to MQ err.", zap.Error(err))
			}
		}()
		// cache中查不到的话直接查db
		logger.InfoT(stdCtx, "cache miss and then query db")
		qb := article_repo.NewQueryBuilder()
		qb.WhereUid(db.EqualPredicate, req.Uid)
		qb.WhereId(db.EqualPredicate, req.Id)

		info, err := qb.QueryOne(s.rp.GetDb().GetDbW().WithContext(ctx))
		if err != nil {
			return nil, err
		}
		if info == nil {
			return nil, configs.ErrNotExist
		}

		res.Article = &grpc_proto.Article{
			Title:      info.Title,
			Content:    info.Content,
			CreateTime: info.CreateTime.Format(time_parse.CSTLayout),
		}

		return res, nil
	}
}

func New(logger *zap.Logger, rp repo.Repo) Service {
	return &service{
		logger: logger,
		rp:     rp,
	}
}
