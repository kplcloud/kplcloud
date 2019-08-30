/**
 * @Time : 2019-07-12 11:28
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package permission

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"io/ioutil"
	"net/http"
	"strconv"
)

type endpoints struct {
	MenuEndpoint   endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	PostEndpoint   endpoint.Endpoint
	DragEndpoint   endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		MenuEndpoint:   makeMenuEndpoint(svc),
		DeleteEndpoint: makeDeleteEndpoint(svc),
		ListEndpoint:   makeListEndpoint(svc),
		PostEndpoint:   makePostEndpoint(svc),
		DragEndpoint:   makeDragEndpoint(svc),
		UpdateEndpoint: makeUpdateEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Menu":   ems,
		"Delete": ems,
		"List":   ems,
		"Post":   ems,
		"Drag":   ems,
		"Update": ems,
	}

	for _, m := range mw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mw["Menu"] {
		eps.MenuEndpoint = m(eps.MenuEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["Drag"] {
		eps.DragEndpoint = m(eps.DragEndpoint)
	}
	for _, m := range mw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/permission/menu", kithttp.NewServer(
		eps.MenuEndpoint,
		func(ctx context.Context, req *http.Request) (request interface{}, err error) {
			return
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/permission/list", kithttp.NewServer(
		eps.ListEndpoint,
		func(ctx context.Context, req *http.Request) (request interface{}, err error) {
			return
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/permission/{id:[0-9]+}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeletePermissionRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/permission/{id:[0-9]+}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodePermissionRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/permission/create", kithttp.NewServer(
		eps.PostEndpoint,
		decodePermissionRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/permission/drag", kithttp.NewServer(
		eps.DragEndpoint,
		decodeDragRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	return r
}

func decodeDeletePermissionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	permissionId, _ := strconv.ParseInt(id, 10, 64)

	return deletePermissionRequest{Id: permissionId}, nil
}

func decodePermissionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req permissionRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	// todo 需要参数校验

	return req, nil
}

func decodeDragRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dragPermissionRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	// todo 需要参数校验

	return req, nil
}
