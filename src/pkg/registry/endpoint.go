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
	"time"
)

type (
	createRequest struct {
		Name     string `json:"name" valid:"required"`
		Host     string `json:"host" valid:"required"`
		Username string `json:"username" valid:"required"`
		Password string `json:"password" valid:"required"`
		Remark   string `json:"remark"`
	}

	result struct {
		Name      string    `json:"name"`
		Host      string    `json:"host"`
		Username  string    `json:"username"`
		Password  string    `json:"password"`
		Remark    string    `json:"remark"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	listRequest struct {
		query          string
		page, pageSize int
	}
)

type Endpoints struct {
	CreateEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateEndpoint: makeCreateEndpoint(s),
		ListEndpoint:   makeListEndpoint(s),
	}

	for _, m := range dmw["Create"] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	return eps
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, total, err := s.List(ctx, req.query, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
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
