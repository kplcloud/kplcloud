/**
 * @Time : 7/5/21 5:40 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package audit

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

func (l *logging) List(ctx context.Context, query string, page, pageSize int) (res []types.Audit, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"query", query,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, query, page, pageSize)
}

func (l *logging) Save(ctx context.Context, data *types.Audit) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Save",
			"err", err,
		)
	}(time.Now())
	return l.next.Save(ctx, data)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.audit", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
