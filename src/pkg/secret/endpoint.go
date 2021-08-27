/**
 * @Time : 8/11/21 4:24 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package secret

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
	imageSecret struct {
		Name     string `json:"name" valid:"required"`
		Host     string `json:"host" valid:"required"`
		Username string `json:"username" valid:"required"`
		Password string `json:"password" valid:"required"`
	}
	deleteRequest struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
)

type Endpoints struct {
	SyncEndpoint        endpoint.Endpoint
	ImageSecretEndpoint endpoint.Endpoint
	DeleteEndpoint      endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SyncEndpoint:        makeSyncEndpoint(s),
		ImageSecretEndpoint: makeImageSecretEndpoint(s),
		DeleteEndpoint:      makeDeleteEndpoint(s),
	}

	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range dmw["ImageSecret"] {
		eps.ImageSecretEndpoint = m(eps.ImageSecretEndpoint)
	}
	for _, m := range dmw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	return eps
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, ok := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		req := request.(deleteRequest)
		err = s.Delete(ctx, clusterId, ns, req.Name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeImageSecretEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, ok := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		req := request.(imageSecret)
		err = s.ImageSecret(ctx, clusterId, ns, req.Name, req.Host, req.Username, req.Password)
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
