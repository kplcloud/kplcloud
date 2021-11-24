/**
 * @Time : 5/25/21 3:54 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package syspermission

import (
	"context"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Save(ctx context.Context, data *types.SysPermission) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.syspermission",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Save(ctx, data)
}

func (s *tracing) Delete(ctx context.Context, id int64) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.syspermission",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, id)
}

func (s *tracing) FindByIds(ctx context.Context, ids []int64) (res []types.SysPermission, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Find", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.syspermission",
	})
	defer func() {
		span.LogKV("ids", ids, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindByIds(ctx, ids)
}

func (s *tracing) Find(ctx context.Context, id int64) (res types.SysPermission, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByIds", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.syspermission",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Find(ctx, id)
}

func (s *tracing) FindAll(ctx context.Context) (res []types.SysPermission, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindAll", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.syspermission",
	})
	defer func() {
		span.LogKV("err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindAll(ctx)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
