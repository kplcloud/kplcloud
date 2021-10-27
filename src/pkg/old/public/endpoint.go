/**
 * @Time : 2019-06-27 13:34
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package public

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

func makeGitPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(gitlabHook)
		err := s.GitPost(ctx, req.Namespace, req.Name, req.Token, req.KeyWord, req.Branch, req)
		return encode.Response{Err: err}, err
	}
}

func makePrometheusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*prometheusAlerts)
		err = s.PrometheusAlert(ctx, req)
		return encode.Response{Err: err}, err
	}
}

func makeConfigEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.Config(ctx)
		return encode.Response{Data: res, Err: err}, err
	}
}
