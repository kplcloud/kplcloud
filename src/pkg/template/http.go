/**
 * @Time : 3/9/21 5:58 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package template

import (
	"context"
	"encoding/json"
	valid "github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
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
		"Add":    ems,
		"Delete": ems,
		"Update": ems,
		"Locked": ems,
	})

	r := mux.NewRouter()

	r.Handle("/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/add", kithttp.NewServer(
		eps.AddEndpoint,
		decodeAddRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	req.email = r.URL.Query().Get("email")

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

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Error()
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
