/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package storageclass

import (
	"context"
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// 链路追踪中间件
type tracing struct {
	next   Service
	tracer stdopentracing.Tracer
}

func (s *tracing) List(ctx context.Context, clusterId int64, page, pageSize int) (res []listResult, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.storageclass",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "page", page, "pageSize", pageSize, "total", total, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, page, pageSize)
}

func (s *tracing) Delete(ctx context.Context, clusterId int64, storageName string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.storageclass",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "storageName", storageName, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, clusterId, storageName)
}

func (s *tracing) Update(ctx context.Context, clusterId int64, storageName, provisioner string, reclaimPolicy *v1.PersistentVolumeReclaimPolicy, volumeBindingMode *storagev1.VolumeBindingMode) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.storageclass",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "storageName", storageName, "provisioner", provisioner, "reclaimPolicy", reclaimPolicy, "volumeBindingMode", volumeBindingMode, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, clusterId, storageName, provisioner, reclaimPolicy, volumeBindingMode)
}

func (s *tracing) Create(ctx context.Context, clusterId int64, ns, name, provisioner string, reclaimPolicy *v1.PersistentVolumeReclaimPolicy, volumeBindingMode *storagev1.VolumeBindingMode) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SyncPv", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.StorageClass",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"ns", ns,
			"name", name,
			"provisioner", provisioner,
			"reclaimPolicy", reclaimPolicy,
			"volumeBindingMode", volumeBindingMode,
			"err", err)
		span.Finish()
	}()
	return s.next.Create(ctx, clusterId, ns, name, provisioner, reclaimPolicy, volumeBindingMode)
}

func (s *tracing) CreateProvisioner(ctx context.Context, clusterId int64) (err error) {
	panic("implement me")
}

func (s *tracing) SyncPv(ctx context.Context, clusterId int64, storageName string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SyncPv", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.StorageClass",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"storageName", storageName,
			"err", err)
		span.Finish()
	}()
	return s.next.SyncPv(ctx, clusterId, storageName)
}

func (s *tracing) SyncPvc(ctx context.Context, clusterId int64, ns string, storageName string) (err error) {
	panic("implement me")
}

func (s *tracing) Sync(ctx context.Context, clusterId int64) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.StorageClass",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"err", err)
		span.Finish()
	}()
	return s.next.Sync(ctx, clusterId)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
