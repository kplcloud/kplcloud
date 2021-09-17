package auth

import (
	"context"
	"encoding/json"
	valid "github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	captcha "github.com/icowan/kit-captcha"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/pkg/errors"
	"net/http"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption, captchaSvc captcha.Service) http.Handler {
	var ems = []endpoint.Middleware{
		checkCaptchaMiddleware(captchaSvc),
	}
	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Login":    ems,
		"Register": ems, // 验证码的中间件
	})

	r := mux.NewRouter()

	r.Handle("/login", kithttp.NewServer(
		eps.LoginEndpoint,
		decodeLoginRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/register", kithttp.NewServer(
		eps.RegisterEndpoint,
		decodeRegisterRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	validResult, err := valid.ValidateStruct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if !validResult {
		return nil, encode.InvalidParams.Wrap(errors.New("valid false"))
	}
	return req, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	return req, nil
}
