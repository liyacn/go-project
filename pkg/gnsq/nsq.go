package gnsq

import (
	"github.com/nsqio/go-nsq"
	"log"
	"project/pkg/core"
)

func NewProducer(addr string) *nsq.Producer {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	producer.SetLogger(nil, nsq.LogLevelDebug)
	err = producer.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return producer
}

type HandlerFunc func(*core.Context, *nsq.Message) error

func NewConsumer(addr, topic, channel string, concurrency int, handler HandlerFunc) *nsq.Consumer {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = concurrency
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Fatal(err)
	}
	consumer.SetLogger(nil, nsq.LogLevelDebug)
	var fn nsq.HandlerFunc = func(msg *nsq.Message) (err error) {
		ctx := core.NewContext(string(msg.ID[:]), "", "", "")
		defer core.RecoverE(ctx, &err)
		return handler(ctx, msg)
	}
	consumer.AddConcurrentHandlers(fn, concurrency)
	err = consumer.ConnectToNSQLookupd(addr)
	if err != nil {
		log.Fatal(err)
	}
	return consumer
}
