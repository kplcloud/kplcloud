/**
 * @Time : 8/11/21 11:25 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package nodes

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	Save(ctx context.Context, data *types.Nodes) (err error)
	FindByName(ctx context.Context, clusterId int64, name string) (res types.Nodes, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindByName(ctx context.Context, clusterId int64, name string) (res types.Nodes, err error) {
	err = s.db.Model(&types.Nodes{}).Where("cluster_id = ?", clusterId).Where("name = ?", name).First(&res).Error
	return
}

func (s *service) Save(ctx context.Context, data *types.Nodes) (err error) {
	return s.db.Model(data).Save(data).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
