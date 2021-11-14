/**
 * @Time : 8/19/21 1:32 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package registry

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

func (s *tracing) Delete(ctx context.Context, id int64, call Call) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.registry",
	})
	defer func() {
		span.LogKV("id", id, "call", call, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, id, call)
}

func (s *tracing) SaveCall(ctx context.Context, reg *types.Registry, call Call) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SaveCall", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.registry",
	})
	defer func() {
		span.LogKV("reg", reg, "call", call, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.SaveCall(ctx, reg, call)
}

func (s *tracing) List(ctx context.Context, query string, page, pageSize int) (res []types.Registry, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Registry",
	})
	defer func() {
		span.LogKV(
			"query", query,
			"page", page,
			"pageSize", pageSize,
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, query, page, pageSize)
}

func (s *tracing) FindByNames(ctx context.Context, names []string) (res []types.Registry, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByNames", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Registry",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindByNames(ctx, names)
}

func (s *tracing) Save(ctx context.Context, reg *types.Registry) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Registry",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Save(ctx, reg)
}

func (s *tracing) FindByName(ctx context.Context, name string) (res types.Registry, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Registry",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindByName(ctx, name)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
