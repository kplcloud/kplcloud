/**
 * @Time : 2019-06-26 14:58
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package persistentvolumeclaim

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) Sync(ctx context.Context, ns string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "sync",
			"namespace", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Sync(ctx, ns)
}

func (s *loggingService) Get(ctx context.Context, ns, name string) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "sync",
			"namespace", ns,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, ns, name)
}

func (s *loggingService) List(ctx context.Context, ns string, page, limit int) (resp map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "list",
			"namespace", ns,
			"name", page,
			"limit", limit,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, ns, page, limit)
}

func (s *loggingService) Delete(ctx context.Context, ns, name string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "sync",
			"namespace", ns,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, ns, name)
}

func (s *loggingService) Post(ctx context.Context, ns, name, storage, storageClassName string, accessModes []string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Post",
			"namespace", ns,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, ns, name, storage, storageClassName, accessModes)
}
