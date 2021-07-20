/**
 * @Time : 3/10/21 3:41 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package sysrole

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// 链路追踪中间件
type traceServer struct {
	next   Service
	tracer stdopentracing.Tracer
}

func (s *traceServer) Delete(ctx context.Context, id int64) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysRole",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"err", err)
		span.Finish()
	}()
	return s.next.Delete(ctx, id)
}

func (s *traceServer) Update(ctx context.Context, id int64, alias, name, description string, enabled bool) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysRole",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"name", name,
			"description", description,
			"enabled", enabled,
			"err", err)
		span.Finish()
	}()
	return s.next.Update(ctx, id, alias, name, description, enabled)
}

func (s *traceServer) User(ctx context.Context, id int64, userId []int64) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "User", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysRole",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"userId", userId,
			"err", err)
		span.Finish()
	}()
	return s.next.User(ctx, id, userId)
}

func (s *traceServer) Permissions(ctx context.Context, id int64, page, pageSize int) (res []listResult, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Permissions", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysRole",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"page", page,
			"pageSize", pageSize,
			"err", err)
		span.Finish()
	}()
	return s.next.Permissions(ctx, id, page, pageSize)
}

func (s *traceServer) Permission(ctx context.Context, id int64, permIds []int64) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Permission", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysRole",
	})
	defer func() {
		span.LogKV(
			"id", id,
			"permIds", permIds,
			"err", err)
		span.Finish()
	}()
	return s.next.Permission(ctx, id, permIds)
}

//func (s *traceServer) Users(ctx context.Context, id int64, page, pageSize int) (res []listResult, total int, err error) {
//	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Users", stdopentracing.Tag{
//		Key:   string(ext.Component),
//		Value: "SysRole",
//	})
//	defer func() {
//		span.LogKV(
//			"id", id,
//			"page", page,
//			"pageSize", pageSize,
//			"err", err)
//		span.Finish()
//	}()
//	return s.next.Users(ctx, id, page, pageSize)
//}

func (s *traceServer) Add(ctx context.Context, alias, name, description string, enabled bool) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Add", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysRole",
	})
	defer func() {
		span.LogKV(
			"alias", alias,
			"name", name,
			"description", description,
			"enabled", enabled,
			"err", err)
		span.Finish()
	}()
	return s.next.Add(ctx, alias, name, description, enabled)
}

func (s *traceServer) List(ctx context.Context, page, pageSize int) (res []listResult, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysRole",
	})
	defer func() {
		span.LogKV(
			"err", err)
		span.Finish()
	}()
	return s.next.List(ctx, page, pageSize)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &traceServer{
			next:   next,
			tracer: otTracer,
		}
	}
}
