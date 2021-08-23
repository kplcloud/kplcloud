/**
 * @Time : 8/11/21 4:26 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package storageclass

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Sync":   ems,
		"SyncPv": ems,
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
