/**
 * @Time : 2019-06-20 17:52
 * @Author : solacowa@gmail.com
 * @File : namespace
 * @Software: GoLand
 */

package namespace

import (
	"context"
	"encoding/json"
	"errors"
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

var errBadRoute = errors.New("bad route")

//type NamespaceEndpoints struct {
//	GetEndpoint  endpoint.Endpoint
//	PostEndpoint endpoint.Endpoint
//}

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerBefore(middleware.CasbinToContext()),
	}

	epsMap := map[string]endpoint.Endpoint{
		"get":    makeGetEndpoint(svc),
		"post":   makePostEndpoint(svc),
		"sync":   makeSyncEndpoint(svc),
		"update": makeUpdateEndpoint(svc),
		"list":   makeListEndpoint(svc),
	}

	for key, val := range epsMap {
		epsMap[key] = kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory)(val)
		epsMap[key] = middleware.CheckAuthMiddleware(logger)(val)
	}

	get := kithttp.NewServer(
		epsMap["get"],
		decodeGetRequest,
		encodeResponse,
		opts...,
	)

	create := kithttp.NewServer(
		epsMap["post"],
		decodePostRequest,
		encodeResponse,
		opts...,
	)

	sync := kithttp.NewServer(
		epsMap["sync"],
		decodeSyncRequest,
		encodeResponse,
		opts...,
	)

	update := kithttp.NewServer(
		epsMap["update"],
		decodeGetRequest,
		encodeResponse,
		opts...,
	)

	list := kithttp.NewServer(
		epsMap["list"],
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/namespace/{name}", get).Methods("GET")
	r.Handle("/namespace/{name}", update).Methods("PUT")
	r.Handle("/namespace", create).Methods("POST")
	r.Handle("/namespace", list).Methods("GET")
	r.Handle("/namespace/sync/all", sync).Methods("GET")

	//r.Handle("/namespace/{name}", ns).Methods("DELETE")

	return r
}

func decodeSyncRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nsRequest{}, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, errBadRoute
	}
	return nsRequest{
		Name: name,
	}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req nsRequest

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case kitjwt.ErrTokenContextMissing, kitjwt.ErrTokenExpired, middleware.ErrorASD:
		w.WriteHeader(http.StatusForbidden)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  -1,
		"error": err.Error(),
	})
}
