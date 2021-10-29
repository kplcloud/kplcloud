/**
 * @Time : 2019-07-23 18:50
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package tools

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

func (s *loggingService) Duplication(ctx context.Context, sourceNamespace, sourceAppName, destinationNamespace string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Duplication",
			"sourceNamespace", sourceNamespace,
			"sourceAppName", sourceAppName,
			"destinationNamespace", destinationNamespace,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Duplication(ctx, sourceNamespace, sourceAppName, destinationNamespace)
}

func (s *loggingService) FakeTime(ctx context.Context, fakeTime time.Time, method FakeTimeMethod) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "FakeTime",
			"fakeTime", fakeTime.String(),
			"method", method,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.FakeTime(ctx, fakeTime, method)
}
