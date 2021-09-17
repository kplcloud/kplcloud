/**
 * @Time : 2021/9/17 3:28 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package account

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

func (s *tracing) UserInfo(ctx context.Context, userId int64) (res userInfoResult, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserInfo", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Account",
	})
	defer func() {
		span.LogKV(
			"userId", userId,
			"err", err)
		span.Finish()
	}()
	return s.next.UserInfo(ctx, userId)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
