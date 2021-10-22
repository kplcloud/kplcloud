/**
 * @Time : 8/11/21 4:26 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package nodes

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
		"Sync":   ems,
		"Info":   ems,
		"List":   ems,
		"Cordon": ems,
		"Drain":  ems,
	})

	r := mux.NewRouter()

	r.Handle("/{cluster}/sync", kithttp.NewServer(
		eps.SyncEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/info/{name}", kithttp.NewServer(
		eps.InfoEndpoint,
		decodeInfoRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/cordon/{name}", kithttp.NewServer(
		eps.CordonEndpoint,
		decodeInfoRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)
	r.Handle("/{cluster}/drain/{name}", kithttp.NewServer(
		eps.DrainEndpoint,
		decodeDrainRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPut)

	return r
}

func decodeDrainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req drainRequest

	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}
	req.Name = name

	req.Force, _ = strconv.ParseBool(r.URL.Query().Get("force"))

	return req, nil
}
func decodeInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req infoRequest

	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}
	req.Name = name

	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 1 {
		pageSize = 10
	}
	req.page = page
	req.pageSize = pageSize
	req.query = r.URL.Query().Get("query")
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
