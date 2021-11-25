/**
 * @Time : 3/5/21 2:43 PM
 * @Author : solacowa@gmais.com
 * @File : logging
 * @Software: GoLand
 */

package namespace

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

func (s *logging) FindByCluster(ctx context.Context, clusterId int64) (res []types.Namespace, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByCluster", "clusterId", clusterId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByCluster(ctx, clusterId)
}

func (s *logging) FindByNames(ctx context.Context, clusterId int64, names []string) (res []types.Namespace, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByNames", "clusterId", clusterId, "names", names,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByNames(ctx, clusterId, names)
}

func (s *logging) List(ctx context.Context, clusterId int64, names []string, query string, page, pageSize int) (res []types.Namespace, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List",
			"clusterId", clusterId,
			"query", query,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, clusterId, names, query, page, pageSize)
}

func (s *logging) SaveCall(ctx context.Context, data *types.Namespace, call Call) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "SaveCall",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.SaveCall(ctx, data, call)
}

func (s *logging) FindByName(ctx context.Context, clusterId int64, name string) (res types.Namespace, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByName",
			"clusterId", clusterId,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByName(ctx, clusterId, name)
}

func (s *logging) FindByIds(ctx context.Context, ids []int64) (res []types.Namespace, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByIds",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByIds(ctx, ids)
}

func (s *logging) Save(ctx context.Context, data *types.Namespace) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, data)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "namespace", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
