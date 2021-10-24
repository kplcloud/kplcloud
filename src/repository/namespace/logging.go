/**
 * @Time : 3/5/21 2:43 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package namespace

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type loggingServer struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *loggingServer) List(ctx context.Context, clusterId int64, names []string, query string, page, pageSize int) (res []types.Namespace, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"clusterId", clusterId,
			"query", query,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, clusterId, names, query, page, pageSize)
}

func (l *loggingServer) SaveCall(ctx context.Context, data *types.Namespace, call Call) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "SaveCall",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.SaveCall(ctx, data, call)
}

func (l *loggingServer) FindByName(ctx context.Context, clusterId int64, name string) (res types.Namespace, err error) {
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

func (l *loggingServer) FindByIds(ctx context.Context, ids []int64) (res []types.Namespace, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByIds",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByIds(ctx, ids)
}

func (l *loggingServer) Save(ctx context.Context, data *types.Namespace) (err error) {
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

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "namespace", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
