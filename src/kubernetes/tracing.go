/**
 * @Time : 8/11/21 4:28 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package kubernetes

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// 链路追踪中间件
type tracing struct {
	next   K8sClient
	tracer stdopentracing.Tracer
}

func (s *tracing) Do(ctx context.Context) *kubernetes.Clientset {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Do", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "Kubernetes",
	})
	defer func() {
		span.Finish()
	}()
	return s.next.Do(ctx)
}

func (s *tracing) Config(ctx context.Context) *rest.Config {
	return s.next.Config(ctx)
}

func (s *tracing) Reload(ctx context.Context) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Reload", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "Kubernetes",
	})
	defer func() {
		span.LogKV(
			"err", err)
		span.Finish()
	}()
	return s.next.Reload(ctx)
}

func (s *tracing) Connect(ctx context.Context, name, configData string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Connect", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "Kubernetes",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"err", err)
		span.Finish()
	}()
	return s.next.Connect(ctx, name, configData)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next K8sClient) K8sClient {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
