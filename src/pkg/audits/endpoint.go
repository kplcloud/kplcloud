/**
 * @Time : 2021/10/27 5:18 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package audits

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"time"
)

type (
	listRequest struct {
		page, pageSize int
		query, status  string
	}
	auditResult struct {
		Cluster        string    `json:"cluster"`
		Name           string    `json:"name"`
		Namespace      string    `json:"namespace"`
		Username       string    `json:"username"`
		Method         string    `json:"method"`
		Remark         string    `json:"remark"`
		PermissionName string    `json:"permissionName"`
		Request        string    `json:"request"`
		Response       string    `json:"response"`
		Headers        string    `json:"headers"`
		TimeSince      string    `json:"timeSince"`
		Status         string    `json:"status"`
		Url            string    `json:"url"`
		TraceId        string    `json:"traceId"`
		CreatedAt      time.Time `json:"createdAt"` // 创建时间
	}
)

type Endpoints struct {
	ListEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint: makeListEndpoint(s),
	}
	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	return eps
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, total, err := s.List(ctx, req.query, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}
