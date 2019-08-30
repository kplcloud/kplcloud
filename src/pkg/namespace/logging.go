package namespace

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) Get(ctx context.Context, name string) (resp *types.Namespace, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, name)
}

func (s *loggingService) Post(ctx context.Context, name, displayName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "post",
			"name", name,
			"displayName", displayName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, name, displayName)
}

func (s *loggingService) Update(ctx context.Context, name, displayName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "update",
			"name", name,
			"displayName", displayName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Update(ctx, name, displayName)
}

func (s *loggingService) Sync(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "sync",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Sync(ctx)
}

func (s *loggingService) List(ctx context.Context) (res []*types.Namespace, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "List",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx)
}
