/**
 * @Time : 8/11/21 4:26 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package secret

import (
	"context"
	"encoding/json"
	valid "github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Sync":        ems,
		"ImageSecret": ems,
		"Delete":      ems,
	})

	r := mux.NewRouter()

	r.Handle("/{cluster}/sync/{namespace}", kithttp.NewServer(
		eps.SyncEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/image-secret/{namespace}", kithttp.NewServer(
		eps.ImageSecretEndpoint,
		decodeImageSecretRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/{cluster}/delete/{namespace}/name/{name}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteSecretRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodDelete)

	return r
}

func decodeDeleteSecretRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req deleteRequest
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}
	req.Name = name
	return req, nil
}

func decodeImageSecretRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req imageSecret
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

func decodeSyncRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req syncRequest
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
