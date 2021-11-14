/**
 * @Time : 8/11/21 4:26 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package registry

import (
	"context"
	"encoding/json"
	"errors"
	valid "github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Create": ems,
	})

	r := mux.NewRouter()

	r.Handle("/create", kithttp.NewServer(
		eps.CreateEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/update/{name}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeUpdateRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)
	r.Handle("/delete/{name}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodDelete)
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
	req.query = r.URL.Query().Get("query")

	return req, nil
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createRequest
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

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req updateRequest
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.Invalid.Error()
	}
	req.Name = name
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
	var req registryRequest
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.Invalid.Error()
	}
	req.Name = name
	return req, nil
}
