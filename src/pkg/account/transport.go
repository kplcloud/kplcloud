package account

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"io/ioutil"
	"net/http"
)

type endpoints struct {
	DetailEndpoint        endpoint.Endpoint
	GetReceiveEndpoint    endpoint.Endpoint
	GetProjectEndpoint    endpoint.Endpoint
	UpdateReceiveEndpoint endpoint.Endpoint
	UpdateBaseEndpoint    endpoint.Endpoint
	UnWechatBindEndpoint  endpoint.Endpoint
}

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	eps := endpoints{
		DetailEndpoint:        makeDetailEndpoint(svc),
		GetReceiveEndpoint:    makeGetReceiveEndpoint(svc),
		UpdateReceiveEndpoint: makeUpdateReceiveEndpoint(svc),
		UpdateBaseEndpoint:    makeUpdateBaseEndpoint(svc),
		UnWechatBindEndpoint:  makeUnWechatBindEndpoint(svc),
		GetProjectEndpoint:    makeGetProjectEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		//middleware.ProjectMiddleware(logger, projectRepository),                                   // 4
		//middleware.NamespaceMiddleware(logger),                                                    // 3
		middleware.CheckAuthMiddleware(logger),                                                    // 2
		kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Detail":        ems,
		"GetReceive":    ems,
		"UpdateReceive": ems,
		"UpdateBase":    ems,
		"UnWechatBind":  ems,
		"GetProject":    ems,
	}

	for _, m := range mw["Detail"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}
	for _, m := range mw["GetReceive"] {
		eps.GetReceiveEndpoint = m(eps.GetReceiveEndpoint)
	}
	for _, m := range mw["UpdateReceive"] {
		eps.UpdateReceiveEndpoint = m(eps.UpdateReceiveEndpoint)
	}
	for _, m := range mw["UpdateBase"] {
		eps.UpdateBaseEndpoint = m(eps.UpdateBaseEndpoint)
	}
	for _, m := range mw["UnWechatBind"] {
		eps.UnWechatBindEndpoint = m(eps.UnWechatBindEndpoint)
	}
	for _, m := range mw["GetProject"] {
		eps.GetProjectEndpoint = m(eps.GetProjectEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/account/current", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeCountReadRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/account/notice/receive", kithttp.NewServer(
		eps.GetReceiveEndpoint,
		decodeGetReceiveRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/account/notice/update", kithttp.NewServer(
		eps.UpdateReceiveEndpoint,
		decodeUpdateReceiveRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/account/base/update", kithttp.NewServer(
		eps.UpdateBaseEndpoint,
		decodeUpdateBaseRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/account/unBindWechat", kithttp.NewServer(
		eps.UnWechatBindEndpoint,
		decodeUnWechatBindRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/account/project", kithttp.NewServer(
		eps.GetProjectEndpoint,
		decodeGetProjectRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeCountReadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return accountRequest{MemberId: 0}, nil
}

func decodeGetReceiveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeUpdateReceiveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req accountReceiveRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateBaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req accountBaseRequest
	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUnWechatBindRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetProjectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
