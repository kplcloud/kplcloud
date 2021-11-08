/**
 * @Time : 3/10/21 3:35 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package storageclass

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

func (s *logging) List(ctx context.Context, clusterId int64, page, pageSize int) (res []types.StorageClass, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "clusterId", clusterId, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, clusterId, page, pageSize)
}

func (s *logging) FindName(ctx context.Context, clusterId int64, name string) (res types.StorageClass, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindName",
			"clusterId", clusterId,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindName(ctx, clusterId, name)
}

func (s *logging) Find(ctx context.Context, id int64) (res types.StorageClass, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Find",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Find(ctx, id)
}

func (s *logging) FirstInsert(ctx context.Context, data *types.StorageClass) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FirstInsert",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FirstInsert(ctx, data)
}

func (s *logging) Save(ctx context.Context, data *types.StorageClass, call Call) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, data, call)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "storageClass", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
