/**
 * @Time : 2019-07-09 18:50
 * @Author : soupzhb@gmail.com
 * @File : endpoint.go
 * @Software: GoLand
 */

package wechat

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

func (s loggingService) Receive(ctx context.Context) (str, contentType string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Receive(ctx)
}

func (s loggingService) GetQr(ctx context.Context, req qrRequest) (res interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "post",
			"email", req.Email,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetQr(ctx, req)
}
