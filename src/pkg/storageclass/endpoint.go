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
	coreV1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"time"
)

type (
	syncRequest struct {
		ClusterName string `json:"clusterName" valid:"required"`
	}
	syncPvRequest struct {
		Name string `json:"name"`
	}
	syncPvcRequest struct {
		Name string `json:"name"`
	}

	listRequest struct {
		page, pageSize int
	}
	listResult struct {
		Name          string    `json:"name"`
		Namespace     string    `json:"namespace"`
		Provisioner   string    `json:"provisioner"`
		VolumeMode    string    `json:"volumeMode"`
		ReclaimPolicy string    `json:"reclaimPolicy"`
		Remark        string    `json:"remark"`
		CreatedAt     time.Time `json:"createdAt"`
		UpdatedAt     time.Time `json:"updatedAt"`
	}

	infoResult struct {
		Name            string    `json:"name"`
		Provisioner     string    `json:"provisioner"`
		VolumeMode      string    `json:"volumeMode"`
		ReclaimPolicy   string    `json:"reclaimPolicy"`
		Remark          string    `json:"remark"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
		ResourceVersion string    `json:"resourceVersion"`
		Detail          string    `json:"detail"`
		ClusterAlias    string    `json:"clusterAlias"`
		ClusterName     string    `json:"clusterName"`
	}

	createRequest struct {
		Name          string                               `json:"name" valid:"required"`
		Namespace     string                               `json:"namespace"`
		ReclaimPolicy coreV1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy" valid:"required"`
		VolumeMode    storagev1.VolumeBindingMode          `json:"volumeMode" valid:"required"`
		Provisioner   string                               `json:"provisioner" valid:"required"`
		Remark        string                               `json:"remark"`
	}
)

type Endpoints struct {
	SyncEndpoint    endpoint.Endpoint
	SyncPvEndpoint  endpoint.Endpoint
	SyncPvcEndpoint endpoint.Endpoint
	CreateEndpoint  endpoint.Endpoint
	ListEndpoint    endpoint.Endpoint
	DeleteEndpoint  endpoint.Endpoint
	UpdateEndpoint  endpoint.Endpoint
	RecoverEndpoint endpoint.Endpoint
	InfoEndpoint    endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SyncEndpoint:    makeSyncEndpoint(s),
		SyncPvEndpoint:  makeSyncPvEndpoint(s),
		SyncPvcEndpoint: makeSyncPvcEndpoint(s),
		CreateEndpoint:  makeCreateEndpoint(s),
		ListEndpoint:    makeListEndpoint(s),
		DeleteEndpoint:  makeDeleteEndpoint(s),
		UpdateEndpoint:  makeUpdateEndpoint(s),
		RecoverEndpoint: makeRecoverEndpoint(s),
		InfoEndpoint:    makeInfoEndpoint(s),
	}

	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
		eps.SyncPvcEndpoint = m(eps.SyncPvcEndpoint)
	}
	for _, m := range dmw["SyncPv"] {
		eps.SyncPvEndpoint = m(eps.SyncPvEndpoint)
	}
	for _, m := range dmw["Create"] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
		eps.RecoverEndpoint = m(eps.RecoverEndpoint)
	}
	for _, m := range dmw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
		eps.InfoEndpoint = m(eps.InfoEndpoint)
	}
	return eps
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(listRequest)
		res, total, err := s.List(ctx, clusterId, req.page, req.pageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(createRequest)
		err = s.Update(ctx, clusterId, req.Name, req.Provisioner, &req.ReclaimPolicy, &req.VolumeMode, req.Remark)
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
		req := request.(createRequest)
		err = s.Create(ctx, clusterId, req.Namespace, req.Name, req.Provisioner, &req.ReclaimPolicy, &req.VolumeMode, req.Remark)
		return encode.Response{
			Error: err,
		}, err
	}
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

func makeSyncPvcEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(syncPvRequest)
		err = s.SyncPvc(ctx, clusterId, req.Name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(syncPvRequest)
		err = s.Delete(ctx, clusterId, req.Name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeRecoverEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(syncPvRequest)
		err = s.Recover(ctx, clusterId, req.Name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		name, ok := ctx.Value(contextKeyStorageClassName).(string)
		if !ok {
			return nil, encode.ErrStorageClassNotfound.Error()
		}
		res, err := s.Info(ctx, clusterId, name)
		return encode.Response{
			Data:  res,
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
