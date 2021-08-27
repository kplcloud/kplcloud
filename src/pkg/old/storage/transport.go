/**
 * @Time : 2019-06-25 19:24
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package storage

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
	SyncEndpoint   endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
	PostEndpoint   endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
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
		SyncEndpoint:   makeSyncEndpoint(svc),
		GetEndpoint:    makeGetEndpoint(svc),
		PostEndpoint:   makePostEndpoint(svc),
		DeleteEndpoint: makeDeleteEndpoint(svc),
		ListEndpoint:   makeListEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Sync":   ems,
		"Get":    ems,
		"Post":   ems,
		"Delete": ems,
		"List":   ems,
	}

	for _, m := range mw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range mw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}

	sync := kithttp.NewServer(
		eps.SyncEndpoint,
		decodeRequest,
		encode.EncodeResponse,
		opts...,
	)

	get := kithttp.NewServer(
		eps.GetEndpoint,
		decodeGetRequest,
		encode.EncodeResponse,
		opts...,
	)

	post := kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)

	deleteHandle := kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeGetRequest,
		encode.EncodeResponse,
		opts...,
	)

	list := kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/storageclass/sync/all", sync).Methods("GET")
	r.Handle("/storageclass/{name}", get).Methods("GET")
	r.Handle("/storageclass/{name}", deleteHandle).Methods("DELETE")
	r.Handle("/storageclass", post).Methods("POST")
	r.Handle("/storageclass", list).Methods("GET")

	return r
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return storageRequest{}, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 15
	}
	return storageListRequest{Offset: offset, Limit: limit}, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	return storageRequest{Name: name}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req storageRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// todo 需要校验参数是否正确

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}
