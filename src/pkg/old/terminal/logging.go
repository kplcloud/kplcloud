/**
 * @Time : 2019-06-27 18:14
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package terminal

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

func (s *loggingService) Index(ctx context.Context, podName, container string) (*IndexData, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Index",
			"podName", podName,
			"container", container,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Index(ctx, podName, container)
}
