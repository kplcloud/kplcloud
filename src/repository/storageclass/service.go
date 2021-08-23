/**
 * @Time : 2021/8/23 10:22 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package storageclass

import (
	"context"
	"github.com/jinzhu/gorm"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	Save(ctx context.Context, data *types.StorageClass) (err error)
	FirstInsert(ctx context.Context, data *types.StorageClass) (err error)
	Find(ctx context.Context, id int64) (res types.StorageClass, err error)
	FindName(ctx context.Context, clusterId int64, name string) (res types.StorageClass, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindName(ctx context.Context, clusterId int64, name string) (res types.StorageClass, err error) {
	err = s.db.Model(&types.StorageClass{}).
		Where("cluster_id = ?", clusterId).
		Where("name = ?", name).
		First(&res).Error
	return
}

func (s *service) Find(ctx context.Context, id int64) (res types.StorageClass, err error) {
	err = s.db.Model(&types.StorageClass{}).Where("id = ?", id).First(&res).Error
	return
}

func (s *service) FirstInsert(ctx context.Context, data *types.StorageClass) (err error) {
	return s.db.Model(data).FirstOrCreate(data, "cluster_id = ? AND name = ?", data.ClusterId, data.Name).Error
}

func (s *service) Save(ctx context.Context, data *types.StorageClass) (err error) {
	return s.db.Model(data).Save(data).Error
}

func New(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}
