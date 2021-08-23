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

func (l *logging) FindName(ctx context.Context, clusterId int64, name string) (res types.StorageClass, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindName",
			"clusterId", clusterId,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindName(ctx, clusterId, name)
}

func (l *logging) Find(ctx context.Context, id int64) (res types.StorageClass, err error) {
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

func (l *logging) FirstInsert(ctx context.Context, data *types.StorageClass) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FirstInsert",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FirstInsert(ctx, data)
}

func (l *logging) Save(ctx context.Context, data *types.StorageClass) (err error) {
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
