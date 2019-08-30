/**
 * @Time : 2019-07-17 19:03
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package persistentvolume

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
	GetEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
	}

	eps := endpoints{
		GetEndpoint: makeGetEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Get": ems,
	}

	for _, m := range mw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/persistentvolume/{namespace}/pv/{pvName}", kithttp.NewServer(
		eps.GetEndpoint,
		decodeGetRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["pvName"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	return getPvRequest{
		name,
	}, nil
}
