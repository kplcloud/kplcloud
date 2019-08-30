/**
 * @Time : 2019/7/5 11:02 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package configmap

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
	Desc string `json:"desc"`
	Data []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"data"`
}

type createConfigMapRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Desc      string `json:"desc"`
	Type      int64  `json:"type"`
}

type createConfigMapDataRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	ConfigMapId int64  `json:"config_map_id"`
}

type configMapDataRequest struct {
	ConfigMapId     int64  `json:"config_map_id"`
	ConfigMapDataId int64  `json:"id"`
	Key             string `json:"key"`
	Value           string `json:"value"`
	Path            string `json:"path"`
}

type configEnvRequest struct {
	EnvDesc   string `json:"env_desc"`
	EnvKey    string `json:"env_key"`
	EnvVar    string `json:"env_var"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Id        int64  `json:"id"`
}

func makeDeleteConfigEnvEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(configEnvRequest)
		err = s.ConfigEnvDel(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeUpdateConfigEnvEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(configEnvRequest)
		err = s.ConfigEnvUpdate(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeCreateConfigEnvEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(configEnvRequest)
		err = s.CreateConfigEnv(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeGetConfigEnvEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, err := s.GetConfigEnv(ctx, req.Name, req.Namespace, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeGetOneEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		res, err := s.GetOne(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeGetOnePullEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		res, err := s.GetOnePull(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, err := s.List(ctx, req)
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

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Update(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		err = s.Delete(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err}, err
	}
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		err = s.Sync(ctx, req.Namespace)
		return encode.Response{Err: err}, err
	}
}

func makeCreateConfigMapEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createConfigMapRequest)
		err = s.CreateConfigMap(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeGetConfigMapEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		res, err := s.GetConfigMap(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeGetConfigMapDataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, err := s.GetConfigMapData(ctx, req.Namespace, req.Name, req.Page, req.Limit)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeCreateConfigMapDataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createConfigMapDataRequest)
		err = s.CreateConfigMapData(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeUpdateConfigMapDataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(configMapDataRequest)
		err = s.UpdateConfigMapData(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeDeleteConfigMapDataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(configMapDataRequest)
		err = s.DeleteConfigMapData(ctx, req)
		return encode.Response{Err: err}, err
	}
}
