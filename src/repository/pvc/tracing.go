/**
 * @Time : 2021/11/25 4:30 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package pvc

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

func (s *tracing) SavePv(ctx context.Context, pv *types.PersistentVolume, call Call) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SavePv", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.pvc",
	})
	defer func() {
		span.LogKV("pv", pv, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.SavePv(ctx, pv, call)
}

func (s *tracing) Delete(ctx context.Context, pvcId int64, call ...Call) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.pvc",
	})
	defer func() {
		span.LogKV("pvcId", pvcId, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, pvcId, call...)
}

func (s *tracing) FindByName(ctx context.Context, clusterId int64, ns, name string) (res types.PersistentVolumeClaim, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByName", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.pvc",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "ns", ns, "name", name, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindByName(ctx, clusterId, ns, name)
}

func (s *tracing) List(ctx context.Context, clusterId int64, storageClassIds []int64, ns, name string, page, pageSize int) (res []types.PersistentVolumeClaim, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.pvc",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "storageClassIds", storageClassIds, "page", page, "namespace", ns, "name", name, "pageSize", pageSize, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, storageClassIds, ns, name, page, pageSize)
}

func (s *tracing) Save(ctx context.Context, app *types.PersistentVolumeClaim, call Call) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.application",
	})
	defer func() {
		span.LogKV("app", app, "call", call, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Save(ctx, app, call)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
