// Package message_queue
// @Author      : lilinzhen
// @Time        : 2022/4/11 16:07:01
// @Description :
package message_queue

import (
	"encoding/json"
	"time"

	"blueblog/internal/pkg/configs"
	"blueblog/pkg/errors"
	"blueblog/pkg/time_parse"
	"blueblog/pkg/trace"

	"github.com/streadway/amqp"
)

type RabbitMQOption func(*option)

type RabbitMQTrace = trace.T

type option struct {
	Trace    *trace.Trace
	RabbitMQ *trace.RabbitMQ
}

func newOption() *option {
	return &option{}
}

type RabbitMQRepo interface {
	i()
	Close() error
	IsClosed() bool

	Publish(exchange, routingKey string, data interface{}, options ...RabbitMQOption) (err error)
	Consume(queue string, handleFunc func(delivery amqp.Delivery), options ...RabbitMQOption) (err error)
}

type rabbitmqRepo struct {
	client *amqp.Connection
}

func (r rabbitmqRepo) i() {}

func (r rabbitmqRepo) Close() error {
	return r.client.Close()
}

func (r rabbitmqRepo) IsClosed() bool {
	return r.client.IsClosed()
}

func (r rabbitmqRepo) Publish(exchange, routingKey string, data interface{}, options ...RabbitMQOption) (err error) {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.RabbitMQ.Timestamp = time_parse.CSTLayoutString()
			opt.RabbitMQ.Exchange = exchange
			opt.RabbitMQ.RoutingKey = routingKey
			opt.RabbitMQ.Data = data
			opt.RabbitMQ.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRabbitMQ(opt.RabbitMQ)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marshal data error")
	}

	ch, err := r.client.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to open a channel")
	}
	defer ch.Close()

	err = ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        b,
	})

	if err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	return err
}

func (r rabbitmqRepo) Consume(queue string, handleFunc func(delivery amqp.Delivery), options ...RabbitMQOption) (err error) {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.RabbitMQ.Timestamp = time_parse.CSTLayoutString()
			opt.RabbitMQ.Queue = queue
			opt.RabbitMQ.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRabbitMQ(opt.RabbitMQ)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	ch, err := r.client.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to open a channel")
	}
	defer func() {
		_ = ch.Close()
	}()

	//????????????
	msgs, err := ch.Consume(
		queue, // ????????????
		"",    // ???????????????
		true,  // ???????????????,????????????????????????????????????????????????
		false, // ???????????????,?????????????????????????????????????????????.??????????????????????????????????????????????????????.??????????????????
		false, //??????
		false, // ?????????????????????,??????????????????????????????,??????????????????????????????,?????????????????????????????????
		nil,   // ????????????
	)
	if err != nil {
		return errors.Wrap(err, "Failed to register a consumer")
	}

	closeChan := make(chan *amqp.Error, 1)
	notifyClose := ch.NotifyClose(closeChan)
	for{
		select{
		case d := <-msgs:
			handleFunc(d)
		case e:=<-notifyClose:
			return e
		}
	}
}

func RabbitMQNew() (RabbitMQRepo, error) {
	client, err := rabbitmqConnect()
	if err != nil {
		return nil, err
	}

	return &rabbitmqRepo{
		client: client,
	}, nil
}

func rabbitmqConnect() (*amqp.Connection, error) {
	cfg := configs.Get().RabbitMQ
	conn, err := amqp.Dial(cfg.Url)
	if err != nil {
		return &amqp.Connection{}, errors.Wrap(err, "parse rabbitmq url err")
	}
	return conn, nil
}

// WithRabbitMQTrace ??????trace??????
func WithRabbitMQTrace(t RabbitMQTrace) RabbitMQOption {
	return func(opt *option) {
		if t != nil {
			opt.Trace = t.(*trace.Trace)
			opt.RabbitMQ = new(trace.RabbitMQ)
		}
	}
}
