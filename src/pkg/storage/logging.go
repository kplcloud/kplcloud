/**
 * @Time : 2019-06-25 19:24
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package storage

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

func (s *loggingService) Get(ctx context.Context, name string) (rs interface{}, err error) {
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

func (s *loggingService) Post(ctx context.Context, name, provisioner, reclaimRolicy, volumeBindingMode string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "post",
			"name", name,
			"provisioner", provisioner,
			"reclaimRolicy", reclaimRolicy,
			"volumeBindingMode", volumeBindingMode,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, name, provisioner, reclaimRolicy, volumeBindingMode)
}

func (s *loggingService) Delete(ctx context.Context, name string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "delete",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, name)
}

func (s *loggingService) List(ctx context.Context, offset, limit int) (res []*types.StorageClass, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "list",
			"offset", offset,
			"limit", limit,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, offset, limit)
}
