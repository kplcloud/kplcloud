/**
 * @Time : 3/10/21 3:35 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package nodes

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

func (l *logging) Save(ctx context.Context, data *types.Nodes) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Save(ctx, data)
}

func (l *logging) FindByName(ctx context.Context, clusterId int64, name string) (res types.Nodes, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByName",
			"clusterId", clusterId,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByName(ctx, clusterId, name)
}

func (l *logging) List(ctx context.Context, clusterId int64, query string, page, pageSize int) (res []types.Nodes, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"clusterId", clusterId,
			"page", page,
			"pageSize", pageSize,
			"query", query,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, clusterId, query, page, pageSize)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "nodes", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
