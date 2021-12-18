/**
 * @Time : 8/19/21 1:29 PM
 * @Author : solacowa@gmais.com
 * @File : logging
 * @Software: GoLand
 */

package secrets

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

func (s *logging) List(ctx context.Context, clusterId int64, ns, name string, page, pageSize int) (res []types.Secret, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List", "clusterId", clusterId, "ns", ns, "name", name, "page", page, "pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List(ctx, clusterId, ns, name, page, pageSize)
}

func (s *logging) FindNsByNames(ctx context.Context, clusterId int64, ns string, names []string) (res []types.Secret, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindNsByNames", "clusterId", clusterId, "ns", ns, "names", names,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindNsByNames(ctx, clusterId, ns, names)
}

func (s *logging) FindByName(ctx context.Context, name string) (res []types.Secret, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindByName", "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindByName(ctx, name)
}

func (s *logging) Delete(ctx context.Context, clusterId int64, ns, name string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete",
			"clusterId", clusterId,
			"ns", ns,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, clusterId, ns, name)
}

func (s *logging) FindBy(ctx context.Context, clusterId int64, ns, name string) (res types.Secret, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "FindBy",
			"clusterId", clusterId,
			"ns", ns,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindBy(ctx, clusterId, ns, name)
}

func (s *logging) Save(ctx context.Context, configMap *types.Secret, data []types.Data) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Save(ctx, configMap, data)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "secrets", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
