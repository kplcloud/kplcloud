/**
 * @Time : 3/10/21 3:36 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package nodes

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Save(ctx context.Context, data *types.Nodes) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Nodes",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.Finish()
	}()
	return s.next.Save(ctx, data)
}

func (s *tracing) FindByName(ctx context.Context, clusterId int64, name string) (res types.Nodes, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Nodes",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"name", name,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.FindByName(ctx, clusterId, name)
}

func (s *tracing) List(ctx context.Context, clusterId int64, page, pageSize int) (res []types.Nodes, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Nodes",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, page, pageSize)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
