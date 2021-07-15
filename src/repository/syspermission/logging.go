/**
 * @Time : 5/25/21 3:53 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package syspermission

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/icowan/kit-admin/src/repository/types"
)

type loggingServer struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *loggingServer) Find(ctx context.Context, id int64) (res types.SysPermission, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Find",
			"id", id,
			"err", err,
		)
	}(time.Now())
	return l.next.Find(ctx, id)
}

func (l *loggingServer) FindByIds(ctx context.Context, ids []int64) (res []types.SysPermission, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByIds",
			"ids", ids,
			"err", err,
		)
	}(time.Now())
	return l.next.FindByIds(ctx, ids)
}

func (l *loggingServer) Delete(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Delete",
			"id", id,
			"err", err,
		)
	}(time.Now())
	return l.next.Delete(ctx, id)
}

func (l *loggingServer) Save(ctx context.Context, data *types.SysPermission) (err error) {
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
	logger = log.With(logger, "syspermission", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
