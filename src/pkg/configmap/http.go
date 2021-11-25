/**
 * @Time: 2021/8/18 23:14
 * @Author: solacowa@gmail.com
 * @File: http
 * @Software: GoLand
 */

package configmap

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Sync": ems,
		"List": ems,
	})

	r := mux.NewRouter()

	r.Handle("/{cluster}/sync/{namespace}", kithttp.NewServer(
		eps.SyncEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/list/{namespace}", kithttp.NewServer(
		eps.ListEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}
