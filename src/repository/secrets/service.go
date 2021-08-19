/**
 * @Time : 8/19/21 2:04 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package secrets

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Call func() error

type Middleware func(Service) Service

type Service interface {
	FindBy(ctx context.Context, clusterId int64, ns, name string) (res types.Secret, err error)
	Save(ctx context.Context, secret *types.Secret, data []types.Data) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindBy(ctx context.Context, clusterId int64, ns, name string) (res types.Secret, err error) {
	err = s.db.Model(&types.Secret{}).
		Preload("Data").
		Where("cluster_id = ? AND namespace = ? AND `name` = ?", clusterId, ns, name).
		First(&res).Error

	return
}

func (s *service) Save(ctx context.Context, secret *types.Secret, data []types.Data) (err error) {
	tx := s.db.Model(secret).Begin()
	if err = tx.FirstOrCreate(secret, "cluster_id = ? AND namespace = ? AND `name` = ?", secret.ClusterId, secret.Namespace, secret.Name).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range data {
		if err = tx.FirstOrCreate(&v, "target_id = ? AND style = ? AND `key` = ?", secret.Id, types.DataStyleSecret, v.Key).Error; err != nil {
			tx.Rollback()
			return err
		}
		v.TargetId = secret.Id
		v.Style = types.DataStyleSecret
		if err = tx.Model(&v).Save(&v).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
