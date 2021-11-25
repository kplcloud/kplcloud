/**
 * @Time : 3/10/21 3:35 PM
 * @Author : solacowa@gmais.com
 * @File : logging
 * @Software: GoLand
 */

package cluster

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

func (s *logging) Find(ctx context.Context, id int64) (res types.Cluster, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Find", "id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Find(ctx, id)
}

func (s *logging) FindByNames(ctx context.Context, names []string) (res []types.Cluster, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByNames",
			"names", names,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByNames(ctx, names)
}

func (s *logging) FindByIds(ctx context.Context, ids []int64) (res []types.Cluster, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByIds",
			"ids", ids,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByIds(ctx, ids)
}

func (s *logging) List(ctx context.Context, name string, status, page, pageSize int) (res []types.Cluster, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List",
			"name", name,
			"status", status,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, name, status, page, pageSize)
}

func (s *logging) SaveRole(ctx context.Context, clusterRole *types.ClusterRole, roles []types.PolicyRule) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SaveRole",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.SaveRole(ctx, clusterRole, roles)
}

func (s *logging) FindAll(ctx context.Context, status int) (res []types.Cluster, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindAll",
			"status", status,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindAll(ctx, status)
}

func (s *logging) FindByName(ctx context.Context, name string) (res types.Cluster, err error) {
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

func (s *logging) Save(ctx context.Context, data *types.Cluster, calls ...Call) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, data, calls...)
}

func (s *logging) Delete(ctx context.Context, id int64, unscoped bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete",
			"id", id,
			"unscoped", unscoped,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, id, unscoped)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.cluster", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
