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
	"time"

	"github.com/kplcloud/kplcloud/src/encode"
)

type (
	listRequest struct {
		key            string
		page, pageSize int
	}
	listResult struct {
		Section   string    `json:"section"`
		Key       string    `json:"key"`
		Value     string    `json:"value"`
		Id        int64     `json:"id"`
		Remark    string    `json:"remark"`
		Enable    bool      `json:"enable"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	addRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	getRequest struct {
		Section string `json:"section" valid:"required"`
		Key     string `json:"key" valid:"required"`
		Value   string `json:"value" valid:"required"`
		Id      int64  `json:"id"`
		Remark  string `json:"remark"`
		Enable  bool   `json:"enable"`
	}
)

type Endpoints struct {
	ListEndpoint   endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	AddEndpoint    endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint:   makeListEndpoint(s),
		DeleteEndpoint: makeDeleteEndpoint(s),
		AddEndpoint:    makeAddEndpoint(s),
		UpdateEndpoint: makeUpdateEndpoint(s),
	}

	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range dmw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range dmw["Add"] {
		eps.AddEndpoint = m(eps.AddEndpoint)
	}
	for _, m := range dmw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}

	return eps
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		err = s.Update(ctx, req.Id, req.Section, req.Key, req.Value, req.Remark, req.Enable)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		err = s.Add(ctx, req.Section, req.Key, req.Value, req.Remark, req.Enable)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		err = s.Delete(ctx, req.Id)
		return encode.Response{
			Error: err,
		}, err
	}
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
