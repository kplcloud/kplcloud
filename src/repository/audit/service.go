/**
 * @Time : 2021/9/8 2:39 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package audit

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	Save(ctx context.Context, audit *types.Audit) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Save(ctx context.Context, audit *types.Audit) (err error) {
	return s.db.Model(audit).Save(audit).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
