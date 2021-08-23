/**
 * @Time : 8/19/21 1:29 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package secrets

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

func (l *logging) FindBy(ctx context.Context, clusterId int64, ns, name string) (res types.Secret, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindBy",
			"clusterId", clusterId,
			"ns", ns,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindBy(ctx, clusterId, ns, name)
}

func (l *logging) Save(ctx context.Context, configMap *types.Secret, data []types.Data) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Save(ctx, configMap, data)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "secrets", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
