/**
 * @Time : 7/20/21 5:44 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package syssetting

import (
	"context"
	"encoding/json"
	"errors"
	valid "github.com/asaskevich/govalidator"
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
		"Add":    ems,
		"Update": ems,
	})

	r := mux.NewRouter()

	r.Handle("/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{id:[0-9]+}/delete", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodDelete)
	r.Handle("/add", kithttp.NewServer(
		eps.AddEndpoint,
		decodeAddRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/{id:[0-9]+}/update", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeUpdateRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)

	return r
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
	req.Id = settingId

	return req, nil
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getRequest
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
