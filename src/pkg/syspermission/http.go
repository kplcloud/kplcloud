/**
 * @Time : 5/25/21 3:12 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package syspermission

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		//"Add": ems,
		"All": ems,
	})

	r := mux.NewRouter()

	r.Handle("/add", kithttp.NewServer(
		eps.AddEndpoint,
		decodeAddRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/all", kithttp.NewServer(
		eps.AllEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{id:[0-9]+}/update", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeUpdateRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)
	r.Handle("/{id:[0-9]+}/delete", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodDelete)
	r.Handle("/drag", kithttp.NewServer(
		eps.DragEndpoint,
		decodeDragRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)

	return r
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}
	permId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, encode.InvalidParams.Error()
	}
	var req deleteRequest
	req.Id = permId
	return req, nil
}

func decodeDragRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dragRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
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

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
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
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}
	permId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, encode.InvalidParams.Error()
	}
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	validResult, err := valid.ValidateStruct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if !validResult {
		return nil, encode.InvalidParams.Wrap(errors.New("valid false"))
	}
	req.Id = permId
	return req, nil
}
