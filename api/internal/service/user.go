package service

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"project/model/cache"
	"project/model/entity"
	"project/model/queue"
	"project/pkg/logger"
	"time"
)

func (s *Service) UserSave(ctx context.Context, data *entity.User) (int, error) {
	err := s.orm.WithContext(ctx).FirstOrCreate(data, "openid=?", data.Openid).Error
	if err != nil {
		b, _ := json.Marshal(data)
		s.redis.Set(ctx, cache.UserInfoKey(data.ID), b, time.Hour)
	}
	return data.ID, err
}

func (s *Service) UserTokenSet(ctx context.Context, data *cache.UserToken) (string, error) {
	b, _ := json.Marshal(data)
	token := cache.GenerateToken()
	err := s.redis.Set(ctx, cache.UserTokenKey(token), b, time.Hour).Err()
	return token, err
}

func (s *Service) UserTokenGet(ctx context.Context, token string) (*cache.UserToken, error) {
	key := cache.UserTokenKey(token)
	b, err := s.redis.Get(ctx, key).Bytes()
	//b, err := s.redis.GetEx(ctx, key, time.Hour).Bytes() // redis>=v6.2使用原子性命令GETEX合并GET和EXPIRE操作
	if err != nil && err != redis.Nil {
		return nil, err
	}
	var result cache.UserToken
	if len(b) == 0 {
		return &result, nil
	}
	if err = s.redis.Expire(ctx, key, time.Hour).Err(); err != nil {
		logger.FromContext(ctx).Error("redis.Expire", key, err)
	}
	err = json.Unmarshal(b, &result)
	return &result, err
}

func (s *Service) UserFindByID(ctx context.Context, id int) (*entity.User, error) {
	var result entity.User
	err := s.redis.FetchJSON(ctx, cache.UserInfoKey(id), &result, func() error {
		return s.orm.WithContext(ctx).Where("id=?", id).Take(&result).Error
	}, time.Hour)
	return &result, err
}

func (s *Service) UserUpdate(ctx context.Context, data *entity.User) error {
	if err := s.orm.WithContext(ctx).Updates(data).Error; err != nil {
		return err
	}
	return s.redis.Del(ctx, cache.UserInfoKey(data.ID)).Err()
}

func (s *Service) AvatarToCdnAsync(data *queue.AvatarToCdn) error {
	b, _ := json.Marshal(data)
	return s.producer.Publish(queue.AvatarToCdnTP, b)
	//return s.rabbit.Publish(queue.DefaultEX, queue.AvatarToCdnQN, false, false, amqp.Publishing{
	//	MessageId: random.UUID(),
	//	Timestamp: time.Now(),
	//	Body:      b,
	//})
}
