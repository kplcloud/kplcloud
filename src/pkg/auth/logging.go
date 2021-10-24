/**
 * @Time : 8/11/21 4:21 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package auth

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

func (s *logging) Register(ctx context.Context, username, email, password, mobile, remark string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Register",
			"username", username,
			"mobile", mobile,
			"remark", remark,
			"email", email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Register(ctx, username, email, password, mobile, remark)
}

func (s *logging) Login(ctx context.Context, username, password string) (rs string, sessionTimeout int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Login",
			"username", username,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Login(ctx, username, password)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "auth", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
