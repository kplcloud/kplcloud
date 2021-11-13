/**
 * @Time : 2021/9/17 3:28 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package account

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) Namespaces(ctx context.Context, userId, clusterId int64) (res []nsResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Namespaces", "userId", userId, "clusterId", clusterId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Namespaces(ctx, userId, clusterId)
}

func (s *logging) Logout(ctx context.Context, userId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Logout",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Logout(ctx, userId)
}

func (s *logging) Menus(ctx context.Context, userId int64) (res []userMenuResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Menus",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Menus(ctx, userId)
}

func (s *logging) UserInfo(ctx context.Context, userId int64) (res userInfoResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UserInfo",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.UserInfo(ctx, userId)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "account", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
