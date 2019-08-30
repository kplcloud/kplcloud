/**
 * @Time : 2019-07-22 14:24
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package event

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

func makeAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.All(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
