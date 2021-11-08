/**
 * @Time : 3/10/21 3:36 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package storageclass

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

func (s *tracing) List(ctx context.Context, clusterId int64, page, pageSize int) (res []types.StorageClass, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.storageclass",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, page, pageSize)
}

func (s *tracing) FindName(ctx context.Context, clusterId int64, name string) (res types.StorageClass, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.StorageClass",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"name", name,
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindName(ctx, clusterId, name)
}

func (s *tracing) Find(ctx context.Context, id int64) (res types.StorageClass, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Find", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.StorageClass",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Find(ctx, id)
}

func (s *tracing) FirstInsert(ctx context.Context, data *types.StorageClass) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FirstInsert", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.StorageClass",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FirstInsert(ctx, data)
}

func (s *tracing) Save(ctx context.Context, data *types.StorageClass, call Call) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.StorageClass",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Save(ctx, data, call)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
