/**
 * @Time : 3/9/21 5:58 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package template

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

func (s *tracing) Add(ctx context.Context, kind, alias, rules, content string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Add", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Template",
	})
	defer func() {
		span.LogKV(
			"kind", kind,
			"alias", alias,
			"err", err)
		span.Finish()
	}()
	return s.next.Add(ctx, kind, alias, rules, content)
}

func (s *tracing) Delete(ctx context.Context, kind string) (err error) {
	panic("implement me")
}

func (s *tracing) Update(ctx context.Context, kind, alias, rules, content string) (err error) {
	panic("implement me")
}

func (s *tracing) List(ctx context.Context, searchValue string, page, pageSize int64) (res []interface{}, total int, err error) {
	panic("implement me")
}

func (s *tracing) Info(ctx context.Context, kind string) (res interface{}, err error) {
	panic("implement me")
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
