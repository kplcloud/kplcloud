/**
 * @Time : 8/13/21 10:01 AM
 * @Author : solacowa@gmail.com
 * @File : cluster
 * @Software: GoLand
 */

package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitcache "github.com/icowan/kit-cache"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"strings"
)

func ClusterMiddleware(store repository.Repository, cache kitcache.Service, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if tracer != nil {
				var span opentracing.Span
				span, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "ClusterMiddleware", opentracing.Tag{
					Key:   string(ext.Component),
					Value: "Middleware",
				}, opentracing.Tag{
					Key:   "Cluster",
					Value: ctx.Value(ContextKeyClusterName),
				})
				defer func() {
					span.LogKV("err", err)
					span.Finish()
				}()
			}

			clusterName, ok := ctx.Value(ContextKeyClusterName).(string)
			if !ok {
				return nil, encode.ErrClusterParams.Error()
			}

			var clusters []string
			rs, err := cache.Get(ctx, fmt.Sprintf("user:%d:clusters", ctx.Value(ContextUserId).(int64)), clusters)
			if err != nil {
				return nil, encode.ErrClusterNotfound.Wrap(err)
			}
			// 判断用户是否有权限访问该集群
			var pass bool
			if err = json.Unmarshal([]byte(rs), &clusters); err != nil {
				return nil, encode.ErrClusterNotfound.Wrap(err)
			}
			for _, v := range clusters {
				if strings.EqualFold(v, clusterName) {
					pass = true
					break
				}
			}
			if !pass {
				return nil, encode.ErrClusterNotPermission.Error()
			}

			cluster, err := store.Cluster(ctx).FindByName(ctx, clusterName)
			if err != nil {
				return nil, encode.ErrClusterNotfound.Error()
			}

			ctx = context.WithValue(ctx, ContextKeyClusterName, clusterName)
			ctx = context.WithValue(ctx, ContextKeyClusterId, cluster.Id)

			return next(ctx, request)
		}
	}
}
