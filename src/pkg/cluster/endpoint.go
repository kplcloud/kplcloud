/**
 * @Time : 8/11/21 4:24 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package cluster

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
)

type (
	addRequest struct {
		Name  string `json:"name"  valid:"required"`
		Alias string `json:"alias"  valid:"required"`
		Data  string `json:"data"  valid:"required"`
	}
)

type Endpoints struct {
	AddEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		AddEndpoint: makeAddEndpoint(s),
	}

	for _, m := range dmw["Add"] {
		eps.AddEndpoint = m(eps.AddEndpoint)
	}
	return eps
}

func makeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addRequest)
		err = s.Add(ctx, req.Name, req.Alias, req.Data)
		return encode.Response{
			Error: err,
		}, err
	}
}
