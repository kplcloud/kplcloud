/**
 * @Time : 7/20/21 5:41 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package syssetting

import (
	"context"
	"github.com/go-kit/kit/endpoint"

	"github.com/kplcloud/kplcloud/src/encode"
)

type (
	listRequest struct {
		key            string
		page, pageSize int
	}
	listResult struct {
		Key string `json:"key"`
	}

	addRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

type Endpoints struct {
	ListEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint: makeListEndpoint(s),
	}

	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}

	return eps
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, total, err := s.List(ctx, req.key, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}
