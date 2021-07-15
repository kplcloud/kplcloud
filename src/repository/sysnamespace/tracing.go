/**
 * @Time : 3/5/21 2:43 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package sysnamespace

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/icowan/kit-admin/src/repository/types"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) FindByIds(ctx context.Context, ids []int64) (res []types.SysNamespace, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByIds", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysNamespace",
	})
	defer func() {
		span.LogKV("error", err)
		span.Finish()
	}()
	return s.next.FindByIds(ctx, ids)
}

func (s *tracing) Save(ctx context.Context, data *types.SysNamespace) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysNamespace",
	})
	defer func() {
		span.LogKV("error", err)
		span.Finish()
	}()
	return s.next.Save(ctx, data)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
