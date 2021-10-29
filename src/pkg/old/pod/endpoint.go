/**
 * @Time : 2019-07-04 16:04
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package pod

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type detailRequest struct {
	PodName string `json:"pod_name"`
}

type getLogsRequest struct {
	PodName   string `json:"pod_name"`
	Container string `json:"container"`
	Previous  bool   `json:"previous"`
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(detailRequest)
		res, err := s.Detail(ctx, req.PodName)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeProjectPodsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.ProjectPods(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeGetLogEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getLogsRequest)
		res, err := s.GetLog(ctx, req.PodName, req.Container, req.Previous)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeDownloadLogEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getLogsRequest)
		res, err := s.DownloadLog(ctx, req.PodName, req.Container, req.Previous)
		return response{Err: err, Data: res}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(detailRequest)
		err := s.Delete(ctx, req.PodName)

		return encode.Response{Err: err}, err
	}
}

func makePodsMetricsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.PodsMetrics(ctx)

		return encode.Response{Err: err, Data: res}, err
	}
}
