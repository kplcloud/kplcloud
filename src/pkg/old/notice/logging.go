/**
 * @Time : 2019-07-02 10:41
 * @Author : soupzhb@gmail.com
 * @File : logging.go
 * @Software: GoLand
 */

package notice

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

func (s loggingService) List(ctx context.Context, param map[string]string, page, limit int) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"page", page,
			"limit", limit,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, param, page, limit)
}

func (s loggingService) Tips(ctx context.Context) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Tips(ctx)
}

func (s loggingService) CountRead(ctx context.Context, param map[string]string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CountRead(ctx, param)
}

func (s loggingService) ClearAll(ctx context.Context, noticeType int) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.ClearAll(ctx, noticeType)
}

func (s loggingService) Detail(ctx context.Context, noticeId int) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, noticeId)
}
