/**
 * @Time : 8/13/21 10:01 AM
 * @Author : solacowa@gmail.com
 * @File : cluster
 * @Software: GoLand
 */

package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/repository"
)

func ClusterMiddleware(store repository.Repository) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			clusterName, ok := ctx.Value(ContextKeyClusterName).(string)
			if !ok {
				return nil, encode.ErrClusterParams.Error()
			}

			cluster, err := store.Cluster(ctx).FindByName(ctx, clusterName)
			if err != nil {
				return nil, encode.ErrClusterNotfound.Error()
			}

			// TODO 判断用户是否有权限访问该集群

			ctx = context.WithValue(ctx, ContextKeyClusterName, clusterName)
			ctx = context.WithValue(ctx, ContextKeyClusterId, cluster.Id)

			return next(ctx, request)
		}
	}
}
