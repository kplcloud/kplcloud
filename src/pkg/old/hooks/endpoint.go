/**
 * @Time : 2019/6/27 10:10 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package hooks

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type hookRequest struct {
	Id        int    `json:"id,omitempty" json:"id"`
	Name      string `json:"name,omitempty" json:"name"`
	Url       string `json:"url,omitempty" json:"url"`
	Token     string `json:"token,omitempty" json:"token"`
	Target    string `json:"target,omitempty" json:"target"`
	AppName   string `json:"app_name,omitempty" json:"app_name"`
	Namespace string `json:"namespace,omitempty" json:"namespace"`
	Status    int64  `json:"status,omitempty" json:"status"`
	Events    []int  `json:"events,omitempty" json:"events"`
}

type hookListRequest struct {
	Name, AppName, Namespace string
	Page, Limit              int
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hookRequest)
		rs, err := s.Get(ctx, req.Id)
		return encode.Response{Err: err, Data: rs}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hookListRequest)
		rs, err := s.List(ctx, req.Name, req.AppName, req.Namespace, req.Page, req.Limit)
		return encode.Response{Err: err, Data: rs}, err
	}
}

func makeListNoAppEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(hookListRequest)
		rs, err := s.List(ctx, req.Name, "", req.Namespace, req.Page, req.Limit)
		return encode.Response{Err: err, Data: rs}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hookRequest)
		err := s.Post(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hookRequest)
		err := s.Update(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hookRequest)
		err := s.Delete(ctx, req.Id)
		return encode.Response{Err: err}, err
	}
}

func makeTestSendEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hookRequest)
		err := s.TestSend(ctx, req.Id)
		return encode.Response{Err: err}, err
	}
}
