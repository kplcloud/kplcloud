/**
 * @Time : 2021/9/3 12:07 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package registry

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"strings"
)

type Middleware func(Service) Service

type Call func() error

type Service interface {
	Save(ctx context.Context, reg *types.Registry) (err error)
	FindByName(ctx context.Context, name string) (res types.Registry, err error)
	FindByNames(ctx context.Context, names []string) (res []types.Registry, err error)
	List(ctx context.Context, query string, page, pageSize int) (res []types.Registry, total int, err error)
	SaveCall(ctx context.Context, reg *types.Registry, call Call) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) SaveCall(ctx context.Context, reg *types.Registry, call Call) (err error) {
	tx := s.db.Model(reg).Begin()
	if err = tx.Save(reg).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = call(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *service) List(ctx context.Context, query string, page, pageSize int) (res []types.Registry, total int, err error) {
	q := s.db.Model(&types.Registry{})
	if !strings.EqualFold(query, "") {
		q = q.Where("name LIKE '%?%'", query)
	}
	err = q.Count(&total).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&res).Error
	return
}

func (s *service) FindByNames(ctx context.Context, names []string) (res []types.Registry, err error) {
	err = s.db.Model(&types.Registry{}).Where("name IN (?)", names).Find(&res).Error
	return
}

func (s *service) Save(ctx context.Context, reg *types.Registry) (err error) {
	return s.db.Model(reg).Save(reg).Error
}

func (s *service) FindByName(ctx context.Context, name string) (res types.Registry, err error) {
	err = s.db.Model(&types.Registry{}).Where("name = ?", name).First(&res).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}
