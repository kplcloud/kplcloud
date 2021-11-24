/**
 * @Time : 5/25/21 3:53 PM
 * @Author : solacowa@gmais.com
 * @File : logging
 * @Software: GoLand
 */

package syspermission

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

func (s *logging) FindAll(ctx context.Context) (res []types.SysPermission, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindAll",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindAll(ctx)
}

func (s *logging) Find(ctx context.Context, id int64) (res types.SysPermission, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Find",
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.next.Find(ctx, id)
}

func (s *logging) FindByIds(ctx context.Context, ids []int64) (res []types.SysPermission, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByIds",
			"ids", ids,
			"err", err,
		)
	}(time.Now())
	return s.next.FindByIds(ctx, ids)
}

func (s *logging) Delete(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete",
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, id)
}

func (s *logging) Save(ctx context.Context, data *types.SysPermission) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save",
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, data)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "syspermission", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
