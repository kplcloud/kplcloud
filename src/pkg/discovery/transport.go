/**
 * @Time : 2019-07-08 10:42
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package discovery

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
	DeleteEndpoint endpoint.Endpoint
	DetailEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	CreateEndpoint endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
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
		DeleteEndpoint: makeDeleteEndpoint(svc),
		DetailEndpoint: makeDetailEndpoint(svc),
		ListEndpoint:   makeListEndpoint(svc),
		CreateEndpoint: makeCreateEndpoint(svc),
		UpdateEndpoint: makeUpdateEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		//middleware.ProjectMiddleware(logger, projectRepository),                                   // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Delete": ems,
		"Detail": ems,
		"List":   ems,
		"Create": ems,
	}

	for _, m := range mw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mw["Detail"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Create"] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
	}
	for _, m := range mw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/discovery/{namespace}/services/{serviceName}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeGetRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("DELETE")

	r.Handle("/discovery/{namespace}/services/{serviceName}", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeGetRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/discovery/{namespace}/services", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/discovery/{namespace}/services", kithttp.NewServer(
		eps.CreateEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/discovery/{namespace}/services/{serviceName}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	return r
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createRequest
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

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	return listRequest{
		page,
		limit,
	}, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	serviceName, ok := vars["serviceName"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	return getRequest{
		ServiceName: serviceName,
	}, nil
}
