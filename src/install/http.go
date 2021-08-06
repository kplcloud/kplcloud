/**
 * @Time : 7/21/21 2:26 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package install

import (
	"context"
	"encoding/json"
	"github.com/kplcloud/kplcloud/src/util"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"InitDb": ems,
	})

	r := mux.NewRouter()

	r.Handle("/init-db", kithttp.NewServer(
		eps.InitDbEndpoint,
		decodeInitDbRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/init-platform", kithttp.NewServer(
		eps.InitPlatformEndpoint,
		decodeInitPlatformRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/init-logo", kithttp.NewServer(
		eps.InitLogoEndpoint,
		decodeInitLogoPlatformRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/init-cors", kithttp.NewServer(
		eps.InitCorsEndpoint,
		decodeInitCorsRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/init-redis", kithttp.NewServer(
		eps.InitRedisEndpoint,
		decodeInitRedisRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodeInitRedisRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req initRedisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
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

func decodeInitCorsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req initCorsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Allow {
		validResult, err := valid.ValidateStruct(req)
		if err != nil {
			return nil, encode.InvalidParams.Wrap(err)
		}
		if !validResult {
			return nil, encode.InvalidParams.Wrap(errors.New("valid false"))
		}
		req.Method = strings.Join(req.Methods, ",")
	}

	return req, nil
}

func decodeInitDbRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req initDbRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
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

func decodeInitPlatformRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req initPlatformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	validResult, err := valid.ValidateStruct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if !validResult {
		return nil, encode.InvalidParams.Wrap(errors.New("valid false"))
	}
	req.AppKey = string(util.Krand(12, util.KC_RAND_KIND_ALL))
	return req, nil
}

func decodeInitLogoPlatformRequest(_ context.Context, r *http.Request) (interface{}, error) {
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, errors.Wrap(err, "r.MultipartReader")
	}

	form, err := reader.ReadForm(32 << 10)
	if err != nil {
		return nil, errors.Wrap(err, "reader.ReadForm")
	}

	if form.File == nil || form.File["logo"] == nil {
		return nil, errors.New("请选择所需要上传的logo")
	}

	return initLogoRequest{Files: form.File["logo"]}, nil

}
