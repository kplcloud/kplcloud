/**
 * @Time : 2021/9/6 5:15 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package pvc

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"strings"
)

type Middleware func(service Service) Service

type Call func() error

type Service interface {
	Save(ctx context.Context, pvc *types.PersistentVolumeClaim, call Call) (err error)
	List(ctx context.Context, clusterId int64, storageClassIds []int64, ns, name string, page, pageSize int) (res []types.PersistentVolumeClaim, total int, err error)
	FindByName(ctx context.Context, clusterId int64, ns, name string) (res types.PersistentVolumeClaim, err error)
	Delete(ctx context.Context, pvcId int64, call ...Call) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Delete(ctx context.Context, pvcId int64, call ...Call) (err error) {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&types.PersistentVolumeClaim{}).Delete(&types.PersistentVolumeClaim{}, "id = ?", pvcId).Error; err != nil {
			return err
		}
		if call == nil {
			return nil
		}
		for _, c := range call {
			if err = c(); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *service) FindByName(ctx context.Context, clusterId int64, ns, name string) (res types.PersistentVolumeClaim, err error) {
	err = s.db.Model(res).
		Preload("Cluster").
		Where("cluster_id = ?", clusterId).
		Where("namespace = ?", ns).
		Where("name = ?", name).First(&res).Error
	return
}

func (s *service) List(ctx context.Context, clusterId int64, storageClassIds []int64, ns, name string, page, pageSize int) (res []types.PersistentVolumeClaim, total int, err error) {
	q := s.db.Model(&types.PersistentVolumeClaim{}).
		Preload("StorageClass").
		Where("cluster_id = ?", clusterId)
	if len(storageClassIds) > 0 {
		q = q.Where("storage_class_id IN (?)", storageClassIds)
	}
	if !strings.EqualFold(ns, "") {
		q = q.Where("namespace = ?", ns)
	}
	if !strings.EqualFold(name, "") {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	err = q.Order("created_at DESC").
		Count(&total).
		Offset((page - 1) * pageSize).
		Limit(total).Find(&res).Error
	return
}

func (s *service) Save(ctx context.Context, pvc *types.PersistentVolumeClaim, call Call) (err error) {
	if pvc.Id == 0 {
		var p types.PersistentVolumeClaim
		if err = s.db.Model(pvc).Where("name = ? AND storage_class_id = ? AND namespace = ?", pvc.Name, pvc.StorageClassId, pvc.Namespace).First(&p).Error; err == nil {
			pvc.Id = p.Id
		}
	}
	return s.db.Model(pvc).Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(pvc).Error; err != nil {
			return err
		}
		if call != nil {
			err = call()
		}
		return err
	})
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
