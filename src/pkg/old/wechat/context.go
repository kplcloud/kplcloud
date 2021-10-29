/**
 * @Time : 2019-07-10 16:30
 * @Author : soupzhb@gmail.com
 * @File : http-options
 * @Software: GoLand
 */

package wechat

import (
	"context"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
)

const httpRequestContext = "http-request-context"

func httpToContest() kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, httpRequestContext, r)
	}
}
