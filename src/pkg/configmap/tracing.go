/**
 * @Time: 2021/8/18 23:13
 * @Author: solacowa@gmail.com
 * @File: tracing
 * @Software: GoLand
 */

package configmap

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

func (s *tracing) Sync(ctx context.Context, ns string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.ConfigMap",
	})
	defer func() {
		span.LogKV(
			"ns", ns,
			"err", err,
		)
		span.Finish()
	}()
	return s.next.Sync(ctx, ns)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
