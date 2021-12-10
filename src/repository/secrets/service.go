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
	"strings"
)

type Call func() error

type Middleware func(Service) Service

type Service interface {
	FindBy(ctx context.Context, clusterId int64, ns, name string) (res types.Secret, err error)
	Save(ctx context.Context, secret *types.Secret, data []types.Data) (err error)
	Delete(ctx context.Context, clusterId int64, ns, name string) (err error)
	FindByName(ctx context.Context, name string) (res []types.Secret, err error)
	FindNsByNames(ctx context.Context, clusterId int64, ns string, names []string) (res []types.Secret, err error)
	List(ctx context.Context, clusterId int64, ns, name string, page, pageSize int) (res []types.Secret, total int, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) List(ctx context.Context, clusterId int64, ns, name string, page, pageSize int) (res []types.Secret, total int, err error) {
	query := s.db.Model(&types.Secret{}).
		Where("cluster_id = ?", clusterId).Where("namespace = ?", ns)
	if !strings.EqualFold(name, "") {
		query = query.Where("name = ?", name)
	}
	err = query.Order("updated_at DESC").
		Count(&total).
		Offset((page - 1) * pageSize).Limit(total).Find(&res).Error

	return
}

func (s *service) FindNsByNames(ctx context.Context, clusterId int64, ns string, names []string) (res []types.Secret, err error) {
	err = s.db.Model(&types.Secret{}).
		Preload("Data").
		Where("cluster_id = ?", clusterId).
		Where("namespace = ?", ns).
		Where("name IN (?)", names).Find(&res).Error
	return
}

func (s *service) FindByName(ctx context.Context, name string) (res []types.Secret, err error) {
	err = s.db.Model(&types.Secret{}).
		Preload("Data").
		Where("name = ?", name).Find(&res).Error
	return
}

func (s *service) Delete(ctx context.Context, clusterId int64, ns, name string) (err error) {
	tx := s.db.Begin()
	var secret types.Secret
	if err = tx.Model(&types.Secret{}).
		Preload("Data", func(db *gorm.DB) *gorm.DB {
			return db.Where("style = ?", types.DataStyleSecret)
		}).
		Where("cluster_id = ? AND namespace = ? AND name = ?", clusterId, ns, name).First(&secret).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, v := range secret.Data {
		if tx.Model(v).Unscoped().Delete(v, "id = ?", v.Id).Error != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Model(&secret).Unscoped().
		Where("id = ?", secret.Id).
		Delete(&secret).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
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
