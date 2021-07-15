/**
 * @Time : 3/10/21 3:35 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package sysrole

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

func (l *loggingServer) Delete(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Delete",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Delete(ctx, id)
}

func (l *loggingServer) AddPermissions(ctx context.Context, role *types.SysRole, permissions []types.SysPermission) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "AddPermissions",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.AddPermissions(ctx, role, permissions)
}

func (l *loggingServer) Find(ctx context.Context, id int64) (res types.SysRole, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Find",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Find(ctx, id)
}

func (l *loggingServer) FindByIds(ctx context.Context, ids []int64) (res []types.SysRole, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByIds",
			"ids", ids,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByIds(ctx, ids)
}

func (l *loggingServer) Save(ctx context.Context, data *types.SysRole) (err error) {
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

func (l *loggingServer) List(ctx context.Context, page, pageSize int) (res []types.SysRole, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, page, pageSize)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "sysrole", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
