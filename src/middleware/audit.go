/**
 * @Time : 2021/9/9 10:45 AM
 * @Author : solacowa@gmail.com
 * @File : audit
 * @Software: GoLand
 */

package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func AuditMiddleware(store repository.Repository, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			begin := time.Now()
			method := ctx.Value(kithttp.ContextKeyRequestMethod).(string)
			if !strings.EqualFold(method, http.MethodPost) &&
				!strings.EqualFold(method, http.MethodDelete) &&
				!strings.EqualFold(method, http.MethodPatch) &&
				!strings.EqualFold(method, http.MethodPut) {
				return next(ctx, request)
			}

			u, _ := url.Parse(ctx.Value(kithttp.ContextKeyRequestURI).(string))

			if tracer != nil {
				var span opentracing.Span
				span, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "AuditMiddleware", opentracing.Tag{
					Key:   string(ext.Component),
					Value: "Middleware",
				}, opentracing.Tag{
					Key:   "Method",
					Value: method,
				}, opentracing.Tag{
					Key:   "URI",
					Value: u.RawPath,
				})
				defer func() {
					span.LogKV(
						"err", err,
					)
					span.Finish()
				}()
			}

			defer func() {
				go func(ctx context.Context, begin time.Time, request interface{}, response interface{}, err error) {
					ns, _ := ctx.Value(ContextKeyNamespaceName).(string)
					svc, _ := ctx.Value(ContextKeyName).(string)
					clusterId, _ := ctx.Value(ContextKeyClusterId).(int64)
					userId, _ := ctx.Value(ContextUserId).(int64)
					permId, _ := ctx.Value(ContextPermissionId).(int64)
					var remark string
					var status = types.AuditStatusSuccess
					if err != nil {
						status = types.AuditStatusFailed
						remark = err.Error()
					}
					req, _ := json.Marshal(request)
					res, _ := json.Marshal(response)
					headers, _ := json.Marshal(ctx.Value(kithttp.ContextKeyResponseHeaders))
					u, _ := url.Parse(ctx.Value(kithttp.ContextKeyRequestURI).(string))
					if e := store.Audit(ctx).Save(ctx, &types.Audit{
						ClusterId:    clusterId,
						Namespace:    ns,
						Name:         svc,
						UserId:       userId,
						PermissionId: permId,
						Request:      string(req),
						Response:     string(res),
						Headers:      string(headers),
						TimeSince:    time.Since(begin).String(),
						Status:       status,
						Remark:       remark,
						Url:          u.String(),
						TraceId:      ctx.Value("traceId").(string),
					}); e != nil {
						log.Println(errors.Wrap(err, "middleware.store.Audit.Save"))
					}
				}(ctx, begin, request, response, err)
			}()

			return next(ctx, request)
		}
	}
}
