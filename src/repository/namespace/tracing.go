/**
 * @Time : 3/5/21 2:43 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package namespace

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

func (s *tracing) SaveCall(ctx context.Context, data *types.Namespace, call Call) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SaveCall", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Namespace",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.Finish()
	}()
	return s.next.SaveCall(ctx, data, call)
}

func (s *tracing) FindByName(ctx context.Context, clusterId int64, name string) (res types.Namespace, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Namespace",
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

func (s *tracing) FindByIds(ctx context.Context, ids []int64) (res []types.Namespace, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByIds", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Namespace",
	})
	defer func() {
		span.LogKV("error", err)
		span.Finish()
	}()
	return s.next.FindByIds(ctx, ids)
}

func (s *tracing) Save(ctx context.Context, data *types.Namespace) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Namespace",
	})
	defer func() {
		span.LogKV("error", err)
		span.Finish()
	}()
	return s.next.Save(ctx, data)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
