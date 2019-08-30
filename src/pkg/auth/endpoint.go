package auth

import (
	"context"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/util/encode"
)

type authRequest struct {
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	LoginType string `json:"login_type,omitempty"`
}

type authResponse struct {
	//Success bool   `json:"success"`
	Code  int                    `json:"code"`
	Token string                 `json:"token"`
	Data  map[string]interface{} `json:"data"`
	Err   error                  `json:"error"`
}

func (r authResponse) error() error { return r.Err }

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(authRequest)
		rs, err := s.Login(ctx, req.Username, req.Password)
		var code int
		if err != nil {
			code = -1
		}
		res := map[string]interface{}{}
		if code == 0 {
			ctx = context.WithValue(ctx, jwt.JWTTokenContextKey, rs)
			res, err = s.ParseToken(ctx, rs)
			res["token"] = rs
		}
		return authResponse{Code: code, Err: err, Data: res, Token: rs}, err
	}
}

func makeAuthLoginTypeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		rs := s.AuthLoginType(ctx)
		return encode.Response{Code: 0, Err: nil, Data: rs}, nil
	}
}
