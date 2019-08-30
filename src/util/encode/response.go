package encode

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/auth/casbin"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/kplcloud/kplcloud/src/middleware"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Err  error       `json:"error,omitempty"`
}

func (r Response) error() error { return r.Err }

var (
	ErrBadRoute      = errors.New("bad route")
	ErrToken         = errors.New("Token 失效, 请重新登陆")
	ErrParamsRefused = errors.New("参数校验未通过")
	ErrCasbin        = errors.New("您无权访问该模块,请联系管理员添加权限")
)

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		EncodeError(ctx, e.error(), w)
		return nil
	}
	resp := response.(Response)
	if resp.Err != nil {
		resp.Code = -1
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case kitjwt.ErrTokenContextMissing, kitjwt.ErrTokenExpired:
		err = ErrToken
		w.WriteHeader(http.StatusUnauthorized)
	case casbin.ErrUnauthorized:
		err = ErrCasbin
		w.WriteHeader(http.StatusForbidden)
	case middleware.ErrorASD:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  -1,
		"error": err.Error(),
	})
}
