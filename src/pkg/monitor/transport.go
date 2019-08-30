/**
 * @Time : 2019-07-29 15:22
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package monitor

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"net/http"
)

type endpoints struct {
	QueryNetworkEndpoint endpoint.Endpoint
	OpsEndpoint          endpoint.Endpoint
	MetricsEndpoint      endpoint.Endpoint
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
		QueryNetworkEndpoint: makeQueryNetworkEndpoint(svc),
		OpsEndpoint:          makeOpsEndpoint(svc),
		MetricsEndpoint:      makeMetricsEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"QueryNetwork": ems,
		"Ops":          ems,
		"Metrics":      ems,
	}

	for _, m := range mw["QueryNetwork"] {
		eps.QueryNetworkEndpoint = m(eps.QueryNetworkEndpoint)
	}
	for _, m := range mw["Ops"] {
		eps.OpsEndpoint = m(eps.OpsEndpoint)
	}
	for _, m := range mw["Metrics"] {
		eps.MetricsEndpoint = m(eps.MetricsEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/monitor/prometheus/query/network", kithttp.NewServer(
		eps.QueryNetworkEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/monitor/prometheus/query/ops", kithttp.NewServer(
		eps.OpsEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/monitor/metrics", kithttp.NewServer(
		eps.MetricsEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}
