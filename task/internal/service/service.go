package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"project/pkg/db"
	"project/pkg/gredis"
	"project/pkg/process"
)

type Config struct {
	Mysql db.Config
	Redis gredis.Config
	Nsq   struct {
		Producer string
		Consumer string
	}
}

type Service struct {
	orm   *gorm.DB
	redis *redis.Client
	//producer *nsq.Producer
}

type Option func(*Service)

func Orm(cfg *db.Config) Option {
	return func(s *Service) {
		if s.orm == nil {
			s.orm = db.NewMysql(cfg)
		}
	}
}

func Redis(cfg *gredis.Config) Option {
	return func(s *Service) {
		if s.redis == nil {
			s.redis = gredis.NewClient(cfg)
		}
	}
}

//func Producer(addr string) Option {
//	return func(s *Service) {
//		if s.producer == nil {
//			s.producer = gnsq.NewProducer(addr)
//		}
//	}
//}

func New(options ...Option) *Service {
	s := &Service{}
	for _, opt := range options {
		opt(s)
	}
	process.RegisterCleanup(func() {
		if s.orm != nil {
			db.Close(s.orm)
		}
		if s.redis != nil {
			s.redis.Close()
		}
		//if s.producer != nil {
		//	s.producer.Stop()
		//}
	})
	return s
}
