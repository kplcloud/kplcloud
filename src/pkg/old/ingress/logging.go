/**
 * @Time : 2019/7/2 2:10 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package ingress

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
	return &loggingService{logger: level.Info(logger), Service: s}
}

func (s loggingService) Get(ctx context.Context, ns string, name string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"namespace", ns,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, ns, name)
}

func (s loggingService) List(ctx context.Context, ns string, page, limit int) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"namespace", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, ns, page, limit)
}

func (s loggingService) Post(ctx context.Context, req postRequest) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "post",
			"namespace", req.Namespace,
			"name", req.Name,
			"rules", req.Rules,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, req)
}

func (s loggingService) GetNoIngressProject(ctx context.Context, ns string) (res []map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"namespace", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetNoIngressProject(ctx, ns)
}

func (s loggingService) Sync(ctx context.Context, ns string) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"namespace", ns,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Sync(ctx, ns)
}

func (s loggingService) Generate(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Generate",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Generate(ctx)
}
