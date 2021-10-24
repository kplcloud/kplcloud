package registry

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(Service) Service

// Service 镜像仓库管理模块
type Service interface {
	// Create 创建空间
	Create(ctx context.Context, name, host, username, password, remark string) (err error)
	// List 仓库列表
	List(ctx context.Context, query string, page, pageSize int) (res []result, total int, err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
}

func (s *service) List(ctx context.Context, query string, page, pageSize int) (res []result, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	list, total, err := s.repository.Registry(ctx).List(ctx, query, page, pageSize)
	if err != nil {
		_ = level.Error(logger).Log("repository.Registry", "List", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, result{
			Name:      v.Name,
			Host:      v.Host,
			Username:  v.Username,
			Password:  "",
			Remark:    v.Remark,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return
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
