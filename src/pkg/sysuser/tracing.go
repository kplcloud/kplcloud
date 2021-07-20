/**
 * @Time : 3/9/21 5:58 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package sysuser

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

func (s *traceServer) Locked(ctx context.Context, userId int64) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Locked", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysUser",
	})
	defer func() {
		span.LogKV(
			"userId", userId,
			"err", err)
		span.Finish()
	}()
	return s.next.Locked(ctx, userId)
}

func (s *traceServer) Delete(ctx context.Context, userId int64, unscoped bool) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysUser",
	})
	defer func() {
		span.LogKV(
			"userId", userId,
			"unscoped", unscoped,
			"err", err)
		span.Finish()
	}()
	return s.next.Delete(ctx, userId, unscoped)
}

func (s *traceServer) Update(ctx context.Context, userId int64, username, email string, locked bool, namespaceIds, roleIds []int64) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysUser",
	})
	defer func() {
		span.LogKV(
			"userId", userId,
			"username", username,
			"email", email,
			"locked", locked,
			"err", err)
		span.Finish()
	}()
	return s.next.Update(ctx, userId, username, email, locked, namespaceIds, roleIds)
}

func (s *traceServer) Add(ctx context.Context, username, email string, locked bool, namespaceIds, roleIds []int64) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Add", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysUser",
	})
	defer func() {
		span.LogKV(
			"username", username,
			"email", email,
			"locked", locked,
			"err", err)
		span.Finish()
	}()
	return s.next.Add(ctx, username, email, locked, namespaceIds, roleIds)
}

func (s *traceServer) List(ctx context.Context, email string, page, pageSize int) (res []listResult, total int, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "SysUser",
	})
	defer func() {
		span.LogKV(
			"email", email,
			"err", err)
		span.Finish()
	}()
	return s.next.List(ctx, email, page, pageSize)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &traceServer{
			next:   next,
			tracer: otTracer,
		}
	}
}
