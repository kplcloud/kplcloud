/**
 * @Time : 2019-07-29 11:30
 * @Author : soupzhb@gmail.com
 * @File : endpoint.go
 * @Software: GoLand
 */

package statistics

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type buildRequest struct {
	Namespace   string
	Name        string
	GroupID     int
	STime       string
	ETime       string
	Page, Limit int
}

func makeBuildEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(buildRequest)
		res, err := s.Build(ctx, req)
		return encode.Response{Err: err, Data: res}, err
	}
}
