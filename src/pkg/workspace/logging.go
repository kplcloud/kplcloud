/**
 * @Time : 2019-07-25 15:11
 * @Author : soupzhb@gmail.com
 * @File : logging.go
 * @Software: GoLand
 */

package workspace

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

func (s loggingService) Metrice(ctx context.Context, ns string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Metrice",
			"namespace", ns,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Metrice(ctx, ns)
}

func (s loggingService) Active(ctx context.Context) (res []map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"path", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Active",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Active(ctx)
}
