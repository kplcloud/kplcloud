/**
 * @Time : 3/10/21 3:35 PM
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

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *logging) List(ctx context.Context, name string, status, page, pageSize int) (res []types.Cluster, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"name", name,
			"status", status,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, name, status, page, pageSize)
}

func (l *logging) SaveRole(ctx context.Context, clusterRole *types.ClusterRole, roles []types.PolicyRule) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "SaveRole",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.SaveRole(ctx, clusterRole, roles)
}

func (l *logging) FindAll(ctx context.Context, status int) (res []types.Cluster, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindAll",
			"status", status,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindAll(ctx, status)
}

func (l *logging) FindByName(ctx context.Context, name string) (res types.Cluster, err error) {
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

func (l *logging) Save(ctx context.Context, data *types.Cluster, calls ...Call) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Save(ctx, data, calls...)
}

func (l *logging) Delete(ctx context.Context, id int64, unscoped bool) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Delete",
			"id", id,
			"unscoped", unscoped,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Delete(ctx, id, unscoped)
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
