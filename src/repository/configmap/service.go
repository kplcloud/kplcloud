/**
 * @Time : 8/19/21 11:24 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package configmap

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Call func() error

type Service interface {
	Save(ctx context.Context, configMap *types.ConfigMap, data []types.Data) (err error)
	SaveData(ctx context.Context, configMapId int64, key, value string) (err error)
	FindBy(ctx context.Context, clusterId int64, ns, name string) (res types.ConfigMap, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindBy(ctx context.Context, clusterId int64, ns, name string) (res types.ConfigMap, err error) {
	err = s.db.Model(&types.ConfigMap{}).
		Preload("Data").
		Where("cluster_id = ? AND namespace = ? AND `name` = ?", clusterId, ns, name).
		First(&res).Error

	return
}

func (s *service) SaveData(ctx context.Context, configMapId int64, key, value string) (err error) {
	return s.db.Model(&types.Data{}).Save(&types.Data{
		TargetId: configMapId,
		Key:      key,
		Value:    value,
	}).Error
}

func (s *service) Save(ctx context.Context, configMap *types.ConfigMap, data []types.Data) (err error) {
	tx := s.db.Model(configMap).Begin()
	if err = tx.FirstOrCreate(configMap, "cluster_id = ? AND namespace = ? AND `name` = ?", configMap.ClusterId, configMap.Namespace, configMap.Name).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Save(configMap).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range data {
		if err = tx.FirstOrCreate(&v, "target_id = ? AND style = ? AND `key` = ?", configMap.Id, types.DataStyleConfigMap, v.Key).Error; err != nil {
			tx.Rollback()
			return err
		}
		v.TargetId = configMap.Id
		v.Style = types.DataStyleConfigMap
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
