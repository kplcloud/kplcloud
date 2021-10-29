/**
 * @Time : 2019-07-29 11:30
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package market

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

func (s *loggingService) Post(ctx context.Context, name, language, version, detail, desc, dockerfile, fullPath string, status int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Post",
			"took", time.Since(begin),
			"name", name,
			"language", language,
			"version", version,
			"fullPath", fullPath,
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, name, language, version, detail, desc, dockerfile, fullPath, status)
}

func (s *loggingService) Detail(ctx context.Context, id int64) (res *types.Dockerfile, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Detail",
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, id)
}

func (s *loggingService) List(ctx context.Context, page, limit int, language []string, status int, name string) (res []*types.Dockerfile, count int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "List",
			"took", time.Since(begin),
			"page", page,
			"limit", limit,
			"status", status,
			"name", name,
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, page, limit, language, status, name)
}

func (s *loggingService) Put(ctx context.Context, id int64, name, language, version, detail, desc, dockerfile, fullPath string, status int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Put",
			"took", time.Since(begin),
			"id", id,
			"name", name,
			"language", language,
			"version", version,
			"fullPath", fullPath,
			"err", err,
		)
	}(time.Now())
	return s.Service.Put(ctx, id, name, language, version, detail, desc, dockerfile, fullPath, status)
}

func (s *loggingService) Delete(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Detail",
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, id)
}
