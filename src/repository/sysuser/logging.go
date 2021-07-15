/**
 * @Time : 3/5/21 2:19 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package sysuser

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

func (l *loggingServer) FindByRoleId(ctx context.Context, roleId int64, page, pageSize int) (res []types.SysUser, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByRoleId",
			"roleId", roleId,
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByRoleId(ctx, roleId, page, pageSize)
}

func (l *loggingServer) FindByIds(ctx context.Context, ids []int64, preload ...string) (res []types.SysUser, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByIds",
			"ids", ids,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByIds(ctx, ids, preload...)
}

func (l *loggingServer) AddRoles(ctx context.Context, user *types.SysUser, roles []types.SysRole) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "AddRoles",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.AddRoles(ctx, user, roles)
}

func (l *loggingServer) Find(ctx context.Context, userId int64, preload ...string) (res types.SysUser, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Find",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Find(ctx, userId, preload...)
}

func (l *loggingServer) List(ctx context.Context, namespaceIds []int64, email string, page, pageSize int) (res []types.SysUser, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"email", email,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, namespaceIds, email, page, pageSize)
}

func (l *loggingServer) Save(ctx context.Context, user *types.SysUser) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Save(ctx, user)
}

func (l *loggingServer) FindByEmail(ctx context.Context, email string) (res types.SysUser, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByEmail",
			"email", email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByEmail(ctx, email)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "sysuser", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
