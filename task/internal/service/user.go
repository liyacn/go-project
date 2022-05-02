package service

import (
	"context"
	"project/model/entity"
)

func (s *Service) UserAvatarUpdate(ctx context.Context, id int, path string) error {
	return s.orm.WithContext(ctx).Model(&entity.User{}).
		Where("id=?", id).Update("avatar_url", path).Error
}
