// Package article_job
// @Author      : lilinzhen
// @Time        : 2022/5/8 20:40:51
// @Description :
package article_job

import (
	"blueblog/internal/app/blueblog-service/repository/article_repo"
	"blueblog/internal/pkg/configs"
	"blueblog/internal/pkg/db"
	"blueblog/internal/pkg/repo"
	"blueblog/pkg/time_parse"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

type readArticleBody struct {
	Id  string
	Uid string
}

func ReadArticle(logger *zap.Logger, rp repo.Repo) {
	logger.Debug("start job:ReadArticle, queue:" + fmt.Sprintf(configs.MQRoutingKey4Article, "read"))
	for {
		err := rp.GetRabbitMQ().Consume(
			fmt.Sprintf(configs.MQRoutingKey4Article, "read"),
			func(delivery amqp.Delivery) {
				logger.Debug(fmt.Sprintf("Received a message for article read: %s", delivery.Body))
				var rab readArticleBody
				err := json.Unmarshal(delivery.Body, &rab)
				if err != nil {
					logger.Error(fmt.Sprintf("Message Body Error: %s", delivery.Body))
				}
				key := fmt.Sprintf(configs.CacheKey4Article, rab.Uid, rab.Id)
				qb := article_repo.NewQueryBuilder()
				qb.WhereUid(db.EqualPredicate, rab.Uid)
				qb.WhereId(db.EqualPredicate, rab.Id)

				info, err := qb.QueryOne(rp.GetDb().GetDbW())
				if err != nil {
					logger.Error("Query DB Error.", zap.Error(err))
				}
				if info == nil {
					logger.Info(fmt.Sprintf("%s not exist", key))
					_, err = rp.GetCache().HSet(key, []interface{}{"exist", "false"})
					if err != nil {
						logger.Error(fmt.Sprintf("cache %s error.", key), zap.Error(err))
					}
					if rp.GetCache().Expire(key, time.Second*5) {
						logger.Debug(fmt.Sprintf("cache %s will expire 5s", key))
					}
				} else {
					logger.Info(fmt.Sprintf("%s founded", key))
					_, err = rp.GetCache().HSet(key,
						[]interface{}{
							"exist", "true",
							"title", info.Title,
							"content", info.Content,
							"create_time", info.CreateTime.Format(time_parse.CSTLayout),
						})
					if err != nil {
						logger.Error(fmt.Sprintf("cache %s error.", key), zap.Error(err))
					}
					if rp.GetCache().Expire(key, time.Second*5) {
						logger.Debug(fmt.Sprintf("cache %s will expire 5s", key))
					}

				}
			},
		)
		if err != nil {
			logger.Error("read article err.", zap.Error(err))
		}
		time.Sleep(time.Second * 5)
	}
}
