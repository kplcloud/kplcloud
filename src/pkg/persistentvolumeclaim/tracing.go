/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package persistentvolumeclaim

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

func (s *tracing) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.PersistentVolumeClaim",
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

func (s *tracing) Get(ctx context.Context, clusterId int64, ns, name string) (rs interface{}, err error) {
	panic("implement me")
}

func (s *tracing) Delete(ctx context.Context, clusterId int64, ns, name string) (err error) {
	panic("implement me")
}

func (s *tracing) Create(ctx context.Context, clusterId int64, ns, name, storage, storageClassName string, accessModes []string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.PersistentVolumeClaim",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"ns", ns,
			"name", name, "storage", storage,
			"storageClassName", storageClassName,
			"accessModes", accessModes,
			"err", err)
		span.Finish()
	}()
	return s.next.Create(ctx, clusterId, ns, name, storage, storageClassName, accessModes)
}

func (s *tracing) List(ctx context.Context, clusterId int64, ns string, page, pageSize int) (resp map[string]interface{}, err error) {
	panic("implement me")
}

func (s *tracing) All(ctx context.Context, clusterId int64) (resp map[string]interface{}, err error) {
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
