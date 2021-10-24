/**
 * @Time: 2021/10/24 14:38
 * @Author: solacowa@gmail.com
 * @File: namespace
 * @Software: GoLand
 */

package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func NamespaceMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var namespace, name string
			namespace, _ = ctx.Value(ContextKeyNamespaceName).(string)
			name, _ = ctx.Value(ContextKeyName).(string)

			var permission bool
			var namespaces []string
			if ctx.Value(ContextKeyNamespaceList) != nil {
				namespaces = ctx.Value(ContextKeyNamespaceList).([]string)
			}

			for _, v := range namespaces {
				if v == namespace {
					permission = true
					break
				}
			}

			if !permission {
				_ = level.Error(logger).Log("name", name, "namespace", namespace, "permission", permission)
				return nil, ErrorASD
			}

			return next(ctx, request)
		}
	}
}
