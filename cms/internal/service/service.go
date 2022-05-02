package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"project/pkg/db"
	"project/pkg/gredis"
)

type Service struct {
	orm   *gorm.DB
	redis *redis.Client
	//producer *nsq.Producer
}

type Config struct {
	Mysql db.Config
	Redis gredis.Config
	Nsq   struct {
		Producer string
	}
}

func New(cfg *Config) *Service {
	return &Service{
		orm:   db.NewMysql(&cfg.Mysql),
		redis: gredis.NewClient(&cfg.Redis),
		//producer: gnsq.NewProducer(cfg.Nsq.Producer),
	}
}

func (s *Service) Stop() {
	db.Close(s.orm)
	s.redis.Close()
	//s.producer.Stop()
}
