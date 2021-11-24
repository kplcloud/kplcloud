/**
 * @Time : 8/11/21 4:21 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package namespace

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

func (s *logging) IssueSecret(ctx context.Context, clusterId int64, name, regName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "IssueSecret", "clusterId", clusterId, "name", name, "regName", regName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.IssueSecret(ctx, clusterId, name, regName)
}

func (s *logging) ReloadSecret(ctx context.Context, clusterId int64, name, regName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ReloadSecret", "clusterId", clusterId, "name", name, "regName", regName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ReloadSecret(ctx, clusterId, name, regName)
}

func (s *logging) Info(ctx context.Context, clusterId int64, name string) (res result, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Info", "clusterId", clusterId, "name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Info(ctx, clusterId, name)
}

func (s *logging) List(ctx context.Context, clusterId int64, names []string, query string, page, pageSize int) (res []result, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "List",
			"clusterId", clusterId,
			"names", names,
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

func (s *logging) Delete(ctx context.Context, clusterId int64, name string, force bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Delete",
			"clusterId", clusterId,
			"name", name,
			"force", force,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Delete(ctx, clusterId, name, force)
}

func (s *logging) Update(ctx context.Context, clusterId int64, name string, alias, remark, status string, imageSecrets []string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Update",
			"clusterId", clusterId,
			"name", name,
			"alias", alias,
			"remark", remark,
			"status", status,
			"imageSecrets", imageSecrets,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Update(ctx, clusterId, name, alias, remark, status, imageSecrets)
}

func (s *logging) Create(ctx context.Context, clusterId int64, name, alias, remark string, imageSecrets []string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Create",
			"clusterId", clusterId,
			"name", name,
			"alias", alias,
			"remark", remark,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Create(ctx, clusterId, name, alias, remark, imageSecrets)
}

func (s *logging) Sync(ctx context.Context, clusterId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Sync",
			"clusterId", clusterId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Sync(ctx, clusterId)
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
