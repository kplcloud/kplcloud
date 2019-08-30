/**
 * @Time : 2019/7/24 3:46 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package audit

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type detailRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type stepRequest struct {
	detailRequest
	Kind string `json:"kind"`
}

func makeAccessAuditEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(detailRequest)
		err = s.AccessAudit(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err}, err
	}
}

func makeAduitStepEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(stepRequest)
		err = s.AuditStep(ctx, req.Namespace, req.Name, req.Kind)
		return encode.Response{Err: err}, err
	}
}

func makeRefusedEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(detailRequest)
		err = s.Refused(ctx, req.Namespace, req.Name)
		return encode.Response{Err: err}, err
	}
}
