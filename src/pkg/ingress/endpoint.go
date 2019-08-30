/**
 * @Time : 2019/7/2 2:10 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package ingress

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type getRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type listRequest struct {
	getRequest
	Page, Limit int
}

type postRequest struct {
	getRequest
	Rules []*ruleStruct `json:"rules"`
}

type ruleStruct struct {
	Domain string `json:"domain"`
	Paths  []struct {
		Path        string `json:"path"`
		ServiceName string `json:"serviceName"`
		PortName    int    `json:"port"`
	} `json:"paths"`
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		res, err := s.Get(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		res, err := s.List(ctx, req.getRequest.Namespace, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Post(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeGetNoIngressProjectEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		res, err := s.GetNoIngressProject(ctx, req.Namespace)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		err = s.Sync(ctx, req.Namespace)
		return encode.Response{Err: err}, err
	}
}

func makeGenerateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = s.Generate(ctx)
		return encode.Response{Err: err}, err
	}
}
