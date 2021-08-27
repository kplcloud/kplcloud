/**
 * @Time : 2021/8/25 9:55 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package k8stpl

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Service interface {
	FindByKind(ctx context.Context, kind types.Kind) (res types.K8sTemplate, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindByKind(ctx context.Context, kind types.Kind) (res types.K8sTemplate, err error) {
	err = s.db.Model(&res).Where("kind = ?", kind).First(&res).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
