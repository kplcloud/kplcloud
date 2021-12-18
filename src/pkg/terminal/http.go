/**
 * @Time : 2021/12/6 10:43 AM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package terminal

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kplcloud/kplcloud/src/encode"
	"net/http"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Token": ems,
	})

	r := mux.NewRouter()

	r.Handle("/{cluster}/namespace/{namespace}/service/{svcName}/pods/{podName}/token", kithttp.NewServer(
		eps.TokenEndpoint,
		decodeTokenRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}

func decodeTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req tokenRequest

	vars := mux.Vars(r)
	podName, ok := vars["podName"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}
	svcName, ok := vars["svcName"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}

	req.PodName = podName
	req.ServiceName = svcName

	return req, nil
}
