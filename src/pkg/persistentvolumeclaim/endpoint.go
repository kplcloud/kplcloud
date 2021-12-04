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
		Storage     string   `json:"requestStorage" valid:"required"`
		StorageName string   `json:"storageClass" valid:"required"`
		AccessModes []string `json:"accessModes" valid:"required"`
		Remark      string   `json:"remark"`
	}

	deleteRequest struct {
		name string
	}

	listRequest struct {
		page, pageSize   int
		namespace, name  string
		cluster, storage string
	}
	result struct {
		Name           string            `json:"name"`
		Namespace      string            `json:"namespace"`
		StorageClass   string            `json:"storageClass"`
		CreatedAt      time.Time         `json:"createdAt"`
		UpdatedAt      time.Time         `json:"updatedAt"`
		AccessModes    []string          `json:"accessModes"`
		Remark         string            `json:"remark"`
		RequestStorage string            `json:"requestStorage"`
		LimitStorage   string            `json:"limitStorage"`
		Status         string            `json:"status"`
		Annotations    map[string]string `json:"annotations,omitempty"`
		Finalizers     []string          `json:"finalizers,omitempty"`
		Labels         map[string]string `json:"labels,omitempty"`
		VolumeName     string            `json:"volumeName,omitempty"`
		ClusterName    string            `json:"clusterName,omitempty"`
		ClusterAlias   string            `json:"clusterAlias,omitempty"`
	}
)

type Endpoints struct {
	SyncEndpoint   endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	CreateEndpoint endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateEndpoint: makeCreateEndpoint(s),
		SyncEndpoint:   makeSyncEndpoint(s),
		ListEndpoint:   makeListEndpoint(s),
		GetEndpoint:    makeGetEndpoint(s),
		DeleteEndpoint: makeDeleteEndpoint(s),
	}

	for _, m := range dmw["Create"] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	return eps
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, _ := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		//name, _ := ctx.Value(middleware.ContextKeyName).(string)
		req := request.(listRequest)
		err = s.Delete(ctx, clusterId, ns, req.name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, _ := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		req := request.(listRequest)
		res, err := s.Get(ctx, clusterId, ns, req.name)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
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
