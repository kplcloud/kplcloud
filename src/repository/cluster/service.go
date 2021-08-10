/**
 * @Time : 8/9/21 6:25 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package cluster

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Call func() error

type Service interface {
	FindAll(ctx context.Context, status int) (res []types.Cluster, err error)
	FindByName(ctx context.Context, name string) (res types.Cluster, err error)
	Save(ctx context.Context, data *types.Cluster, calls ...Call) (err error)
	Delete(ctx context.Context, id int64, unscoped bool) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Delete(ctx context.Context, id int64, unscoped bool) (err error) {
	query := s.db.Model(&types.Cluster{}).Where("id = ?", id)
	if unscoped {
		query = query.Unscoped()
	}
	err = query.Delete(&types.Cluster{}).Error
	return
}

func (s *service) FindByName(ctx context.Context, name string) (res types.Cluster, err error) {
	err = s.db.Model(&types.Cluster{}).Where("name = ?", name).First(&res).Error
	return
}

func (s *service) Save(ctx context.Context, data *types.Cluster, calls ...Call) (err error) {
	tx := s.db.Begin()
	if err = tx.Model(data).Save(data).Error; err != nil {
		return tx.Rollback().Error
	}

	for _, call := range calls {
		if err = call(); err != nil {
			return tx.Rollback().Error
		}
	}

	return tx.Commit().Error
}

func (s *service) FindAll(ctx context.Context, status int) (res []types.Cluster, err error) {
	err = s.db.Model(&types.Cluster{}).Where("status = ?", 1).Find(&res).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
