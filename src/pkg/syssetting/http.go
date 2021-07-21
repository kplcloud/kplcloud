/**
 * @Time : 7/20/21 5:44 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package syssetting

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"List": ems,
	})

	r := mux.NewRouter()

	r.Handle("/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	req.key = r.URL.Query().Get("key")

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		req.page = p
	} else {
		req.page = 1
	}
	if p, err := strconv.Atoi(r.URL.Query().Get("pageSize")); err == nil {
		req.pageSize = p
	} else {
		req.pageSize = 10
	}

	return req, nil
}
