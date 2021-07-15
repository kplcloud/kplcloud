/**
 * @Time : 3/10/21 3:36 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package sysrole

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/icowan/kit-admin/src/repository/types"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Delete(ctx context.Context, id int64) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysRole",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.Delete(ctx, id)
}

func (s *tracing) AddPermissions(ctx context.Context, role *types.SysRole, permissions []types.SysPermission) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AddPermissions", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysRole",
	})
	defer func() {
		span.LogKV(
			"error", err,
		)
		span.Finish()
	}()
	return s.next.AddPermissions(ctx, role, permissions)
}

func (s *tracing) Find(ctx context.Context, id int64) (res types.SysRole, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Find", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysRole",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"error", err)
		span.Finish()
	}()
	return s.next.Find(ctx, id)
}

func (s *tracing) FindByIds(ctx context.Context, ids []int64) (res []types.SysRole, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByIds", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysRole",
	})
	defer func() {
		span.LogKV(
			"error", err)
		span.Finish()
	}()
	return s.next.FindByIds(ctx, ids)
}

func (s *tracing) Save(ctx context.Context, data *types.SysRole) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysRole",
	})
	defer func() {
		span.LogKV(
			"error", err)
		span.Finish()
	}()
	return s.next.Save(ctx, data)
}

func (s *tracing) List(ctx context.Context, page, pageSize int) (res []types.SysRole, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysRole",
	})
	defer func() {
		span.LogKV(
			"error", err)
		span.Finish()
	}()
	return s.next.List(ctx, page, pageSize)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
