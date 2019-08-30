/**
 * @Time : 2019-07-17 14:19
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package member

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
	NamespacesEndpoint endpoint.Endpoint
	DetailEndpoint     endpoint.Endpoint
	PostEndpoint       endpoint.Endpoint
	UpdateEndpoint     endpoint.Endpoint
	ListEndpoint       endpoint.Endpoint
	MeEndpoint         endpoint.Endpoint
	AllEndpoint        endpoint.Endpoint
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
		NamespacesEndpoint: makeNamespacesEndpoint(svc),
		DetailEndpoint:     makeDetailEndpoint(svc),
		PostEndpoint:       makePostEndpoint(svc),
		UpdateEndpoint:     makeUpdateEndpoint(svc),
		MeEndpoint:         makeMeEndpoint(svc),
		ListEndpoint:       makeListEndpoint(svc),
		AllEndpoint:        makeAllEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Namespaces": ems,
		"Detail":     ems,
		"Post":       ems,
		"Update":     ems,
		"List":       ems,
		"Me":         ems,
		"All":        ems,
	}

	for _, m := range mw["Namespaces"] {
		eps.NamespacesEndpoint = m(eps.NamespacesEndpoint)
	}
	for _, m := range mw["Detail"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}
	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Me"] {
		eps.MeEndpoint = m(eps.MeEndpoint)
	}
	for _, m := range mw["All"] {
		eps.AllEndpoint = m(eps.AllEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/member/namespaces", kithttp.NewServer(
		eps.NamespacesEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/member", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListMemberRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/member/me", kithttp.NewServer(
		eps.MeEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/member/{id:[0-9]+}", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeGetMemberRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/member", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostMemberRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/member/{id:[0-9]+}", kithttp.NewServer(
		eps.UpdateEndpoint,
		decodeUpdateMemberRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/member/all", kithttp.NewServer(
		eps.AllEndpoint,
		decodeAllRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeGetMemberRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	memberId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	return getMemberRequest{Id: memberId}, nil
}

func decodePostMemberRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := memberRequest{
		State: 1,
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateMemberRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := memberRequest{
		State: 1,
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	memberId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	req.Id = memberId

	return req, nil
}

func decodeListMemberRequest(_ context.Context, r *http.Request) (interface{}, error) {
	pageStr := r.URL.Query().Get("p")
	//limitStr := r.URL.Query().Get("limit")
	email := r.URL.Query().Get("email")

	page, _ := strconv.Atoi(pageStr)
	//limit, _ := strconv.Atoi(limitStr)

	return listMemberRequest{
		Page:  page,
		Limit: 10,
		Email: email,
	}, nil
}

func decodeAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
