/**
 * @Time : 2019-07-17 19:02
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package persistentvolume

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

func (s *loggingService) Get(ctx context.Context, name string) (rs interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Get",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, name)
}
