/**
 * @Time : 8/11/21 3:19 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package kubernetes

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    K8sClient
	traceId string
}

func (s *logging) Do(ctx context.Context) *kubernetes.Clientset {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Do",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.Do(ctx)
}

func (s *logging) Config(ctx context.Context) *rest.Config {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Config",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.Config(ctx)
}

func (s *logging) Reload(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Reload",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Reload(ctx)
}

func (s *logging) Connect(ctx context.Context, name, configData string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Connect",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Connect(ctx, name, configData)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "kubernetes", "logging")
	return func(next K8sClient) K8sClient {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
