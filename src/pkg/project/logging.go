/**
 * @Time : 2019-07-02 17:18
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package project

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/util/pods"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) Post(ctx context.Context, ns, name, displayName, desc string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Post",
			"namespace", ns,
			"name", name,
			"displayName", displayName,
			"desc", desc,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, ns, name, displayName, desc)
}

func (s *loggingService) BasicPost(ctx context.Context, name string, req basicRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "BasicPost",
			"namespace", req.Namespace,
			"name", req.Name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.BasicPost(ctx, name, req)
}

func (s *loggingService) List(ctx context.Context, page, limit int, name string, groupId int64) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "List",
			"name", name,
			"page", page,
			"limit", limit,
			"groupId", groupId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, page, limit, name, groupId)
}

func (s *loggingService) ListByNs(ctx context.Context) (res []map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "ListByNs",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.ListByNs(ctx)
}

func (s *loggingService) PomFile(ctx context.Context, pomFile string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "PomFile",
			"took", time.Since(begin),
			"pomFile", pomFile,
			"err", err,
		)
	}(time.Now())
	return s.Service.PomFile(ctx, pomFile)
}

func (s *loggingService) GitAddr(ctx context.Context, gitAddr string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "GitAddr",
			"took", time.Since(begin),
			"gitAddr", gitAddr,
			"err", err,
		)
	}(time.Now())
	return s.Service.GitAddr(ctx, gitAddr)
}

func (s *loggingService) Detail(ctx context.Context) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Detail",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx)
}

func (s *loggingService) Update(ctx context.Context, displayName, desc string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Update",
			"took", time.Since(begin),
			"displayName", displayName,
			"desc", desc,
			"err", err,
		)
	}(time.Now())
	return s.Service.Update(ctx, displayName, desc)
}

func (s *loggingService) Sync(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Sync",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Sync(ctx)
}

func (s *loggingService) Delete(ctx context.Context, ns, name, code string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Delete",
			"name", name,
			"namespace", ns,
			"code", code,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, ns, name, code)
}

func (s *loggingService) Config(ctx context.Context) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Config",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Config(ctx)
}

func (s *loggingService) Monitor(ctx context.Context, metrics, podName, container string) (res map[string]map[string]map[string][]pods.XYRes, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Monitor",
			"metrics", metrics,
			"podName", podName,
			"container", container,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Monitor(ctx, metrics, podName, container)
}

func (s *loggingService) Alerts(ctx context.Context) (res alertsResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Alerts",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Alerts(ctx)
}
