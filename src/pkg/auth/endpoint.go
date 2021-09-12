package auth

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
)

type authRequest struct {
	Username    string `json:"username" valid:"required"`
	Password    string `json:"password" valid:"required"`
	LoginType   string `json:"loginType,omitempty"`
	CaptchaId   string `json:"captchaId,omitempty"`
	CaptchaCode string `json:"captchaCode,omitempty"`
}

type Endpoints struct {
	LoginEndpoint         endpoint.Endpoint
	GithubEndpoint        endpoint.Endpoint
	GoogleEndpoint        endpoint.Endpoint
	AuthLoginTypeEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		LoginEndpoint: makeLoginEndpoint(s),
	}

	for _, m := range dmw["Login"] {
		eps.LoginEndpoint = m(eps.LoginEndpoint)
	}
	return eps
}

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(authRequest)
		tk, err := s.Login(ctx, req.Username, req.Password)
		return encode.Response{
			Data:  tk,
			Error: err,
		}, err
	}
}
