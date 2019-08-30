package proclaim

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

func (s loggingService) Get(ctx context.Context, id int) (resp *types.Notices, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, id)
}

func (s loggingService) Post(ctx context.Context, req proclaimRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "post",
			"request", req,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, req)
}

func (s loggingService) List(ctx context.Context, name string, page, limit int) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "delete",
			"page", page,
			"limit", limit,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, name, page, limit)
}
