/**
 * @Time : 2019-06-25 19:24
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package storage

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type storageRequest struct {
	Name              string `json:"name,omitempty"`
	Provisioner       string `json:"provisioner"`
	ReclaimPolicy     string `json:"reclaim_policy"`
	VolumeBindingMode string `json:"volume_binding_mode"`
}

type storageListRequest struct {
	Offset, Limit int
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.Sync(ctx)
		return encode.Response{Err: err}, err
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(storageRequest)
		rs, err := s.Get(ctx, req.Name)
		return encode.Response{Err: err, Data: rs}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(storageRequest)
		err := s.Post(ctx, req.Name, req.Provisioner, req.ReclaimPolicy, req.VolumeBindingMode)
		return encode.Response{Err: err}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(storageRequest)
		err := s.Delete(ctx, req.Name)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(storageListRequest)
		res, err := s.List(ctx, req.Offset, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}
