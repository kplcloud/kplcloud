/**
 * @Time : 2019-07-04 16:12
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package pod

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/middleware"
	"io"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) Detail(ctx context.Context, podName string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Detail",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"podName", podName,
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, podName)
}

func (s *loggingService) ProjectPods(ctx context.Context) (res []map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "ProjectPods",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"err", err,
		)
	}(time.Now())
	return s.Service.ProjectPods(ctx)
}

func (s *loggingService) GetLog(ctx context.Context, podName, container string, previous bool) (res *LogDetails, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "GetLog",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"podName", podName,
			"container", container,
			"previous", previous,
			"err", err,
		)
	}(time.Now())
	return s.Service.GetLog(ctx, podName, container, previous)
}

func (s *loggingService) DownloadLog(ctx context.Context, podName, container string, previous bool) (res io.ReadCloser, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "DownloadLog",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"podName", podName,
			"container", container,
			"previous", previous,
			"err", err,
		)
	}(time.Now())
	return s.Service.DownloadLog(ctx, podName, container, previous)
}

func (s *loggingService) Delete(ctx context.Context, podName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Delete",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"podName", podName,
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, podName)
}

func (s *loggingService) PodsMetrics(ctx context.Context) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Metrics",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"err", err,
		)
	}(time.Now())
	return s.Service.PodsMetrics(ctx)
}
