/**
 * @Time : 7/20/21 5:44 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package syssetting

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// 链路追踪中间件
type traceServer struct {
	next   Service
	tracer stdopentracing.Tracer
}

func (s *traceServer) Delete(ctx context.Context, id int64) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.SysSetting",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"err", err)
		span.Finish()
	}()
	return s.next.Delete(ctx, id)
}

func (s *traceServer) List(ctx context.Context, key string, page, pageSize int) (res []listResult, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.SysSetting",
	})
	defer func() {
		span.LogKV(
			"key", key,
			"err", err)
		span.Finish()
	}()
	return s.next.List(ctx, key, page, pageSize)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &traceServer{
			next:   next,
			tracer: otTracer,
		}
	}
}
