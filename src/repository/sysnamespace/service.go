/**
 * @Time : 3/5/21 2:41 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package sysnamespace

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(next Service) Service

type Service interface {
	// 保存过更新
	Save(ctx context.Context, data *types.SysNamespace) (err error)
	FindByIds(ctx context.Context, ids []int64) (res []types.SysNamespace, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindByIds(ctx context.Context, ids []int64) (res []types.SysNamespace, err error) {
	err = s.db.Model(&types.SysNamespace{}).Where("id IN (?)", ids).Find(&res).Error
	return
}

func (s *service) Save(ctx context.Context, data *types.SysNamespace) (err error) {
	return s.db.Model(data).Save(data).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
