/**
 * @Time : 2021/11/25 4:27 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package application

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(service Service) Service

type Call func() error

type Service interface {
	Save(ctx context.Context, app *types.Application, call Call) (err error)
	List(ctx context.Context, clusterId int64, namespace string, names []string, page, pageSize int) (res []types.Application, total int, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) List(ctx context.Context, clusterId int64, namespace string, names []string, page, pageSize int) (res []types.Application, total int, err error) {
	q := s.db.Model(&types.Application{}).Where("cluster_id = ? AND namespace = ?", clusterId, namespace)
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

func (s *service) Save(ctx context.Context, app *types.Application, call Call) (err error) {
	return s.db.Model(app).Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(app).Error; err != nil {
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
