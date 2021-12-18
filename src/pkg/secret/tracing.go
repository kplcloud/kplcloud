/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package secret

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

func (s *tracing) Add(ctx context.Context, clusterId int64, namespace, name string) (err error) {
	panic("implement me")
}

func (s *tracing) List(ctx context.Context, clusterId int64, namespace, name string, page, pageSize int) (res []secretResult, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.secret",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "namespace", namespace, "name", name, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, namespace, name, page, pageSize)
}

func (s *tracing) Delete(ctx context.Context, clusterId int64, ns, name string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
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
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ImageSecret", opentracing.Tag{
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
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", opentracing.Tag{
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

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
