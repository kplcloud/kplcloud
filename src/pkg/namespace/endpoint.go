package namespace

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type nsRequest struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

type nsResponse struct {
	Code int              `json:"code"`
	Data *types.Namespace `json:"data,omitempty"`
	Err  error            `json:"error,omitempty"`
}

func (r nsResponse) error() error { return r.Err }

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(nsRequest)
		rs, err := s.Get(ctx, req.Name)
		return encode.Response{Err: err, Data: rs}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(nsRequest)
		err := s.Post(ctx, req.Name, req.DisplayName)
		return encode.Response{Err: err}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(nsRequest)
		err := s.Update(ctx, req.Name, req.DisplayName)
		return encode.Response{Err: err}, err
	}
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.Sync(ctx)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.List(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
