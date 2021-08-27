/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package secret

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

func (s *tracing) Delete(ctx context.Context, clusterId int64, ns, name string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Secret",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"ns", ns,
			"name", name,
			"err", err,
		)
		span.Finish()
	}()
	return s.next.Delete(ctx, clusterId, ns, name)
}

func (s *tracing) ImageSecret(ctx context.Context, clusterId int64, ns, name, host, username, password string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ImageSecret", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Secret",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"ns", ns,
			"name", name,
			"host", host,
			"username", username,
			"password", password,
			"err", err,
		)
		span.Finish()
	}()
	return s.next.ImageSecret(ctx, clusterId, ns, name, host, username, password)
}

func (s *tracing) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Secret",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"ns", ns,
			"err", err)
		span.Finish()
	}()
	return s.next.Sync(ctx, clusterId, ns)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
