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
)

type Middleware func(service Service) Service

type Call func() error

type Service interface {
	Save(ctx context.Context, pvc *types.PersistentVolumeClaim, call Call) (err error)
}

type service struct {
	db *gorm.DB
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
