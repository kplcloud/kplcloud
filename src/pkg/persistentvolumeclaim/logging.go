/**
 * @Time : 8/11/21 4:21 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package persistentvolumeclaim

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) List(ctx context.Context, clusterId int64, storageClass, ns string, page, pageSize int) (resp []result, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "clusterId", clusterId, "storageClass", storageClass, "ns", ns, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, clusterId, storageClass, ns, page, pageSize)
}

func (s *logging) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Sync",
			"clusterId", clusterId,
			"ns", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Sync(ctx, clusterId, ns)
}

func (s *logging) Get(ctx context.Context, clusterId int64, ns, name string) (res result, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Get", "clusterId", clusterId, "ns", ns, "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Get(ctx, clusterId, ns, name)
}

func (s *logging) Delete(ctx context.Context, clusterId int64, ns, name string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete", "clusterId", clusterId, "ns", ns, "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, clusterId, ns, name)
}

func (s *logging) Create(ctx context.Context, clusterId int64, ns, name, storage, storageClassName string, accessModes []string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create",
			"clusterId", clusterId,
			"ns", ns,
			"name", name, "storage", storage,
			"storageClassName", storageClassName,
			"accessModes", accessModes,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, clusterId, ns, name, storage, storageClassName, accessModes)
}

func (s *logging) All(ctx context.Context, clusterId int64) (resp map[string]interface{}, err error) {
	panic("implement me")
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "persistentvolumeclaim", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
