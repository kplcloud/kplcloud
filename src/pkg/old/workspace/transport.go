/**
 * @Time : 2019-07-25 15:10
 * @Author : soupzhb@gmail.com
 * @File : transport.go
 * @Software: GoLand
 */

package workspace

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"net/http"
)

var errBadRoute = errors.New("bad route")

type endpoints struct {
	MetriceEndpoint endpoint.Endpoint
	ActiveEndpoint  endpoint.Endpoint
}

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		MetriceEndpoint: makeMetriceEndpoint(svc),
		ActiveEndpoint:  makeActiveEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		//middleware.ProjectMiddleware(logger, projectRepository),                                   // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Metrice": ems,
		"Active":  ems,
	}

	for _, m := range mw["Metrice"] {
		eps.MetriceEndpoint = m(eps.MetriceEndpoint)
	}
	for _, m := range mw["Active"] {
		eps.ActiveEndpoint = m(eps.ActiveEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/workspace/{namespace}/metrics", kithttp.NewServer(
		eps.MetriceEndpoint,
		decodeMetriceRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/workspace/{namespace}/active", kithttp.NewServer(
		eps.ActiveEndpoint,
		decodeActiveRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeMetriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, errBadRoute
	}
	return nsRequest{
		Namespace: ns,
	}, nil
}

func decodeActiveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
