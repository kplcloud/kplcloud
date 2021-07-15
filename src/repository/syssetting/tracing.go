/**
 * @Time : 2020/12/31 4:56 PM
 * @Author : solacowa@gmail.com
 * @File : setting_middleware
 * @Software: GoLand
 */

package syssetting

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(next Service) Service

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Delete(ctx context.Context, key string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV("key", key, "error", err)
		span.Finish()
	}()
	return s.next.Delete(ctx, key)
}

func (s *tracing) Update(ctx context.Context, data *types.SysSetting) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV("error", err)
		span.Finish()
	}()
	return s.next.Update(ctx, data)
}

func (s *tracing) Find(ctx context.Context, key string) (res types.SysSetting, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Find", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV("key", key, "error", err)
		span.Finish()
	}()
	return s.next.Find(ctx, key)
}

func (s *tracing) Add(ctx context.Context, key, val, desc string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Add", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV("key", key, "val", val, "desc", desc, "error", err)
		span.Finish()
	}()
	return s.next.Add(ctx, key, val, desc)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
