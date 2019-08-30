/**
 * @Time : 2019/7/24 3:46 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package audit

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

func (s *loggingService) AccessAudit(ctx context.Context, ns, name string) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "AccessAudit",
			"took", time.Since(begin),
			"namespace", ns,
			"name", name,
		)
	}(time.Now())
	return s.Service.AccessAudit(ctx, ns, name)
}

func (s *loggingService) AuditStep(ctx context.Context, ns, name, kind string) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "AuditStep",
			"took", time.Since(begin),
			"namespace", ns,
			"name", name,
			"kind", kind,
		)
	}(time.Now())
	return s.Service.AuditStep(ctx, ns, name, kind)
}

func (s *loggingService) Refused(ctx context.Context, ns, name string) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "Build",
			"took", time.Since(begin),
			"namespace", ns,
			"name", name,
		)
	}(time.Now())
	return s.Service.Refused(ctx, ns, name)
}
