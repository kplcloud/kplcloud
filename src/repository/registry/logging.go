/**
 * @Time : 8/19/21 1:29 PM
 * @Author : solacowa@gmais.com
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

func (s *logging) Delete(ctx context.Context, id int64, call Call) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete", "id", id, "call", call,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, id, call)
}

func (s *logging) SaveCall(ctx context.Context, reg *types.Registry, call Call) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SaveCall", "reg", reg, "call", call,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.SaveCall(ctx, reg, call)
}

func (s *logging) List(ctx context.Context, query string, page, pageSize int) (res []types.Registry, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List",
			"query", query,
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, query, page, pageSize)
}

func (s *logging) FindByNames(ctx context.Context, names []string) (res []types.Registry, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save",
			"names", names,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByNames(ctx, names)
}

func (s *logging) Save(ctx context.Context, reg *types.Registry) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, reg)
}

func (s *logging) FindByName(ctx context.Context, name string) (res types.Registry, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByName",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByName(ctx, name)
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
