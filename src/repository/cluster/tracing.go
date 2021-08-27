/**
 * @Time : 3/10/21 3:36 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package cluster

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

func (s *tracing) List(ctx context.Context, name string, status int, page, pageSize int) (res []types.Cluster, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Cluster",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"status", status,
			"page", page,
			"pageSize", pageSize,
			"total", total,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.List(ctx, name, status, page, pageSize)
}

func (s *tracing) SaveRole(ctx context.Context, clusterRole *types.ClusterRole, roles []types.PolicyRule) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SaveRole", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Cluster",
	})
	defer func() {
		span.LogKV(
			"clusterRole", clusterRole,
			"roles", roles,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.SaveRole(ctx, clusterRole, roles)
}

func (s *tracing) FindAll(ctx context.Context, status int) (res []types.Cluster, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindAll", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Cluster",
	})
	defer func() {
		span.LogKV(
			"status", status,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.FindAll(ctx, status)
}

func (s *tracing) FindByName(ctx context.Context, name string) (res types.Cluster, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Cluster",
	})
	defer func() {
		span.LogKV(
			"name", name,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.FindByName(ctx, name)
}

func (s *tracing) Save(ctx context.Context, data *types.Cluster, calls ...Call) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Cluster",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.Finish()
	}()
	return s.next.Save(ctx, data, calls...)
}

func (s *tracing) Delete(ctx context.Context, id int64, unscoped bool) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.Cluster",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"unscoped", unscoped,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.Delete(ctx, id, unscoped)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
