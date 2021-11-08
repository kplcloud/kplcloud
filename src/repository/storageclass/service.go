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

type Call func() error

type Service interface {
	Save(ctx context.Context, data *types.StorageClass, call Call) (err error)
	FirstInsert(ctx context.Context, data *types.StorageClass) (err error)
	Find(ctx context.Context, id int64) (res types.StorageClass, err error)
	FindName(ctx context.Context, clusterId int64, name string) (res types.StorageClass, err error)
	List(ctx context.Context, clusterId int64, page, pageSize int) (res []types.StorageClass, total int, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) List(ctx context.Context, clusterId int64, page, pageSize int) (res []types.StorageClass, total int, err error) {
	err = s.db.Model(&types.StorageClass{}).
		Where("cluster_id = ?", clusterId).
		Count(&total).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&res).Error
	return
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

func (s *service) Save(ctx context.Context, data *types.StorageClass, call Call) (err error) {
	tx := s.db.Model(data).Begin()
	if err = tx.Save(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	if call != nil {
		if err = call(); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func New(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}
