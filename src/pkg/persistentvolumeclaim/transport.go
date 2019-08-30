/**
 * @Time : 2019-06-26 14:49
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package persistentvolumeclaim

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
	DeleteEndpoint endpoint.Endpoint
	PostEndpoint   endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	AllEndpoint    endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		SyncEndpoint:   makeSyncEndpoint(svc),
		GetEndpoint:    makeGetEndpoint(svc),
		DeleteEndpoint: makeDeleteEndpoint(svc),
		PostEndpoint:   makePostEndpoint(svc),
		ListEndpoint:   makeListEndpoint(svc),
		AllEndpoint:    makeAllEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Sync":   ems,
		"Get":    ems,
		"Delete": ems,
		"Post":   ems,
		"List":   ems,
	}

	for _, m := range mw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range mw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["All"] {
		eps.AllEndpoint = m(eps.AllEndpoint)
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

	deletePvc := kithttp.NewServer(
		eps.DeleteEndpoint,
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

	list := kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)

	all := kithttp.NewServer(
		eps.AllEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/persistentvolumeclaim/{namespace}/sync", sync).Methods("GET")
	r.Handle("/persistentvolumeclaim/{namespace}/pvc/{name}", get).Methods("GET")
	r.Handle("/persistentvolumeclaim/{namespace}/pvc/{name}", deletePvc).Methods("DELETE")
	r.Handle("/persistentvolumeclaim/{namespace}", post).Methods("POST")
	r.Handle("/persistentvolumeclaim/{namespace}", list).Methods("GET")
	r.Handle("/persistentvolumeclaim/{namespace}/all", all).Methods("GET")

	return r
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	return pvcRequest{Namespace: ns}, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	return listRequest{getRequest: getRequest{Namespace: ns}, Page: page, Limit: pageSize}, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	return pvcRequest{Namespace: ns, Name: name}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pvcRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// todo 需要校验参数是否正确

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	req.Storage = req.Storage + req.Unit.String()

	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	req.Namespace = ns

	return req, nil
}
