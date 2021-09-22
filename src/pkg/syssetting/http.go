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
		"List":   ems,
		"Delete": ems,
	})

	r := mux.NewRouter()

	r.Handle("/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/delete/{id:[0-9]+}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodDelete)

	return r
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getRequest
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}
	settingId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	req.Id = settingId

	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	req.key = r.URL.Query().Get("query")

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
