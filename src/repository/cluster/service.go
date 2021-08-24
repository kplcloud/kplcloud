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
	"github.com/pkg/errors"
)

type Call func(tx *gorm.DB) error

type Middleware func(Service) Service

type Service interface {
	FindAll(ctx context.Context, status int) (res []types.Cluster, err error)
	FindByName(ctx context.Context, name string) (res types.Cluster, err error)
	Save(ctx context.Context, data *types.Cluster, calls ...Call) (err error)
	Delete(ctx context.Context, id int64, unscoped bool) (err error)
	SaveRole(ctx context.Context, clusterRole *types.ClusterRole, roles []types.PolicyRule) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) SaveRole(ctx context.Context, clusterRole *types.ClusterRole, roles []types.PolicyRule) (err error) {
	tx := s.db.Model(clusterRole).Begin()
	if err = tx.FirstOrCreate(clusterRole, "cluster_id = ? AND `name` = ?", clusterRole.ClusterId, clusterRole.Name).Error; err != nil {
		tx.Rollback()
		return err
	}
	// TODO 直接删除所有
	for _, v := range roles {
		if err = tx.FirstOrCreate(&v, "target_id = ? AND kind = ?", clusterRole.Id, types.KindClusterRole).Error; err != nil {
			tx.Rollback()
			return err
		}
		v.TargetId = clusterRole.Id
		v.Kind = types.KindClusterRole.String()
		if err = tx.Model(&v).Save(&v).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
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
		if err = call(tx); err != nil {
			tx.Rollback()
			return errors.Wrap(err, "call")
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
