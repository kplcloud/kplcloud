/**
 * @Time : 3/5/21 2:24 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package sysuser

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/kplcloud/kplcloud/src/repository/types"
)

type Middleware func(next Service) Service

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) FindByRoleId(ctx context.Context, roleId int64, page, pageSize int) (res []types.SysUser, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByRoleId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysUser",
	})
	defer func() {
		span.LogKV(
			"roleId", roleId,
			"page", page,
			"pageSize", pageSize,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.FindByRoleId(ctx, roleId, page, pageSize)
}

func (s *tracing) FindByIds(ctx context.Context, ids []int64, preload ...string) (res []types.SysUser, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByIds", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysUser",
	})
	defer func() {
		span.LogKV(
			"ids", ids,
			"error", err,
		)
		span.Finish()
	}()
	return s.next.FindByIds(ctx, ids, preload...)
}

func (s *tracing) AddRoles(ctx context.Context, user *types.SysUser, roles []types.SysRole) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "AddRoles", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysUser",
	})
	defer func() {
		span.LogKV(
			"error", err)
		span.Finish()
	}()
	return s.next.AddRoles(ctx, user, roles)
}

func (s *tracing) Find(ctx context.Context, userId int64, preload ...string) (res types.SysUser, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Find", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysUser",
	})
	defer func() {
		span.LogKV(
			"userId", userId,
			"preload", preload,
			"error", err)
		span.Finish()
	}()
	return s.next.Find(ctx, userId, preload...)
}

func (s *tracing) List(ctx context.Context, namespaceIds []int64, email string, page, pageSize int) (res []types.SysUser, total int, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysUser",
	})
	defer func() {
		span.LogKV(
			"namespaceIds", namespaceIds,
			"email", email,
			"error", err)
		span.Finish()
	}()
	return s.next.List(ctx, namespaceIds, email, page, pageSize)
}

func (s *tracing) Save(ctx context.Context, user *types.SysUser) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Save", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysUser",
	})
	defer func() {
		span.LogKV("error", err)
		span.Finish()
	}()
	return s.next.Save(ctx, user)
}

func (s *tracing) FindByEmail(ctx context.Context, email string) (res types.SysUser, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "FindByEmail", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.SysUser",
	})
	defer func() {
		span.LogKV("email", email, "error", err)
		span.Finish()
	}()
	return s.next.FindByEmail(ctx, email)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
