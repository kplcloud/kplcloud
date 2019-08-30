/**
 * @Time : 2019-07-12 11:55
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package permission

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

func (s *loggingService) Delete(ctx context.Context, id int64) (res []*types.Permission, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "sync",
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, id)
}

func (s *loggingService) Update(ctx context.Context, id int64, icon, keyType string, menu bool, name, path, method string) (res []*types.Permission, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Update",
			"took", time.Since(begin),
			"id", id,
			"icon", icon,
			"keyType", keyType,
			"menu", menu,
			"name", name,
			"path", path,
			"m", method,
			"err", err,
		)
	}(time.Now())
	return s.Service.Update(ctx, id, icon, keyType, menu, name, path, method)
}

func (s *loggingService) Post(ctx context.Context, name, path, method, icon string, isMenu bool, parentId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Post",
			"took", time.Since(begin),
			"icon", icon,
			"menu", isMenu,
			"name", name,
			"path", path,
			"m", method,
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, name, path, method, icon, isMenu, parentId)
}

func (s *loggingService) Drag(ctx context.Context, dragKey, dropKey int64) (res []*types.Permission, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Drag",
			"took", time.Since(begin),
			"dragKey", dragKey,
			"dropKey", dropKey,
			"err", err,
		)
	}(time.Now())
	return s.Service.Drag(ctx, dragKey, dropKey)
}

func (s *loggingService) Menu(ctx context.Context) (res []*types.Permission, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Menu",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Menu(ctx)
}

func (s *loggingService) List(ctx context.Context) (res []*types.Permission, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "List",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx)
}
