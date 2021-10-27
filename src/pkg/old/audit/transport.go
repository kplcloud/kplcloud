/**
 * @Time : 2019/7/24 3:46 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package audit

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
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"net/http"
)

type endpoints struct {
	RefusedEndpoint     endpoint.Endpoint
	AccessAuditEndpoint endpoint.Endpoint
	AuditStepEndpoint   endpoint.Endpoint
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
		RefusedEndpoint:     makeRefusedEndpoint(svc),
		AuditStepEndpoint:   makeAduitStepEndpoint(svc),
		AccessAuditEndpoint: makeAccessAuditEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()),           // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Refused":   ems,
		"Access":    ems,
		"AuditStep": ems,
	}

	for _, m := range mw["Refused"] {
		eps.RefusedEndpoint = m(eps.RefusedEndpoint)
		eps.AccessAuditEndpoint = m(eps.AccessAuditEndpoint)
		eps.AuditStepEndpoint = m(eps.AuditStepEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/audit/{namespace}/name/{name}", kithttp.NewServer(
		eps.AccessAuditEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/audit/{namespace}/refused/{name}", kithttp.NewServer(
		eps.RefusedEndpoint,
		decodeDetailRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/audit/{namespace}/step/{name}/kind/{kind}", kithttp.NewServer(
		eps.AuditStepEndpoint,
		decodeStepRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")
	return r
}

func decodeDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	return detailRequest{
		Name:      name,
		Namespace: ns,
	}, nil
}

func decodeStepRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["namespace"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	name, ok := vars["name"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	kind, ok := vars["kind"]
	if !ok {
		return nil, encode.ErrBadRoute
	}
	return stepRequest{
		Kind: kind,
		detailRequest: detailRequest{
			Name:      name,
			Namespace: ns,
		},
	}, nil
}
