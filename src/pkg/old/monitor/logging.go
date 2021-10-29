/**
 * @Time : 2019-07-29 15:18
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package monitor

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

func (s *loggingService) QueryNetwork(ctx context.Context) (data []map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "QueryNetwork",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.QueryNetwork(ctx)
}

func (s *loggingService) Ops(ctx context.Context) (rs interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Ops",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Ops(ctx)
}

func (s *loggingService) Metrics(ctx context.Context) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Metrics",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Metrics(ctx)
}
