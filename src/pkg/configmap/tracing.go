/**
 * @Time: 2021/8/18 23:13
 * @Author: solacowa@gmail.com
 * @File: tracing
 * @Software: GoLand
 */

package configmap

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// 链路追踪中间件
type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) List(ctx context.Context, clusterId int64, ns, name string, page, pageSize int) (res []configMapResult, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.configmap",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "ns", ns, "name", name, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, ns, name, page, pageSize)
}

func (s *tracing) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.ConfigMap",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"ns", ns,
			"err", err,
		)
		span.Finish()
	}()
	return s.next.Sync(ctx, clusterId, ns)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
