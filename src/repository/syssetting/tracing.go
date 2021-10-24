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

func (s *tracing) FindById(ctx context.Context, id int64) (res types.SysSetting, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindById", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.FindById(ctx, id)
}

func (s *tracing) List(ctx context.Context, key string, page, pageSize int) (res []types.SysSetting, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV(
			"key", key,
			"page", page,
			"pageSize", pageSize,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.List(ctx, key, page, pageSize)
}

func (s *tracing) FindAll(ctx context.Context) (res []types.SysSetting, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindAll", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.Finish()
	}()
	return s.next.FindAll(ctx)
}

func (s *tracing) Delete(ctx context.Context, section, key string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV(
			"key", key,
			"section", section,
			"error", err)
		span.Finish()
	}()
	return s.next.Delete(ctx, section, key)
}

func (s *tracing) Update(ctx context.Context, data *types.SysSetting) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV(
			"error", err)
		span.Finish()
	}()
	return s.next.Update(ctx, data)
}

func (s *tracing) Find(ctx context.Context, section, key string) (res types.SysSetting, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Find", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV(
			"section", section,
			"key", key,
			"error", err)
		span.Finish()
	}()
	return s.next.Find(ctx, section, key)
}

func (s *tracing) Add(ctx context.Context, section, key, val, desc string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Add", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Setting",
	})
	defer func() {
		span.LogKV(
			"section", section,
			"key", key,
			"val", val,
			"desc", desc, "error", err)
		span.Finish()
	}()
	return s.next.Add(ctx, section, key, val, desc)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
