/**
 * @Time : 2019-07-08 10:39
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package discovery

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type getRequest struct {
	ServiceName string
}

type listRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type createRequest struct {
	Name           string `json:"name"`
	EditState      bool   `json:"edit_state"`
	ResourceType   string `json:"resource_type"`   //service endpoint
	ServiceProject string `json:"service_project"` //集群内访问项目名称
	Routes         []struct {
		Port       int    `json:"port"`
		Protocol   string `json:"protocol"`
		TargetPort int    `json:"target_port"`
		Name       string `json:"name"`
	} `json:"routes"`
	Address struct {
		Ips   []string `json:"ips"`
		Ports []struct {
			Port int    `json:"port"`
			Name string `json:"name"`
		} `json:"ports"`
	} `json:"address"`
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		err := s.Delete(ctx, req.ServiceName)
		return encode.Response{Err: err}, err
	}
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		rs, err := s.Detail(ctx, req.ServiceName)
		return encode.Response{Err: err, Data: rs}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		res, err := s.List(ctx, req.Page, req.Limit)
		return encode.Response{Err: err, Data: map[string]interface{}{
			"items": res,
			"page":  req.Page,
			"limit": req.Limit,
		}}, err
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.Create(ctx, request.(createRequest))
		return encode.Response{Err: err}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.Update(ctx, request.(createRequest))
		return encode.Response{Err: err}, err
	}
}
