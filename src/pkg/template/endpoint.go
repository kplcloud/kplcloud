/**
 * @Time : 2019/6/25 4:07 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package template

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type templateListRequest struct {
	Name        string
	Page, Limit int
}

type templateRequest struct {
	Id     int    `json:"-"`
	Name   string `json:"name,omitempty"`
	Kind   string `json:"kind,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(templateRequest)
		rs, err := s.Get(ctx, req.Id)
		return encode.Response{Err: err, Data: rs}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(templateRequest)
		err := s.Post(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(templateRequest)
		err := s.Update(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(templateRequest)
		err := s.Delete(ctx, req.Id)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(templateListRequest)
		res, err := s.List(ctx, req.Name, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}
