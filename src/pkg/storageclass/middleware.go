/**
 * @Time: 2021/11/14 22:24
 * @Author: solacowa@gmail.com
 * @File: middleware
 * @Software: GoLand
 */

package storageclass

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	captcha "github.com/icowan/kit-captcha"
)

type contextKey string

const (
	contextKeyStorageClassName contextKey = "storageclass-context-key-name"
)

// 获取名称
func getStorageNameMiddleware(captcha captcha.Service) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return next(ctx, request)
		}
	}
}
