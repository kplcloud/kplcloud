/**
 * @Time : 8/11/21 4:21 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package deployment

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) PutImage(ctx context.Context, clusterId int64, ns, name, image string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "PutImage",
			"clusterId", clusterId,
			"ns", ns,
			"name", name,
			"image", image,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.PutImage(ctx, clusterId, ns, name, image)
}

func (s *logging) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Sync",
			"clusterId", clusterId,
			"ns", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Sync(ctx, clusterId, ns)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "deployment", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
