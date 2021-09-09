package auth

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kplcloud/kplcloud/src/encode"
	"net/http"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	var ems = []endpoint.Middleware{
		//checkCaptchaMiddleware(captchaSvc),
	}
	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Login": ems,
	})

	r := mux.NewRouter()

	r.Handle("/login", kithttp.NewServer(
		eps.LoginEndpoint,
		decodeLoginRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}
