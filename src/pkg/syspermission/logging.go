/**
 * @Time : 5/25/21 3:11 PM
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
)

type loggingServer struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *loggingServer) All(ctx context.Context) (res []result, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "All",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.All(ctx)
}

func (s *loggingServer) Drag(ctx context.Context, dragKey, dropKey int64) (res []result, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Drag",
			"dragKey", dragKey,
			"dropKey", dropKey,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Drag(ctx, dragKey, dropKey)
}

func (s *loggingServer) Add(ctx context.Context, name, alias, icon, path, method, desc string, parentId int64, menu bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Add",
			"name", name,
			"alias", alias,
			"icon", icon,
			"path", path,
			"method", method,
			"desc", desc,
			"parentId", parentId,
			"menu", menu,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Add(ctx, name, alias, icon, path, method, desc, parentId, menu)
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

func (s *loggingServer) Update(ctx context.Context, id int64, name, alias, icon, path, method, desc string, parentId int64, menu bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Update",
			"id", id,
			"name", name,
			"alias", alias,
			"icon", icon,
			"path", path,
			"method", method,
			"desc", desc,
			"parentId", parentId,
			"menu", menu,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Update(ctx, id, name, alias, icon, path, method, desc, parentId, menu)
}

func (s *loggingServer) List(ctx context.Context, page, pageSize int) (res []result, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List",
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, page, pageSize)
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
