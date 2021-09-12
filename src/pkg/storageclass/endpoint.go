/**
 * @Time : 8/11/21 4:24 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package storageclass

import (
	"context"
	coreV1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"

	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/middleware"
)

type (
	syncRequest struct {
		ClusterName string `json:"clusterName" valid:"required"`
	}
	syncPvRequest struct {
		Name string `json:"name"`
	}

	listRequest struct {
		page, pageSize int
	}

	createRequest struct {
		Name          string                                `json:"name" valid:"required"`
		Namespace     string                                `json:"namespace" valid:"required"`
		ReclaimPolicy *coreV1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy" valid:"required"`
		VolumeMode    *storagev1.VolumeBindingMode          `json:"volumeMode" valid:"required"`
		Provisioner   string                                `json:"provisioner" valid:"required"`
	}
)

type Endpoints struct {
	SyncEndpoint   endpoint.Endpoint
	SyncPvEndpoint endpoint.Endpoint
	CreateEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SyncEndpoint:   makeSyncEndpoint(s),
		SyncPvEndpoint: makeSyncPvEndpoint(s),
		CreateEndpoint: makeCreateEndpoint(s),
	}

	for _, m := range dmw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range dmw["SyncPv"] {
		eps.SyncPvEndpoint = m(eps.SyncPvEndpoint)
	}
	for _, m := range dmw["Create"] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
	}
	return eps
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		req := request.(createRequest)
		err = s.Create(ctx, clusterId, req.Namespace, req.Name, req.Provisioner, req.ReclaimPolicy, req.VolumeMode)
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
