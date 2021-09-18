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
	"github.com/kplcloud/kplcloud/src/repository"
)

type Middleware func(Service) Service

type Service interface {
	List(ctx context.Context, key string, page, pageSize int) (res []listResult, total int, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
}

func (s *service) List(ctx context.Context, key string, page, pageSize int) (res []listResult, total int, err error) {
	list, total, err := s.repository.SysSetting().List(ctx, key, page, pageSize)
	if err != nil {
		return
	}
	for _, v := range list {
		res = append(res, listResult{
			Section:   v.Section,
			Key:       v.Key,
			Value:     v.Value,
			Id:        v.Id,
			Remark:    v.Description,
			Enable:    v.Enable,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return
}

func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "syssetting", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}
