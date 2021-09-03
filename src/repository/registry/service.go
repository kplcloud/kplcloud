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
)

type Middleware func(Service) Service

type Service interface {
	Save(ctx context.Context, reg *types.Registry) (err error)
	FindByName(ctx context.Context, name string) (res types.Registry, err error)
	FindByNames(ctx context.Context, names []string) (res []types.Registry, err error)
}

type service struct {
	db *gorm.DB
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
