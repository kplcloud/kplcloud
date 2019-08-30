/**
 * @Time : 2019-07-16 17:59
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package role

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) PermissionSelected(ctx context.Context, id int64) (ids []int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "PermissionSelected",
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.PermissionSelected(ctx, id)
}

func (s *loggingService) Detail(ctx context.Context, id int64) (ids *types.Role, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Detail",
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, id)
}

func (s *loggingService) Post(ctx context.Context, name, desc string, level int) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Post",
			"took", time.Since(begin),
			"name", name,
			"desc", desc,
			"level", level,
		)
	}(time.Now())
	return s.Service.Post(ctx, name, desc, level)
}

func (s *loggingService) Update(ctx context.Context, id int64, name, desc string, level int) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Update",
			"took", time.Since(begin),
			"name", name,
			"id", id,
			"desc", desc,
			"level", level,
		)
	}(time.Now())
	return s.Service.Update(ctx, id, name, desc, level)
}

func (s *loggingService) All(ctx context.Context) ([]*types.Role, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "All",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.All(ctx)
}

func (s *loggingService) Delete(ctx context.Context, id int64) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "All",
			"id", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Delete(ctx, id)
}

func (s *loggingService) RolePermission(ctx context.Context, id int64, permIds []int64) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "All",
			"id", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.RolePermission(ctx, id, permIds)
}
