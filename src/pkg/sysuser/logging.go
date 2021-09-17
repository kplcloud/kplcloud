/**
 * @Time : 3/9/21 5:58 PM
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
)

type loggingServer struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *loggingServer) Locked(ctx context.Context, userId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Locked",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Locked(ctx, userId)
}

func (s *loggingServer) Delete(ctx context.Context, userId int64, unscoped bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, userId, unscoped)
}

func (s *loggingServer) Update(ctx context.Context, userId int64, username, email string, locked bool, clusterIds, roleIds []int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Update",
			"userId", userId,
			"username", username,
			"email", email,
			"locked", locked,
			"clusterIds", clusterIds,
			"roleIds", roleIds,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Update(ctx, userId, username, email, locked, clusterIds, roleIds)
}

func (s *loggingServer) Add(ctx context.Context, username, email string, locked bool, clusterIds, roleIds []int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Add",
			"username", username,
			"email", email,
			"locked", locked,
			"clusterIds", clusterIds,
			"roleIds", roleIds,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Add(ctx, username, email, locked, clusterIds, roleIds)
}

func (s *loggingServer) List(ctx context.Context, email string, page, pageSize int) (res []listResult, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List",
			"email", email,
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, email, page, pageSize)
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
