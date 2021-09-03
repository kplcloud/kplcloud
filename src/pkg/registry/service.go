package registry

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	// Create 创建空间
	Create(ctx context.Context, name, host, username, password, remark string) (err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
}

func (s *service) Create(ctx context.Context, name, host, username, password, remark string) (err error) {
	return s.repository.Registry(ctx).Save(ctx, &types.Registry{
		Name:     name,
		Host:     host,
		Username: username,
		Password: password,
		Remark:   remark,
	})
}

func New(logger log.Logger, traceId string, store repository.Repository) Service {
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: store,
	}
}
