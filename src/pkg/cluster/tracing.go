/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package cluster

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

func (s *tracing) List(ctx context.Context, name string, page, pageSize int) (res []listResult, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Cluster",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"err", err,
		)
		span.Finish()
	}()
	return s.next.List(ctx, name, page, pageSize)
}

func (s *tracing) SyncRoles(ctx context.Context, clusterId int64) (err error) {
	panic("implement me")
}

func (s *tracing) Add(ctx context.Context, name, alias, data string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Add", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Cluster",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"alias", alias,
			"err", err)
		span.Finish()
	}()
	return s.next.Add(ctx, name, alias, data)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
