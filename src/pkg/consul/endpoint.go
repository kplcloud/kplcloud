/**
 * @Time : 2019/7/17 2:16 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package consul

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type detailRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type listRequest struct {
	detailRequest
	Page, Limit int
}
type rules struct {
	Name   string `json:"name"`
	Policy string `json:"policy"`
}
type consulRule struct {
	Key     []*rules `json:"key"`
	Query   []*rules `json:"query"`
	Service []*rules `json:"service"`
	Event   []*rules `json:"event"`
}

type postRequest struct {
	Namespace string      `json:"namespace"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Rules     *consulRule `json:"rules"`
	JsonRule  string
}

type kvDetailRequest struct {
	detailRequest
	Prefix      string `json:"prefix"`
	FolderState bool   `json:"folder"`
}
type kvPostRequest struct {
	detailRequest
	Key   string `json:"key"`
	Value string `json:"value"`
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = s.Sync(ctx)
		return encode.Response{Err: err}, err
	}
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(detailRequest)
		data, err := s.Detail(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err, Data: data}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		data, err := s.List(ctx, req.Namespace, req.Name, req.Page, req.Limit)
		return encode.Response{Err: err, Data: data}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Post(ctx, req.Namespace, req.Name, req.Type, req.JsonRule)
		return encode.Response{Err: err}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Update(ctx, req.Namespace, req.Name, req.Type, req.JsonRule)
		return encode.Response{Err: err}, err

	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(detailRequest)
		err = s.Delete(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err}, err
	}
}

func makeKVDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(kvDetailRequest)
		data, err := s.KVDetail(ctx, req.Namespace, req.Name, req.Prefix)
		return encode.Response{Err: err, Data: data}, err
	}
}

func makeKVListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(kvDetailRequest)
		data, err := s.KVList(ctx, req.Namespace, req.Name, req.Prefix)
		return encode.Response{Err: err, Data: data}, err
	}
}

func makeKVPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(kvPostRequest)
		err = s.KVPost(ctx, req.Namespace, req.Name, req.Key, req.Value)
		return encode.Response{Err: err}, err
	}
}

func makeKVDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(kvDetailRequest)
		err = s.KVDelete(ctx, req.Namespace, req.Name, req.Prefix, req.FolderState)
		return encode.Response{Err: err}, err
	}
}
