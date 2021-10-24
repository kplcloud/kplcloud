package account

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

func (s loggingService) Detail(ctx context.Context) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx)
}

func (s loggingService) GetReceive(ctx context.Context) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetReceive(ctx)
}

func (s loggingService) UpdateReceive(ctx context.Context, req accountReceiveRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UpdateReceive(ctx, req)
}

func (s loggingService) UnWechatBind(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UnWechatBind(ctx)
}

func (s loggingService) GetProject(ctx context.Context) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetProject(ctx)
}
