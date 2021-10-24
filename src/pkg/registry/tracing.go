/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package registry

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

func (s *tracing) List(ctx context.Context, query string, page, pageSize int) (res []result, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", stdopentracing.Tag{
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
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", stdopentracing.Tag{
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

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
