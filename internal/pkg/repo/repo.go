// Package repo
// @Author      : lilinzhen
// @Time        : 2022/5/7 14:40:01
// @Description :
package repo

import (
	"blueblog/internal/pkg/cache"
	"blueblog/internal/pkg/db"
	"blueblog/internal/pkg/message_queue"
	"go.uber.org/zap"
)

type repo struct {
	db       db.Repo
	cache    cache.Repo
	rabbitmq message_queue.RabbitMQRepo
}

type Repo interface {
	i()
	GetDb() db.Repo
	GetCache() cache.Repo
	GetRabbitMQ() message_queue.RabbitMQRepo
	Close(logger *zap.Logger)
}

type Option func(*option)

type option struct {
	needDb       bool
	needCache    bool
	needRabbitMQ bool
}

func WithDb() Option {
	return func(o *option) {
		o.needDb = true
	}
}

func WithCache() Option {
	return func(o *option) {
		o.needCache = true
	}
}

func WithRabbitMQ() Option {
	return func(o *option) {
		o.needRabbitMQ = true
	}
}

func (r *repo) GetDb() db.Repo {
	return r.db
}

func (r *repo) GetCache() cache.Repo {
	return r.cache
}

func (r *repo) GetRabbitMQ() message_queue.RabbitMQRepo {
	if r.rabbitmq.IsClosed() {
		rabbitmqRepo, _ := message_queue.RabbitMQNew()
		r.rabbitmq = rabbitmqRepo
	}
	return r.rabbitmq
}

func (r *repo) Close(logger *zap.Logger) {
	if r.db != nil {
		if err := r.db.DbWClose(); err != nil {
			logger.Error("dbw close err", zap.Error(err))
		}
		if err := r.db.DbRClose(); err != nil {
			logger.Error("dbr close err", zap.Error(err))
		}
	}

	if r.cache != nil {
		if err := r.cache.Close(); err != nil {
			logger.Error("cache close err", zap.Error(err))
		}
	}

	if r.rabbitmq != nil {
		if err := r.rabbitmq.Close(); err != nil {
			logger.Error("rabbitmq close err", zap.Error(err))
		}
	}
}

func (r *repo) i() {}

func NewRepo(logger *zap.Logger, opts ...Option) Repo {
	opt := &option{}
	rp := &repo{}
	for _, f := range opts {
		f(opt)
	}

	if opt.needDb {
		// 初始化 DB
		dbRepo, err := db.New()
		if err != nil {
			logger.Fatal("new db err", zap.Error(err))
		}
		logger.Info("connect to DB success")
		rp.db = dbRepo
	}

	if opt.needCache {
		// 初始化 Cache
		cacheRepo, err := cache.New()
		if err != nil {
			logger.Fatal("new cache err", zap.Error(err))
		}
		logger.Info("connect to Cache success")
		rp.cache = cacheRepo
	}

	if opt.needRabbitMQ {
		// 初始化 RabbitMQ
		rabbitmqRepo, err := message_queue.RabbitMQNew()
		if err != nil {
			logger.Fatal("new rabbitmq err", zap.Error(err))
		}
		logger.Info("connect to rabbitmq success")
		rp.rabbitmq = rabbitmqRepo
	}

	return rp
}
