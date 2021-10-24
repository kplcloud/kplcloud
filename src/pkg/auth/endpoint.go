package auth

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
)

type authRequest struct {
	Username    string `json:"name" valid:"required"`
	Email       string `json:"email" valid:"required,email"`
	Password    string `json:"password" valid:"required"`
	LoginType   string `json:"loginType,omitempty"`
	Mobile      string `json:"mobile,omitempty"`
	Remark      string `json:"remark,omitempty"`
	CaptchaId   string `json:"captchaId,omitempty"`
	CaptchaCode string `json:"captchaCode,omitempty"`
}

type Endpoints struct {
	LoginEndpoint         endpoint.Endpoint
	GithubEndpoint        endpoint.Endpoint
	GoogleEndpoint        endpoint.Endpoint
	AuthLoginTypeEndpoint endpoint.Endpoint
	RegisterEndpoint      endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		LoginEndpoint:    makeLoginEndpoint(s),
		RegisterEndpoint: makeRegisterEndpoint(s),
	}

	for _, m := range dmw["Login"] {
		eps.LoginEndpoint = m(eps.LoginEndpoint)
	}
	for _, m := range dmw["Register"] {
		eps.RegisterEndpoint = m(eps.RegisterEndpoint)
	}
	return eps
}

func makeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(authRequest)
		err = s.Register(ctx, req.Username, req.Email, req.Password, req.Mobile, req.Remark)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(authRequest)
		tk, sessionTimeout, err := s.Login(ctx, req.Email, req.Password)
		return encode.Response{
			Data: map[string]interface{}{
				"sessionTimeout": sessionTimeout,
				"token":          tk,
			},
			Error: err,
		}, err
	}
}
