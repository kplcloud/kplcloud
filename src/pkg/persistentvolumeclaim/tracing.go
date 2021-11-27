/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package persistentvolumeclaim

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

func (s *tracing) List(ctx context.Context, clusterId int64, storageClass, ns string, page, pageSize int) (resp []result, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.persistentvolumeclaim",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "storageClass", storageClass, "ns", ns, "page", page, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, storageClass, ns, page, pageSize)
}

func (s *tracing) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", opentracing.Tag{
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
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{
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

func (s *tracing) All(ctx context.Context, clusterId int64) (resp map[string]interface{}, err error) {
	panic("implement me")
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
