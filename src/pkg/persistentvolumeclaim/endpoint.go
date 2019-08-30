/**
 * @Time : 2019-06-26 14:48
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package persistentvolumeclaim

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type ResourceUnit string

const (
	ResourceMi ResourceUnit = "Mi"
	ResourceGi ResourceUnit = "Gi"
)

func (c ResourceUnit) String() string {
	return string(c)
}

type getRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type pvcRequest struct {
	Namespace        string       `json:"namespace"`
	Name             string       `json:"name"`
	AccessModes      []string     `json:"access_modes"`
	Storage          string       `json:"storage"`
	Unit             ResourceUnit `json:"unit"`
	StorageClassName string       `json:"storage_class_name"`
}

type listRequest struct {
	getRequest
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pvcRequest)
		err := s.Sync(ctx, req.Namespace)
		return encode.Response{Err: err}, err
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pvcRequest)
		rs, err := s.Get(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err, Data: rs}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pvcRequest)
		err := s.Delete(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pvcRequest)
		err := s.Post(ctx, req.Namespace, req.Name, req.Storage, req.StorageClassName, req.AccessModes)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		res, err := s.List(ctx, req.Namespace, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		res, err := s.List(ctx, req.Namespace, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}
