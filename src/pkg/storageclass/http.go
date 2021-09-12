/**
 * @Time : 8/11/21 4:26 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package storageclass

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	coreV1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"net/http"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Sync":   ems,
		"SyncPv": ems,
		"Create": ems,
	})

	r := mux.NewRouter()

	r.Handle("/{cluster}/sync", kithttp.NewServer(
		eps.SyncEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/sync/{storage}/pv", kithttp.NewServer(
		eps.SyncPvEndpoint,
		decodeSyncPvRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/create", kithttp.NewServer(
		eps.CreateEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodeSyncPvRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req syncPvRequest
	vars := mux.Vars(r)
	storageName, ok := vars["storage"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}
	req.Name = storageName

	return req, nil
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}

	var reclaimPolicy *coreV1.PersistentVolumeReclaimPolicy
	switch *req.ReclaimPolicy {
	case coreV1.PersistentVolumeReclaimRecycle:
		policy := coreV1.PersistentVolumeReclaimRecycle
		reclaimPolicy = &policy
	case coreV1.PersistentVolumeReclaimDelete:
		policy := coreV1.PersistentVolumeReclaimDelete
		reclaimPolicy = &policy
	case coreV1.PersistentVolumeReclaimRetain:
		policy := coreV1.PersistentVolumeReclaimRetain
		reclaimPolicy = &policy
	}
	var volumeMode *storagev1.VolumeBindingMode
	switch *req.VolumeMode {
	case storagev1.VolumeBindingImmediate:
		policy := storagev1.VolumeBindingImmediate
		volumeMode = &policy
	case storagev1.VolumeBindingWaitForFirstConsumer:
		policy := storagev1.VolumeBindingWaitForFirstConsumer
		volumeMode = &policy
	}

	if reclaimPolicy == nil {
		return nil, encode.InvalidParams.Wrap(errors.New("reclaimPolicy not null"))
	}
	if volumeMode == nil {
		return nil, encode.InvalidParams.Wrap(errors.New("volumeMode not null"))
	}

	req.VolumeMode = volumeMode
	req.ReclaimPolicy = reclaimPolicy

	return req, nil
}
