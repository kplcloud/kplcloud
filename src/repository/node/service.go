/**
 * @Time: 2021/6/29 下午10:07
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package node

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	Save(ctx context.Context, dat *types.Node) (err error)
	List(ctx context.Context, page, pageSize int) (res []types.Node, total int, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) List(ctx context.Context, page, pageSize int) (res []types.Node, total int, err error) {
	panic("implement me")
}

func (s *service) Save(ctx context.Context, dat *types.Node) (err error) {
	return s.db.Model(dat).Save(dat).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
