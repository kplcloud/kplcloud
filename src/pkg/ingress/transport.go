/**
 * @Time : 2019/7/2 2:11 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package ingress

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
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"io/ioutil"
	"net/http"
	"strconv"
)

type endpoints struct {
	GetEndpoint       endpoint.Endpoint
	ListEndpoint      endpoint.Endpoint
	PostEndpoint      endpoint.Endpoint
	NoIngressEndpoint endpoint.Endpoint
	SyncEndpoint      endpoint.Endpoint
	GenerateEndpoint  endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger, repository repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		GetEndpoint:       makeGetEndpoint(svc),
		ListEndpoint:      makeListEndpoint(svc),
		PostEndpoint:      makePostEndpoint(svc),
		NoIngressEndpoint: makeGetNoIngressProjectEndpoint(svc),
		SyncEndpoint:      makeSyncEndpoint(svc),
		GenerateEndpoint:  makeGenerateEndpoint(svc),
	}
	ems := []endpoint.Middleware{
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}
	pEms := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()),
		middleware.NamespaceMiddleware(logger),
		middleware.CheckAuthMiddleware(logger),
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory),
	}
	mw := map[string][]endpoint.Middleware{
		"Get":       pEms,
		"List":      ems,
		"Post":      pEms,
		"NoIngress": ems,
		"Sync":      ems,
		"Generate":  pEms,
	}
	for _, m := range mw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["NoIngress"] {
		eps.NoIngressEndpoint = m(eps.NoIngressEndpoint)
	}
	for _, m := range mw["Sync"] {
		eps.SyncEndpoint = m(eps.SyncEndpoint)
	}
	for _, m := range mw["Generate"] {
		eps.GenerateEndpoint = m(eps.GenerateEndpoint)
	}

	get := kithttp.NewServer(eps.GetEndpoint, decodeGetRequest, encode.EncodeResponse, opts...)
	list := kithttp.NewServer(eps.ListEndpoint, decodeListRequest, encode.EncodeResponse, opts...)
	post := kithttp.NewServer(eps.PostEndpoint, decodePostRequest, encode.EncodeResponse, opts...)
	noIngress := kithttp.NewServer(eps.NoIngressEndpoint, decodeSyncRequest, encode.EncodeResponse, opts...)
	sync := kithttp.NewServer(eps.SyncEndpoint, decodeSyncRequest, encode.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/ingress/{namespace}/detail/{name}", get).Methods("GET")
	r.Handle("/ingress/{namespace}", list).Methods("GET")
	r.Handle("/ingress/{namespace}/detail/{name}", post).Methods("POST")
	r.Handle("/ingress/{namespace}/project", noIngress).Methods("GET")
	r.Handle("/ingress/{namespace}/sync", sync).Methods("GET")
	r.Handle("/ingress/{namespace}/generate/{name}", kithttp.NewServer(
		eps.GenerateEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...)).Methods("PUT")
	return r
}

func decodeGenerateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	name, nameOk := vars["name"]
	if !ok || !nameOk {
		return nil, encode.ErrBadRoute
	}
	return getRequest{
		Namespace: ns,
		Name:      name,
	}, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	pageStr := r.URL.Query().Get("p")
	limitStr := r.URL.Query().Get("limit")
	p, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 15
	}
	return listRequest{
		getRequest: getRequest{Namespace: ns},
		Page:       p,
		Limit:      limit,
	}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	name, nameOk := vars["name"]
	if !ok || !nameOk {
		return nil, encode.ErrBadRoute
	}

	// get request body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req postRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	req.Namespace = ns
	req.Name = name
	return req, nil
}

func decodeSyncRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	return getRequest{Namespace: ns}, nil
}
