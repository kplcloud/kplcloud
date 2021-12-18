/**
 * @Time : 2021/12/13 5:41 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package hpa

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(service Service) Service

type Call func() error

type Service interface {
	Save(ctx context.Context, dat *types.HorizontalPodAutoscaler, calls ...Call) (err error)
	List(ctx context.Context, clusterId int64, namespace string, names []string, page, pageSize int) (res []types.HorizontalPodAutoscaler, total int, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) List(ctx context.Context, clusterId int64, namespace string, names []string, page, pageSize int) (res []types.HorizontalPodAutoscaler, total int, err error) {
	q := s.db.Model(&types.HorizontalPodAutoscaler{}).
		Where("cluster_id = ?", clusterId).
		Where("namespace = ?", namespace)
	if len(names) > 0 {
		q = q.Where("name IN (?)", names)
	}
	err = q.Count(&total).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&res).Error
	return
}

func (s *service) Save(ctx context.Context, dat *types.HorizontalPodAutoscaler, calls ...Call) (err error) {
	if dat.Id == 0 {
		var hpaData types.HorizontalPodAutoscaler
		err = s.db.Model(dat).Where("cluster_id = ? AND namespace = ? AND name = ?", dat.ClusterId, dat.Namespace, dat.Name).First(&hpaData).Error
		if !gorm.IsRecordNotFoundError(err) {
			dat.Id = hpaData.Id
		}
	}
	return s.db.Model(dat).Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(dat).Error; err != nil {
			return err
		}
		if calls != nil {
			for _, call := range calls {
				if err = call(); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
