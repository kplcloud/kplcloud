/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package nodes

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

func (s *tracing) Cordon(ctx context.Context, clusterId int64, nodeName string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Cordon", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Nodes",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"nodeName", nodeName,
			"err", err)
		span.Finish()
	}()
	return s.next.Cordon(ctx, clusterId, nodeName)
}

func (s *tracing) Drain(ctx context.Context, clusterId int64, nodeName string, force bool) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Drain", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Nodes",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"nodeName", nodeName,
			"force", force,
			"err", err)
		span.Finish()
	}()
	return s.next.Drain(ctx, clusterId, nodeName, force)
}

func (s *tracing) Info(ctx context.Context, clusterId int64, nodeName string) (res infoResult, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Info", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Nodes",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"nodeName", nodeName,
			"err", err)
		span.Finish()
	}()
	return s.next.Info(ctx, clusterId, nodeName)
}

func (s *tracing) List(ctx context.Context, clusterId int64, query string, page, pageSize int) (res []nodeResult, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Nodes",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"page", page,
			"pageSize", pageSize,
			"query", query,
			"err", err)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, query, page, pageSize)
}

func (s *tracing) Sync(ctx context.Context, clusterName string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Sync", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Nodes",
	})
	defer func() {
		span.LogKV(
			"clusterName", clusterName,
			"err", err)
		span.Finish()
	}()
	return s.next.Sync(ctx, clusterName)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
