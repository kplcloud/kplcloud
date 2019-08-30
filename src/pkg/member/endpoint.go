/**
 * @Time : 2019-07-17 14:17
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package member

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type getMemberRequest struct {
	Id int64 `json:"id"`
}

type memberRequest struct {
	Id         int64    `json:"id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	State      int64    `json:"state"`
	Namespaces []string `json:"namespace"`
	Roles      []int64  `json:"role"`
}

type listMemberRequest struct {
	Page  int
	Limit int
	Email string
}

func makeNamespacesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Namespaces(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getMemberRequest)
		res, err := s.Detail(ctx, req.Id)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(memberRequest)
		err := s.Post(ctx, req.Username, req.Email, req.Password, req.State, req.Namespaces, req.Roles)
		return encode.Response{Err: err}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(memberRequest)
		err := s.Update(ctx, req.Id, req.Username, req.Email, req.Password, req.State, req.Namespaces, req.Roles)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listMemberRequest)
		res, err := s.List(ctx, req.Page, req.Limit, req.Email)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeMeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Me(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.All(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
