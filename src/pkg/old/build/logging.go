/**
 * @Time : 2019-07-09 16:03
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package build

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/middleware"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) Build(ctx context.Context, gitType, version, buildEnv, buildEnvDesc, buildTime string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Build",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"gitType", gitType,
			"version", version,
			"buildEnv", buildEnv,
			"buildEnvDesc", buildEnvDesc,
			"buildTime", buildTime,
			"err", err,
		)
	}(time.Now())
	return s.Service.Build(ctx, gitType, version, buildEnv, buildEnvDesc, buildTime)
}

func (s *loggingService) BuildConsole(ctx context.Context, number, start int) (out string, end int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "BuildConsole",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"number", number,
			"start", start,
			"end", end,
			"err", err,
		)
	}(time.Now())
	return s.Service.BuildConsole(ctx, number, start)
}

func (s *loggingService) AbortBuild(ctx context.Context, jenkinsBuildId int) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "AbortBuild",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"jenkinsBuildId", jenkinsBuildId,
		)
	}(time.Now())
	return s.Service.AbortBuild(ctx, jenkinsBuildId)
}

func (s *loggingService) History(ctx context.Context, page, limit int) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "History",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"page", page,
			"limit", limit,
			"err", err,
		)
	}(time.Now())
	return s.Service.History(ctx, page, limit)
}

func (s *loggingService) Rollback(ctx context.Context, buildId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "RollBack",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"buildId", buildId,
			"err", err,
		)
	}(time.Now())
	return s.Service.Rollback(ctx, buildId)
}

func (s *loggingService) BuildConf(ctx context.Context, ns, name string) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "RollBack",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"name", name,
			"err", err,
		)
	}(time.Now())
	return s.Service.BuildConf(ctx, ns, name)
}

func (s *loggingService) CronHistory(ctx context.Context, page, limit int) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "CronHistory",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"page", page,
			"limit", limit,
			"err", err,
		)
	}(time.Now())
	return s.Service.CronHistory(ctx, page, limit)
}

func (s *loggingService) CronBuildConsole(ctx context.Context, number, start int) (out string, end int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "CronBuildConsole",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"number", number,
			"start", start,
			"end", end,
			"err", err,
		)
	}(time.Now())
	return s.Service.CronBuildConsole(ctx, number, start)
}
