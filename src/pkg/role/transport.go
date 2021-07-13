/**
 * @Time : 2019-07-16 18:01
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package role

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
	PermissionSelectedEndpoint endpoint.Endpoint
	UpdateEndpoint             endpoint.Endpoint
	PostEndpoint               endpoint.Endpoint
	AllEndpoint                endpoint.Endpoint
	DetailEndpoint             endpoint.Endpoint
	DeleteEndpoint             endpoint.Endpoint
	RolePermissionEndpoint     endpoint.Endpoint
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
		PermissionSelectedEndpoint: makePermissionSelectedEndpoint(svc),
		UpdateEndpoint:             makeUpdateEndpoint(svc),
		PostEndpoint:               makePostEndpoint(svc),
		AllEndpoint:                makeAllEndpoint(svc),
		DetailEndpoint:             makeDetailEndpoint(svc),
		DeleteEndpoint:             makeDeleteEndpoint(svc),
		RolePermissionEndpoint:     makeRolePermissionEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"PermissionSelected": ems,
		"UpdateEndpoint":     ems,
		"PostEndpoint":       ems,
		"AllEndpoint":        ems,
		"DetailEndpoint":     ems,
		"DeleteEndpoint":     ems,
		"RolePermission":     ems,
	}

	for _, m := range mw["PermissionSelected"] {
		eps.PermissionSelectedEndpoint = m(eps.PermissionSelectedEndpoint)
	}
	for _, m := range mw["UpdateEndpoint"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}
	for _, m := range mw["PostEndpoint"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["AllEndpoint"] {
		eps.AllEndpoint = m(eps.AllEndpoint)
	}
	for _, m := range mw["DetailEndpoint"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}
	for _, m := range mw["DeleteEndpoint"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mw["RolePermission"] {
		eps.RolePermissionEndpoint = m(eps.RolePermissionEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/role/{id:[0-9]+}/permission", kithttp.NewServer(
		eps.PermissionSelectedEndpoint,
		decodeGetRolePermissionRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/role/{id:[0-9]+}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeGetRolePermissionRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/role/{id:[0-9]+}", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeGetRolePermissionRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/role", kithttp.NewServer(
		eps.AllEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/role", kithttp.NewServer(
		eps.PostEndpoint,
		decodeRoleRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/role/{id:[0-9]+}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeRoleRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/role/{id:[0-9]+}/permission", kithttp.NewServer(
		eps.RolePermissionEndpoint,
		decodeRolePermissionRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	return r
}

func decodeGetRolePermissionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	roleId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	return getRoleRequest{
		roleId,
	}, nil
}

func decodeRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req roleRequest
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

func decodeRolePermissionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req permissionIds
	//var req rolePermissionsRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	roleId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var permIds []int64
	for _, v := range req.PermissionIds {
		if permId, err := strconv.ParseInt(v, 10, 64); err == nil {
			permIds = append(permIds, permId)
		}
	}

	return rolePermissionsRequest{
		RoleId:        roleId,
		PermissionIds: permIds,
	}, nil
}
