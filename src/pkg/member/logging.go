/**
 * @Time : 2019-07-17 14:18
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package member

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

func (s *loggingService) Namespaces(ctx context.Context) (list []map[string]string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Namespaces",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Namespaces(ctx)
}

func (s *loggingService) Detail(ctx context.Context, id int64) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Detail",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, id)
}

func (s *loggingService) Post(ctx context.Context, username, email, password string, state int64, namespaces []string, roleIds []int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Post",
			"username", username,
			"email", email,
			"state", state,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, username, email, password, state, namespaces, roleIds)
}

func (s *loggingService) Update(ctx context.Context, id int64, username, email, password string, state int64, namespaces []string, roleIds []int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Update",
			"id", id,
			"username", username,
			"email", email,
			"state", state,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Update(ctx, id, username, email, password, state, namespaces, roleIds)
}

func (s *loggingService) List(ctx context.Context, page, limit int, email string) (rs map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Update",
			"page", page,
			"limit", limit,
			"email", email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, page, limit, email)
}

func (s *loggingService) Me(ctx context.Context) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Me",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Me(ctx)
}

func (s *loggingService) All(ctx context.Context) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "All",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.All(ctx)
}
