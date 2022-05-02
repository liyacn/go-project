package service

import (
	"context"
	"project/model/entity"
)

func (s *Service) UserCount(ctx context.Context, begin, end string) (count int64, err error) {
	err = s.orm.WithContext(ctx).Model(&entity.User{}).
		Where("created_at>? AND created_at<?", begin, end).Count(&count).Error
	return
}

func (s *Service) UserAvatarUpdate(ctx context.Context, id int, path string) error {
	return s.orm.WithContext(ctx).Model(&entity.User{}).
		Where("id=?", id).Update("avatar_url", path).Error
}
