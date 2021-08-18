/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package nodes

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

func (s *tracing) List(ctx context.Context, clusterId int64, page, pageSize int) (res []nodeResult, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Nodes",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"page", page,
			"pageSize", pageSize,
			"err", err)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, page, pageSize)
}

func (s *tracing) Sync(ctx context.Context, clusterName string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Nodes",
	})
	defer func() {
		span.LogKV(
			"clusterName", clusterName,
			"err", err)
		span.Finish()
	}()
	return s.next.Sync(ctx, clusterName)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}