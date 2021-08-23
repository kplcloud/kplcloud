/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package storageclass

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
