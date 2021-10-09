/**
 * @Time : 8/11/21 4:26 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package cluster

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Add":    ems,
		"List":   ems,
		"Delete": ems,
		"Update": ems,
		"Info":   ems,
	})

	r := mux.NewRouter()

	r.Handle("/add", kithttp.NewServer(
		eps.AddEndpoint,
		decodeAddRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{name}/delete", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodDelete)
	r.Handle("/{name}/update", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeAddRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)
	r.Handle("/{name}/info", kithttp.NewServer(
		eps.InfoEndpoint,
		decodeInfoRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}

func decodeInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req infoRequest
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.Invalid.Error()
	}
	req.Name = name
	return req, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req addRequest
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.Invalid.Error()
	}
	req.Name = name
	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	query := r.URL.Query().Get("query")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	req.pageSize = pageSize
	req.page = page
	req.name = query
	return req, nil
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, encode.InvalidParams.Wrap(err)
	}
	validResult, err := valid.ValidateStruct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if !validResult {
		return nil, encode.InvalidParams.Wrap(errors.New("valid false"))
	}
	return req, nil
}
