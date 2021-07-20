/**
 * @Time : 3/10/21 3:40 PM
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
)

type loggingServer struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *loggingServer) Delete(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, id)
}

func (s *loggingServer) Update(ctx context.Context, id int64, alias, name, description string, enabled bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Update",
			"id", id,
			"alias", alias,
			"name", name,
			"description", description,
			"enabled", enabled,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Update(ctx, id, alias, name, description, enabled)
}

func (s *loggingServer) User(ctx context.Context, id int64, userId []int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "User",
			"id", id,
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.User(ctx, id, userId)
}

func (s *loggingServer) Permissions(ctx context.Context, id int64, page, pageSize int) (res []listResult, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Permissions",
			"id", id,
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Permissions(ctx, id, page, pageSize)
}

func (s *loggingServer) Permission(ctx context.Context, id int64, permIds []int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Permission",
			"id", id,
			"permIds", permIds,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Permission(ctx, id, permIds)
}

//func (s *loggingServer) Users(ctx context.Context, id int64, page, pageSize int) (res []listResult, total int, err error) {
//	defer func(begin time.Time) {
//		_ = s.logger.Log(
//			s.traceId, ctx.Value(s.traceId),
//			"method", "Users",
//			"id", id,
//			"page", page,
//			"pageSize", pageSize,
//			"took", time.Since(begin),
//			"err", err,
//		)
//	}(time.Now())
//	return s.next.Users(ctx, id, page, pageSize)
//}

func (s *loggingServer) Add(ctx context.Context, alias, name, description string, enabled bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Add",
			"alias", alias,
			"name", name,
			"description", description,
			"enabled", enabled,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Add(ctx, alias, name, description, enabled)
}

func (s *loggingServer) List(ctx context.Context, page, pageSize int) (res []listResult, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List",
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, page, pageSize)
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
