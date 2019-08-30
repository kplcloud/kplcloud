/**
 * @Time : 2019-06-27 13:34
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package public

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) GitPost(ctx context.Context, namespace, name, token, keyWord, branch string, req gitlabHook) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "sync",
			"namespace", namespace,
			"name", name,
			"token", token,
			"keyword", keyWord,
			"branch", branch,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GitPost(ctx, namespace, name, token, keyWord, branch, req)
}

func (s *loggingService) PrometheusAlert(ctx context.Context, req *prometheusAlerts) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "sync",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.PrometheusAlert(ctx, req)
}

func (s *loggingService) Config(ctx context.Context) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "config",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Config(ctx)
}
