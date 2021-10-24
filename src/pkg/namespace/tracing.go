/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package namespace

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

func (s *tracing) List(ctx context.Context, clusterId int64, names []string, query string, page, pageSize int) (res []result, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", stdopentracing.Tag{
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
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", stdopentracing.Tag{
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
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", stdopentracing.Tag{
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
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", stdopentracing.Tag{
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
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", stdopentracing.Tag{
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

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
