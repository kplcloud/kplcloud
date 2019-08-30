/**
 * @Time : 2019-07-29 11:31
 * @Author : soupzhb@gmail.com
 * @File : logging.go
 * @Software: GoLand
 */

package statistics

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

func (s loggingService) Build(ctx context.Context, req buildRequest) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Metrice",
			"namespace", req.Namespace,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Build(ctx, req)
}
