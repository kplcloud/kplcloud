/**
 * @Time : 8/11/21 4:21 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package nodes

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

func (s *logging) Cordon(ctx context.Context, clusterId int64, nodeName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Cordon",
			"clusterId", clusterId,
			"nodeName", nodeName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Cordon(ctx, clusterId, nodeName)
}

func (s *logging) Drain(ctx context.Context, clusterId int64, nodeName string, force bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Drain",
			"clusterId", clusterId,
			"nodeName", nodeName,
			"force", force,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Drain(ctx, clusterId, nodeName, force)
}

func (s *logging) Info(ctx context.Context, clusterId int64, nodeName string) (res infoResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Info",
			"clusterId", clusterId,
			"nodeName", nodeName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Info(ctx, clusterId, nodeName)
}

func (s *logging) List(ctx context.Context, clusterId int64, page, pageSize int) (res []nodeResult, total int, err error) {
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
	return s.next.List(ctx, clusterId, page, pageSize)
}

func (s *logging) Sync(ctx context.Context, clusterName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Sync",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Sync(ctx, clusterName)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "cluster", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
