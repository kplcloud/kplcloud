/**
 * @Time : 8/19/21 1:32 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package configmap

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

func (s *tracing) List(ctx context.Context, clusterId int64, ns, name string, page, pageSize int) (res []types.ConfigMap, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.configmap",
	})
	defer func() {
		span.LogKV("clusterId", clusterId, "ns", ns, "name", name, "page", page, "pageSize", pageSize, "total", total, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, clusterId, ns, name, page, pageSize)
}

func (s *tracing) FindBy(ctx context.Context, clusterId int64, ns, name string) (res types.ConfigMap, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindBy", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.ConfigMap",
	})
	defer func() {
		span.LogKV(
			"clusterId", clusterId,
			"namespace", ns,
			"name", name,
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.FindBy(ctx, clusterId, ns, name)
}

func (s *tracing) Save(ctx context.Context, configMap *types.ConfigMap, data []types.Data) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.ConfigMap",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Save(ctx, configMap, data)
}

func (s *tracing) SaveData(ctx context.Context, configMapId int64, key, value string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SaveData", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.ConfigMap",
	})
	defer func() {
		span.LogKV(
			"configMapId", configMapId,
			"key", key,
			"value", value,
			"error", err,
		)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.SaveData(ctx, configMapId, key, value)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
