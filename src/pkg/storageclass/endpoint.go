/**
 * @Time : 8/11/21 4:24 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package storageclass

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
	syncPvRequest struct {
		Name string `json:"name"`
	}

	listRequest struct {
		page, pageSize int
	}
)

type Endpoints struct {
	SyncEndpoint   endpoint.Endpoint
	SyncPvEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SyncEndpoint:   makeSyncEndpoint(s),
		SyncPvEndpoint: makeSyncPvEndpoint(s),
	}

	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range dmw["SyncPv"] {
		eps.SyncPvEndpoint = m(eps.SyncPvEndpoint)
	}
	return eps
}

func makeSyncPvEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(syncPvRequest)
		err = s.SyncPv(ctx, clusterId, req.Name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeSyncEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		err = s.Sync(ctx, clusterId)
		return encode.Response{
			Error: err,
		}, err
	}
}
