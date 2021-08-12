/**
 * @Time : 8/11/21 4:21 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package cluster

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

func (s *logging) Add(ctx context.Context, name, alias, data string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Add",
			"name", name,
			"alias", alias,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Add(ctx, name, alias, data)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "cluster", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
