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

func (s *tracing) Info(ctx context.Context, name string) (res infoResult, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Info", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Cluster",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"err", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Info(ctx, name)
}

func (s *tracing) Update(ctx context.Context, name, alias, data, remark string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Cluster",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"alias", alias,
			"remark", remark,
			"err", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, name, alias, data, remark)
}

func (s *tracing) Delete(ctx context.Context, name string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Cluster",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"err", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, name)
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
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, name, page, pageSize)
}

func (s *tracing) SyncRoles(ctx context.Context, clusterId int64) (err error) {
	panic("implement me")
}

func (s *tracing) Add(ctx context.Context, name, alias, data, remark string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Add", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Cluster",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"alias", alias,
			"remark", remark,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Add(ctx, name, alias, data, remark)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
