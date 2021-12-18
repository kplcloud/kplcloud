/**
 * @Time : 2021/12/13 5:56 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package hpa

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

func (s *tracing) List(ctx context.Context, clusterId int64, namespace string, names []string, page, pageSize int) (res []types.HorizontalPodAutoscaler, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.hpa",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "namespace", namespace, "names", names, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, namespace, names, page, pageSize)
}

func (s *tracing) Save(ctx context.Context, dat *types.HorizontalPodAutoscaler, calls ...Call) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.hpa",
	})
	defer func() {
		span.LogKV("dat", dat, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Save(ctx, dat, calls...)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
