/**
 * @Time : 2019-07-09 18:50
 * @Author : soupzhb@gmail.com
 * @File : endpoint.go
 * @Software: GoLand
 */

package wechat

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type qrRequest struct {
	Email string `json:"email"`
}

type receiveResponse struct {
	Code        int    `json:"code"`
	Data        string `json:"data"`
	Err         error  `json:"error,omitempty"`
	ContentType string `json:"content_type"`
}

func makeReceiveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, contentType, err := s.Receive(ctx)
		return receiveResponse{Err: err, Data: data, ContentType: contentType}, err

	}
}

func makeGetQrEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(qrRequest)
		res, err := s.GetQr(ctx, req)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeTestSendEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.TestSend(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeMenuEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Menu(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
