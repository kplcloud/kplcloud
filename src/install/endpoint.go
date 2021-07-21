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
		Drive    string `json:"drive"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	}
)

type Endpoints struct {
	InitEndpoint   endpoint.Endpoint
	InitDbEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		InitDbEndpoint: makeInitDbEndpoint(s),
		InitEndpoint:   nil,
	}

	for _, m := range dmw["Init"] {
		eps.InitEndpoint = m(eps.InitEndpoint)
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
