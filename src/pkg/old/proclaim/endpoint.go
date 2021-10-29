package proclaim

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type proclaimListRequest struct {
	Name        string
	Page, Limit int
}

type proclaimRequest struct {
	Id           int      `json:"id"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	ProclaimType string   `json:"proclaim_type"`
	Namespace    []string `json:"namespace"`
	UserList     []string `json:"userlist"`
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(proclaimRequest)
		res, err := s.Get(ctx, req.Id)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(proclaimRequest)
		err := s.Post(ctx, req)
		return encode.Response{Err: err, Data: nil}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(proclaimListRequest)
		res, err := s.List(ctx, req.Name, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}
