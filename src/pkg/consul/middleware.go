/**
 * @Time : 2019/7/19 5:15 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package consul

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/icowan/config"
)

var (
	ErrConsulVisit = errors.New("抱歉,您未开启consul模块,请在app.cfg把consul_kv设置为true")
)

func consulVisitMiddleware(logger log.Logger, cf *config.Config) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if cf.GetString("server", "consul_kv") != "true" {
				_ = logger.Log("consulVisitMiddleware", "consul refused")
				return nil, ErrConsulVisit
			}
			return next(ctx, request)
		}
	}
}
