/**
 * @Time : 7/20/21 5:41 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package syssetting

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Middleware func(Service) Service

type Service interface {
	List(ctx context.Context, key string, page, pageSize int) (res []listResult, total int, err error)
}

type service struct {
	logger  log.Logger
	traceId string
}

func (s *service) List(ctx context.Context, key string, page, pageSize int) (res []listResult, total int, err error) {
	panic("implement me")
}

func New(logger log.Logger, traceId string) Service {
	logger = log.With(logger, "syssetting", "service")
	return &service{
		logger:  logger,
		traceId: traceId,
	}
}
