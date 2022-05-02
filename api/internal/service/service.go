package service

import (
	"context"
	"github.com/nsqio/go-nsq"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"project/model/cache"
	"project/pkg/db"
	"project/pkg/gnsq"
	"project/pkg/gredis"
	"project/pkg/process"
)

type Service struct {
	orm      *gorm.DB
	redis    *gredis.Client
	producer *nsq.Producer
	single   singleflight.Group
}

type Config struct {
	Mysql db.Config
	Redis gredis.Config
	Nsq   struct {
		Producer string
	}
}

func New(cfg *Config) *Service {
	s := &Service{
		orm:      db.NewMysql(&cfg.Mysql),
		redis:    gredis.NewPlusClient(&cfg.Redis),
		producer: gnsq.NewProducer(cfg.Nsq.Producer),
	}
	process.RegisterCleanup(func() {
		db.Close(s.orm)
		s.redis.Close()
		s.producer.Stop()
	})
	return s
}

func (s *Service) WechatToken(ctx context.Context) (string, error) {
	val, err, _ := s.single.Do("WechatToken", func() (any, error) {
		return s.redis.Get(ctx, cache.WechatTokenKey).Result()
	})
	return val.(string), err
}
