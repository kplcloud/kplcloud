/**
 * @Time : 2021/12/13 5:55 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package hpa

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) List(ctx context.Context, clusterId int64, namespace string, names []string, page, pageSize int) (res []types.HorizontalPodAutoscaler, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "clusterId", clusterId, "namespace", namespace, "names", names, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, clusterId, namespace, names, page, pageSize)
}

func (s *logging) Save(ctx context.Context, dat *types.HorizontalPodAutoscaler, calls ...Call) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save", "dat", dat,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, dat, calls...)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "hpa", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
