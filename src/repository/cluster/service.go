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
)

type Service interface {
	FindAll(ctx context.Context, status int) (res []types.Cluster, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindAll(ctx context.Context, status int) (res []types.Cluster, err error) {
	err = s.db.Model(&types.Cluster{}).Where("status = ?", 1).Find(&res).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
