/**
 * @Time : 3/5/21 2:43 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package k8stpl

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type loggingServer struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *loggingServer) EncodeTemplate(ctx context.Context, kind types.Kind, paramContent map[string]interface{}, data interface{}) (tpl []byte, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "EncodeTemplate",
			"kind", kind,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.EncodeTemplate(ctx, kind, paramContent, data)
}

func (l *loggingServer) FindByKind(ctx context.Context, kind types.Kind) (tpl types.K8sTemplate, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByKind",
			"kind", kind,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByKind(ctx, kind)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "template", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
