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
	"github.com/kplcloud/kplcloud/src/middleware"
)

type (
	syncRequest struct {
		ClusterName string `json:"clusterName"  valid:"required"`
	}

	nodeResult struct {
		Name   string `json:"name"`
		Memory int64  `json:"memory"`
		Cpu    int64  `json:"cpu"`
	}

	listRequest struct {
		page, pageSize int
	}
)

type Endpoints struct {
	SyncEndpoint endpoint.Endpoint
	ListEndpoint endpoint.Endpoint
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
		clusterName, ok := ctx.Value(middleware.ContextKeyClusterName).(string)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		err = s.Sync(ctx, clusterName)
		return encode.Response{
			Error: err,
		}, err
	}
}
