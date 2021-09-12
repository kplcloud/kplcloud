/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package namespace

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// 链路追踪中间件
type tracing struct {
	next   Service
	tracer stdopentracing.Tracer
}

func (s *tracing) Create(ctx context.Context, clusterId int64, name, alias, remark string, imageSecrets []string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Namespace",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"name", name, "alias", alias, "remark", remark,
			"err", err)
		span.Finish()
	}()
	return s.next.Create(ctx, clusterId, name, alias, remark, imageSecrets)
}

func (s *tracing) Sync(ctx context.Context, clusterId int64) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Namespace",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"err", err)
		span.Finish()
	}()
	return s.next.Sync(ctx, clusterId)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
