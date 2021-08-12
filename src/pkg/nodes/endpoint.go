/**
 * @Time : 8/11/21 4:24 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package nodes

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
)

type (
	syncRequest struct {
		ClusterName string `json:"clusterName"  valid:"required"`
	}
)

type Endpoints struct {
	SyncEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SyncEndpoint: makeSyncEndpoint(s),
	}

	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	return eps
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(syncRequest)
		err = s.Sync(ctx, req.ClusterName)
		return encode.Response{
			Error: err,
		}, err
	}
}
