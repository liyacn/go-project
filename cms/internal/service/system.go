package service

import (
	"context"
	"gorm.io/gorm"
	"project/cms/internal/proto"
	"project/model/entity"
)

func (s *Service) SystemActionKeyNames(ctx context.Context, onlyRoute bool) ([]string, error) {
	query := s.orm.WithContext(ctx).Model(&entity.SystemAction{}).Select("key_name")
	if onlyRoute {
		query.Where("key_name LIKE '/%'")
	}
	var result []string
	err := query.Scan(&result).Error
	return result, err
}

func (s *Service) SystemActionList(ctx context.Context) ([]*entity.SystemAction, error) {
	var result []*entity.SystemAction
	err := s.orm.WithContext(ctx).Find(&result).Error
	return result, err
}

func (s *Service) SystemActionUpdate(ctx context.Context, data *entity.SystemAction) error {
	return s.orm.WithContext(ctx).Updates(data).Error
}

func (s *Service) SystemActionSave(ctx context.Context, p *proto.SystemActionSyncData) error {
	return s.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		if len(p.Delete) > 0 {
			if err = tx.Delete(&entity.SystemAction{}, p.Delete).Error; err != nil {
				return
			}
		}
		if len(p.Create) > 0 {
			err = tx.Create(p.Create).Error
		}
		return
	})
}
