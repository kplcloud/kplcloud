/**
 * @Time : 2019-07-22 13:51
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package event

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Service interface {
	// 获取所有events
	All(ctx context.Context) (res []*types.Event, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, store repository.Repository) Service {
	return &service{logger, store}
}

func (c *service) All(ctx context.Context) (res []*types.Event, err error) {
	return c.repository.Event().FindAllEvents()
}
