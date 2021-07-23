/**
 * @Time : 7/21/21 2:26 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package install

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
)

type (
	initDbRequest struct {
		Drive    string `json:"drive" valid:"required"`
		Host     string `json:"host" valid:"required"`
		Port     int    `json:"port" valid:"required"`
		User     string `json:"user" valid:"required"`
		Password string `json:"password" valid:"required"`
		Database string `json:"database" valid:"required"`
	}
)

type Endpoints struct {
	InitDbEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		InitDbEndpoint: makeInitDbEndpoint(s),
	}

	for _, m := range dmw["InitDb"] {
		eps.InitDbEndpoint = m(eps.InitDbEndpoint)
	}

	return eps
}

func makeInitDbEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(initDbRequest)
		err = s.InitDb(ctx, req.Drive, req.Host, req.Port, req.User, req.Password, req.Database)
		return encode.Response{
			Error: err,
		}, err
	}
}
