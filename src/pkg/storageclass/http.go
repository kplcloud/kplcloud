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
	valid "github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	coreV1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"net/http"
	"strconv"
	"strings"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opt []kithttp.ServerOption) http.Handler {
	var ems []endpoint.Middleware

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Sync":   ems,
		"SyncPv": ems,
		"Create": ems,
		"Info":   ems,
		"List":   ems,
		"Update": ems,
	})
	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			var storageName string
			vars := mux.Vars(request)
			storageName, ok := vars["storage"]
			if !ok {
				return ctx
			}
			ctx = context.WithValue(ctx, contextKeyStorageClassName, storageName)
			return ctx
		}),
	}

	opts = append(opts, opt...)

	r := mux.NewRouter()

	r.Handle("/{cluster}/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
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
	r.Handle("/{cluster}/sync/{storage}/pvc", kithttp.NewServer(
		eps.SyncPvcEndpoint,
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
	r.Handle("/{cluster}/delete/{storage}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeSyncPvRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodDelete)
	r.Handle("/{cluster}/update/{storage}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)
	r.Handle("/{cluster}/recover/{storage}", kithttp.NewServer(
		eps.RecoverEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)
	r.Handle("/{cluster}/info/{storage}", kithttp.NewServer(
		eps.InfoEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	req.page = page
	req.pageSize = pageSize
	//req.query = r.URL.Query().Get("query")

	return req, nil
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

	var reclaimPolicy coreV1.PersistentVolumeReclaimPolicy
	switch req.ReclaimPolicy {
	case coreV1.PersistentVolumeReclaimRecycle:
		reclaimPolicy = coreV1.PersistentVolumeReclaimRecycle
	case coreV1.PersistentVolumeReclaimDelete:
		reclaimPolicy = coreV1.PersistentVolumeReclaimDelete
	case coreV1.PersistentVolumeReclaimRetain:
		reclaimPolicy = coreV1.PersistentVolumeReclaimRetain
	default:
		reclaimPolicy = ""
	}
	var volumeMode storagev1.VolumeBindingMode
	switch req.VolumeMode {
	case storagev1.VolumeBindingImmediate:
		volumeMode = storagev1.VolumeBindingImmediate
	case storagev1.VolumeBindingWaitForFirstConsumer:
		volumeMode = storagev1.VolumeBindingWaitForFirstConsumer
	default:
		volumeMode = ""
	}

	if strings.EqualFold(string(reclaimPolicy), "") {
		return nil, encode.InvalidParams.Wrap(errors.New("reclaimPolicy not null"))
	}
	if strings.EqualFold(string(volumeMode), "") {
		return nil, encode.InvalidParams.Wrap(errors.New("volumeMode not null"))
	}

	req.VolumeMode = volumeMode
	req.ReclaimPolicy = reclaimPolicy

	validResult, err := valid.ValidateStruct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if !validResult {
		return nil, encode.InvalidParams.Wrap(errors.New("valid false"))
	}

	return req, nil
}
