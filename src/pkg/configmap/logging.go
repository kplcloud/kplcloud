/**
 * @Time: 2021/8/18 23:11
 * @Author: solacowa@gmail.com
 * @File: logging
 * @Software: GoLand
 */

package configmap

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

func (s *logging) Sync(ctx context.Context, ns string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Sync",
			"ns", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Sync(ctx, ns)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "configmap", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
