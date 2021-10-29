/**
 * @Time : 2019-07-22 14:23
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package event

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/middleware"
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

func (s *loggingService) All(ctx context.Context) (res []*types.Event, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "All",
			"took", time.Since(begin),
			"namespace", ctx.Value(middleware.NamespaceContext),
			"err", err,
		)
	}(time.Now())
	return s.Service.All(ctx)
}
