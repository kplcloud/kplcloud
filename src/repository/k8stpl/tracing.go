/**
 * @Time : 3/5/21 2:43 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package k8stpl

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

func (s *tracing) Save(ctx context.Context, tpl *types.K8sTemplate) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Template",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.Finish()
	}()
	return s.next.Save(ctx, tpl)
}

func (s *tracing) EncodeTemplate(ctx context.Context, kind types.Kind, paramContent map[string]interface{}, data interface{}) (tpl []byte, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "EncodeTemplate", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Template",
	})
	defer func() {
		span.LogKV(
			"kind", kind,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.EncodeTemplate(ctx, kind, paramContent, data)
}

func (s *tracing) FindByKind(ctx context.Context, kind types.Kind) (tpl types.K8sTemplate, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByKind", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Template",
	})
	defer func() {
		span.LogKV(
			"kind", kind,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.FindByKind(ctx, kind)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
