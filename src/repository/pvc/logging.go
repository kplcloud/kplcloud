/**
 * @Time : 2021/11/25 4:30 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package pvc

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

func (s *logging) Delete(ctx context.Context, pvcId int64, call ...Call) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete", "pvcId", pvcId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, pvcId, call...)
}

func (s *logging) FindByName(ctx context.Context, clusterId int64, ns, name string) (res types.PersistentVolumeClaim, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByName", "clusterId", clusterId, "ns", ns, "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByName(ctx, clusterId, ns, name)
}

func (s *logging) List(ctx context.Context, clusterId int64, storageClassIds []int64, ns, name string, page, pageSize int) (res []types.PersistentVolumeClaim, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "clusterId", clusterId, "storageClassIds", storageClassIds, "namespace", ns, "name", name, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, clusterId, storageClassIds, ns, name, page, pageSize)
}

func (s *logging) Save(ctx context.Context, app *types.PersistentVolumeClaim, call Call) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save", "app", app, "call", call,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, app, call)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "repository.pvc", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
