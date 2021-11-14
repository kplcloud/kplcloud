/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package registry

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

func (s *tracing) Secret(ctx context.Context, name string) (err error) {
	panic("implement me")
}

func (s *tracing) Update(ctx context.Context, name, host, username, password, remark string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.registry",
	})
	defer func() {
		span.LogKV("name", name, "host", host, "username", username, "password", password, "remark", remark, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, name, host, username, password, remark)
}

func (s *tracing) Delete(ctx context.Context, name string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.registry",
	})
	defer func() {
		span.LogKV("name", name, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, name)
}

func (s *tracing) Password(ctx context.Context, name string) (res string, err error) {
	panic("implement me")
}

func (s *tracing) Info(ctx context.Context, name string) (res result, err error) {
	panic("implement me")
}

func (s *tracing) List(ctx context.Context, query string, page, pageSize int) (res []result, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Registry",
	})
	defer func() {
		span.LogKV(
			"query", query,
			"page", page,
			"pageSize", pageSize,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, query, page, pageSize)
}

func (s *tracing) Create(ctx context.Context, name, host, username, password, remark string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Registry",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"host", host,
			"username", username,
			"password", password,
			"remark", remark,
			"err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, name, host, username, password, remark)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
