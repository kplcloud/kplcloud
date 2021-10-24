/**
 * @Time: 2020/12/28 23:23
 * @Author: solacowa@gmail.com
 * @File: tracing
 * @Software: GoLand
 */

package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

func TracingServerBefore(tracer stdopentracing.Tracer) kithttp.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		if tracer == nil {
			return ctx
		}
		reqPath := ctx.Value(kithttp.ContextKeyRequestURI).(string)
		u, _ := url.Parse(reqPath)
		span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, tracer, u.Path, stdopentracing.Tag{
			Key:   string(ext.Component),
			Value: "ServerBefore",
		}, stdopentracing.Tag{
			Key:   string(ext.HTTPMethod),
			Value: ctx.Value(kithttp.ContextKeyRequestMethod),
		}, stdopentracing.Tag{
			Key:   string(ext.HTTPUrl),
			Value: ctx.Value(kithttp.ContextKeyRequestURI),
		}, stdopentracing.Tag{
			Key:   "user_agent",
			Value: ctx.Value(kithttp.ContextKeyRequestUserAgent),
		}, stdopentracing.Tag{
			Key:   "x-request-id",
			Value: ctx.Value(kithttp.ContextKeyRequestXRequestID),
		}, stdopentracing.Tag{
			Key:   "Authorization",
			Value: ctx.Value(kithttp.ContextKeyRequestAuthorization),
		}, stdopentracing.Tag{
			Key:   "referer",
			Value: ctx.Value(kithttp.ContextKeyRequestReferer),
		}, stdopentracing.Tag{
			Key:   "x-forwarded-for",
			Value: ctx.Value(kithttp.ContextKeyRequestXForwardedFor),
		})
		traceId := span.Context().(jaeger.SpanContext).TraceID().String()
		ctx = context.WithValue(ctx, "TraceId", traceId)
		ctx = context.WithValue(ctx, "traceId", traceId)
		span = span.SetTag("TraceId", span.Context().(jaeger.SpanContext).TraceID().String())

		defer func() {
			b, _ := json.Marshal(request)
			span.LogKV("request", string(b))
			span.Finish()
		}()
		return ctx
	}
}

func TracingMiddleware(tracer stdopentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if tracer == nil {
				return next(ctx, request)
			}
			span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, tracer, "TracingMiddleware", stdopentracing.Tag{
				Key:   string(ext.Component),
				Value: "Middleware",
			})
			defer func() {
				b, _ := json.Marshal(request)
				res, _ := json.Marshal(response)
				span.LogKV("request", string(b), "response", string(res), "err", err)
				span.Finish()
			}()
			return next(ctx, request)
		}
	}
}
