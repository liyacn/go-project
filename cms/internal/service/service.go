package service

import (
	"gorm.io/gorm"
	"project/pkg/db"
	"project/pkg/gredis"
	"project/pkg/process"
)

type Service struct {
	orm   *gorm.DB
	redis *gredis.Client
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
	s := &Service{
		orm:   db.NewMysql(&cfg.Mysql),
		redis: gredis.NewPlusClient(&cfg.Redis),
		//producer: gnsq.NewProducer(cfg.Nsq.Producer),
	}
	process.RegisterCleanup(func() {
		db.Close(s.orm)
		s.redis.Close()
		//s.producer.Stop()
	})
	return s
}
