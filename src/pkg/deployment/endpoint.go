/**
 * @Time : 8/11/21 4:24 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package deployment

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/middleware"
)

type (
	syncRequest struct {
		ClusterId string `json:"clusterId"`
		Namespace string `json:"namespace"`
	}

	nodeResult struct {
		Name   string `json:"name"`
		Memory int64  `json:"memory"`
		Cpu    int64  `json:"cpu"`
	}

	listRequest struct {
		page, pageSize int
	}

	putImageRequest struct {
		Image string `json:"image" valid:"required"`
	}
)

type Endpoints struct {
	SyncEndpoint     endpoint.Endpoint
	ListEndpoint     endpoint.Endpoint
	PutImageEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SyncEndpoint:     makeSyncEndpoint(s),
		PutImageEndpoint: makePutImageEndpoint(s),
	}

	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range dmw["PutImage"] {
		eps.PutImageEndpoint = m(eps.PutImageEndpoint)
	}
	return eps
}

func makePutImageEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, ok := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		name, ok := ctx.Value(middleware.ContextKeyName).(string)
		if !ok {
			return nil, encode.ErrNameNotfound.Error()
		}
		req := request.(putImageRequest)
		err = s.PutImage(ctx, clusterId, ns, name, req.Image)
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
		ns, ok := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		err = s.Sync(ctx, clusterId, ns)
		return encode.Response{
			Error: err,
		}, err
	}
}
