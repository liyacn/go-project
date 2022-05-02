package service

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"project/cms/internal/proto"
	"project/model/cache"
	"project/model/entity"
	"project/pkg/logger"
	"time"
)

func (s *Service) AdminTokenSet(ctx context.Context, data *cache.AdminToken) (string, error) {
	ssoKey := cache.AdminSsoKey(data.ID)
	oldToken, err := s.redis.Get(ctx, ssoKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	b, _ := json.Marshal(data)
	newToken := cache.GenerateToken()
	tx := s.redis.TxPipeline()
	if oldToken != "" {
		tx.Unlink(ctx, cache.AdminTokenKey(oldToken))
	}
	tx.Set(ctx, cache.AdminTokenKey(newToken), b, time.Hour)
	tx.Set(ctx, ssoKey, newToken, time.Hour)
	_, err = tx.Exec(ctx)
	return newToken, err
}

func (s *Service) AdminTokenGet(ctx context.Context, token string) (*cache.AdminToken, error) {
	key := cache.AdminTokenKey(token)
	b, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	var result cache.AdminToken
	if len(b) == 0 {
		return &result, nil
	}
	err = json.Unmarshal(b, &result)
	tx := s.redis.TxPipeline()
	tx.Expire(ctx, key, time.Hour)
	tx.Expire(ctx, cache.AdminSsoKey(result.ID), time.Hour)
	if _, err := tx.Exec(ctx); err != nil {
		logger.FromContext(ctx).Error("redis.TxPipeline.Exec error", err)
	}
	return &result, err
}

func (s *Service) AdminUserLogout(ctx context.Context, id int) error {
	token, err := s.redis.Get(ctx, cache.AdminSsoKey(id)).Result() //redis>=v6.2改用GetDel方法更佳
	if err != nil && err != redis.Nil {
		return err
	}
	if token != "" {
		return s.redis.Unlink(ctx, cache.AdminTokenKey(token)).Err()
	}
	return nil
}

func (s *Service) AdminUserFindByName(ctx context.Context, name string) (*entity.AdminUser, error) {
	var result entity.AdminUser
	err := s.orm.WithContext(ctx).Where("username=?", name).Take(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &result, nil
}

func (s *Service) AdminUserFindByID(ctx context.Context, id int) (*entity.AdminUser, error) {
	var result entity.AdminUser
	err := s.orm.WithContext(ctx).Where("id=?", id).Take(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &result, nil
}

func (s *Service) AdminUserUpdate(ctx context.Context, data *entity.AdminUser) error {
	if err := s.orm.WithContext(ctx).Updates(data).Error; err != nil {
		return err
	}
	if data.Status == entity.StatusDisabled || data.RoleID != 0 {
		return s.AdminUserLogout(ctx, data.ID) //禁用或变更角色强制退出登录
	}
	return nil
}

func (s *Service) AdminRoleList(ctx context.Context,
	p *proto.ListArgs) (total int64, list []*entity.AdminRole, err error) {
	query := s.orm.WithContext(ctx).Model(&entity.AdminRole{})
	err = query.Count(&total).Error
	offset := p.Size * (p.Page - 1)
	if err != nil || total == 0 || offset >= int(total) {
		return
	}
	err = query.Order("id DESC").Limit(p.Size).Offset(offset).Find(&list).Error
	return
}

func (s *Service) AdminRoles(ctx context.Context) ([]*entity.AdminRole, error) {
	var result []*entity.AdminRole
	err := s.orm.WithContext(ctx).Select("id", "name").
		Order("id").Find(&result).Error
	return result, err
}

func (s *Service) AdminRoleFindByID(ctx context.Context, id int) (*entity.AdminRole, error) {
	var result entity.AdminRole
	err := s.orm.WithContext(ctx).Where("id=?", id).Take(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &result, nil
}

func (s *Service) AdminRoleActions(ctx context.Context, id int) ([]string, error) {
	var role entity.AdminRole
	err := s.redis.FetchJSON(ctx, cache.AdminRoleActionsKey(id), &role.Actions, func() error {
		return s.orm.WithContext(ctx).Select("actions").Where("id=?", id).Take(&role).Error
	}, time.Hour)
	return role.Actions, err
}

func (s *Service) AdminRoleSave(ctx context.Context, data *entity.AdminRole) error {
	if data.ID == 0 {
		return s.orm.WithContext(ctx).Create(data).Error
	}
	if err := s.orm.WithContext(ctx).Updates(data).Error; err != nil {
		return err
	}
	return s.redis.Unlink(ctx, cache.AdminRoleActionsKey(data.ID)).Err()
}

func (s *Service) AdminUserList(ctx context.Context,
	p *proto.AdminUserListArgs) (total int64, list []*entity.AdminUser, err error) {
	query := s.orm.WithContext(ctx).Model(&entity.AdminUser{})
	if p.RoleID > 0 {
		query.Where("role_id=?", p.RoleID)
	}
	if p.Username != "" {
		query.Where("username LIKE ?", p.Username+"%")
	}
	if p.Status != 0 {
		query.Where("status=?", p.Status)
	}
	err = query.Count(&total).Error
	offset := p.Size * (p.Page - 1)
	if err != nil || total == 0 || offset >= int(total) {
		return
	}
	err = query.Order("id DESC").Limit(p.Size).Offset(offset).Find(&list).Error
	return
}

func (s *Service) AdminUserCreate(ctx context.Context, data *entity.AdminUser) (bool, error) {
	opt := s.orm.WithContext(ctx).FirstOrCreate(data, "username=?", data.Username)
	return opt.RowsAffected > 0, opt.Error
}
