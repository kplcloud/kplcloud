/**
 * @Time : 2019/7/17 2:16 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package consul

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

func (s *loggingService) Detail(ctx context.Context, ns, name string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Detail",
			"namespace", ns,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, ns, name)
}

func (s *loggingService) List(ctx context.Context, ns, name string, page, limit int) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "List",
			"namespace", ns,
			"name", name,
			"page", page,
			"limit", limit,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, ns, name, page, limit)
}

func (s *loggingService) Post(ctx context.Context, ns, name, clientType, rules string) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Post",
			"namespace", ns,
			"name", name,
			"clientType", clientType,
			"rules", rules,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Post(ctx, ns, name, clientType, rules)
}

func (s *loggingService) Update(ctx context.Context, ns, name, clientType, rules string) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Update",
			"namespace", ns,
			"name", name,
			"clientType", clientType,
			"rules", rules,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Update(ctx, ns, name, clientType, rules)
}

func (s *loggingService) Delete(ctx context.Context, ns, name string) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Delete",
			"namespace", ns,
			"name", name,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Delete(ctx, ns, name)
}

func (s *loggingService) KVDetail(ctx context.Context, ns, name, prefix string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "KVDetail",
			"namespace", ns,
			"name", name,
			"prefix", prefix,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.KVDetail(ctx, ns, name, prefix)
}

func (s *loggingService) KVList(ctx context.Context, ns, name, prefix string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "KVList",
			"namespace", ns,
			"name", name,
			"prefix", prefix,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.KVList(ctx, ns, name, prefix)
}

func (s *loggingService) KVPost(ctx context.Context, ns, name, key, value string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "KVPost",
			"namespace", ns,
			"name", name,
			"key", key,
			"value", value,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.KVPost(ctx, ns, name, key, value)
}

func (s *loggingService) KVDelete(ctx context.Context, ns, name, prefix string, filderState bool) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "KVDelete",
			"namespace", ns,
			"name", name,
			"prefix", prefix,
			"filderState", filderState,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.KVDelete(ctx, ns, name, prefix, filderState)
}
