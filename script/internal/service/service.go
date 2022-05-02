package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"project/pkg/db"
	"project/pkg/gredis"
)

type Config struct {
	Mysql db.Config
	//Pgsql db.Config
	Redis gredis.Config
	Nsq   struct {
		Producer string
		Consumer string
	}
	//Rabbitmq rabbitmq.Config
}

type Service struct {
	orm   *gorm.DB
	redis *redis.Client
	//producer *nsq.Producer
	//rabbit   *amqp.Channel
}

type Option func(*Service)

func Orm(cfg *db.Config) Option {
	return func(s *Service) {
		if s.orm == nil {
			s.orm = db.NewMysql(cfg)
			//s.orm = db.NewPostgres(cfg)
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
//
//func Rabbit(cfg *rabbitmq.Config) Option {
//	return func(s *Service) {
//		if s.rabbit == nil {
//			s.rabbit = rabbitmq.NewChannel(cfg)
//		}
//	}
//}

func New(options ...Option) *Service {
	s := &Service{}
	for _, opt := range options {
		opt(s)
	}
	return s
}
