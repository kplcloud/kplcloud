/**
 * @Time : 2019-07-29 15:20
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package monitor

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

func makeQueryNetworkEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.QueryNetwork(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeOpsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Ops(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeMetricsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Metrics(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
