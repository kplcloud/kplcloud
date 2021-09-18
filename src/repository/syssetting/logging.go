/**
 * @Time : 3/4/21 11:14 AM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package syssetting

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

func (l *loggingServer) List(ctx context.Context, key string, page, pageSize int) (res []types.SysSetting, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"key", key,
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, key, page, pageSize)
}

func (l *loggingServer) FindAll(ctx context.Context) (res []types.SysSetting, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindAll",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindAll(ctx)
}

func (l *loggingServer) Add(ctx context.Context, section, key, val, desc string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Add",
			"section", section,
			"key", key,
			"val", val,
			"desc", desc,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Add(ctx, section, key, val, desc)
}

func (l *loggingServer) Delete(ctx context.Context, section, key string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Delete",
			"key", key,
			"section", section,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Delete(ctx, section, key)
}

func (l *loggingServer) Update(ctx context.Context, data *types.SysSetting) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Update",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Update(ctx, data)
}

func (l *loggingServer) Find(ctx context.Context, section, key string) (res types.SysSetting, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Find",
			"section", section,
			"key", key,
			"err", err,
		)
	}(time.Now())
	return l.next.Find(ctx, section, key)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "syssetting", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
