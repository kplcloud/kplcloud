/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package namespace

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// 链路追踪中间件
type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Info(ctx context.Context, clusterId int64, name string) (res result, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Info", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.namespace",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "name", name, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Info(ctx, clusterId, name)
}

func (s *tracing) List(ctx context.Context, clusterId int64, names []string, query string, page, pageSize int) (res []result, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Namespace",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"names", names,
			"query", query,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, names, query, page, pageSize)

}

func (s *tracing) Delete(ctx context.Context, clusterId int64, name string, force bool) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Namespace",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"name", name,
			"force", force,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, clusterId, name, force)
}

func (s *tracing) Update(ctx context.Context, clusterId int64, name string, alias, remark, status string, imageSecrets []string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Namespace",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"name", name,
			"alias", alias,
			"remark", remark,
			"status", status,
			"imageSecrets", imageSecrets,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, clusterId, name, alias, remark, status, imageSecrets)
}

func (s *tracing) Create(ctx context.Context, clusterId int64, name, alias, remark string, imageSecrets []string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Namespace",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"name", name, "alias", alias, "remark", remark,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, clusterId, name, alias, remark, imageSecrets)
}

func (s *tracing) Sync(ctx context.Context, clusterId int64) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Namespace",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Sync(ctx, clusterId)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
