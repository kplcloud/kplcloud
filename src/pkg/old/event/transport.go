/**
 * @Time : 2019-07-22 14:25
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package event

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
	AllEndpoint endpoint.Endpoint
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
		AllEndpoint: makeAllEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"All": ems,
	}

	for _, m := range mw["All"] {
		eps.AllEndpoint = m(eps.AllEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/event/all", kithttp.NewServer(
		eps.AllEndpoint,
		func(ctx context.Context, req *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}
