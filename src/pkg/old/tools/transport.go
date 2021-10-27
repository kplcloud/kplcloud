/**
 * @Time : 2019-07-23 18:51
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package tools

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
)

type endpoints struct {
	FakeTimeEndpoint    endpoint.Endpoint
	DuplicationEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger, repository repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
		kithttp.ServerBefore(middleware.NamespaceToContext()),
	}

	eps := endpoints{
		FakeTimeEndpoint:    makeFakeTimeEndpoint(svc),
		DuplicationEndpoint: makeDuplicationEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.ProjectMiddleware(logger, repository.Project(), repository.Groups()),           // 4
		middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"FakeTime":    ems,
		"Duplication": ems[2:],
	}

	for _, m := range mw["FakeTime"] {
		eps.FakeTimeEndpoint = m(eps.FakeTimeEndpoint)
	}
	for _, m := range mw["Duplication"] {
		eps.DuplicationEndpoint = m(eps.DuplicationEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/tools/faketime/{namespace}/project/{name}", kithttp.NewServer(
		eps.FakeTimeEndpoint,
		decodeFakeTimeRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/tools/duplication/single", kithttp.NewServer(
		eps.DuplicationEndpoint,
		decodeDuplicationRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	return r
}

func decodeFakeTimeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req fakeTimeRequest

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	if req.FakeTime.Second() < 0 {
		return nil, ErrToolFakeTimeErr
	}

	return req, nil
}

func decodeDuplicationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req duplicationRequest

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}
