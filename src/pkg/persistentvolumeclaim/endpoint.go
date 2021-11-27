/**
 * @Time : 8/11/21 4:24 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package persistentvolumeclaim

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/middleware"
	"time"
)

type (
	createRequest struct {
		Name        string   `json:"name" valid:"required"`
		Storage     string   `json:"storage" valid:"required"`
		StorageName string   `json:"storageName" valid:"required"`
		AccessModes []string `json:"accessModes" valid:"required"`
	}

	listRequest struct {
		page, pageSize   int
		namespace, name  string
		cluster, storage string
	}
	result struct {
		Name           string    `json:"name"`
		Namespace      string    `json:"namespace"`
		StorageClass   string    `json:"storageClass"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
		AccessModes    []string  `json:"accessModes"`
		Remark         string    `json:"remark"`
		RequestStorage string    `json:"requestStorage"`
		LimitStorage   string    `json:"limitStorage"`
		Status         string    `json:"status"`
	}
)

type Endpoints struct {
	SyncEndpoint   endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	CreateEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateEndpoint: makeCreateEndpoint(s),
		SyncEndpoint:   makeSyncEndpoint(s),
		ListEndpoint:   makeListEndpoint(s),
	}

	for _, m := range dmw["Create"] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
	}
	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	return eps
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, _ := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		req := request.(listRequest)
		res, total, err := s.List(ctx, clusterId, req.storage, ns, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
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

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, ok := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		req := request.(createRequest)
		err = s.Create(ctx, clusterId, ns, req.Name, req.Storage, req.StorageName, req.AccessModes)
		return encode.Response{
			Error: err,
		}, err
	}
}
