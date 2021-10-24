/**
 * @Time : 8/19/21 1:29 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package registry

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

func (l *logging) List(ctx context.Context, query string, page, pageSize int) (res []types.Registry, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"query", query,
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, query, page, pageSize)
}

func (l *logging) FindByNames(ctx context.Context, names []string) (res []types.Registry, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Save",
			"names", names,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByNames(ctx, names)
}

func (l *logging) Save(ctx context.Context, reg *types.Registry) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Save(ctx, reg)
}

func (l *logging) FindByName(ctx context.Context, name string) (res types.Registry, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByName",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByName(ctx, name)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "registry", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
