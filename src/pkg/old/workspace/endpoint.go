/**
 * @Time : 2019-07-25 15:11
 * @Author : soupzhb@gmail.com
 * @File : endpoint.go
 * @Software: GoLand
 */

package workspace

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type nsRequest struct {
	Namespace string
}

func makeMetriceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(nsRequest)
		res, err := s.Metrice(ctx, req.Namespace)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeActiveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Active(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
