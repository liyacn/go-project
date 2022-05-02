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
	//rabbit   *amqp.Channel
}

type Config struct {
	Mysql db.Config
	//Pgsql db.Config
	Redis gredis.Config
	//Nsq   struct {
	//	Producer string
	//}
	//Rabbitmq rabbitmq.Config
}

func New(cfg *Config) *Service {
	return &Service{
		orm: db.NewMysql(&cfg.Mysql),
		//orm:   db.NewPostgres(&cfg.Pgsql),
		redis: gredis.NewClient(&cfg.Redis),
		//producer: gnsq.NewProducer(cfg.Nsq.Producer),
		//rabbit:   rabbitmq.NewChannel(&cfg.Rabbitmq),
	}
}
