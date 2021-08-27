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
		Name  string `json:"name" valid:"required"`
		Alias string `json:"alias" valid:"required"`
		Data  string `json:"data" valid:"required"`
	}

	listRequest struct {
		name           string
		page, pageSize int
	}
	listResult struct {
		Name    string `json:"name"`
		Alias   string `json:"alias"`
		Remark  string `json:"remark"` // 备注
		Label   string `json:"label"`  // 标签
		Status  int    `json:"status"` // 状态
		NodeNum int    `json:"node_num"`
	}
)

type Endpoints struct {
	AddEndpoint  endpoint.Endpoint
	ListEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		AddEndpoint:  makeAddEndpoint(s),
		ListEndpoint: makeListEndpoint(s),
	}

	for _, m := range dmw["Add"] {
		eps.AddEndpoint = m(eps.AddEndpoint)
	}
	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	return eps
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, total, err := s.List(ctx, req.name, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"total": total,
				"list":  res,
			},
			Error: err,
		}, err
	}
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
