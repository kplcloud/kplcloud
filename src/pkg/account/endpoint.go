package account

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type accountRequest struct {
	MemberId int
}

type accountReceiveRequest struct {
	Action string `json:"notice_action"`
	Site   int    `json:"site"`
	Wechat int    `json:"wechat"`
	Email  int    `json:"email"`
	Sms    int    `json:"sms"`
	Bee    int    `json:"bee"`
}

type accountBaseRequest struct {
	Name       string `json:"name"`
	City       string `json:"city"`
	Department string `json:"department"`
	Phone      string `json:"phone"`
}

type receiveResponse struct {
	Action     string `json:"action"`
	ActionDesc string `json:"action_desc"`
	Site       int    `json:"site"`
	Wechat     int    `json:"wechat"`
	Email      int    `json:"email"`
	Sms        int    `json:"sms"`
	Bee        int    `json:"bee"`
}

type accountResponse struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error                  `json:"error,omitempty"`
}

func (r accountResponse) error() error { return r.Err }

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Detail(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeGetReceiveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetReceive(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeUpdateReceiveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(accountReceiveRequest)
		err := s.UpdateReceive(ctx, req)
		return encode.Response{Err: err, Data: nil}, err
	}
}

func makeUpdateBaseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(accountBaseRequest)
		err := s.UpdateBase(ctx, req)
		return encode.Response{Err: err, Data: nil}, err
	}
}

func makeUnWechatBindEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.UnWechatBind(ctx)
		return encode.Response{Err: err, Data: nil}, err
	}
}

func makeGetProjectEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetProject(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
