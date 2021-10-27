/**
 * @Time : 2019-06-27 18:02
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package terminal

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type indexRequest struct {
	PodName   string
	Container string
	Token     string
}

type indexResponse struct {
	Code int        `json:"code"`
	Data *IndexData `json:"data,omitempty"`
	Err  error      `json:"error,omitempty"`
}

func (r indexResponse) error() error { return r.Err }

type errorer interface {
	error() error
}

func makeIndexEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(indexRequest)
		var code int
		res, err := s.Index(ctx, req.PodName, req.Container)
		if err != nil {
			code = -1
		}
		return indexResponse{Code: code, Err: err, Data: res}, err
	}
}
