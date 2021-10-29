/**
 * @Time : 2019-07-08 10:41
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package discovery

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/middleware"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) Delete(ctx context.Context, svcName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Delete",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"svcName", svcName,
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, svcName)
}

func (s *loggingService) Detail(ctx context.Context, svcName string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Detail",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"svcName", svcName,
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, svcName)
}

func (s *loggingService) List(ctx context.Context, page, limit int) (res []*serviceList, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "List",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"page", page,
			"limit", limit,
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, page, limit)
}

func (s *loggingService) Create(ctx context.Context, req createRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Create",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"name", req.Name,
			"err", err,
		)
	}(time.Now())
	return s.Service.Create(ctx, req)
}

func (s *loggingService) Update(ctx context.Context, req createRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Update",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"name", req.Name,
			"err", err,
		)
	}(time.Now())
	return s.Service.Update(ctx, req)
}
