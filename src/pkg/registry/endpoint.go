/**
 * @Time : 8/13/21 2:55 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package registry

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
)

type (
	createRequest struct {
		Name     string `json:"name" valid:"required"`
		Host     string `json:"host" valid:"required"`
		Username string `json:"username" valid:"required"`
		Password string `json:"password" valid:"required"`
		Remark   string `json:"remark"`
	}
)

type Endpoints struct {
	CreateEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateEndpoint: makeCreateEndpoint(s),
	}

	for _, m := range dmw["Create"] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
	}
	return eps
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createRequest)
		err = s.Create(ctx, req.Name, req.Host, req.Username, req.Password, req.Remark)
		return encode.Response{
			Error: err,
		}, err
	}
}
